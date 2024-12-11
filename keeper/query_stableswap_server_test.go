package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"swap.noble.xyz/keeper"
	"swap.noble.xyz/types/stableswap"
	"swap.noble.xyz/utils"
	"swap.noble.xyz/utils/mocks"
)

func TestPositionsByProvider(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	stableswapQueryServer := keeper.NewStableSwapQueryServer(k.Stableswap)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	user := utils.TestAccount()

	// ARRANGE: create pools
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdc", math.NewInt(100)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdn", math.NewInt(100)))
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
		InitialA: 100,
		FutureA:  100,
	})
	assert.NoError(t, err)

	// ACT: Attempt to query with an invalid request.
	_, err = stableswapQueryServer.PositionsByProvider(ctx, &stableswap.QueryPositionsByProvider{})
	assert.Error(t, err)

	// ACT: Attempt to query with a valid request but 0 positions.
	res, err := stableswapQueryServer.PositionsByProvider(ctx, &stableswap.QueryPositionsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(res.Positions))

	// ARRANGE: Create a provider position by adding liquidity.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(10)),
			sdk.NewCoin("uusdn", math.NewInt(10)),
		),
	})
	assert.NoError(t, err)

	// ACT: Attempt to query with a valid request and 1 current position.
	res, err = stableswapQueryServer.PositionsByProvider(ctx, &stableswap.QueryPositionsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Positions))
}

func TestUnbondingPositionsByProvider(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	stableswapQueryServer := keeper.NewStableSwapQueryServer(k.Stableswap)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	user := utils.TestAccount()

	// ARRANGE: create pools
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdc", math.NewInt(100)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdn", math.NewInt(100)))
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
		InitialA: 100,
		FutureA:  100,
	})
	assert.NoError(t, err)

	// ACT: Attempt to query with an invalid request.
	_, err = stableswapQueryServer.UnbondingPositionsByProvider(ctx, &stableswap.QueryUnbondingPositionsByProvider{})
	assert.Error(t, err)

	// ACT: Attempt to query with a valid request but 0 positions.
	res, err := stableswapQueryServer.UnbondingPositionsByProvider(ctx, &stableswap.QueryUnbondingPositionsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(res.UnbondingPositions))
}
