package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/header"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"swap.noble.xyz/types"
	"swap.noble.xyz/types/stableswap"
	"swap.noble.xyz/utils"
	"swap.noble.xyz/utils/mocks"
)

func TestStableSwapBeginBlocker(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	user := utils.TestAccount()

	// ARRANGE: Trigger the BeginBlocker with an invalid height.
	ctx = ctx.WithHeaderInfo(header.Info{Height: 1})
	executed := k.StableSwapBeginBlocker(ctx)

	// ASSERT: The execution should've been skipped.
	assert.False(t, executed)

	// ARRANGE: Trigger the BeginBlocker with a valid height.
	ctx = ctx.WithHeaderInfo(header.Info{Height: 10})
	executed = k.StableSwapBeginBlocker(ctx)

	// ASSERT: The BeginBlocker should still succeed but logging the error.
	assert.True(t, executed)

	// ARRANGE: Add an invalid StableSwap Pool triggering errors on the StableSwap BeginBlocker execution, and a paused Pool.
	err := k.Stableswap.Pools.Set(ctx, 0, stableswap.Pool{
		ProtocolFeePercentage: 0,
	})
	assert.NoError(t, err)
	err = k.Stableswap.Pools.Set(ctx, 1, stableswap.Pool{
		ProtocolFeePercentage: 0,
	})
	assert.NoError(t, err)
	err = k.Pools.Set(ctx, 1, types.Pool{
		Id:           1,
		Address:      user.Address,
		Algorithm:    types.STABLESWAP,
		Pair:         "uusdc",
		CreationTime: time.Time{},
	})
	assert.NoError(t, err)
	err = k.Paused.Set(ctx, 1, true)
	assert.NoError(t, err)

	executed = k.StableSwapBeginBlocker(ctx)
	// ASSERT: The BeginBlocker should still succeed but logging the error.
	assert.True(t, executed)

	// ARRANGE: Add an invalid unbonding position
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)})
	err = k.Stableswap.UnbondingPositions.Set(ctx, collections.Join3(time.Time{}.Unix(), user.Address, uint64(1)), stableswap.UnbondingPosition{
		Shares:  math.LegacyDec{},
		Amount:  nil,
		EndTime: time.Time{},
	})
	assert.NoError(t, err)
	err = k.Pools.Remove(ctx, 0)
	assert.NoError(t, err)
	err = k.Paused.Set(ctx, 1, false)
	assert.NoError(t, err)
	k.Stableswap.UsersTotalBondedShares = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.StableSwapUsersTotalBondedSharesPrefix, "stableswap_users_total_bonded_shares", collections.PairKeyCodec(collections.Uint64Key, collections.StringKey), sdk.LegacyDecValue,
	)
	executed = k.StableSwapBeginBlocker(ctx)
	// ASSERT: The BeginBlocker should still succeed but logging the error.
	assert.True(t, executed)
}
