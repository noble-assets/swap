package keeper_test

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/header"
	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	stableswapkeeper "swap.noble.xyz/keeper"
	"swap.noble.xyz/types"
	"swap.noble.xyz/types/stableswap"
	"swap.noble.xyz/utils"
	"swap.noble.xyz/utils/mocks"
)

func TestCreateStableSwapPool(t *testing.T) {
	// ARRANGE: Test cases validating each message attribute and collections access.
	tests := []struct {
		name        string
		msg         *stableswap.MsgCreatePool
		error       error
		mockStateFn func(k *stableswapkeeper.Keeper, ctx sdk.Context)
	}{
		{
			"Invalid authority",
			&stableswap.MsgCreatePool{
				Signer: "user",
				Pair:   "",
			},
			sdkerrors.Wrapf(types.ErrInvalidAuthority, "expected authority, got user"),
			nil,
		},
		{
			"Missing pair",
			&stableswap.MsgCreatePool{
				Signer: "authority",
				Pair:   "",
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "missing pair value"),
			nil,
		},
		{
			"Non existing Pair",
			&stableswap.MsgCreatePool{
				Signer: "authority",
				Pair:   "invalid-pair",
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "invalid-pair does not exists on chain"),
			nil,
		},
		{
			"Same pair denom",
			&stableswap.MsgCreatePool{
				Signer: "authority",
				Pair:   "uusdn",
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "pair denom must be different from uusdn"),
			nil,
		},
		{
			"Missing InitialA value",
			&stableswap.MsgCreatePool{
				Signer: "authority",
				Pair:   "uusdc",
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "invalid InitialA value"),
			nil,
		},
		{
			"Invalid InitialA value",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: -1,
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "invalid InitialA value"),
			nil,
		},
		{
			"Invalid RateMultipliers, missing pair value",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: 100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
				),
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "RateMultipliers length must be 2, got 1"),
			nil,
		},
		{
			"Invalid Rate Multipliers, invalid pair value",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(0)),
				),
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "RateMultipliers length must be 2, got 1"),
			nil,
		},
		{
			"Invalid Rate Multipliers, too many values",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: 100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
					sdk.NewCoin("ueure", math.NewInt(1000000000000000000)),
				),
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "RateMultipliers length must be 2, got 3"),
			nil,
		},
		{
			"Invalid Rate Multipliers, invalid base pair denom",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdx", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "uusdn rate multiplier must be positive, got 0"),
			nil,
		},
		{
			"Invalid Rate Multipliers, invalid pair denom",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdx", math.NewInt(1000000000000000000)),
				),
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "uusdc rate multiplier must be positive, got 0"),
			nil,
		},
		{
			"Invalid MaxFee value",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
				MaxFee: -1,
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "MaxFee cannot be negative"),
			nil,
		},
		{
			"Invalid ProtocolFee value (<0)",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
				ProtocolFeePercentage: -1,
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "invalid ProtocolFeePercentage value"),
			nil,
		},
		{
			"Invalid ProtocolFee value (>100)",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
				ProtocolFeePercentage: 104,
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "invalid ProtocolFeePercentage value"),
			nil,
		},
		{
			"Invalid RewardsFee value",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
				RewardsFee: -1,
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "RewardsFee cannot be negative"),
			nil,
		},
		{
			"Invalid MaxFee value",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
				MaxFee: -1,
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "MaxFee cannot be negative"),
			nil,
		},
		{
			"[collections] Failing collection on Set NextPoolId",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
			},
			sdkerrors.Wrapf(mocks.ErrorStoreAccess, "unable to set next pool id"),
			func(k *stableswapkeeper.Keeper, ctx sdk.Context) {
				k.NextPoolID = collections.NewSequence(
					collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
					types.NextPoolIDPrefix, "next_pool_id",
				)
			},
		},
		{
			"[collections] Failing collection on Set Pool",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
			},
			sdkerrors.Wrapf(mocks.ErrorStoreAccess, "unable to set pool"),
			func(k *stableswapkeeper.Keeper, ctx sdk.Context) {
				k.Pools = collections.NewMap(
					collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
					types.PoolsPrefix, "pools_generic", collections.Uint64Key, codec.CollValue[types.Pool](mocks.MakeTestEncodingConfig("noble").Codec),
				)
			},
		},
		{
			"[collections] Failing collection on Set StableSwap Pool",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
			},
			sdkerrors.Wrapf(mocks.ErrorStoreAccess, "unable to set stableswap pool"),
			func(k *stableswapkeeper.Keeper, ctx sdk.Context) {
				k.Stableswap.Pools = collections.NewMap(
					collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
					types.StableSwapPoolsPrefix, "stableswap_pools", collections.Uint64Key, codec.CollValue[stableswap.Pool](mocks.MakeTestEncodingConfig("noble").Codec),
				)
			},
		},
		{
			"[collections] Failing collection on Set Paused",
			&stableswap.MsgCreatePool{
				Signer:   "authority",
				Pair:     "uusdc",
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
			},
			sdkerrors.Wrapf(mocks.ErrorStoreAccess, "unable to create paused entry"),
			func(k *stableswapkeeper.Keeper, ctx sdk.Context) {
				k.Paused = collections.NewMap(
					collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
					types.PausedPrefix, "paused", collections.Uint64Key, codec.BoolValue,
				)
			},
		},
	}

	// ASSERT: Execute each test case and expect the related error.
	for _, tt := range tests {
		account := mocks.AccountKeeper{
			Accounts: make(map[string]sdk.AccountI),
		}
		bank := mocks.BankKeeper{
			Balances:    make(map[string]sdk.Coins),
			Restriction: mocks.NoOpSendRestrictionFn,
		}
		k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
		server := stableswapkeeper.NewStableSwapMsgServer(k)

		t.Run(tt.name, func(t *testing.T) {
			if tt.mockStateFn != nil {
				tt.mockStateFn(k, ctx)
			}
			_, err := server.CreatePool(ctx, tt.msg)
			assert.NotNil(t, err)
			require.Equal(t, tt.error.Error(), err.Error())
		})
	}

	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	server := stableswapkeeper.NewStableSwapMsgServer(k)

	// ARRANGE: Create a valid Pool.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)})
	_, err := server.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer: "authority",
		Pair:   "uusdc",
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
		InitialA:              100,
		FutureA:               200,
		FutureATime:           10000000,
		MaxFee:                10000,
		RewardsFee:            10,
		ProtocolFeePercentage: 1,
	})
	assert.Nil(t, err)

	// ARRANGE: Retrieve the Generic Pool from state.
	pool, err := k.Pools.Get(ctx, 0)
	assert.Nil(t, err)

	// ARRANGE: Compute the expected Pool account.
	pool0Account := authtypes.NewEmptyModuleAccount(fmt.Sprintf("%s/pool/0", types.ModuleName))

	// ASSERT: Expect matching values in the state.
	assert.Equal(t, pool, types.Pool{
		Id:           0,
		Address:      pool0Account.Address,
		Algorithm:    types.STABLESWAP,
		Pair:         "uusdc",
		CreationTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	// ARRANGE: Retrieve the StableSwap Pool from state.
	ssPool, err := k.Stableswap.Pools.Get(ctx, 0)
	assert.Nil(t, err)

	// ASSERT: Expect matching values in the StableSwap state.
	assert.Equal(t, stableswap.Pool{
		ProtocolFeePercentage: 1,
		RewardsFee:            10,
		MaxFee:                10000,
		InitialA:              100,
		FutureA:               200,
		InitialATime:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		FutureATime:           10000000,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
		TotalShares: math.LegacyZeroDec(),
	}, ssPool)

	// ASSERT: Expect empty liquidity amount
	poolLiquidity := bank.GetAllBalances(ctx, sdk.AccAddress(pool0Account.Address))
	assert.Equal(t, poolLiquidity, sdk.Coins(nil))

	// ARRANGE: Create a Second Pool, with unordered rate multipliers.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)})
	_, err = server.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "ueure",
		RewardsFee:            4e3,
		ProtocolFeePercentage: 1,
		MaxFee:                1,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           10000000,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("ueure", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
		),
	})
	assert.Nil(t, err)

	// ARRANGE: Compute the expected Pool account.
	pool1Account := authtypes.NewEmptyModuleAccount(fmt.Sprintf("%s/pool/1", types.ModuleName))

	// ASSERT: Expect matching values in the state.
	pools := k.GetPools(ctx)
	assert.Equal(t, map[uint64]types.Pool{
		0: {
			Id:           0,
			Address:      pool0Account.Address,
			Algorithm:    types.STABLESWAP,
			Pair:         "uusdc",
			CreationTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		1: {
			Id:           1,
			Address:      pool1Account.Address,
			Algorithm:    types.STABLESWAP,
			Pair:         "ueure",
			CreationTime: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
		},
	}, pools)

	// ASSERT: Expect matching values in the StableSwap state.
	ssPools := k.Stableswap.GetPools(ctx)
	assert.Equal(t, map[uint64]stableswap.Pool{
		0: {
			RateMultipliers: sdk.NewCoins(
				sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
				sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
			),
			InitialA:              100,
			InitialATime:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			FutureA:               200,
			FutureATime:           10000000,
			MaxFee:                10000,
			RewardsFee:            10,
			ProtocolFeePercentage: 1,
			TotalShares:           math.LegacyZeroDec(),
		},
		1: {
			RateMultipliers: sdk.NewCoins(
				sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
				sdk.NewCoin("ueure", math.NewInt(1000000000000000000)),
			),
			InitialATime:          time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
			RewardsFee:            4e3,
			ProtocolFeePercentage: 1,
			MaxFee:                1,
			InitialA:              100,
			FutureA:               100,
			FutureATime:           10000000,
			TotalShares:           math.LegacyZeroDec(),
		},
	}, ssPools)

	_, err = server.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "ueure",
		RewardsFee:            4e3,
		ProtocolFeePercentage: 1,
		MaxFee:                1,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           10000000,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("ueure", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
		),
	})
	assert.Equal(t, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "pool with pair ueure and STABLESWAP algorithm already exists").Error(), err.Error())
}

