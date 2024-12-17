package stableswap_test

import (
	"testing"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	stableswapkeeper "swap.noble.xyz/keeper/stableswap"
	types2 "swap.noble.xyz/types"
	"swap.noble.xyz/types/stableswap"
	"swap.noble.xyz/utils"
	"swap.noble.xyz/utils/mocks"
)

func TestGetPoolsTotalUnbondingShares(t *testing.T) {
	keeper, ctx := mocks.SwapKeeper(t)

	// ACT: Get the PoolsTotalUnbondingShares with no state.
	value := keeper.Stableswap.GetPoolsTotalUnbondingShares(ctx)
	assert.Equal(t, []stableswap.PoolsTotalUnbondingSharesEntry(nil), value)

	// ARRANGE: Set an entry in the state.
	err := keeper.Stableswap.SetPoolTotalUnbondingShares(ctx, 0, math.LegacyNewDec(1))
	assert.NoError(t, err)

	// ACT: Get the PoolsTotalUnbondingShares.
	value = keeper.Stableswap.GetPoolsTotalUnbondingShares(ctx)

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, 1, len(value))
	assert.Equal(t, []stableswap.PoolsTotalUnbondingSharesEntry{
		{
			PoolId: 0,
			Shares: math.LegacyNewDec(1),
		},
	}, value)

	// ARRANGE: Set another entry in the state.
	err = keeper.Stableswap.SetPoolTotalUnbondingShares(ctx, 1, math.LegacyNewDec(2))
	assert.NoError(t, err)

	// ACT: Get the PoolsTotalUnbondingShares.
	value = keeper.Stableswap.GetPoolsTotalUnbondingShares(ctx)

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, 2, len(value))
	assert.Equal(t, []stableswap.PoolsTotalUnbondingSharesEntry{
		{
			PoolId: 0,
			Shares: math.LegacyNewDec(1),
		},
		{
			PoolId: 1,
			Shares: math.LegacyNewDec(2),
		},
	}, value)
}

func TestGetUsersTotalBondedShares(t *testing.T) {
	keeper, ctx := mocks.SwapKeeper(t)

	// ACT: Get the UsersTotalBondedShares with no state.
	value := keeper.Stableswap.GetUsersTotalBondedShares(ctx)
	assert.Equal(t, []stableswap.UsersTotalBondedSharesEntry(nil), value)

	// ARRANGE: Set an entry in the state.
	err := keeper.Stableswap.SetUserTotalBondedShares(ctx, 0, "address1", math.LegacyNewDec(1))
	assert.NoError(t, err)

	// ACT: Get the UsersTotalBondedShares.
	value = keeper.Stableswap.GetUsersTotalBondedShares(ctx)

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, 1, len(value))
	assert.Equal(t, []stableswap.UsersTotalBondedSharesEntry{
		{
			PoolId:  0,
			Address: "address1",
			Shares:  math.LegacyNewDec(1),
		},
	}, value)

	// ARRANGE: Set another entry in the state.
	err = keeper.Stableswap.SetUserTotalBondedShares(ctx, 1, "address2", math.LegacyNewDec(2))
	assert.NoError(t, err)

	// ACT: Get the UsersTotalBondedShares.
	value = keeper.Stableswap.GetUsersTotalBondedShares(ctx)

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, 2, len(value))
	assert.Equal(t, []stableswap.UsersTotalBondedSharesEntry{
		{
			PoolId:  0,
			Address: "address1",
			Shares:  math.LegacyNewDec(1),
		},
		{
			PoolId:  1,
			Address: "address2",
			Shares:  math.LegacyNewDec(2),
		},
	}, value)
}

