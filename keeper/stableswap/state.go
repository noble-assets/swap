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

package stableswap

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"swap.noble.xyz/types/stableswap"
)

//

// GetPools retrieves all StableSwap pools from the state.
func (k *Keeper) GetPools(ctx context.Context) map[uint64]stableswap.Pool {
	pools := map[uint64]stableswap.Pool{}
	_ = k.Pools.Walk(ctx, nil, func(key uint64, value stableswap.Pool) (stop bool, err error) {
		pools[key] = value
		return false, nil
	})
	return pools
}

// GetPool retrieves a StableSwap pool by its ID.
func (k *Keeper) GetPool(ctx context.Context, poolId uint64) (stableswap.Pool, error) {
	return k.Pools.Get(ctx, poolId)
}

// SetPool stores a StableSwap pool in the state by its ID.
func (k *Keeper) SetPool(ctx context.Context, poolId uint64, pool stableswap.Pool) error {
	return k.Pools.Set(ctx, poolId, pool)
}

//

// GetPoolsTotalUnbondingShares retrieves the total unbonding shares for all pools.
func (k *Keeper) GetPoolsTotalUnbondingShares(ctx context.Context) []stableswap.PoolsTotalUnbondingSharesEntry {
	var pools []stableswap.PoolsTotalUnbondingSharesEntry
	_ = k.PoolsTotalUnbondingShares.Walk(ctx, nil, func(key uint64, value math.LegacyDec) (stop bool, err error) {
		pools = append(pools, stableswap.PoolsTotalUnbondingSharesEntry{
			PoolId: key,
			Shares: value,
		})
		return false, nil
	})
	return pools
}

// HasPoolTotalUnbondingShares checks if total unbonding shares exist for a specific pool.
func (k *Keeper) HasPoolTotalUnbondingShares(ctx context.Context, poolId uint64) bool {
	has, _ := k.PoolsTotalUnbondingShares.Has(ctx, poolId)
	return has
}

// GetPoolTotalUnbondingShares retrieves the total unbonding shares for a specific pool by ID.
func (k *Keeper) GetPoolTotalUnbondingShares(ctx context.Context, poolId uint64) math.LegacyDec {
	totalUnbondingShares, _ := k.PoolsTotalUnbondingShares.Get(ctx, poolId)
	return totalUnbondingShares
}

// SetPoolTotalUnbondingShares updates the total unbonding shares for a specific pool.
func (k *Keeper) SetPoolTotalUnbondingShares(ctx context.Context, poolId uint64, value math.LegacyDec) error {
	return k.PoolsTotalUnbondingShares.Set(ctx, poolId, value)
}

//

// GetUsersTotalBondedShares retrieves total bonded shares for all users.
func (k *Keeper) GetUsersTotalBondedShares(ctx context.Context) []stableswap.UsersTotalBondedSharesEntry {
	var entries []stableswap.UsersTotalBondedSharesEntry
	_ = k.UsersTotalBondedShares.Walk(ctx, nil, func(key collections.Pair[uint64, string], value math.LegacyDec) (stop bool, err error) {
		entries = append(entries, stableswap.UsersTotalBondedSharesEntry{
			PoolId:  key.K1(),
			Address: key.K2(),
			Shares:  value,
		})
		return false, nil
	})
	return entries
}

// HasUserTotalBondedShares checks if bonded shares exist for a user in a specific pool.
func (k *Keeper) HasUserTotalBondedShares(ctx context.Context, poolId uint64, address string) bool {
	has, _ := k.UsersTotalBondedShares.Has(ctx, collections.Join(poolId, address))
	return has
}

// GetUserTotalBondedShares retrieves the bonded shares for a user in a specific pool.
func (k *Keeper) GetUserTotalBondedShares(ctx context.Context, poolId uint64, address string) math.LegacyDec {
	userTotalBondedShares, _ := k.UsersTotalBondedShares.Get(ctx, collections.Join(poolId, address))
	return userTotalBondedShares
}

