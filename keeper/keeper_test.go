package keeper_test

import (
	"testing"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/stretchr/testify/require"
	modulev1 "swap.noble.xyz/api/module/v1"
	"swap.noble.xyz/keeper"
	"swap.noble.xyz/types"
	"swap.noble.xyz/utils/mocks"
)

func TestNewKeeper(t *testing.T) {
	// ARRANGE: Set the PoolsPrefix to an already existing key
	types.PoolsPrefix = types.NextPoolIDPrefix

	// ACT: Attempt to create a new Keeper with overlapping prefixes
	require.Panics(t, func() {
		cfg := mocks.MakeTestEncodingConfig("noble")
		keeper.NewKeeper(
			cfg.Codec,
			mocks.FailingStore(mocks.Set, nil),
			runtime.ProvideEventService(),
			runtime.ProvideHeaderInfoService(&runtime.AppBuilder{}),
			log.NewNopLogger(),
			"authority",
			"uusdn",
			&modulev1.StableSwap{},
			address.NewBech32Codec("noble"),
			mocks.AccountKeeper{},
			mocks.BankKeeper{},
		)
	})
	// ASSERT: The function should've panicked.

	// ARRANGE: Restore the original PoolsPrefix
	types.PoolsPrefix = []byte("pools_generic")

	// ACT: Test the logger
	k, _ := mocks.SwapKeeper(t)
	k.Logger()
}
