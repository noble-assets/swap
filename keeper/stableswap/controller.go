package stableswap

import (
	"context"
	"fmt"
	"time"

	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	anyproto "github.com/cosmos/gogoproto/types/any"
	"swap.noble.xyz/types"
	"swap.noble.xyz/types/stableswap"
	stableswaptypes "swap.noble.xyz/types/stableswap"
)

// Controller for `StableSwap` Pools.
type Controller struct {
	baseDenom  string
	bankKeeper *types.BankKeeper
	pool       *types.Pool
	paused     bool

	stableswapPool   *stableswap.Pool
	stableswapKeeper *Keeper
}

// NewController initializes and returns the appropriate controller for a `StableSwap` Pool.
func NewController(
	bankKeeper *types.BankKeeper,
	baseDenom string,
	pool *types.Pool,
	paused bool,
	stableswapPool *stableswap.Pool,
	stableswapKeeper *Keeper,
) Controller {
	return Controller{
		baseDenom:        baseDenom,
		bankKeeper:       bankKeeper,
		stableswapPool:   stableswapPool,
		paused:           paused,
		pool:             pool,
		stableswapKeeper: stableswapKeeper,
	}
}

func (c *Controller) GetId() uint64 {
	return c.pool.Id
}

func (c *Controller) GetAlgorithm() types.Algorithm {
	return c.pool.Algorithm
}

func (c *Controller) GetPair() string {
	return c.pool.Pair
}

func (c *Controller) GetAddress() string {
	return c.pool.Address
}

// PoolDetails returns the underlying detailed information about the `StableSwap` Pool as an `Any` object.
func (c *Controller) PoolDetails() *anyproto.Any {
	details, _ := anyproto.NewAnyWithCacheWithValue(c.stableswapPool)
	return details
}

// IsPaused returns true if the Pool is paused.
func (c *Controller) IsPaused() bool {
	return c.paused
}

// Swap computes a token swap within the Pool using the `StableSwap` algorithm, returning a commitment.
func (c *Controller) Swap(
	ctx context.Context,
	currentTime int64,
	coin sdk.Coin,
	denomTo string,
) (*types.SwapCommitment, error) {
	// Ensure that the Pool has liquidity.
	poolLiquidity := c.GetLiquidity(ctx)
	if !poolLiquidity.IsAllPositive() {
		return nil, fmt.Errorf("pool liquidity must be positive")
	}

	// Calculate the liquidity adjusted to the pool rates.
	adjustedLiquidity, err := calculateAdjustedBalancesInRates(c.stableswapPool.RateMultipliers, poolLiquidity)
	if err != nil {
		return nil, err
	}

	// Get the current amplification coefficient.
	amp := getAmplificationCoefficient(
		currentTime,
		math.LegacyNewDec(c.stableswapPool.InitialA),
		math.LegacyNewDec(c.stableswapPool.FutureA),
		c.stableswapPool.InitialATime,
		c.stableswapPool.FutureATime,
	)

	// Calculate new token balance after adding input amount
	x := computeNewAdjustedBalance(
		adjustedLiquidity.AmountOf(coin.Denom),
		coin.Amount.ToLegacyDec(),
		c.stableswapPool.RateMultipliers.AmountOf(coin.Denom).ToLegacyDec(),
		DecimalPrecision,
	)

	// Compute the result of the swap.
	swapResult, err := performSwap(
		sdk.NewCoin(coin.Denom, x.TruncateInt()),
		adjustedLiquidity,
		amp,
		denomTo,
		c.stableswapPool.RewardsFee,
		c.stableswapPool.ProtocolFeePercentage,
		c.stableswapPool.MaxFee,
		c.stableswapPool.RateMultipliers,
	)
	if err != nil {
		return nil, err
	}

	return &types.SwapCommitment{
		In:  sdk.NewCoin(coin.Denom, coin.Amount),
		Out: sdk.NewCoin(denomTo, swapResult.Dy.TruncateInt()),
		Fees: []types.Receiver{
			{ // protocol fees
				Amount:  sdk.NewCoin(coin.Denom, swapResult.ProtocolFee.TruncateInt()),
				Address: authtypes.NewModuleAddress(fmt.Sprintf("%s/pool/%d/protocol_fees", types.ModuleName, c.GetId())),
			},
			{ // rewards fees
				Amount:  sdk.NewCoin(coin.Denom, swapResult.RewardsFee.TruncateInt()),
				Address: authtypes.NewModuleAddress(fmt.Sprintf("%s/pool/%d/rewards_fees", types.ModuleName, c.GetId())),
			},
		},
	}, nil
}

