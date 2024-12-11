package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"swap.noble.xyz/types/stableswap"
	"swap.noble.xyz/utils/mocks"
)

func TestBeginBlocker(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)

	// ARRANGE: Execute the BeginBlocker.
	err := k.BeginBlocker(ctx)
	// ASSERT: Correct execution.
	assert.Nil(t, err)

	// ARRANGE: Add an invalid StableSwap Pool triggering errors on the StableSwap BeginBlocker execution.
	err = k.Stableswap.Pools.Set(ctx, 0, stableswap.Pool{
		ProtocolFeePercentage: 0,
	})
	assert.NoError(t, err)

	// ASSERT: The BeginBlocker should still succeed but logging the error.
	err = k.BeginBlocker(ctx)
	assert.NoError(t, err)
}
