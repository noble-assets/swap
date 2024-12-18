package types

import "cosmossdk.io/errors"

var (
	ErrInvalidAuthority        = errors.Register(ModuleName, 1, "signer is not authority")
	ErrInvalidSwapRoutingPlan  = errors.Register(ModuleName, 2, "invalid swap routing plan")
	ErrInvalidAmount           = errors.Register(ModuleName, 3, "invalid amount")
	ErrInsufficientBalance     = errors.Register(ModuleName, 4, "insufficient balance")
	ErrInvalidPool             = errors.Register(ModuleName, 5, "invalid pool")
	ErrInvalidPoolParams       = errors.Register(ModuleName, 6, "invalid pool params")
	ErrInvalidAlgorithm        = errors.Register(ModuleName, 7, "invalid algorithm")
	ErrPoolActivityPaused      = errors.Register(ModuleName, 8, "pool activity is paused")
	ErrInvalidUnbondAmount     = errors.Register(ModuleName, 9, "invalid unbond amount")
	ErrInvalidUnbondPercentage = errors.Register(ModuleName, 10, "invalid unbond percentage")
	ErrInvalidUnbondPosition   = errors.Register(ModuleName, 11, "invalid unbond position")
)