// AddLiquidity creates a new user bonded position, providing liquidity to the given `StableSwap` Pool, setting the state.
func (c *Controller) AddLiquidity(
	ctx context.Context,
	currentTime time.Time,
	msg *stableswap.MsgAddLiquidity,
) (*types.AddLiquidityCommitment, error) {
	// Get current amplification coefficient.
	amp := getAmplificationCoefficient(
		currentTime.Unix(),
		math.LegacyNewDec(c.stableswapPool.InitialA),
		math.LegacyNewDec(c.stableswapPool.FutureA),
		c.stableswapPool.InitialATime,
		c.stableswapPool.FutureATime,
	)

	// Retrieve the Pool liquidity.
	liquidity := c.GetLiquidity(ctx)

	// Compute the total coins amount.
	totalSupply := math.LegacyZeroDec()
	for _, coin := range liquidity {
		totalSupply = totalSupply.Add(liquidity.AmountOf(coin.Denom).ToLegacyDec())
	}

	// Calculate pre-deposit invariant.
	xp, err := calculateAdjustedBalancesInRates(c.stableswapPool.RateMultipliers, liquidity)
	if err != nil {
		return nil, err
	}
	D0, err := calculateInvariant(xp, amp)
	if err != nil {
		return nil, err
	}

	// Calculate new liquidity balances after adding the user provided amounts.
	for _, coin := range msg.Amount {
		liquidity = liquidity.Add(sdk.NewCoin(coin.Denom, math.NewInt(liquidity.AmountOf(coin.Denom).Int64()+msg.Amount.AmountOf(coin.Denom).Int64())))
	}

	// Calculate post-deposit invariant.
	xp, err = calculateAdjustedBalancesInRates(c.stableswapPool.RateMultipliers, liquidity)
	if err != nil {
		return nil, err
	}
	D1, err := calculateInvariant(xp, amp)
	if err != nil {
		return nil, err
	}
	if D1.LTE(D0) {
		return nil, fmt.Errorf("d1 must be greater than d0")
	}

	// Calculate how many LP tokens to mint.
	newTotalMint := math.LegacyZeroDec()
	if totalSupply.GT(math.LegacyZeroDec()) {
		newTotalMint = totalSupply.Mul(D1.Sub(D0))
		newTotalMint = newTotalMint.Quo(D0)
	} else {
		// First liquidity provider mints tokens equal to D1.
		newTotalMint.Set(D1)
	}

	bondedPosition := stableswaptypes.BondedPosition{
		Balance:            newTotalMint.Sub(c.stableswapPool.TotalShares),
		Timestamp:          currentTime,
		RewardsPeriodStart: currentTime,
	}

	// Add the new BondedPosition on the StableSwap State.
	if c.stableswapKeeper.HasBondedPosition(ctx, msg.PoolId, msg.Signer, currentTime.Unix()) {
		return nil, fmt.Errorf("cannot create multiple user positions in the same block")
	}
	if err = c.stableswapKeeper.SetBondedPosition(ctx, msg.PoolId, msg.Signer, currentTime.Unix(), bondedPosition); err != nil {
		return nil, sdkerrors.Wrap(err, "unable to set updated user pool position")
	}

	// Update the new user total bonded shares on the State.
	newUserTotal := math.LegacyZeroDec()
	if c.stableswapKeeper.HasUserTotalBondedShares(ctx, c.GetId(), msg.Signer) {
		currentUserTotal := c.stableswapKeeper.GetUserTotalBondedShares(ctx, c.GetId(), msg.Signer)
		newUserTotal = newUserTotal.Add(currentUserTotal)
	}
	newUserTotal = newUserTotal.Add(bondedPosition.Balance)
	if err = c.stableswapKeeper.SetUserTotalBondedShares(ctx, c.GetId(), msg.Signer, newUserTotal); err != nil {
		return nil, sdkerrors.Wrap(err, "unable to set stableswap user pool total bonded shares")
	}

	// Update the Pool TotalShares on the State.
	c.stableswapPool.TotalShares = c.stableswapPool.TotalShares.Add(bondedPosition.Balance)
	if c.GetLiquidity(ctx).IsZero() {
		c.stableswapPool.InitialRewardsTime = currentTime
	}
	if err = c.stableswapKeeper.SetPool(ctx, msg.PoolId, *c.stableswapPool); err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to set stableswap pool total bonded shares")
	}

	return &types.AddLiquidityCommitment{
		BondedPosition: bondedPosition,
	}, nil
}

