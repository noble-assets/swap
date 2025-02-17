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

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/stretchr/testify/require"
	modulev1 "swap.noble.xyz/api/module/v1"
	"swap.noble.xyz/keeper"
	"swap.noble.xyz/types"
	"swap.noble.xyz/utils/mocks"
)

func TestNewKeeper(t *testing.T) {
	// ARRANGE: Set the PoolsPrefix to an already existing key.
	types.PoolsPrefix = types.NextPoolIDPrefix

	// ACT: Attempt to create a new Keeper with overlapping prefixes.
	require.Panics(t, func() {
		cfg := mocks.MakeTestEncodingConfig("noble")
		keeper.NewKeeper(
			cfg.Codec,
			mocks.FailingStore(mocks.Set, nil),
			runtime.ProvideEventService(),
			runtime.ProvideHeaderInfoService(&runtime.AppBuilder{}),
			log.NewNopLogger(),
			"authority",
			"uusdn",
			1e6,
			0.5e4,
			&modulev1.StableSwap{},
			address.NewBech32Codec("noble"),
			mocks.AccountKeeper{},
			mocks.BankKeeper{},
		)
	})
	// ASSERT: The function should've panicked.

	// ARRANGE: Restore the original PoolsPrefix.
	types.PoolsPrefix = []byte("pools_generic")

	// ACT: Test the logger.
	k, _ := mocks.SwapKeeper(t)
	k.Logger()
}
