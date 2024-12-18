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
