package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"swap.noble.xyz/utils/mocks"
)

func TestGetNextPoolId(t *testing.T) {
	keeper, ctx := mocks.SwapKeeper(t)

	// ACT: Retrieve the NextPoolId with no state.
	nextPoolId := keeper.GetNextPoolID(ctx)

	// ASSERT: Expect 0.
	require.Equal(t, uint64(0), nextPoolId)

	// ACT: Increase the NextPoolId.
	_, err := keeper.IncreaseNextPoolID(ctx)

	// ASSERT: Expect 1.
	require.NoError(t, err)
	require.Equal(t, uint64(1), keeper.GetNextPoolID(ctx))
}