func TestGetUsersTotalUnbondingShares(t *testing.T) {
	keeper, ctx := mocks.SwapKeeper(t)

	// ACT: Get the UsersTotalUnbondingShares with no state.
	value := keeper.Stableswap.GetUsersTotalUnbondingShares(ctx)
	assert.Equal(t, []stableswap.UsersTotalUnbondingSharesEntry(nil), value)

	// ARRANGE: Set an entry in the state.
	err := keeper.Stableswap.SetUserTotalUnbondingShares(ctx, 0, "address1", math.LegacyNewDec(1))
	assert.NoError(t, err)

	// ACT: Get the UsersTotalUnbondingShares.
	value = keeper.Stableswap.GetUsersTotalUnbondingShares(ctx)

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, 1, len(value))
	assert.Equal(t, []stableswap.UsersTotalUnbondingSharesEntry{
		{
			PoolId:  0,
			Address: "address1",
			Shares:  math.LegacyNewDec(1),
		},
	}, value)

	// ARRANGE: Set another entry in the state.
	err = keeper.Stableswap.SetUserTotalUnbondingShares(ctx, 1, "address2", math.LegacyNewDec(2))
	assert.NoError(t, err)

	// ACT: Get the UsersTotalUnbondingShares.
	value = keeper.Stableswap.GetUsersTotalUnbondingShares(ctx)

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, 2, len(value))
	assert.Equal(t, []stableswap.UsersTotalUnbondingSharesEntry{
		{
			PoolId:  0,
			Address: "address1",
			Shares:  math.LegacyNewDec(1),
		},
		{
			PoolId:  1,
			Address: "address2",
			Shares:  math.LegacyNewDec(2),
		},
	}, value)
}

func TestGetBondedPositions(t *testing.T) {
	keeper, ctx := mocks.SwapKeeper(t)

	// ACT: Get the BondedPositions with no state.
	value := keeper.Stableswap.GetBondedPositions(ctx)
	assert.Equal(t, []stableswap.BondedPositionEntry(nil), value)

	// ARRANGE: Set an entry in the state.
	t1 := time.Time{}.Add(time.Hour)
	err := keeper.Stableswap.SetBondedPosition(ctx, 0, "address1", t1.Unix(), stableswap.BondedPosition{
		Balance:            math.LegacyNewDec(1),
		Timestamp:          t1,
		RewardsPeriodStart: t1,
	})
	assert.NoError(t, err)

	// ACT: Get the BondedPositions.
	value = keeper.Stableswap.GetBondedPositions(ctx)

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, 1, len(value))
	assert.Equal(t, []stableswap.BondedPositionEntry{
		{
			PoolId:    0,
			Address:   "address1",
			Timestamp: t1.Unix(),
			BondedPosition: stableswap.BondedPosition{
				Balance:            math.LegacyNewDec(1),
				Timestamp:          t1,
				RewardsPeriodStart: t1,
			},
		},
	}, value)

	// ARRANGE: Set another entry in the state.
	t2 := t1.Add(time.Hour)
	err = keeper.Stableswap.SetBondedPosition(ctx, 1, "address2", t2.Unix(), stableswap.BondedPosition{
		Balance:            math.LegacyNewDec(2),
		Timestamp:          t2,
		RewardsPeriodStart: t2,
	})
	assert.NoError(t, err)

	// ACT: Get the BondedPositions.
	value = keeper.Stableswap.GetBondedPositions(ctx)

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, 2, len(value))
	assert.Equal(t, []stableswap.BondedPositionEntry{
		{
			PoolId:    0,
			Address:   "address1",
			Timestamp: t1.Unix(),
			BondedPosition: stableswap.BondedPosition{
				Balance:            math.LegacyNewDec(1),
				Timestamp:          t1,
				RewardsPeriodStart: t1,
			},
		},
		{
			PoolId:    1,
			Address:   "address2",
			Timestamp: t2.Unix(),
			BondedPosition: stableswap.BondedPosition{
				Balance:            math.LegacyNewDec(2),
				Timestamp:          t2,
				RewardsPeriodStart: t2,
			},
		},
	}, value)
}

