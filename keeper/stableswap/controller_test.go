package stableswap_test

import (
	"testing"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"swap.noble.xyz/keeper"
	"swap.noble.xyz/keeper/stableswap"
	"swap.noble.xyz/types"
	stableswaptypes "swap.noble.xyz/types/stableswap"
	"swap.noble.xyz/utils"
	"swap.noble.xyz/utils/mocks"
)

func TestController(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	stableswapServer := keeper.NewStableSwapMsgServer(k)
	_, err := stableswapServer.CreatePool(ctx, &stableswaptypes.MsgCreatePool{
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
	assert.NoError(t, err)

	// ACT: Get a valid pool.
	_, err = keeper.GetStableSwapController(ctx, k, 0)
	assert.NoError(t, err)

	// ACT: Attempt to Get a non-existing pool.
	_, err = keeper.GetStableSwapController(ctx, k, 1)
	assert.Error(t, err)

	// ARRANGE: Set up failing collections for the StableSwap Pool.
	k.Stableswap.Pools = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Get, utils.GetKVStore(ctx, types.ModuleName))),
		types.StableSwapPoolsPrefix, "stableswap_pools", collections.Uint64Key, codec.CollValue[stableswaptypes.Pool](mocks.MakeTestEncodingConfig("noble").Codec),
	)

	// ACT: Attempt to Get a non-existing pool.
	_, err = keeper.GetStableSwapController(ctx, k, 0)
	assert.Error(t, err)

	// ARRANGE: Set up failing collections for the Generic Pool.
	k.Pools = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Get, utils.GetKVStore(ctx, types.ModuleName))),
		types.PoolsPrefix, "pools_generic", collections.Uint64Key, codec.CollValue[types.Pool](mocks.MakeTestEncodingConfig("noble").Codec),
	)

	// ACT: Attempt to Get a non-existing pool.
	_, err = keeper.GetStableSwapController(ctx, k, 0)
	assert.Error(t, err)
}

func TestComputeWeightedPoolUnbondingPeriod(t *testing.T) {
	// Test cases.
	tests := []struct {
		name           string
		totalShares    math.LegacyDec
		sharesToUnbond math.LegacyDec
		expectedTime   time.Duration
		expectError    bool
	}{
		{
			name:           "Base duration when percentage is below 0.1%",
			totalShares:    math.LegacyNewDec(10000),
			sharesToUnbond: math.LegacyNewDec(5), // 0.05%
			expectedTime:   1 * time.Minute,
			expectError:    false,
		},
		{
			name:           "Short duration when percentage is between 0.1% and 1%",
			totalShares:    math.LegacyNewDec(10000),
			sharesToUnbond: math.LegacyNewDec(15), // 0.15%
			expectedTime:   30 * time.Minute,
			expectError:    false,
		},
		{
			name:           "Medium duration when percentage is between 1% and 10%",
			totalShares:    math.LegacyNewDec(10000),
			sharesToUnbond: math.LegacyNewDec(150), // 1.5%
			expectedTime:   12 * time.Hour,
			expectError:    false,
		},
		{
			name:           "Long duration when percentage is above 10%",
			totalShares:    math.LegacyNewDec(10000),
			sharesToUnbond: math.LegacyNewDec(1200), // 12%
			expectedTime:   24 * time.Hour,
			expectError:    false,
		},
		{
			name:           "Error on zero total shares",
			totalShares:    math.LegacyNewDec(0),
			sharesToUnbond: math.LegacyNewDec(100),
			expectedTime:   1 * time.Minute,
			expectError:    true,
		},
		{
			name:           "Error on zero shares to unbond",
			totalShares:    math.LegacyNewDec(10000),
			sharesToUnbond: math.LegacyNewDec(0),
			expectedTime:   1 * time.Minute,
			expectError:    true,
		},
		{
			name:           "Max values",
			totalShares:    math.LegacyNewDec(1e18),
			sharesToUnbond: math.LegacyNewDec(1),
			expectedTime:   1 * time.Minute,
			expectError:    false,
		},
		{
			name:           "Max values",
			totalShares:    math.LegacyNewDec(1),
			sharesToUnbond: math.LegacyNewDec(1e18),
			expectedTime:   24 * time.Hour,
			expectError:    false,
		},
		{
			name:           "Max values",
			totalShares:    math.LegacyNewDec(1e18),
			sharesToUnbond: math.LegacyNewDec(1e18),
			expectedTime:   24 * time.Hour,
			expectError:    false,
		},
	}

	// Execute each test case.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT: Compute the unbonding period.
			result, err := stableswap.ComputeWeightedPoolUnbondingPeriod(tt.totalShares, tt.sharesToUnbond)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTime, result)
			}
		})
	}
}

