package keeper

import (
	"context"
	"time"

	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	anyproto "github.com/cosmos/gogoproto/types/any"
	"swap.noble.xyz/types"
)

// Controller defines and abstracts the core operations for interacting with a pool using any algorithm,
// offering a unified interface for adding/removing liquidity, performing swaps, retrieving pool information,
// and handling algorithm-specific actions seamlessly across different pool implementations.
// State changes are managed internally within the controller for consistency and atomicity.
type Controller interface {
	// GetId retrieves the unique identifier of the pool.
	GetId() uint64

	// GetAddress retrieves the address associated with the pool.
	GetAddress() string

	// GetAlgorithm retrieves the algorithm type used by the pool.
	GetAlgorithm() types.Algorithm

	// GetPair retrieves the token pair managed by the pool.
	GetPair() string

	// PoolDetails returns detailed information about the StableSwap pool as a serialized `Any` object.
	PoolDetails() *anyproto.Any

	// IsPaused checks if the pool is currently paused.
	IsPaused() bool

	// GetRates computes exchange rates for tokens in the pool based on liquidity.
	GetRates(ctx context.Context) []types.Rate

	// GetRate computes the single exchange rate for the base token pair in the pool.
	GetRate(ctx context.Context) math.LegacyDec

	// GetLiquidity retrieves the total liquidity in the pool.
	GetLiquidity(ctx context.Context) sdk.Coins

	// GetProtocolFeesAddresses retrieves the addresses where protocol fees are collected.
	GetProtocolFeesAddresses() []sdk.AccAddress

	// Swap performs a coin swap within a specified pool and its underlying algorithm.
	Swap(
		ctx context.Context,
		currentTime int64,
		coin sdk.Coin,
		denomTo string,
	) (*types.SwapCommitment, error)

	// ProcessUserRewards distributes rewards to a user.
	ProcessUserRewards(ctx context.Context, address string, currentTime time.Time) (sdk.Coins, error)
}

// GetGenericController initializes and returns the Generic Controller for the specified Pool ID.
func GetGenericController(ctx context.Context, keeper *Keeper, poolId uint64) (Controller, error) {
	// Check if the Pool exists.
	if !keeper.HasPool(ctx, poolId) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPool, "pool %d does not exists", poolId)
	}

	// Get the pool.
	pool, err := keeper.GetPool(ctx, poolId)
	if err != nil {
		return nil, err
	}

	// Select the appropriate controller based on the pool's algorithm.
	switch pool.Algorithm {
	case types.STABLESWAP:
		return GetStableSwapController(ctx, keeper, poolId)
	default:
		// Return an error for unsupported algorithms.
		return nil, sdkerrors.Wrapf(types.ErrInvalidAlgorithm, "unsupported pool with algorithm: %s", pool.Algorithm)
	}
}
