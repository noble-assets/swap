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