func TestUpdateStableSwapPool(t *testing.T) {
	// ARRANGE: Test cases validating each message attribute.
	tests := []struct {
		name        string
		msg         *stableswap.MsgUpdatePool
		error       error
		mockStateFn func(k *stableswapkeeper.Keeper, ctx sdk.Context)
	}{
		{
			"Invalid authority",
			&stableswap.MsgUpdatePool{
				Signer: "user",
				PoolId: 10,
			},
			sdkerrors.Wrapf(types.ErrInvalidAuthority, "expected authority, got user"),
			nil,
		},
		{
			"Non existing Pool",
			&stableswap.MsgUpdatePool{
				Signer: "authority",
				PoolId: 10,
			},
			sdkerrors.Wrapf(types.ErrInvalidPool, "stableswap pool with Id 10 does not exists"),
			nil,
		},
		{
			"Invalid initial A param",
			&stableswap.MsgUpdatePool{
				Signer:   "authority",
				PoolId:   0,
				InitialA: -1,
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "invalid InitialA value"),
			nil,
		},
		{
			"Invalid RateMultipliers, missing pair value",
			&stableswap.MsgUpdatePool{
				Signer:   "authority",
				PoolId:   0,
				InitialA: 100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
				),
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "RateMultipliers length must be 2, got 1"),
			nil,
		},
		{
			"Invalid Rate Multipliers, invalid base pair denom",
			&stableswap.MsgUpdatePool{
				Signer:   "authority",
				PoolId:   0,
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdx", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "uusdn rate multiplier must be positive, got 0"),
			nil,
		},
		{
			"Invalid Rate Multipliers, invalid pair value",
			&stableswap.MsgUpdatePool{
				Signer:   "authority",
				PoolId:   0,
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(0)),
				),
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "RateMultipliers length must be 2, got 1"),
			nil,
		},
		{
			"Invalid Rate Multipliers, too many values",
			&stableswap.MsgUpdatePool{
				Signer:   "authority",
				PoolId:   0,
				InitialA: 100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
					sdk.NewCoin("ueure", math.NewInt(1000000000000000000)),
				),
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "RateMultipliers length must be 2, got 3"),
			nil,
		},
		{
			"Invalid Rate Multipliers, invalid pair denom",
			&stableswap.MsgUpdatePool{
				Signer:   "authority",
				PoolId:   0,
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdx", math.NewInt(1000000000000000000)),
				),
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "uusdc rate multiplier must be positive, got 0"),
			nil,
		},
		{
			"Invalid MaxFee value",
			&stableswap.MsgUpdatePool{
				Signer:   "authority",
				PoolId:   0,
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
				MaxFee: -1,
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "MaxFee cannot be negative"),
			nil,
		},
		{
			"Invalid ProtocolFee value (<0)",
			&stableswap.MsgUpdatePool{
				Signer:   "authority",
				PoolId:   0,
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
				ProtocolFeePercentage: -1,
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "invalid ProtocolFeePercentage value"),
			nil,
		},
		{
			"Invalid ProtocolFee value (>100)",
			&stableswap.MsgUpdatePool{
				Signer:   "authority",
				PoolId:   0,
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
				ProtocolFeePercentage: 104,
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "invalid ProtocolFeePercentage value"),
			nil,
		},
		{
			"Invalid RewardsFee value",
			&stableswap.MsgUpdatePool{
				Signer:   "authority",
				PoolId:   0,
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
				RewardsFee: -1,
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "RewardsFee cannot be negative"),
			nil,
		},
		{
			"Invalid MaxFee value",
			&stableswap.MsgUpdatePool{
				Signer:   "authority",
				PoolId:   0,
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
				MaxFee: -1,
			},
			sdkerrors.Wrapf(types.ErrInvalidPoolParams, "MaxFee cannot be negative"),
			nil,
		},
		{
			"Invalid Pool",
			&stableswap.MsgUpdatePool{
				Signer:   "authority",
				PoolId:   0,
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
			},
			sdkerrors.Wrapf(types.ErrInvalidPool, "invalid pool algorithm"),
			func(k *stableswapkeeper.Keeper, ctx sdk.Context) {
				pool, _ := k.GetPool(ctx, 0)
				pool.Algorithm = types.PERFECTPRICE
				_ = k.Pools.Set(ctx, 0, pool)
			},
		},
		{
			"[collections] Invalid Set StableSwap Pool",
			&stableswap.MsgUpdatePool{
				Signer:   "authority",
				PoolId:   0,
				InitialA: 100,
				FutureA:  100,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
					sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
				),
			},
			sdkerrors.Wrapf(mocks.ErrorStoreAccess, "unable to set stableswap pool"),
			func(k *stableswapkeeper.Keeper, ctx sdk.Context) {
				k.Stableswap.Pools = collections.NewMap(
					collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
					types.StableSwapPoolsPrefix, "stableswap_pools", collections.Uint64Key, codec.CollValue[stableswap.Pool](mocks.MakeTestEncodingConfig("noble").Codec),
				)
			},
		},
	}

	initBaseState := func() (mocks.AccountKeeper, mocks.BankKeeper, *stableswapkeeper.Keeper, sdk.Context, stableswap.MsgServer) {
		account := mocks.AccountKeeper{
			Accounts: make(map[string]sdk.AccountI),
		}
		bank := mocks.BankKeeper{
			Balances:    make(map[string]sdk.Coins),
			Restriction: mocks.NoOpSendRestrictionFn,
		}
		k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
		server := stableswapkeeper.NewStableSwapMsgServer(k)

		ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)})

		// ARRANGE: Create 2 Pools.
		_, err := server.CreatePool(ctx, &stableswap.MsgCreatePool{
			Signer: "authority",
			Pair:   "uusdc",
			RateMultipliers: sdk.NewCoins(
				sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
				sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
			),
			InitialA:              100,
			FutureA:               200,
			FutureATime:           10000000,
			MaxFee:                10000,
			RewardsFee:            10,
			ProtocolFeePercentage: 1,
		})
		assert.Nil(t, err)
		_, err = server.CreatePool(ctx, &stableswap.MsgCreatePool{
			Signer: "authority",
			Pair:   "ueure",
			RateMultipliers: sdk.NewCoins(
				sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
				sdk.NewCoin("ueure", math.NewInt(1000000000000000000)),
			),
			InitialA:              300,
			FutureA:               400,
			FutureATime:           10000001,
			MaxFee:                10001,
			RewardsFee:            11,
			ProtocolFeePercentage: 2,
		})
		assert.Nil(t, err)

		return account, bank, k, ctx, server
	}

	// ASSERT: Execute each test case and expect the related error.
	for _, tt := range tests {
		_, _, k, ctx, server := initBaseState()
		// ASSERT: Execute each test case and expect the related error.
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockStateFn != nil {
				tt.mockStateFn(k, ctx)
			}
			_, err := server.UpdatePool(ctx, tt.msg)
			assert.NotNil(t, err)
			require.Equal(t, tt.error.Error(), err.Error())
		})
	}

	_, _, k, ctx, server := initBaseState()

	_, err := server.UpdatePool(ctx, &stableswap.MsgUpdatePool{
		Signer:                "authority",
		PoolId:                1,
		InitialA:              101,
		FutureA:               101,
		FutureATime:           10000002,
		RewardsFee:            11,
		ProtocolFeePercentage: 2,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("ueure", math.NewInt(1000000000000000000)),
		),
	})
	assert.Nil(t, err)

	// ASSERT: Expect matching values in the StableSwap state.
	ssPools := k.Stableswap.GetPools(ctx)
	assert.Equal(t, map[uint64]stableswap.Pool{
		0: {
			RateMultipliers: sdk.NewCoins(
				sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
				sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
			),
			InitialA:              100,
			InitialATime:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			FutureA:               200,
			FutureATime:           10000000,
			MaxFee:                10000,
			RewardsFee:            10,
			ProtocolFeePercentage: 1,
			TotalShares:           math.LegacyZeroDec(),
		},
		1: {
			RateMultipliers: sdk.NewCoins(
				sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
				sdk.NewCoin("ueure", math.NewInt(1000000000000000000)),
			),
			InitialATime:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			InitialA:              101,
			FutureA:               101,
			FutureATime:           10000002,
			MaxFee:                0,
			RewardsFee:            11,
			ProtocolFeePercentage: 2,
			TotalShares:           math.LegacyZeroDec(),
		},
	}, ssPools)
}