// SetUserTotalBondedShares updates the bonded shares for a user in a specific pool.
func (k *Keeper) SetUserTotalBondedShares(ctx context.Context, poolId uint64, address string, value math.LegacyDec) error {
	return k.UsersTotalBondedShares.Set(ctx, collections.Join(poolId, address), value)
}

//

// GetUsersTotalUnbondingShares retrieves total unbonding shares for all users.
func (k *Keeper) GetUsersTotalUnbondingShares(ctx context.Context) []stableswap.UsersTotalUnbondingSharesEntry {
	var entries []stableswap.UsersTotalUnbondingSharesEntry
	_ = k.UsersTotalUnbondingShares.Walk(ctx, nil, func(key collections.Pair[uint64, string], value math.LegacyDec) (stop bool, err error) {
		entries = append(entries, stableswap.UsersTotalUnbondingSharesEntry{
			PoolId:  key.K1(),
			Address: key.K2(),
			Shares:  value,
		})
		return false, nil
	})
	return entries
}

// HasUserTotalUnbondingShares checks if unbonding shares exist for a user in a specific pool.
func (k *Keeper) HasUserTotalUnbondingShares(ctx context.Context, poolId uint64, address string) bool {
	has, _ := k.UsersTotalUnbondingShares.Has(ctx, collections.Join(poolId, address))
	return has
}

// GetUserTotalUnbondingShares retrieves the unbonding shares for a user in a specific pool.
func (k *Keeper) GetUserTotalUnbondingShares(ctx context.Context, poolId uint64, address string) math.LegacyDec {
	userTotalUnbondingShares, _ := k.UsersTotalUnbondingShares.Get(ctx, collections.Join(poolId, address))
	return userTotalUnbondingShares
}

// SetUserTotalUnbondingShares updates the unbonding shares for a user in a specific pool.
func (k *Keeper) SetUserTotalUnbondingShares(ctx context.Context, poolId uint64, address string, value math.LegacyDec) error {
	return k.UsersTotalUnbondingShares.Set(ctx, collections.Join(poolId, address), value)
}

//

// GetBondedPositions retrieves all bonded positions from the state.
func (k *Keeper) GetBondedPositions(ctx context.Context) []stableswap.BondedPositionEntry {
	var entries []stableswap.BondedPositionEntry
	_ = k.BondedPositions.Walk(ctx, nil, func(key collections.Triple[uint64, string, int64], value stableswap.BondedPosition) (stop bool, err error) {
		entries = append(entries, stableswap.BondedPositionEntry{
			PoolId:         key.K1(),
			Address:        key.K2(),
			Timestamp:      key.K3(),
			BondedPosition: value,
		})
		return false, nil
	})
	return entries
}

// HasBondedPosition checks if a bonded position exists for a user in a specific pool.
func (k *Keeper) HasBondedPosition(ctx context.Context, poolId uint64, address string, timestamp int64) bool {
	has, _ := k.BondedPositions.Has(ctx, collections.Join3(poolId, address, timestamp))
	return has
}

// GetBondedPositionsByProvider retrieves all bonded positions by a specific provider.
func (k *Keeper) GetBondedPositionsByProvider(ctx context.Context, provider string) []stableswap.BondedPositionEntry {
	var entries []stableswap.BondedPositionEntry
	itr, err := k.BondedPositions.Indexes.ByProvider.MatchExact(ctx, provider)
	if err != nil {
		return nil
	}

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.PrimaryKey()
		value, _ := k.BondedPositions.Get(ctx, key)
		entries = append(entries, stableswap.BondedPositionEntry{
			PoolId:         key.K1(),
			Address:        key.K2(),
			Timestamp:      key.K3(),
			BondedPosition: value,
		})
	}
	return entries
}

// GetBondedPositionsByPoolAndProvider retrieves all bonded positions by a provider within a specific pool.
func (k *Keeper) GetBondedPositionsByPoolAndProvider(ctx context.Context, poolId uint64, provider string) []stableswap.BondedPositionEntry {
	var entries []stableswap.BondedPositionEntry
	itr, err := k.BondedPositions.Indexes.ByPoolAndProvider.MatchExact(ctx, collections.Join(poolId, provider))
	if err != nil {
		return nil
	}

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.PrimaryKey()
		value, _ := k.BondedPositions.Get(ctx, key)
		entries = append(entries, stableswap.BondedPositionEntry{
			PoolId:         key.K1(),
			Address:        key.K2(),
			Timestamp:      key.K3(),
			BondedPosition: value,
		})
	}
	return entries
}