func TestGetBondedPositionsByProvider(t *testing.T) {
	keeper, ctx := mocks.SwapKeeper(t)

	// ACT: Get the BondedPositions with no state.
	value := keeper.Stableswap.GetBondedPositionsByProvider(ctx, "address1")
	assert.Equal(t, []stableswap.BondedPositionEntry(nil), value)

	// ARRANGE: Set an entry in the state.
	t1 := time.Time{}.Add(time.Hour)
	err := keeper.Stableswap.SetBondedPosition(ctx, 0, "address1", t1.Unix(), stableswap.BondedPosition{
		Balance:            math.LegacyNewDec(1),
		Timestamp:          t1,
		RewardsPeriodStart: t1,
	})
	assert.NoError(t, err)

	// ACT: Get the BondedPositions.
	value = keeper.Stableswap.GetBondedPositionsByProvider(ctx, "address1")

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, 1, len(value))
	assert.Equal(t, []stableswap.BondedPositionEntry{
		{
			PoolId:    0,
			Address:   "address1",
			Timestamp: t1.Unix(),
			BondedPosition: stableswap.BondedPosition{
				Balance:            math.LegacyNewDec(1),
				Timestamp:          t1,
				RewardsPeriodStart: t1,
			},
		},
	}, value)

	// ARRANGE: Set another entry in the state.
	t2 := t1.Add(time.Hour)
	err = keeper.Stableswap.SetBondedPosition(ctx, 1, "address2", t2.Unix(), stableswap.BondedPosition{
		Balance:            math.LegacyNewDec(2),
		Timestamp:          t2,
		RewardsPeriodStart: t2,
	})
	assert.NoError(t, err)

	// ACT: Get the BondedPositions.
	value = keeper.Stableswap.GetBondedPositionsByProvider(ctx, "address1")

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, 1, len(value))
	assert.Equal(t, []stableswap.BondedPositionEntry{
		{
			PoolId:    0,
			Address:   "address1",
			Timestamp: t1.Unix(),
			BondedPosition: stableswap.BondedPosition{
				Balance:            math.LegacyNewDec(1),
				Timestamp:          t1,
				RewardsPeriodStart: t1,
			},
		},
	}, value)

	// ARRANGE: Set failing collections for BondingPositions.
	builder := collections.NewSchemaBuilder(mocks.FailingStore(mocks.Iterator, utils.GetKVStore(ctx, types2.ModuleName)))
	keeper.Stableswap.BondedPositions = collections.NewIndexedMap(
		builder, types2.StableSwapBondedPositionsPrefix, "stableswap_bonded_positions",
		collections.TripleKeyCodec(collections.Uint64Key, collections.StringKey, collections.Int64Key),
		codec.CollValue[stableswap.BondedPosition](mocks.MakeTestEncodingConfig("noble").Codec),
		stableswapkeeper.NewBondedPositionIndexes(builder),
	)

	// ACT: Attempt to get the BondedPositions.
	value = keeper.Stableswap.GetBondedPositionsByProvider(ctx, "")
	// ASSERT: No results due to the failing collections.
	assert.Equal(t, []stableswap.BondedPositionEntry(nil), value)
}

func TestGetUnbondingPositions(t *testing.T) {
	keeper, ctx := mocks.SwapKeeper(t)

	// ACT: Get the UnbondingPositions with no state.
	value := keeper.Stableswap.GetUnbondingPositions(ctx)
	assert.Equal(t, []stableswap.UnbondingPositionEntry(nil), value)

	// ARRANGE: Set an entry in the state.
	t1 := time.Time{}.Add(time.Hour)
	err := keeper.Stableswap.SetUnbondingPosition(ctx, t1.Unix(), "address1", 0, stableswap.UnbondingPosition{
		Shares:  math.LegacyNewDec(1),
		Amount:  types.NewCoins(types.NewCoin("uusdc", math.NewInt(1)), types.NewCoin("uusdn", math.NewInt(1))),
		EndTime: t1,
	})
	assert.NoError(t, err)

	// ACT: Get the UnbondingPositions.
	value = keeper.Stableswap.GetUnbondingPositions(ctx)

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, 1, len(value))
	assert.Equal(t, []stableswap.UnbondingPositionEntry{
		{
			PoolId:    0,
			Address:   "address1",
			Timestamp: t1.Unix(),
			UnbondingPosition: stableswap.UnbondingPosition{
				Shares:  math.LegacyNewDec(1),
				Amount:  types.NewCoins(types.NewCoin("uusdc", math.NewInt(1)), types.NewCoin("uusdn", math.NewInt(1))),
				EndTime: t1,
			},
		},
	}, value)

	// ARRANGE: Set another entry in the state.
	t2 := t1.Add(time.Hour)
	err = keeper.Stableswap.SetUnbondingPosition(ctx, t2.Unix(), "address2", 1, stableswap.UnbondingPosition{
		Shares:  math.LegacyNewDec(2),
		Amount:  types.NewCoins(types.NewCoin("uusdc", math.NewInt(2)), types.NewCoin("uusdn", math.NewInt(2))),
		EndTime: t2,
	})
	assert.NoError(t, err)

	// ACT: Get the UnbondingPositions.
	value = keeper.Stableswap.GetUnbondingPositions(ctx)

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, 2, len(value))
	assert.Equal(t, []stableswap.UnbondingPositionEntry{
		{
			PoolId:    0,
			Address:   "address1",
			Timestamp: t1.Unix(),
			UnbondingPosition: stableswap.UnbondingPosition{
				Shares:  math.LegacyNewDec(1),
				Amount:  types.NewCoins(types.NewCoin("uusdc", math.NewInt(1)), types.NewCoin("uusdn", math.NewInt(1))),
				EndTime: t1,
			},
		},
		{
			PoolId:    1,
			Address:   "address2",
			Timestamp: t2.Unix(),
			UnbondingPosition: stableswap.UnbondingPosition{
				Shares:  math.LegacyNewDec(2),
				Amount:  types.NewCoins(types.NewCoin("uusdc", math.NewInt(2)), types.NewCoin("uusdn", math.NewInt(2))),
				EndTime: t2,
			},
		},
	}, value)
}