func TestAddLiquidity(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	stableswapServer := stableswapkeeper.NewStableSwapMsgServer(k)

	user := utils.TestAccount()

	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		InitialA:              100,
		FutureA:               100,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})
	assert.Nil(t, err)

	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 1,
		Amount: sdk.NewCoins(),
	})
	assert.Error(t, err, types.ErrInvalidPool)

	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(),
	})
	assert.Error(t, err, types.ErrInvalidDenom)

	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(),
	})
	assert.Error(t, err, types.ErrInvalidDenom)

	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1000000000000000000))),
	})
	assert.Error(t, err, types.ErrInvalidDenom)

	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(10)),
			sdk.NewCoin("invalid_denom", math.NewInt(10)),
		),
	})
	assert.Error(t, err, types.ErrInvalidDenom)

	pool, err := k.Pools.Get(ctx, 0)
	assert.Nil(t, err)
	poolAddress, err := sdk.AccAddressFromBech32(pool.Address)
	assert.Nil(t, err)
	poolLiquidity := bank.GetAllBalances(ctx, poolAddress)
	assert.Equal(t, len(poolLiquidity), 0)

	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdc", math.NewInt(100)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdn", math.NewInt(100)))
	response, err := stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(10)),
			sdk.NewCoin("uusdn", math.NewInt(10)),
		),
	})
	assert.Nil(t, err)
	assert.Equal(t, response.MintedShares, int64(20))
	pool, err = k.Pools.Get(ctx, 0)
	assert.Nil(t, err)
	assert.Equal(t, bank.Balances[user.Address].AmountOf("uusdc"), math.NewInt(90))
	assert.Equal(t, bank.Balances[user.Address].AmountOf("uusdn"), math.NewInt(90))
	poolLiquidity = bank.GetAllBalances(ctx, poolAddress)
	assert.Equal(t, poolLiquidity.AmountOf("uusdn").Int64(), int64(10))
	assert.Equal(t, poolLiquidity.AmountOf("uusdc").Int64(), int64(10))

	key := collections.Join3(uint64(0), user.Address, int64(-62135596800))
	position, err := k.Stableswap.BondedPositions.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, position.Balance, math.LegacyNewDec(20))

	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)})
	response, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(10)),
			sdk.NewCoin("uusdn", math.NewInt(10)),
		),
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(20), response.MintedShares)

	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 2, 1, 1, time.UTC)})
	response, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(10)),
			sdk.NewCoin("uusdn", math.NewInt(10)),
		),
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(20), response.MintedShares)

	stableswapPool, _ := k.Stableswap.Pools.Get(ctx, 0)
	assert.Equal(t, math.LegacyNewDec(60), stableswapPool.TotalShares)
}

