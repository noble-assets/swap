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

import sdkerrors "cosmossdk.io/errors"

// ValidateMsgSwap checks and ensures that the swap message is valid and with a correct routing plan.
func ValidateMsgSwap(msg *MsgSwap) error {
	// Ensure that the message contains at least 1 route.
	if len(msg.Routes) < 1 {
		return sdkerrors.Wrapf(ErrInvalidSwapRoutingPlan, "expected at least 1 route, got: %d", len(msg.Routes))
	}

	// Ensure that the Message contains at least 1 route.
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrapf(ErrInvalidSwapRoutingPlan, "invalid swap coin, got: %s", msg.Amount.String())
	}

	// Ensure that the Message contains a valid positive amount to swap.
	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrapf(ErrInvalidSwapRoutingPlan, "amount must be positive, got: %s", msg.Amount.String())
	}

	// Ensure that the Message contains a valid min denom.
	if msg.Min.Denom != msg.Routes[len(msg.Routes)-1].DenomTo {
		return sdkerrors.Wrapf(
			ErrInvalidSwapRoutingPlan,
			"inconsistent min denom: expected %s but got %s", msg.Routes[len(msg.Routes)-1].DenomTo, msg.Min.Denom,
		)
	}

	// Ensure no duplicated routes
	seen := make(map[uint64]bool, len(msg.Routes))
	for _, entry := range msg.Routes {
		if seen[entry.PoolId] {
			return sdkerrors.Wrapf(ErrInvalidSwapRoutingPlan, "found duplicated route on Pool: %d", entry.PoolId)
		}
		seen[entry.PoolId] = true
	}

	return nil
}
