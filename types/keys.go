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

package types

const ModuleName = "swap"

// Prefixes must be unique and cannot overlap, meaning one cannot be the start of another (ex. <pools> and <pools_stableswap> are invalid).
// Distinct prefixes with shared roots (ex. <pools_generic> and <pools_stableswap>) are allowed as long as neither fully contains the other.

var (
	NextPoolIDPrefix = []byte("next_pool_id")
	PausedPrefix     = []byte("paused")
	PoolsPrefix      = []byte("pools_generic")

	StableSwapPoolsPrefix                     = []byte("pools_stableswap")
	StableSwapUsersTotalBondedSharesPrefix    = []byte("stableswap_users_total_bonded_shares")
	StableSwapUsersTotalUnbondingSharesPrefix = []byte("stableswap_users_total_unbonding_shares")
	StableSwapPoolsTotalUnbondingSharesPrefix = []byte("stableswap_pools_total_unbonding_shares")
	StableSwapUnbondingPositionsPrefix        = []byte("stableswap_unbonding_positions")
	StableSwapBondedPositionsPrefix           = []byte("stableswap_bonded_positions")
)
