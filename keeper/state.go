package keeper

import (
	"context"

	"swap.noble.xyz/types"
)

//

// GetPaused returns all the pools' pause State.
func (k *Keeper) GetPaused(ctx context.Context) map[uint64]bool {
	paused := map[uint64]bool{}
	_ = k.Paused.Walk(ctx, nil, func(key uint64, value bool) (stop bool, err error) {
		paused[key] = value
		return false, nil
	})
	return paused
}

// IsPaused returns the pause State of a specific pool.
func (k *Keeper) IsPaused(ctx context.Context, poolId uint64) bool {
	isPaused, _ := k.Paused.Get(ctx, poolId)
	return isPaused
}

// SetPaused sets the pause State of a specific pool.
func (k *Keeper) SetPaused(ctx context.Context, poolId uint64, value bool) error {
	return k.Paused.Set(ctx, poolId, value)
}

//

// GetNextPoolID returns the next pool ID in the State.
func (k *Keeper) GetNextPoolID(ctx context.Context) uint64 {
	nextPoolID, _ := k.NextPoolID.Peek(ctx)
	return nextPoolID
}

// IncreaseNextPoolID sets the next pool ID to the next value in the State.
func (k *Keeper) IncreaseNextPoolID(ctx context.Context) (uint64, error) {
	return k.NextPoolID.Next(ctx)
}

// SetNextPoolID sets the next pool ID to a given value in the State.
func (k *Keeper) SetNextPoolID(ctx context.Context, poolId uint64) error {
	return k.NextPoolID.Set(ctx, poolId)
}

//

// GetPools returns all the generic pools in the State.
func (k *Keeper) GetPools(ctx context.Context) map[uint64]types.Pool {
	pools := map[uint64]types.Pool{}
	_ = k.Pools.Walk(ctx, nil, func(key uint64, value types.Pool) (stop bool, err error) {
		pools[key] = value
		return false, nil
	})
	return pools
}

// HasPool checks if a Pool exists and returns true if it does.
func (k *Keeper) HasPool(ctx context.Context, poolId uint64) bool {
	has, _ := k.Pools.Has(ctx, poolId)
	return has
}

// GetPool returns a specific generic pool by ID.
func (k *Keeper) GetPool(ctx context.Context, poolId uint64) (types.Pool, error) {
	return k.Pools.Get(ctx, poolId)
}

// SetPool sets a pool value in the State.
func (k *Keeper) SetPool(ctx context.Context, poolId uint64, pool types.Pool) error {
	return k.Pools.Set(ctx, poolId, pool)
}