func TestGetUnbondingPositionsByProvider(t *testing.T) {
	keeper, ctx := mocks.SwapKeeper(t)

	// ARRANGE: Set an entry in the state.
	t1 := time.Time{}.Add(time.Hour)
	err := keeper.Stableswap.SetUnbondingPosition(ctx, t1.Unix(), "address1", 0, stableswap.UnbondingPosition{
		Shares:  math.LegacyNewDec(1),
		Amount:  types.NewCoins(types.NewCoin("uusdc", math.NewInt(1)), types.NewCoin("uusdn", math.NewInt(1))),
		EndTime: t1,
	})
	assert.NoError(t, err)

	// ACT: Get the UnbondingPositions.
	value := keeper.Stableswap.GetUnbondingPositionsByProvider(ctx, "address1")

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, 1, len(value))
	assert.Equal(t, []stableswap.UnbondingPositionEntry{
		{
			PoolId:    0,
			Address:   "address1",
			Timestamp: t1.Unix(),
			UnbondingPosition: stableswap.UnbondingPosition{
				Shares:  math.LegacyNewDec(1),
				Amount:  types.NewCoins(types.NewCoin("uusdc", math.NewInt(1)), types.NewCoin("uusdn", math.NewInt(1))),
				EndTime: t1,
			},
		},
	}, value)

	// ARRANGE: Set another entry in the state.
	t2 := t1.Add(time.Hour)
	err = keeper.Stableswap.SetUnbondingPosition(ctx, t2.Unix(), "address2", 1, stableswap.UnbondingPosition{
		Shares:  math.LegacyNewDec(2),
		Amount:  types.NewCoins(types.NewCoin("uusdc", math.NewInt(2)), types.NewCoin("uusdn", math.NewInt(2))),
		EndTime: t2,
	})
	assert.NoError(t, err)

	// ACT: Get the UnbondingPositions.
	value = keeper.Stableswap.GetUnbondingPositionsByProvider(ctx, "address1")

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, 1, len(value))
	assert.Equal(t, []stableswap.UnbondingPositionEntry{
		{
			PoolId:    0,
			Address:   "address1",
			Timestamp: t1.Unix(),
			UnbondingPosition: stableswap.UnbondingPosition{
				Shares:  math.LegacyNewDec(1),
				Amount:  types.NewCoins(types.NewCoin("uusdc", math.NewInt(1)), types.NewCoin("uusdn", math.NewInt(1))),
				EndTime: t1,
			},
		},
	}, value)

	// ARRANGE: Set failing collections for UnbondingPositions
	builder := collections.NewSchemaBuilder(mocks.FailingStore(mocks.Iterator, utils.GetKVStore(ctx, types2.ModuleName)))
	keeper.Stableswap.UnbondingPositions = collections.NewIndexedMap(
		builder, types2.StableSwapUnbondingPositionsPrefix, "stableswap_unbonding_positions",
		collections.TripleKeyCodec(collections.Int64Key, collections.StringKey, collections.Uint64Key),
		codec.CollValue[stableswap.UnbondingPosition](mocks.MakeTestEncodingConfig("noble").Codec),
		stableswapkeeper.NewUnbondingPositionIndexes(builder),
	)

	// ACT: Attempt to get the UnbondingPositions.
	value = keeper.Stableswap.GetUnbondingPositionsByProvider(ctx, "")
	// ASSERT: No results due to the failing collections.
	assert.Equal(t, []stableswap.UnbondingPositionEntry(nil), value)
}
