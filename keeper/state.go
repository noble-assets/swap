// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, NASD Inc. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package keeper

import (
	"context"

	"swap.noble.xyz/types"
)

//

// GetPaused retrieves the paused state for all pools.
func (k *Keeper) GetPaused(ctx context.Context) map[uint64]bool {
	paused := map[uint64]bool{}
	_ = k.Paused.Walk(ctx, nil, func(key uint64, value bool) (stop bool, err error) {
		paused[key] = value
		return false, nil
	})
	return paused
}

// IsPaused checks if a specific pool is in a paused state.
func (k *Keeper) IsPaused(ctx context.Context, poolId uint64) bool {
	isPaused, _ := k.Paused.Get(ctx, poolId)
	return isPaused
}

// SetPaused updates the paused state of a specific pool.
func (k *Keeper) SetPaused(ctx context.Context, poolId uint64, value bool) error {
	return k.Paused.Set(ctx, poolId, value)
}

//

// GetNextPoolID retrieves the next pool ID from the state.
func (k *Keeper) GetNextPoolID(ctx context.Context) uint64 {
	nextPoolID, _ := k.NextPoolID.Peek(ctx)
	return nextPoolID
}

// IncreaseNextPoolID increments the next pool ID in the state.
func (k *Keeper) IncreaseNextPoolID(ctx context.Context) (uint64, error) {
	return k.NextPoolID.Next(ctx)
}

// SetNextPoolID sets the next pool ID to a specified value.
func (k *Keeper) SetNextPoolID(ctx context.Context, poolId uint64) error {
	return k.NextPoolID.Set(ctx, poolId)
}

//

// GetPools retrieves all generic pools from the state.
func (k *Keeper) GetPools(ctx context.Context) map[uint64]types.Pool {
	pools := map[uint64]types.Pool{}
	_ = k.Pools.Walk(ctx, nil, func(key uint64, value types.Pool) (stop bool, err error) {
		pools[key] = value
		return false, nil
	})
	return pools
}

// HasPool checks if a specific pool exists in the state.
func (k *Keeper) HasPool(ctx context.Context, poolId uint64) bool {
	has, _ := k.Pools.Has(ctx, poolId)
	return has
}

// GetPool retrieves a specific pool by its ID from the state.
func (k *Keeper) GetPool(ctx context.Context, poolId uint64) (types.Pool, error) {
	return k.Pools.Get(ctx, poolId)
}

// SetPool updates or creates a pool in the state with a given ID.
func (k *Keeper) SetPool(ctx context.Context, poolId uint64, pool types.Pool) error {
	return k.Pools.Set(ctx, poolId, pool)
}
