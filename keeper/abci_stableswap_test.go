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

package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/header"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"swap.noble.xyz/types"
	"swap.noble.xyz/types/stableswap"
	"swap.noble.xyz/utils"
	"swap.noble.xyz/utils/mocks"
)

func TestStableSwapBeginBlocker(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	user := utils.TestAccount()

	// ARRANGE: Trigger the BeginBlocker with an invalid height.
	ctx = ctx.WithHeaderInfo(header.Info{Height: 1})
	executed := k.StableSwapBeginBlocker(ctx)

	// ASSERT: The execution should've been skipped.
	assert.False(t, executed)

	// ARRANGE: Trigger the BeginBlocker with a valid height.
	ctx = ctx.WithHeaderInfo(header.Info{Height: 10})
	executed = k.StableSwapBeginBlocker(ctx)

	// ASSERT: The BeginBlocker should still succeed but logging the error.
	assert.True(t, executed)

	// ARRANGE: Add an invalid StableSwap Pool triggering errors on the StableSwap BeginBlocker execution, and a paused Pool.
	err := k.Stableswap.Pools.Set(ctx, 0, stableswap.Pool{
		ProtocolFeePercentage: 0,
	})
	assert.NoError(t, err)
	err = k.Stableswap.Pools.Set(ctx, 1, stableswap.Pool{
		ProtocolFeePercentage: 0,
	})
	assert.NoError(t, err)
	err = k.Pools.Set(ctx, 0, types.Pool{
		Id:        0,
		Address:   user.Address,
		Algorithm: types.STABLESWAP,
		Pair:      "uusdc",
	})
	assert.NoError(t, err)
	err = k.Paused.Set(ctx, 0, true)
	assert.NoError(t, err)
	err = k.Stableswap.UnbondingPositions.Set(ctx, collections.Join3(time.Time{}.Unix(), user.Address, uint64(0)), stableswap.UnbondingPosition{
		Shares:  math.LegacyDec{},
		Amount:  nil,
		EndTime: time.Time{},
	})
	assert.NoError(t, err)

	executed = k.StableSwapBeginBlocker(ctx)
	// ASSERT: The BeginBlocker should still succeed but logging the error.
	assert.True(t, executed)

	// ARRANGE: Add an invalid unbonding position
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)})
	err = k.Stableswap.UnbondingPositions.Set(ctx, collections.Join3(time.Time{}.Unix(), user.Address, uint64(1)), stableswap.UnbondingPosition{
		Shares:  math.LegacyDec{},
		Amount:  nil,
		EndTime: time.Time{},
	})
	assert.NoError(t, err)
	err = k.Pools.Remove(ctx, 1)
	assert.NoError(t, err)
	err = k.Paused.Set(ctx, 0, false)
	assert.NoError(t, err)

	// ARRANGE: Setup up failing collections on UsersTotalBondedShares
	tmpUsersTotalBondedShares := k.Stableswap.UsersTotalBondedShares
	k.Stableswap.UsersTotalBondedShares = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.StableSwapUsersTotalBondedSharesPrefix, "stableswap_users_total_bonded_shares", collections.PairKeyCodec(collections.Uint64Key, collections.StringKey), sdk.LegacyDecValue,
	)
	executed = k.StableSwapBeginBlocker(ctx)
	// ASSERT: The BeginBlocker should still succeed but logging the error.
	assert.True(t, executed)
	k.Stableswap.UsersTotalBondedShares = tmpUsersTotalBondedShares

	// ARRANGE: Remove the StableSwap Pool in order to test a "failed to access Pool".
	_ = k.Stableswap.Pools.Remove(ctx, 0)

	// ACT: Execute the BeginBlocker.
	executed = k.StableSwapBeginBlocker(ctx)
	assert.True(t, executed)
}