// RemoveLiquidity begins the process for a user to remove liquidity positions from a given `StableSwap` Pool and
// a given percentage amount in the Pool initiating the unbonding period.
func (c *Controller) RemoveLiquidity(
	ctx context.Context,
	currentTime time.Time,
	msg *stableswaptypes.MsgRemoveLiquidity,
) (*types.RemoveLiquidityCommitment, error) {
	// Get the user total shares in the Pool.
	if !c.stableswapKeeper.HasUserTotalBondedShares(ctx, c.GetId(), msg.Signer) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidUnbondPosition, "user %s does not have a UsersTotalBondedShares position", msg.Signer)
	}
	userTotalShares := c.stableswapKeeper.GetUserTotalBondedShares(ctx, c.GetId(), msg.Signer)

	prevUserTotalUnbondingShares := math.LegacyZeroDec()
	if c.stableswapKeeper.HasUserTotalUnbondingShares(ctx, c.GetId(), msg.Signer) {
		userTotalUnbondingShares := c.stableswapKeeper.GetUserTotalUnbondingShares(ctx, c.GetId(), msg.Signer)
		prevUserTotalUnbondingShares = userTotalUnbondingShares
	}

	// Compute the user remaining bonded shares in the Pool, available to be unbonded.
	availableShares := userTotalShares.Sub(prevUserTotalUnbondingShares)

	// Compute the unbonding shares by the percentage.
	sharesToUnbond := availableShares.Mul(msg.Percentage).QuoInt64(100)
	if !sharesToUnbond.IsPositive() {
		return nil, types.ErrInvalidUnbondAmount
	}

	// Get the Pool liquidity.
	liquidity := c.GetLiquidity(ctx)

	// Calculate the proportional amount of each asset in the pool to return to the user.
	coinsToReturn := sdk.NewCoins()
	for _, asset := range liquidity {
		// Determine the amount of tokens to return for this asset.
		amountToReturn := asset.Amount.Mul(sharesToUnbond.TruncateInt()).Quo(c.stableswapPool.TotalShares.TruncateInt())
		if !amountToReturn.IsPositive() {
			continue
		}

		// Create the coin representation of the token to return.
		coinsToReturn = coinsToReturn.Add(sdk.NewCoin(asset.Denom, amountToReturn))
	}

	// Calculate the proportional amount of each asset in the pool to return to the user.
	coinsToReturn = sdk.NewCoins()
	for _, asset := range liquidity {
		// Determine the amount of tokens to return for this asset.
		amountToReturn := asset.Amount.Mul(sharesToUnbond.TruncateInt()).Quo(c.stableswapPool.TotalShares.TruncateInt())
		if !amountToReturn.IsPositive() {
			continue
		}

		// Create the coin representation of the token to return.
		coinsToReturn = coinsToReturn.Add(sdk.NewCoin(asset.Denom, amountToReturn))
	}

	// Compute the unbonding period weighted to the amount of tokens to unbond and the total pool liquidity.
	unbondingPeriod, err := ComputeWeightedPoolUnbondingPeriod(c.stableswapPool.TotalShares, sharesToUnbond)
	if err != nil {
		return nil, err
	}
	unbondingEndTime := currentTime.Add(unbondingPeriod)

	// Create the unbonding Position.
	unbondingPosition := stableswaptypes.UnbondingPosition{
		Amount:  coinsToReturn,
		EndTime: unbondingEndTime,
		Shares:  sharesToUnbond,
	}

	// Add to the Unbonding queue on the State if a record does not already exist.
	if c.stableswapKeeper.HasUnbondingPosition(ctx, unbondingEndTime.Unix(), msg.Signer, c.GetId()) {
		return nil, fmt.Errorf("unbonding key already exists: %d-%s-%d", unbondingEndTime.Unix(), msg.Signer, c.GetId())
	}
	err = c.stableswapKeeper.SetUnbondingPosition(ctx, unbondingEndTime.Unix(), msg.Signer, c.GetId(), unbondingPosition)
	if err != nil {
		return nil, err
	}

	// Add the shares to the pool total unbonding shares on the State.
	totalPoolUnbondingShares := math.LegacyZeroDec()
	if c.stableswapKeeper.HasPoolTotalUnbondingShares(ctx, c.GetId()) {
		totalPoolUnbondingShares = c.stableswapKeeper.GetPoolTotalUnbondingShares(ctx, c.GetId())
	}
	if err = c.stableswapKeeper.SetPoolTotalUnbondingShares(ctx, c.GetId(), totalPoolUnbondingShares.Add(unbondingPosition.Shares)); err != nil {
		return nil, err
	}

	// Add the shares to the user total unbonding shares on the State.
	prevUnbondingShares := math.LegacyZeroDec()
	if c.stableswapKeeper.HasUserTotalUnbondingShares(ctx, c.GetId(), msg.Signer) {
		prevUnbondingShares = c.stableswapKeeper.GetUserTotalUnbondingShares(ctx, c.GetId(), msg.Signer)
	}
	if err = c.stableswapKeeper.SetUserTotalUnbondingShares(ctx, c.GetId(), msg.Signer, prevUnbondingShares.Add(unbondingPosition.Shares)); err != nil {
		return nil, err
	}

	// Calculate unbonding time and add to unbonding queue.
	return &types.RemoveLiquidityCommitment{
		UnbondingPosition: unbondingPosition,
	}, nil
}