func TestCalculatePositionRewards(t *testing.T) {
	// Define common values for all test cases
	currentTime := time.Now()
	poolCreation := currentTime.Add(-time.Hour * 24 * 30) // Pool created 30 days ago
	totalShares := math.LegacyNewDec(1000)                // Total shares in the pool

	// Test cases
	tests := []struct {
		name           string
		position       stableswaptypes.BondedPosition
		poolRewards    sdk.Coins
		expectedReward sdk.Coins
	}{
		{
			name: "User with 10% share for 15 days",
			position: stableswaptypes.BondedPosition{
				Balance:            math.LegacyNewDec(100),                // 10% of total shares
				RewardsPeriodStart: currentTime.Add(-time.Hour * 24 * 15), // BondedPosition held for 15 days
			},
			poolRewards:    sdk.NewCoins(sdk.NewCoin("token", math.NewInt(1000))), // 1000 tokens in pool rewards
			expectedReward: sdk.NewCoins(sdk.NewCoin("token", math.NewInt(50))),   // Expected reward: 50 tokens
		},
		{
			name: "User with 5% share for 10 days",
			position: stableswaptypes.BondedPosition{
				Balance:            math.LegacyNewDec(50),                 // 5% of total shares
				RewardsPeriodStart: currentTime.Add(-time.Hour * 24 * 10), // BondedPosition held for 10 days
			},
			poolRewards:    sdk.NewCoins(sdk.NewCoin("token", math.NewInt(1000))), // 1000 tokens in pool rewards
			expectedReward: sdk.NewCoins(sdk.NewCoin("token", math.NewInt(16))),   // Expected reward: 5 tokens
		},
		{
			name: "User with 20% share for entire pool duration (30 days)",
			position: stableswaptypes.BondedPosition{
				Balance:            math.LegacyNewDec(200), // 20% of total shares
				RewardsPeriodStart: poolCreation,           // BondedPosition held from pool creation
			},
			poolRewards:    sdk.NewCoins(sdk.NewCoin("token", math.NewInt(1000))), // 1000 tokens in pool rewards
			expectedReward: sdk.NewCoins(sdk.NewCoin("token", math.NewInt(200))),  // Expected reward: 200 tokens
		},
		{
			name: "User with full share for partial period",
			position: stableswaptypes.BondedPosition{
				Balance:            totalShares,                          // 100% of total shares
				RewardsPeriodStart: currentTime.Add(-time.Hour * 24 * 5), // BondedPosition held for 5 days
			},
			poolRewards:    sdk.NewCoins(sdk.NewCoin("token", math.NewInt(500))), // 500 tokens in pool rewards
			expectedReward: sdk.NewCoins(sdk.NewCoin("token", math.NewInt(83))),  // Expected reward: 83 tokens
		},
		{
			name: "User with full share for full period",
			position: stableswaptypes.BondedPosition{
				Balance:            totalShares,  // 100% of total shares
				RewardsPeriodStart: poolCreation, // BondedPosition held for 5 days
			},
			poolRewards:    sdk.NewCoins(sdk.NewCoin("token", math.NewInt(500))), // 500 tokens in pool rewards
			expectedReward: sdk.NewCoins(sdk.NewCoin("token", math.NewInt(500))), // Expected reward: 500 tokens
		},
	}

	// Execute each test case.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ACT: Compute the position rewards.
			reward, err := stableswap.CalculatePositionRewards(currentTime, tt.poolRewards, tt.position, totalShares, poolCreation)
			require.NoError(t, err)
			require.Equal(t, tt.expectedReward, reward)
		})
	}
}
