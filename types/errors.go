package types

import "cosmossdk.io/errors"

var (
	ErrInvalidAuthority        = errors.Register(ModuleName, 1, "signer is not authority")
	ErrInvalidSwapRoutingPlan  = errors.Register(ModuleName, 2, "invalid swap routing plan")
	ErrInvalidDenom            = errors.Register(ModuleName, 3, "invalid denom")
	ErrInsufficientBalance     = errors.Register(ModuleName, 5, "insufficient balance")
	ErrInvalidPool             = errors.Register(ModuleName, 6, "invalid pool")
	ErrInvalidPair             = errors.Register(ModuleName, 7, "invalid pair")
	ErrInvalidPoolParams       = errors.Register(ModuleName, 8, "invalid pool params")
	ErrInvalidAlgorithm        = errors.Register(ModuleName, 9, "invalid algorithm")
	ErrPoolActivityPaused      = errors.Register(ModuleName, 12, "pool activity is paused")
	ErrInvalidUnbondAmount     = errors.Register(ModuleName, 13, "invalid unbond amount")
	ErrInvalidUnbondPercentage = errors.Register(ModuleName, 14, "invalid unbond percentage")
	ErrInvalidUnbondPosition   = errors.Register(ModuleName, 15, "invalid unbond position")
	ErrInvalidSwap             = errors.Register(ModuleName, 16, "invalid swap")
)