// GetLiquidity returns the total liquidity of the `StableSwap` Pool.
func (c *Controller) GetLiquidity(ctx context.Context) sdk.Coins {
	address, err := sdk.AccAddressFromBech32(c.GetAddress())
	if err != nil {
		return sdk.Coins{}
	}

	// Get the liquidity of only the wanted tokens.
	liquidity := sdk.Coins{}
	liquidity = liquidity.Add((*c.bankKeeper).GetBalance(ctx, address, c.baseDenom))
	liquidity = liquidity.Add((*c.bankKeeper).GetBalance(ctx, address, c.GetPair()))
	return liquidity
}

// GetRates computes the exchange rates of the tokens in the `StableSwap` Pool.
func (c *Controller) GetRates(ctx context.Context) []types.Rate {
	liquidity := c.GetLiquidity(ctx)
	amount := liquidity.AmountOf(c.GetPair()).ToLegacyDec()
	vsAmount := liquidity.AmountOf(c.baseDenom).ToLegacyDec()

	price := math.LegacyZeroDec()
	vsPrice := math.LegacyZeroDec()
	if vsAmount.GT(math.LegacyZeroDec()) && amount.GT(math.LegacyZeroDec()) {
		price = vsAmount.Quo(amount)
		vsPrice = math.LegacyOneDec().Quo(price)
	}

	return []types.Rate{
		{
			Denom:     c.GetPair(),
			Vs:        c.baseDenom,
			Price:     price,
			Algorithm: c.GetAlgorithm(),
		}, {
			Denom:     c.baseDenom,
			Vs:        c.GetPair(),
			Price:     vsPrice,
			Algorithm: c.GetAlgorithm(),
		},
	}
}

