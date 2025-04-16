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
	"fmt"

	sdkerrors "cosmossdk.io/errors"

	"swap.noble.xyz/keeper/stableswap"
	"swap.noble.xyz/types"
)

// GetStableSwapController initializes and returns the specific `StableSwap` Controller for the specified Pool.
func GetStableSwapController(ctx context.Context, keeper *Keeper, poolId uint64) (*stableswap.Controller, error) {
	// Check if the `StableSwap` Pool exists.
	if !keeper.HasPool(ctx, poolId) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPool, "stableswap pool with Id %d does not exists", poolId)
	}

	// Retrieve the `StableSwap` pool with the given ID.
	pool, err := keeper.GetPool(ctx, poolId)
	if err != nil {
		return nil, err
	}

	// Retrieve if the Pool is paused.
	paused := keeper.IsPaused(ctx, poolId)

	// Retrieve the `StableSwap` Pool.
	stableswapPool, err := keeper.Stableswap.GetPool(ctx, pool.Id)
	if err != nil {
		return nil, err
	}

	// Create and return the `StableSwap` StableswapController.
	stableswapController := stableswap.NewController(
		&keeper.bankKeeper,
		&keeper.addressCodec,
		keeper.baseDenom,
		&pool, paused,
		&stableswapPool,
		keeper.Stableswap,
		keeper.minRemoveLiquidityAmount,
		keeper.maxRemoveLiquidityPositions,
	)
	return &stableswapController, nil
}

// GetStableSwapControllers initializes and returns all the `StableSwap` Controller for the Pools.
func GetStableSwapControllers(ctx context.Context, keeper *Keeper) map[uint64]*stableswap.Controller {
	controllers := map[uint64]*stableswap.Controller{}

	// Iterate through the StableSwap pools.
	for _, pool := range keeper.GetPools(ctx) {
		controller, err := GetStableSwapController(ctx, keeper, pool.Id)
		if err != nil {
			keeper.Stableswap.Logger().Error(fmt.Sprintf("failed to access Pool %d: %s", pool.Id, err.Error()))
			continue
		}

		controllers[pool.GetId()] = controller
	}

	return controllers
}