// SetBondedPosition updates a bonded position for a user in a specific pool.
func (k *Keeper) SetBondedPosition(ctx context.Context, poolId uint64, address string, timestamp int64, value stableswap.BondedPosition) error {
	return k.BondedPositions.Set(ctx, collections.Join3(poolId, address, timestamp), value)
}

// RemoveBondedPosition deletes a bonded position for a user in a specific pool.
func (k *Keeper) RemoveBondedPosition(ctx context.Context, poolId uint64, address string, timestamp int64) error {
	return k.BondedPositions.Remove(ctx, collections.Join3(poolId, address, timestamp))
}

//

// GetUnbondingPositions retrieves all unbonding positions from the state.
func (k *Keeper) GetUnbondingPositions(ctx context.Context) []stableswap.UnbondingPositionEntry {
	var entries []stableswap.UnbondingPositionEntry
	_ = k.UnbondingPositions.Walk(ctx, nil, func(key collections.Triple[int64, string, uint64], value stableswap.UnbondingPosition) (stop bool, err error) {
		entries = append(entries, stableswap.UnbondingPositionEntry{
			Timestamp:         key.K1(),
			Address:           key.K2(),
			PoolId:            key.K3(),
			UnbondingPosition: value,
		})
		return false, nil
	})
	return entries
}

// GetUnbondingPositionsUntil retrieves unbonding positions until a specified timestamp.
func (k *Keeper) GetUnbondingPositionsUntil(ctx context.Context, to int64) []stableswap.UnbondingPositionEntry {
	var entries []stableswap.UnbondingPositionEntry

	_ = k.UnbondingPositions.Walk(
		ctx,
		collections.NewPrefixUntilTripleRange[int64, string, uint64](to),
		func(key collections.Triple[int64, string, uint64], value stableswap.UnbondingPosition) (stop bool, err error) {
			entries = append(entries, stableswap.UnbondingPositionEntry{
				Timestamp:         key.K1(),
				Address:           key.K2(),
				PoolId:            key.K3(),
				UnbondingPosition: value,
			})
			return false, nil
		})
	return entries
}

// HasUnbondingPosition checks if an unbonding position exists for a user in a specific pool.
func (k *Keeper) HasUnbondingPosition(ctx context.Context, timestamp int64, address string, poolId uint64) bool {
	has, _ := k.UnbondingPositions.Has(ctx, collections.Join3(timestamp, address, poolId))
	return has
}

// GetUnbondingPositionsByProvider retrieves unbonding positions by a specific provider.
func (k *Keeper) GetUnbondingPositionsByProvider(ctx context.Context, provider string) []stableswap.UnbondingPositionEntry {
	var entries []stableswap.UnbondingPositionEntry
	itr, err := k.UnbondingPositions.Indexes.ByProvider.MatchExact(ctx, provider)
	if err != nil {
		return nil
	}

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.PrimaryKey()
		value, _ := k.UnbondingPositions.Get(ctx, key)
		entries = append(entries, stableswap.UnbondingPositionEntry{
			Timestamp:         key.K1(),
			Address:           key.K2(),
			PoolId:            key.K3(),
			UnbondingPosition: value,
		})
	}
	return entries
}

// SetUnbondingPosition updates an unbonding position for a user in a specific pool.
func (k *Keeper) SetUnbondingPosition(ctx context.Context, timestamp int64, address string, poolId uint64, value stableswap.UnbondingPosition) error {
	return k.UnbondingPositions.Set(ctx, collections.Join3(timestamp, address, poolId), value)
}

// RemoveUnbondingPosition deletes an unbonding position for a user in a specific pool.
func (k *Keeper) RemoveUnbondingPosition(ctx context.Context, timestamp int64, address string, poolId uint64) error {
	return k.UnbondingPositions.Remove(ctx, collections.Join3(timestamp, address, poolId))
}