func (c *Controller) GetRate(ctx context.Context) math.LegacyDec {
	liquidity := c.GetLiquidity(ctx)
	if liquidity.IsZero() {
		return math.LegacyOneDec()
	}
	vsAmount := liquidity.AmountOf(c.GetPair()).ToLegacyDec()
	amount := liquidity.AmountOf(c.baseDenom).ToLegacyDec()

	if vsAmount.LTE(math.LegacyZeroDec()) || amount.LTE(math.LegacyZeroDec()) {
		return math.LegacyZeroDec()
	}
	return vsAmount.Quo(amount)
}

// UpdatePool updates the StableSwap Pool params.
func (c *Controller) UpdatePool(
	ctx context.Context,
	protocolFeePercentage int64,
	rewardsFee int64,
	maxFee int64,
	initialA int64,
	initialATime int64,
	futureA int64,
	futureATime int64,
	rateMultipliers sdk.Coins,
) error {
	c.stableswapPool.ProtocolFeePercentage = protocolFeePercentage
	c.stableswapPool.RewardsFee = rewardsFee
	c.stableswapPool.InitialA = initialA
	c.stableswapPool.InitialATime = initialATime
	c.stableswapPool.MaxFee = maxFee
	c.stableswapPool.FutureA = futureA
	c.stableswapPool.FutureATime = futureATime
	c.stableswapPool.RateMultipliers = rateMultipliers

	// Update the `StableSwap` data on state.
	if err := c.stableswapKeeper.SetPool(ctx, c.GetId(), *c.stableswapPool); err != nil {
		return sdkerrors.Wrapf(err, "unable to set stableswap pool")
	}

	return nil
}

