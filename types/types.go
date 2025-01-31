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

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"swap.noble.xyz/types/stableswap"
)

// SwapResult represents the outputs of a swap.
type SwapResult struct {
	Dy          math.LegacyDec
	ProtocolFee math.LegacyDec
	RewardsFee  math.LegacyDec
}

// AddLiquidityCommitment commits to adding liquidity (via a bonded position) to a stableswap pool.
type AddLiquidityCommitment struct {
	BondedPosition stableswap.BondedPosition
}

// RemoveLiquidityCommitment commits to removing liquidity (via an unbonding position) from a stableswap pool.
type RemoveLiquidityCommitment struct {
	UnbondingPosition stableswap.UnbondingPosition
}

// ReceiverMulti specifies a recipient and multiple coin amounts for distribution.
type ReceiverMulti struct {
	Amount  sdk.Coins
	Address sdk.AccAddress
	PoolId  uint64
}

// Receiver specifies a recipient and a single coin amount.
type Receiver struct {
	Amount  sdk.Coin
	Address sdk.AccAddress
}

// SwapCommitment details a swap, including input/output coins and fee distributions.
type SwapCommitment struct {
	In   sdk.Coin
	Out  sdk.Coin
	Fees []Receiver
}

// PlanSwapRoute defines one step of a multi-hop swap, referencing a pool and its swap details.
type PlanSwapRoute struct {
	PoolId      uint64
	PoolAddress string
	Commitment  *SwapCommitment
}

// PlanSwapRoutes aggregates multiple PlanSwapRoute steps into a full multi-hop swap plan.
type PlanSwapRoutes struct {
	Swaps []PlanSwapRoute
}
