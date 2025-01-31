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

package swap

import (
	"context"

	stableswaptypes "swap.noble.xyz/types/stableswap"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"swap.noble.xyz/keeper"
	"swap.noble.xyz/types"
)

func InitGenesis(ctx context.Context, k *keeper.Keeper, gen types.GenesisState) {
	// Set the NextPoolID State.
	if err := k.SetNextPoolID(ctx, gen.GetNextPoolId()); err != nil {
		panic(err)
	}

	// Set the Pools State.
	for poolId, pool := range gen.GetPools() {
		if err := k.SetPool(ctx, poolId, pool); err != nil {
			panic(err)
		}
	}

	// Set the Paused State.
	for poolId, paused := range gen.GetPaused() {
		if err := k.SetPaused(ctx, poolId, paused); err != nil {
			panic(err)
		}
	}

	// Set the StableSwap State.
	for poolId, pool := range gen.StableswapState.GetPools() {
		if err := k.Stableswap.SetPool(ctx, poolId, pool); err != nil {
			panic(err)
		}
	}

	// Set the UsersTotalBondedShares State.
	for _, entry := range gen.StableswapState.GetPoolsTotalUnbondingShares() {
		if err := k.Stableswap.SetPoolTotalUnbondingShares(ctx, entry.PoolId, entry.Shares); err != nil {
			panic(err)
		}
	}

	// Set the UsersTotalBondedShares State.
	for _, entry := range gen.StableswapState.GetUsersTotalBondedShares() {
		if err := k.Stableswap.SetUserTotalBondedShares(ctx, entry.PoolId, entry.Address, entry.Shares); err != nil {
			panic(err)
		}
	}

	// Set the UsersTotalBondedShares State.
	for _, entry := range gen.StableswapState.GetUsersTotalUnbondingShares() {
		if err := k.Stableswap.SetUserTotalUnbondingShares(ctx, entry.PoolId, entry.Address, entry.Shares); err != nil {
			panic(err)
		}
	}

	// Set the BondedPositions State.
	for _, entry := range gen.StableswapState.GetBondedPositions() {
		if err := k.Stableswap.SetBondedPosition(ctx, entry.PoolId, entry.Address, entry.Timestamp, entry.BondedPosition); err != nil {
			panic(err)
		}
	}

	// Set the UnbondingPositions State.
	for _, entry := range gen.StableswapState.GetUnbondingPositions() {
		if err := k.Stableswap.SetUnbondingPosition(ctx, entry.Timestamp, entry.Address, entry.PoolId, entry.UnbondingPosition); err != nil {
			panic(err)
		}
	}

	//
}

func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		NextPoolId: k.GetNextPoolID(ctx),
		Pools:      k.GetPools(ctx),
		Paused:     k.GetPaused(ctx),
		StableswapState: stableswaptypes.GenesisState{
			Pools:                     k.Stableswap.GetPools(ctx),
			PoolsTotalUnbondingShares: k.Stableswap.GetPoolsTotalUnbondingShares(ctx),
			UsersTotalBondedShares:    k.Stableswap.GetUsersTotalBondedShares(ctx),
			UsersTotalUnbondingShares: k.Stableswap.GetUsersTotalUnbondingShares(ctx),
			BondedPositions:           k.Stableswap.GetBondedPositions(ctx),
			UnbondingPositions:        k.Stableswap.GetUnbondingPositions(ctx),
		},
	}
}