func (c *Controller) ProcessUnbondings(ctx context.Context, currentTime time.Time) error {
	// iterate over unbonding entries and process those whose unbonding period has ended
	for _, entry := range c.stableswapKeeper.GetUnbondingPositionsUntil(ctx, currentTime.Unix()) {
		// Check if the unbonding period has ended
		if currentTime.After(entry.UnbondingPosition.EndTime) {
			// Send tokens back to the user
			if err := (*c.bankKeeper).SendCoins(
				ctx,
				sdk.MustAccAddressFromBech32(c.GetAddress()),
				sdk.MustAccAddressFromBech32(entry.Address),
				entry.UnbondingPosition.Amount,
			); err != nil {
				return err
			}

			rewards, err := c.ProcessUserRewards(ctx, entry.Address, currentTime)
			if err != nil {
				return err
			}
			if rewards.Len() > 0 {
				if err = c.stableswapKeeper.eventService.EventManager(ctx).Emit(ctx, &types.WithdrawnRewards{
					Signer:  entry.Address,
					Rewards: rewards,
				}); err != nil {
					return err
				}
			}

			cumulativeUnbonded := math.LegacyZeroDec()
			// Iterate through user's positions to unbond the specified amount
			for _, bondedEntry := range c.stableswapKeeper.GetBondedPositionsByProvider(ctx, entry.Address) {
				if cumulativeUnbonded.Add(bondedEntry.BondedPosition.Balance).GT(entry.UnbondingPosition.Shares) {
					remainingShares := entry.UnbondingPosition.Shares.Sub(cumulativeUnbonded)
					bondedEntry.BondedPosition.Balance = entry.UnbondingPosition.Shares.Sub(remainingShares)
					cumulativeUnbonded = entry.UnbondingPosition.Shares

				} else {
					cumulativeUnbonded = cumulativeUnbonded.Add(bondedEntry.BondedPosition.Balance)
					if err := c.stableswapKeeper.RemoveBondedPosition(ctx, bondedEntry.PoolId, bondedEntry.Address, bondedEntry.Timestamp); err != nil {
						return err
					}
				}
			}

			// Final check to ensure the unbonded amount matches the target
			if cumulativeUnbonded.LT(entry.UnbondingPosition.Shares) {
				return fmt.Errorf("%s is smaller then requested: %s", cumulativeUnbonded.String(), entry.UnbondingPosition.Shares.String())
			}

			// Remove entry from the unbonding queue after processing it.
			if err := c.stableswapKeeper.RemoveUnbondingPosition(ctx, entry.Timestamp, entry.Address, entry.PoolId); err != nil {
				return err
			}

			// Update the pool total shares.
			c.stableswapPool.TotalShares = c.stableswapPool.TotalShares.Sub(entry.UnbondingPosition.Shares)
			if err := c.stableswapKeeper.SetPool(ctx, c.GetId(), *c.stableswapPool); err != nil {
				return err
			}

			// Remove the unbonded shares from the user total.
			userTotalBondedShares := math.LegacyZeroDec()
			if c.stableswapKeeper.HasUserTotalBondedShares(ctx, c.GetId(), entry.Address) {
				userTotalBondedShares = c.stableswapKeeper.GetUserTotalBondedShares(ctx, c.GetId(), entry.Address)
			}
			if err := c.stableswapKeeper.SetUserTotalBondedShares(ctx, c.GetId(), entry.Address, userTotalBondedShares.Sub(entry.UnbondingPosition.Shares)); err != nil {
				return err
			}

			// Remove the shares from the pool total unbonding shares.
			totalPoolUnbondingShares := math.LegacyZeroDec()
			if c.stableswapKeeper.HasPoolTotalUnbondingShares(ctx, c.GetId()) {
				totalPoolUnbondingShares = c.stableswapKeeper.GetPoolTotalUnbondingShares(ctx, c.GetId())
			}
			if err := c.stableswapKeeper.SetPoolTotalUnbondingShares(ctx, c.GetId(), totalPoolUnbondingShares.Sub(entry.UnbondingPosition.Shares)); err != nil {
				return err
			}

			// Remove the shares from the user total unbonding shares.
			userTotalUnbondingShares := math.LegacyZeroDec()
			if c.stableswapKeeper.HasUserTotalUnbondingShares(ctx, c.GetId(), entry.Address) {
				userTotalUnbondingShares = c.stableswapKeeper.GetUserTotalUnbondingShares(ctx, c.GetId(), entry.Address)
			}
			if err := c.stableswapKeeper.SetUserTotalUnbondingShares(ctx, c.GetId(), entry.Address, userTotalUnbondingShares.Sub(entry.UnbondingPosition.Shares)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Controller) GetTotalPoolUserRewards(ctx context.Context, address string, currentTime time.Time) ([]types.ReceiverMulti, error) {
	// Get the total Pool rewards.
	poolRewardsAddress := authtypes.NewModuleAddress(fmt.Sprintf("%s/pool/%d/rewards_fees", types.ModuleName, c.GetId()))
	poolRewards := (*c.bankKeeper).GetAllBalances(ctx, poolRewardsAddress)

	var userRewards []types.ReceiverMulti
	for _, entry := range c.stableswapKeeper.GetBondedPositionsByPoolAndProvider(ctx, c.GetId(), address) {
		// Calculate the adjusted rewards for the given user bonded position.
		rewards, err := CalculatePositionRewards(currentTime, poolRewards, entry.BondedPosition, c.stableswapPool.TotalShares, c.stableswapPool.InitialRewardsTime)
		if err != nil {
			return nil, err
		}
		if rewards.IsZero() {
			continue
		}

		cumulativeRewards := sdk.Coins{}
		cumulativeRewards = cumulativeRewards.Add(rewards...)

		userRewards = append(userRewards, types.ReceiverMulti{
			Amount:  cumulativeRewards,
			Address: poolRewardsAddress,
			PoolId:  entry.PoolId,
		})
	}
	return userRewards, nil
}

func (c *Controller) ProcessUserRewards(ctx context.Context, address string, currentTime time.Time) (sdk.Coins, error) {
	// Get the expected amount of user rewards.
	userRewards, err := c.GetTotalPoolUserRewards(ctx, address, currentTime)
	if err != nil {
		return sdk.Coins{}, err
	}

	// If the user does not have any reward exit.
	if len(userRewards) <= 0 {
		return sdk.Coins{}, nil
	}

	// Update the user rewards periods.
	for _, entry := range c.stableswapKeeper.GetBondedPositionsByPoolAndProvider(ctx, c.GetId(), address) {
		// update the entry of the rewards period with the current time
		entry.BondedPosition.RewardsPeriodStart = currentTime
		err = c.stableswapKeeper.SetBondedPosition(ctx, entry.PoolId, entry.Address, entry.Timestamp, entry.BondedPosition)
		if err != nil {
			return nil, err
		}
	}

	// Send the rewards to the user.
	finalRewards := sdk.Coins{}
	for _, poolRewards := range userRewards {
		for _, coin := range poolRewards.Amount {
			finalRewards = finalRewards.Add(coin)
		}
		err = (*c.bankKeeper).SendCoins(ctx, poolRewards.Address, sdk.MustAccAddressFromBech32(address), poolRewards.Amount)
		if err != nil {
			return nil, err
		}
	}

	return finalRewards, nil
}

func (c *Controller) GetProtocolFeesAddresses() []sdk.AccAddress {
	return []sdk.AccAddress{
		authtypes.NewModuleAddress(fmt.Sprintf("%s/pool/%d/protocol_fees", types.ModuleName, c.GetId())),
	}
}

func ComputeWeightedPoolUnbondingPeriod(totalShares math.LegacyDec, sharesToUnbond math.LegacyDec) (time.Duration, error) {
	const (
		thresholdLong   = 10.
		thresholdMedium = 1.
		thresholdShort  = 0.1
		baseDuration    = 1 * time.Minute
		shortDuration   = 30 * time.Minute
		mediumDuration  = 12 * time.Hour
		longDuration    = 24 * time.Hour
	)

	if !sharesToUnbond.IsPositive() || !totalShares.IsPositive() {
		return baseDuration, fmt.Errorf("invalid zero values")
	}

	// Calculate the percentage of total shares that the amount represents
	percentage := sharesToUnbond.Quo(totalShares).MulInt64(100).MustFloat64()

	// Determine the unbonding duration based on the percentage thresholds
	switch {
	case percentage > thresholdLong:
		return longDuration, nil
	case percentage > thresholdMedium:
		return mediumDuration, nil
	case percentage > thresholdShort:
		return shortDuration, nil
	default:
		return baseDuration, nil
	}
}

func CalculatePositionRewards(
	currentTime time.Time,
	poolRewards sdk.Coins,
	position stableswaptypes.BondedPosition,
	totalShares math.LegacyDec,
	initialPoolRewardsTime time.Time,
) (sdk.Coins, error) {
	duration := currentTime.Sub(position.RewardsPeriodStart)
	poolDuration := currentTime.Sub(initialPoolRewardsTime)

	// Ensure that the period is invalid
	if poolDuration.Seconds() == 0 || duration.Seconds() == 0 {
		return nil, fmt.Errorf("period is too short")
	}

	var positionRewards sdk.Coins
	for _, coinRewards := range poolRewards {
		numerator := position.Balance.MulInt64(int64(duration.Seconds()))  // User's shares * position period
		denominator := totalShares.MulInt64(int64(poolDuration.Seconds())) // Total shares * pool period
		userReward := numerator.Quo(denominator).Mul(coinRewards.Amount.ToLegacyDec())
		positionRewards = append(positionRewards, sdk.NewCoin(coinRewards.Denom, userReward.TruncateInt()))
	}
	return positionRewards, nil
}
