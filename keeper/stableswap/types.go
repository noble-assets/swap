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
	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"

	stableswaptypes "swap.noble.xyz/types/stableswap"
)

// BondedPosition Indexes

type BondedPositionIndexes struct {
	ByProvider        *indexes.Multi[string, collections.Triple[uint64, string, int64], stableswaptypes.BondedPosition]
	ByPoolAndProvider *indexes.Multi[collections.Pair[uint64, string], collections.Triple[uint64, string, int64], stableswaptypes.BondedPosition]
}

func (i BondedPositionIndexes) IndexesList() []collections.Index[collections.Triple[uint64, string, int64], stableswaptypes.BondedPosition] {
	return []collections.Index[collections.Triple[uint64, string, int64], stableswaptypes.BondedPosition]{
		i.ByProvider,
		i.ByPoolAndProvider,
	}
}

func NewBondedPositionIndexes(builder *collections.SchemaBuilder) BondedPositionIndexes {
	return BondedPositionIndexes{
		ByProvider: indexes.NewMulti(
			builder, []byte("position_by_provider"), "position_by_provider",
			collections.StringKey,
			collections.TripleKeyCodec(collections.Uint64Key, collections.StringKey, collections.Int64Key),
			func(key collections.Triple[uint64, string, int64], shares stableswaptypes.BondedPosition) (string, error) {
				return key.K2(), nil
			},
		),
		ByPoolAndProvider: indexes.NewMulti(
			builder,
			[]byte("position_by_pool_and_provider"),
			"position_by_pool_and_provider",
			collections.PairKeyCodec(collections.Uint64Key, collections.StringKey),
			collections.TripleKeyCodec(collections.Uint64Key, collections.StringKey, collections.Int64Key),
			func(key collections.Triple[uint64, string, int64], position stableswaptypes.BondedPosition) (collections.Pair[uint64, string], error) {
				return collections.Join(key.K1(), key.K2()), nil
			},
		),
	}
}

// UnbondingPosition Indexes

type UnbondingPositionIndexes struct {
	ByProvider *indexes.Multi[string, collections.Triple[int64, string, uint64], stableswaptypes.UnbondingPosition]
	ByPool     *indexes.Multi[uint64, collections.Triple[int64, string, uint64], stableswaptypes.UnbondingPosition]
}

func (i UnbondingPositionIndexes) IndexesList() []collections.Index[collections.Triple[int64, string, uint64], stableswaptypes.UnbondingPosition] {
	return []collections.Index[collections.Triple[int64, string, uint64], stableswaptypes.UnbondingPosition]{
		i.ByProvider,
		i.ByPool,
	}
}

func NewUnbondingPositionIndexes(builder *collections.SchemaBuilder) UnbondingPositionIndexes {
	return UnbondingPositionIndexes{
		ByProvider: indexes.NewMulti(
			builder, []byte("unbonding_by_provider"), "unbonding_by_provider",
			collections.StringKey,
			collections.TripleKeyCodec(collections.Int64Key, collections.StringKey, collections.Uint64Key),
			func(key collections.Triple[int64, string, uint64], position stableswaptypes.UnbondingPosition) (string, error) {
				return key.K2(), nil
			},
		),
		ByPool: indexes.NewMulti(
			builder, []byte("unbonding_by_pool"), "unbonding_by_pool",
			collections.Uint64Key,
			collections.TripleKeyCodec(collections.Int64Key, collections.StringKey, collections.Uint64Key),
			func(key collections.Triple[int64, string, uint64], position stableswaptypes.UnbondingPosition) (uint64, error) {
				return key.K3(), nil
			},
		),
	}
}
