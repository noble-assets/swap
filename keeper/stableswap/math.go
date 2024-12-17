package stableswap

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"swap.noble.xyz/types"
)

const (
	FeeDenominator                    = 1e10
	DecimalPrecision                  = 1e18
	AmplificationCoefficientPrecision = 1_000_000
)

// getAmplificationCoefficient computes the amplification coefficient (A) based on the current time and ramp settings.
func getAmplificationCoefficient(currentTime int64, initialA math.LegacyDec, futureA math.LegacyDec, initialATime int64, futureATime int64) math.LegacyDec {
	if futureATime <= currentTime {
		return futureA.MulInt64(AmplificationCoefficientPrecision)
	}

	// Linear interpolation during the ramp period
	if futureA.GT(initialA) {
		// A is increasing
		return initialA.Add(
			futureA.Sub(initialA).
				MulInt64(currentTime - initialATime).
				QuoInt64(futureATime - initialATime),
		).MulInt64(AmplificationCoefficientPrecision)
	}

	// A is decreasing
	return initialA.Sub(
		initialA.Sub(futureA).
			MulInt64(currentTime - initialATime).
			QuoInt64(futureATime - initialATime),
	).MulInt64(AmplificationCoefficientPrecision)
}

// calculateInvariant computes the invariant (D) using pool balances and the amplification coefficient (A).
func calculateInvariant(xp sdk.DecCoins, amp math.LegacyDec) (math.LegacyDec, error) {
	// Calculate S, the sum of all balances in xp
	S := math.LegacyZeroDec()
	for _, x := range xp {
		S = S.Add(x.Amount)
	}

	if S.IsZero() {
		return math.LegacyZeroDec(), nil // Return zero if S is zero
	}
	// D is the invariant we are trying to find
	D := S                 // Start with D = S
	Ann := amp.MulInt64(2) // Ann = A * N

	// Newton-Raphson iteration to find D
	for i := 0; i < 255; i++ {
		// D_P = D
		D_P := D

		// Calculate D_P = D_P * D / (x * NCoins) for each x in xp
		for _, x := range xp {
			D_P = D_P.Mul(D).Quo(x.Amount.Mul(math.LegacyNewDec(2)))
		}

		// Save the current D to Dprev
		Dprev := D

		// Calculate numerator: (Ann * S + D_P * NCoins)
		numerator := Ann.Mul(S).Add(D_P.MulInt64(2))

		// Calculate denominator: (Ann - 1) * D + (NCoins + 1) * D_P
		denominator := Ann.Sub(math.LegacyOneDec()).Mul(D).Add(D_P.MulInt64(3))

		// Update D: D = D * numerator / denominator
		D = D.Mul(numerator).Quo(denominator)

		// Check for convergence: |D - Dprev| <= 1 (LegacyDec doesn't directly support absolute values, so check with IsZero or comparing difference)
		if D.Sub(Dprev).Abs().LTE(math.LegacyOneDec()) {
			return D, nil
		}
	}

	// Error the invariant don't converge within 255 iterations.
	return math.LegacyZeroDec(), fmt.Errorf("invariant calculation did not converge")
}

// calculateAdjustedBalancesInRates adjusts balances by their rate multipliers for precision.
func calculateAdjustedBalancesInRates(rates sdk.Coins, balances sdk.Coins) (sdk.DecCoins, error) {
	balance := sdk.DecCoins{}

	for _, rate := range rates {
		// balance: (rate * balance) / prevision
		adjustedAmount := rate.Amount.Mul(balances.AmountOf(rate.Denom)).QuoRaw(DecimalPrecision)
		balance = balance.Add(sdk.NewDecCoin(rate.Denom, adjustedAmount))
	}

	return balance, nil
}

// GetY calculates the y value for the exchange (the amount of the output coin)
func getY(x sdk.Coin, amp, D math.LegacyDec) (math.LegacyDec, error) {
	Ann := amp.Mul(math.LegacyNewDec(2))
	// Last portion of c calculation: c = c * D / (Ann * N_COINS)
	c := D.Mul(D).Quo(Ann.Mul(math.LegacyNewDec(2)))

	// b = S_ + D / Ann
	b := x.Amount.ToLegacyDec().Add(D.Quo(Ann))

	// Initialize y
	y := D // Start with y = D

	// Newton-Raphson iteration to find y
	for _i := 0; _i < 255; _i++ {
		yPrev := y

		// y = (y^2 + c) / (2y + b - D)
		ySq := y.Mul(y)
		numerator := ySq.Add(c)
		denominator := y.MulInt64(2).Add(b.Sub(D))
		y = numerator.Quo(denominator)

		// Check for convergence: |y - yPrev| <= 1
		if y.Sub(yPrev).Abs().LTE(math.LegacyOneDec()) {
			return y, nil
		}
	}
	return math.LegacyZeroDec(), fmt.Errorf("did not converge")
}

// performSwap executes the internal token swap and computes resulting balances and fees.
func performSwap(x sdk.Coin, xp sdk.DecCoins, amp math.LegacyDec, denomTo string,
	rewardsFee int64, protocolFeePercentage int64, maxFee int64, rateMultipliers sdk.Coins,
) (types.SwapResult, error) {
	// Calculate invariant D
	D, err := calculateInvariant(xp, amp)
	if err != nil {
		return types.SwapResult{}, err
	}

	// Calculate the new y value after the exchange
	y, err := getY(x, amp, D)
	if err != nil {
		return types.SwapResult{}, err
	}

	// Calculate dy (amount to be received)
	dy := xp.AmountOf(denomTo).Sub(y)

	// Calculate the fee
	rewards := dy.MulInt64(rewardsFee).Quo(math.LegacyNewDec(FeeDenominator))

	if rewards.GT(math.LegacyNewDec(maxFee)) {
		rewards = math.LegacyNewDec(maxFee)
	}
	protocol := rewards.MulInt64(protocolFeePercentage).Quo(math.LegacyNewDec(100))
	rewards = rewards.Sub(protocol)
	// Subtract the fee from dy
	dy = dy.Sub(rewards)
	dy = dy.Sub(protocol)

	// Convert dy back to the original units
	dy = dy.Mul(math.LegacyNewDec(DecimalPrecision)).Quo(rateMultipliers.AmountOf(denomTo).ToLegacyDec())

	return types.SwapResult{
		Dy:          dy,
		ProtocolFee: protocol,
		RewardsFee:  rewards,
	}, nil
}

// computeNewAdjustedBalance calculates the new balance after adding the delta, adjusted by the rate multiplier and precision.
func computeNewAdjustedBalance(xp math.LegacyDec, dx math.LegacyDec, rateMultiplier math.LegacyDec, PRECISION int64) math.LegacyDec {
	return xp.Add(dx.Mul(rateMultiplier).QuoInt64(PRECISION))
}
