package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/header"
	"cosmossdk.io/core/store"
	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	modulev1 "swap.noble.xyz/api/module/v1"
	"swap.noble.xyz/keeper/stableswap"
	"swap.noble.xyz/types"
)

type Keeper struct {
	authority        string
	baseDenom        string
	stableswapConfig *modulev1.StableSwap

	eventService  event.Service
	headerService header.Service
	logger        log.Logger

	Schema collections.Schema

	// NextPoolID generates and keeps track of the next available unique pool ID for new pools.
	NextPoolID collections.Sequence

	// Paused tracks the paused state of pools, mapped by their unique pool ID (uint64).
	Paused collections.Map[uint64, bool]

	// Pools stores the generic pools, mapped by their unique pool ID (uint64).
	Pools collections.Map[uint64, types.Pool]

	// Stableswap is the sub-keeper responsible for managing StableSwap-specific functionalities.
	Stableswap *stableswap.Keeper

	addressCodec  address.Codec
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	eventService event.Service,
	headerService header.Service,
	logger log.Logger,
	authority string,
	pairDenom string,
	stableswapConfig *modulev1.StableSwap,
	addressCodec address.Codec,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) *Keeper {
	builder := collections.NewSchemaBuilder(storeService)

	keeper := &Keeper{
		authority:        authority,
		baseDenom:        pairDenom,
		stableswapConfig: stableswapConfig,

		eventService:  eventService,
		headerService: headerService,
		logger:        logger,

		NextPoolID: collections.NewSequence(builder, types.NextPoolIDPrefix, "next_pool_id"),
		Paused:     collections.NewMap(builder, types.PausedPrefix, "paused", collections.Uint64Key, collections.BoolValue),
		Pools:      collections.NewMap(builder, types.PoolsPrefix, "pools_generic", collections.Uint64Key, codec.CollValue[types.Pool](cdc)),

		Stableswap: stableswap.NewKeeper(cdc, storeService, eventService, headerService, logger),

		addressCodec:  addressCodec,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}

	schema, err := builder.Build()
	if err != nil {
		panic(err)
	}

	keeper.Schema = schema
	return keeper
}

// SetBankKeeper overwrites the bank keeper used in this module.
func (k *Keeper) SetBankKeeper(bankKeeper types.BankKeeper) {
	k.bankKeeper = bankKeeper
}

func (k *Keeper) Logger() log.Logger {
	return k.logger.With("module", types.ModuleName)
}

// Swap processes a token swap request, validates the message, and executes the swap routes,
// ensuring all conditions are met, including balances, slippage limits, and pool states.
func (k *Keeper) Swap(ctx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
	// Ensure that the signer is valid.
	userAddress, err := k.addressCodec.StringToBytes(msg.Signer)
	if err != nil {
		return nil, fmt.Errorf("unable to decode signer address: %s", msg.Signer)
	}

	// Validate the Swap message.
	if err = types.ValidateMsgSwap(msg); err != nil {
		return nil, err
	}

	// Check if the user has a balance >= than the requested swap amount.
	userBalance := k.bankKeeper.GetBalance(ctx, userAddress, msg.Amount.Denom)
	if userBalance.Amount.LT(msg.Amount.Amount) {
		return nil, sdkerrors.Wrapf(
			types.ErrInsufficientBalance,
			"%s balance of %s is smaller then %s, available: %s",
			msg.Amount.Denom, msg.Signer, msg.Amount.Amount.String(), userBalance.Amount.String(),
		)
	}

	// Prepare the swap plan in order to be executed, ensuring that the requested route pools are not paused.
	swapRoutesPlan, err := k.PrepareSwapPlan(ctx, msg, k.headerService.GetHeaderInfo(ctx).Time.Unix(), k)
	if err != nil {
		return nil, fmt.Errorf("error computing swap routes plan: %s", err.Error())
	}

	// Verify slippage limits.
	out := swapRoutesPlan.Swaps[len(swapRoutesPlan.Swaps)-1].Commitment.Out
	if out.IsLT(msg.Min) {
		return nil, fmt.Errorf("%s is less then min amount %s", out.String(), msg.Min.String())
	}

	// Commit the plan.
	var executedSwaps []*types.Swap
	for _, swap := range swapRoutesPlan.Swaps {
		poolAddr, err := sdk.AccAddressFromBech32(swap.PoolAddress)
		if err != nil {
			return nil, err
		}
		if err := k.bankKeeper.SendCoins(ctx, userAddress, poolAddr, sdk.NewCoins(swap.Commitment.In)); err != nil {
			return nil, sdkerrors.Wrap(err, "unable to transfer from provider to pool")
		}
		if err := k.bankKeeper.SendCoins(ctx, poolAddr, userAddress, sdk.NewCoins(swap.Commitment.Out)); err != nil {
			return nil, sdkerrors.Wrap(err, "unable to transfer from provider to pool")
		}

		totalFees := sdk.Coins{}
		for _, fee := range swap.Commitment.Fees {
			totalFees = totalFees.Add(fee.Amount)
			if err := k.bankKeeper.SendCoins(ctx, poolAddr, fee.Address.Bytes(), sdk.NewCoins(fee.Amount)); err != nil {
				return nil, sdkerrors.Wrap(err, "unable to transfer from provider to pool")
			}
		}
		executedSwaps = append(executedSwaps, &types.Swap{
			PoolId: swap.PoolId,
			In:     swap.Commitment.In,
			Out:    swap.Commitment.Out,
			Fees:   totalFees,
		})
	}

	return &types.MsgSwapResponse{
		Result: swapRoutesPlan.Swaps[len(swapRoutesPlan.Swaps)-1].Commitment.Out,
		Swaps:  executedSwaps,
	}, nil
}

// PrepareSwapPlan prepares a swap route plan from the swap message, containing the details for its execution.
func (k *Keeper) PrepareSwapPlan(ctx context.Context, msg *types.MsgSwap, timestamp int64, s *Keeper) (*types.PlanSwapRoutes, error) {
	var swaps []types.PlanSwapRoute

	swapIn := msg.Amount // Initial swap amount.
	for _, route := range msg.Routes {
		// Retrieve the Pool StableswapController for the requested Pool.
		controller, err := GetGenericController(ctx, s, route.PoolId)
		if err != nil {
			return nil, err
		}

		// Ensure that the Pool is not paused from execution.
		if controller.IsPaused() {
			return nil, sdkerrors.Wrapf(types.ErrPoolActivityPaused, "pool %d is paused", controller.GetId())
		}

		// Ensure that the Pool contains the requested `DenomTo`.
		if route.DenomTo != s.baseDenom && route.DenomTo != controller.GetPair() {
			return nil, sdkerrors.Wrapf(
				types.ErrInvalidSwapRoutingPlan, "pool %d doesn't contain denom %s", controller.GetId(), route.DenomTo,
			)
		}

		// Compute the Swap result.
		swapRes, err := controller.Swap(ctx, timestamp, swapIn, route.DenomTo)
		if err != nil {
			return nil, err
		}

		// Add the Commitment if the Swap is successful.
		swaps = append(swaps, types.PlanSwapRoute{
			PoolId:      controller.GetId(),
			PoolAddress: controller.GetAddress(),
			Commitment:  swapRes,
		})

		// Update the input swap amount after each route.
		swapIn = sdk.NewCoin(swapRes.Out.Denom, swapRes.Out.Amount)
	}
	return &types.PlanSwapRoutes{
		Swaps: swaps,
	}, nil
}
