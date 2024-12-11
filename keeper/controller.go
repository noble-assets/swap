package keeper

import (
	"context"
	"time"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	anyproto "github.com/cosmos/gogoproto/types/any"
	"swap.noble.xyz/types"
)

// Controller defines and abstracts the core operations for interacting with a pool using any algorithm,
// offering a unified interface for adding/removing liquidity, performing swaps, retrieving pool information,
// and handling algorithm-specific actions seamlessly across different pool implementations.
// State changes are managed internally within the controller for consistency and atomicity.
type Controller interface {
	// GetId returns the Pool id.
	GetId() uint64

	// GetAddress returns as string the Pool address.
	GetAddress() string

	// GetAlgorithm returns the Pool algorithm.
	GetAlgorithm() types.Algorithm

	// GetPair returns the Pool pair denom.
	GetPair() string

	// GetCreationTime returns the time of Pool creation.
	GetCreationTime() time.Time

	// PoolDetails returns the underlying detailed information or metadata about the Pool.
	PoolDetails() *anyproto.Any

	// IsPaused returns true if the Pool is paused.
	IsPaused() bool

	// GetRates retrieves current swap rates for all supported denominations within the Pool.
	GetRates(ctx context.Context) []types.Rate

	// GetLiquidity returns the total liquidity available in the pool as a collection of coins.
	GetLiquidity(ctx context.Context) sdk.Coins

	// GetProtocolFeesAddresses returns the array of addresses containing the protocol fees of the Pool.
	GetProtocolFeesAddresses() []sdk.AccAddress

	// Swap performs a coin swap within a specified pool, exchanging one coin type for another.
	Swap(
		ctx context.Context,
		currentTime int64,
		coin sdk.Coin,
		denomTo string,
	) (*types.SwapCommitment, error)

	// ProcessUserRewards processes the rewards for the given user within the Pool.
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