func TestRemoveLiquiditySingleUser(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)

	stableswapServer := stableswapkeeper.NewStableSwapMsgServer(k)
	user := utils.TestAccount()

	_, _ = stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
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

	pool, _ := k.Pools.Get(ctx, 0)
	poolAddress, _ := sdk.AccAddressFromBech32(pool.Address)
	poolLiquidity := bank.GetAllBalances(ctx, poolAddress)
	assert.Equal(t, len(poolLiquidity), 0)

	_, err := stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     user.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(10),
	})
	assert.Error(t, err)

	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)})

	// 1.

	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdc", math.NewInt(100)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdn", math.NewInt(100)))
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(10)),
			sdk.NewCoin("uusdn", math.NewInt(10)),
		),
	})
	assert.NoError(t, err)
	userTotalShares, _ := k.Stableswap.UsersTotalBondedShares.Get(ctx, collections.Join(pool.Id, user.Address))
	stableswapPool, _ := k.Stableswap.Pools.Get(ctx, 0)
	assert.Equal(t, math.LegacyNewDec(20), stableswapPool.TotalShares)
	assert.Equal(t, userTotalShares, userTotalShares)
	assert.Equal(t, math.NewInt(90), bank.Balances[user.Address].AmountOf("uusdc"))
	assert.Equal(t, math.NewInt(90), bank.Balances[user.Address].AmountOf("uusdn"))

	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 2, 1, 1, 1, 1, 1, time.UTC)})
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(20)),
			sdk.NewCoin("uusdn", math.NewInt(20)),
		),
	})
	assert.NoError(t, err)
	userTotalShares, _ = k.Stableswap.UsersTotalBondedShares.Get(ctx, collections.Join(pool.Id, user.Address))
	stableswapPool, _ = k.Stableswap.Pools.Get(ctx, 0)
	assert.Equal(t, math.LegacyNewDec(60), stableswapPool.TotalShares)
	assert.Equal(t, userTotalShares, stableswapPool.TotalShares)
	assert.Equal(t, math.NewInt(70), bank.Balances[user.Address].AmountOf("uusdc"))
	assert.Equal(t, math.NewInt(70), bank.Balances[user.Address].AmountOf("uusdn"))

	_, err = stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     user.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(120),
	})
	assert.ErrorIs(t, err, types.ErrInvalidUnbondPercentage)

	res, err := stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     user.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(100),
	})
	assert.NoError(t, err)

	totalUnbondingShares, err := k.Stableswap.PoolsTotalUnbondingShares.Get(ctx, pool.Id)
	assert.NoError(t, err)
	assert.Equal(t, res.UnbondingShares, totalUnbondingShares)

	userTotalShares, _ = k.Stableswap.UsersTotalBondedShares.Get(ctx, collections.Join(pool.Id, user.Address))
	stableswapPool, _ = k.Stableswap.Pools.Get(ctx, 0)
	assert.Equal(t, userTotalShares, stableswapPool.TotalShares)
	assert.Equal(t, math.NewInt(70), bank.Balances[user.Address].AmountOf("uusdc"))
	assert.Equal(t, math.NewInt(70), bank.Balances[user.Address].AmountOf("uusdn"))
	err = k.BeginBlocker(ctx)
	assert.NoError(t, err)
	assert.Equal(t, math.NewInt(70), bank.Balances[user.Address].AmountOf("uusdc"))
	assert.Equal(t, math.NewInt(70), bank.Balances[user.Address].AmountOf("uusdn"))

	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 2, 10, 1, 1, 1, 1, time.UTC)})
	err = k.BeginBlocker(ctx)
	assert.NoError(t, err)
	assert.Equal(t, math.NewInt(100), bank.Balances[user.Address].AmountOf("uusdc"))
	assert.Equal(t, math.NewInt(100), bank.Balances[user.Address].AmountOf("uusdn"))
	stableswapPool, _ = k.Stableswap.Pools.Get(ctx, 0)
	assert.Equal(t, math.LegacyZeroDec(), stableswapPool.TotalShares)

	totalUnbondingShares, err = k.Stableswap.PoolsTotalUnbondingShares.Get(ctx, pool.Id)
	assert.NoError(t, err)
	assert.Equal(t, math.LegacyZeroDec(), totalUnbondingShares)

	// 2.

	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 3, 1, 1, 1, 1, 1, time.UTC)})
	resAdd, err := stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(20)),
			sdk.NewCoin("uusdn", math.NewInt(20)),
		),
	})
	assert.NoError(t, err)
	userTotalShares, _ = k.Stableswap.UsersTotalBondedShares.Get(ctx, collections.Join(pool.Id, user.Address))
	stableswapPool, _ = k.Stableswap.Pools.Get(ctx, 0)
	assert.Equal(t, math.LegacyNewDec(resAdd.MintedShares), stableswapPool.TotalShares)
	assert.Equal(t, userTotalShares, stableswapPool.TotalShares)
	assert.Equal(t, math.NewInt(80), bank.Balances[user.Address].AmountOf("uusdc"))
	assert.Equal(t, math.NewInt(80), bank.Balances[user.Address].AmountOf("uusdn"))

	_, err = stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     user.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(100),
	})
	assert.NoError(t, err)
	userTotalShares, _ = k.Stableswap.UsersTotalBondedShares.Get(ctx, collections.Join(pool.Id, user.Address))
	stableswapPool, _ = k.Stableswap.Pools.Get(ctx, 0)
	assert.Equal(t, userTotalShares, stableswapPool.TotalShares)
	assert.Equal(t, math.NewInt(80), bank.Balances[user.Address].AmountOf("uusdc"))
	assert.Equal(t, math.NewInt(80), bank.Balances[user.Address].AmountOf("uusdn"))
	err = k.BeginBlocker(ctx)
	assert.NoError(t, err)
	assert.Equal(t, math.NewInt(80), bank.Balances[user.Address].AmountOf("uusdc"))
	assert.Equal(t, math.NewInt(80), bank.Balances[user.Address].AmountOf("uusdn"))

	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 3, 10, 1, 1, 1, 1, time.UTC)})
	err = k.BeginBlocker(ctx)
	assert.NoError(t, err)
	assert.Equal(t, math.NewInt(100), bank.Balances[user.Address].AmountOf("uusdc"))
	assert.Equal(t, math.NewInt(100), bank.Balances[user.Address].AmountOf("uusdn"))
	stableswapPool, _ = k.Stableswap.Pools.Get(ctx, 0)
	assert.Equal(t, math.LegacyZeroDec(), stableswapPool.TotalShares)

	// 3.

	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 4, 1, 1, 1, 1, 1, time.UTC)})
	resAdd, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(20)),
			sdk.NewCoin("uusdn", math.NewInt(20)),
		),
	})
	assert.NoError(t, err)
	userTotalShares, _ = k.Stableswap.UsersTotalBondedShares.Get(ctx, collections.Join(pool.Id, user.Address))
	stableswapPool, _ = k.Stableswap.Pools.Get(ctx, 0)
	assert.Equal(t, math.LegacyNewDec(resAdd.MintedShares), stableswapPool.TotalShares)
	assert.Equal(t, userTotalShares, stableswapPool.TotalShares)
	assert.Equal(t, math.NewInt(80), bank.Balances[user.Address].AmountOf("uusdc"))
	assert.Equal(t, math.NewInt(80), bank.Balances[user.Address].AmountOf("uusdn"))

	expectedRemovedAmount := resAdd.MintedShares * 40 / 100
	expectedRemovedTokens := int64(20 * 40 / 100)
	resRemove, err := stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     user.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(40),
	})
	assert.NoError(t, err)
	userTotalShares, _ = k.Stableswap.UsersTotalBondedShares.Get(ctx, collections.Join(pool.Id, user.Address))
	stableswapPool, _ = k.Stableswap.Pools.Get(ctx, 0)
	assert.Equal(t, math.LegacyNewDec(expectedRemovedAmount), resRemove.UnbondingShares)
	assert.Equal(t, math.LegacyNewDec(resAdd.MintedShares-expectedRemovedAmount), math.LegacyNewDec(24))
	assert.Equal(t, userTotalShares, stableswapPool.TotalShares)
	assert.Equal(t, math.NewInt(80), bank.Balances[user.Address].AmountOf("uusdc"))
	assert.Equal(t, math.NewInt(80), bank.Balances[user.Address].AmountOf("uusdn"))
	err = k.BeginBlocker(ctx)
	assert.NoError(t, err)
	assert.Equal(t, math.NewInt(80), bank.Balances[user.Address].AmountOf("uusdc"))
	assert.Equal(t, math.NewInt(80), bank.Balances[user.Address].AmountOf("uusdn"))

	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 4, 10, 1, 1, 1, 1, time.UTC)})
	err = k.BeginBlocker(ctx)
	assert.NoError(t, err)
	assert.Equal(t, math.NewInt(80+expectedRemovedTokens), bank.Balances[user.Address].AmountOf("uusdc"))
	assert.Equal(t, math.NewInt(80+expectedRemovedTokens), bank.Balances[user.Address].AmountOf("uusdn"))
	stableswapPool, _ = k.Stableswap.Pools.Get(ctx, 0)
	assert.Equal(t, math.LegacyNewDec(40-expectedRemovedAmount), k.Stableswap.GetUserTotalBondedShares(ctx, 0, user.Address))
	assert.Equal(t, math.LegacyNewDec(40-expectedRemovedAmount), stableswapPool.TotalShares)
}

func TestRemoveLiquidityMultiUser(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)

	stableswapServer := stableswapkeeper.NewStableSwapMsgServer(k)
	bob, alice := utils.TestAccount(), utils.TestAccount()
	bobLiquidity, aliceLiquidity := int64(1_000_000_000), int64(500_000_000)
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdc", math.NewInt(bobLiquidity)))
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdn", math.NewInt(bobLiquidity)))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], sdk.NewCoin("uusdc", math.NewInt(aliceLiquidity)))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], sdk.NewCoin("uusdn", math.NewInt(aliceLiquidity)))

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

	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 3, 10, 1, 1, 1, 1, time.UTC)})
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: bob.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(1000000000)),
			sdk.NewCoin("uusdn", math.NewInt(1000000000)),
		),
	})
	assert.NoError(t, err)
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: alice.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(1000000)),
			sdk.NewCoin("uusdn", math.NewInt(1000000)),
		),
	})
	assert.NoError(t, err)

	_, err = stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     bob.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(100),
	})
	assert.NoError(t, err)
	_, err = stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     alice.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(100),
	})
	assert.NoError(t, err)
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 3, 20, 1, 1, 1, 1, time.UTC)})
	err = k.BeginBlocker(ctx)
	assert.NoError(t, err)

	assert.Equal(t, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(bobLiquidity)), sdk.NewCoin("uusdn", math.NewInt(bobLiquidity))), bank.Balances[bob.Address])
	assert.Equal(t, sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(aliceLiquidity)), sdk.NewCoin("uusdn", math.NewInt(aliceLiquidity))), bank.Balances[alice.Address])
}

func TestUnbondings(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	stableswapServer := stableswapkeeper.NewStableSwapMsgServer(k)

	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		RewardsFee:            4e3,
		ProtocolFeePercentage: 1,
		MaxFee:                1,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           1893452400,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})
	assert.Nil(t, err)

	accounts := map[string]utils.Account{}
	N := 1_00
	for i := 0; i < N; i++ {
		user := utils.TestAccount()
		accounts[user.Address] = user
		amount := math.NewInt(rand.Int64N(1_000_000_000000) + 1_000)
		bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdc", amount))
		bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdn", amount))
	}

	balancesClone := map[string]sdk.Coins{}
	for key, value := range bank.Balances {
		balancesClone[key] = value
	}

	totalSupply := math.ZeroInt()
	for address := range accounts {
		totalSupply = totalSupply.Add(bank.Balances[address][0].Amount)
	}

	// ARRANGE: bond the total liquidity of all the users with multiple AddLiquidity transactions
	for _, user := range accounts {
		x := rand.IntN(10-1) + 1
		amount := bank.Balances[user.Address][0].Amount.QuoRaw(int64(x))
		for i := 0; i < x; i++ {
			if i == x-1 {
				amount = bank.Balances[user.Address][0].Amount
			}

			ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC).Add(time.Duration(i) * time.Second)})
			_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
				Signer: user.Address,
				PoolId: 0,
				Amount: sdk.NewCoins(
					sdk.NewCoin("uusdc", amount),
					sdk.NewCoin("uusdn", amount),
				),
			})
			assert.NoError(t, err)
		}
	}

	// ASSERT: All the users have bonded their full balances
	for _, user := range accounts {
		assert.True(t, bank.Balances[user.Address].IsZero())
	}
	// ASSERT: The bonded liquidity matches the total supply
	pool, err := k.Pools.Get(ctx, 0)
	assert.Nil(t, err)
	assert.Equal(t, totalSupply, bank.Balances[pool.Address][0].Amount)

	// ARRANGE: unbond all the users liquidity with multiple RemoveLiquidity transactions using different percentages
	percentages := []math.LegacyDec{
		math.LegacyMustNewDecFromStr("0.1"),
		math.LegacyMustNewDecFromStr("0.9"),
		math.LegacyNewDec(4),
		math.LegacyNewDec(10),
		math.LegacyNewDec(26),
		math.LegacyNewDec(50),
		math.LegacyNewDec(100),
	}
	fakeTime := time.Date(2020, 2, 1, 1, 1, 1, 1, time.UTC)
	for _, user := range accounts {
		for _, percentage := range percentages {
			fakeTime = fakeTime.Add(1 * time.Second) // increase the timer
			ctx = ctx.WithHeaderInfo(header.Info{Time: fakeTime})
			_, err = stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
				Signer:     user.Address,
				PoolId:     0,
				Percentage: percentage,
			})
			assert.NoError(t, err)
		}
	}

	// ASSERT: matching total shares
	cumulativeTotUnbondingShares := math.LegacyZeroDec()
	itr, err := k.Stableswap.UnbondingPositions.Iterate(ctx, nil)
	if err != nil {
		return
	}
	for ; itr.Valid(); itr.Next() {
		key, _ := itr.Key()
		unbondingPosition, _ := k.Stableswap.UnbondingPositions.Get(ctx, key)
		cumulativeTotUnbondingShares = cumulativeTotUnbondingShares.Add(unbondingPosition.Shares)
	}

	itr2, err := k.Stableswap.UsersTotalUnbondingShares.Iterate(ctx, nil)
	if err != nil {
		return
	}
	cumulativeUsersTotUnbondingShares := math.LegacyZeroDec()
	for ; itr2.Valid(); itr2.Next() {
		key, _ := itr2.Key()
		amount, _ := k.Stableswap.UsersTotalUnbondingShares.Get(ctx, key)
		cumulativeUsersTotUnbondingShares = cumulativeUsersTotUnbondingShares.Add(amount)
	}

	itr3, err := k.Stableswap.BondedPositions.Iterate(ctx, nil)
	if err != nil {
		return
	}
	cumulativeTotShares := math.LegacyZeroDec()
	for ; itr3.Valid(); itr3.Next() {
		key, _ := itr3.Key()
		amount, _ := k.Stableswap.BondedPositions.Get(ctx, key)
		cumulativeTotShares = cumulativeTotShares.Add(amount.Balance)
	}
	assert.Equal(t, cumulativeTotUnbondingShares, cumulativeUsersTotUnbondingShares)

	for _, u := range accounts {
		assert.True(t, bank.Balances[u.Address].IsZero())
	}

	// ARRANGE: execute the BeginBlocker
	ctx = ctx.WithHeaderInfo(header.Info{Time: fakeTime.Add(72 * time.Hour)})
	err = k.BeginBlocker(ctx)
	assert.NoError(t, err)

	for _, u := range accounts {
		diff1 := balancesClone[u.Address].AmountOf("uusdc").Sub(bank.Balances[u.Address].AmountOf("uusdc")).Int64()
		diff2 := balancesClone[u.Address].AmountOf("uusdn").Sub(bank.Balances[u.Address].AmountOf("uusdn")).Int64()

		delta := int64(10)
		assert.LessOrEqual(t, diff1, delta)
		assert.LessOrEqual(t, diff2, delta)
	}
}
