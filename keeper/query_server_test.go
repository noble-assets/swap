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

package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"swap.noble.xyz/keeper"
	"swap.noble.xyz/types"
	"swap.noble.xyz/types/stableswap"
	"swap.noble.xyz/utils"
	"swap.noble.xyz/utils/mocks"
)

func TestSimulateSwap(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	msgServer := keeper.NewMsgServer(k)
	queryServer := keeper.NewQueryServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	alice, bob := utils.TestAccount(), utils.TestAccount()
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], sdk.NewCoin("uusdc", math.NewInt(1_000_000*ONE)))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], sdk.NewCoin("uusdn", math.NewInt(1_000_000*ONE)))

	// ARRANGE: Create a Pool.
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		RewardsFee:            0,
		ProtocolFeePercentage: 0,
		InitialA:              100,
		FutureA:               100,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})
	assert.Nil(t, err)

	// ARRANGE: Provide liquidity.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: alice.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(sdk.NewCoin("uusdc", math.NewInt(1_000_000*ONE)), sdk.NewCoin("uusdn", math.NewInt(1_000_000*ONE))),
	})
	assert.NoError(t, err)

	// ARRANGE: Create the route for the tests.
	routes := []types.Route{
		{
			PoolId:  0,
			DenomTo: "uusdn",
		},
	}

	// ACT: Simulate a swap without a valid balance.
	_, err = msgServer.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		Routes: routes,
	})
	// ASSERT: The simulations should've failed.
	assert.Errorf(t, err, "uusdc balance of noble1y96tyawaa9adsu22gwp5vs0llc9glpuh9p8n5z is smaller then 100000000.000000000000000000, available: 0")

	// ARRANGE: provide a valid balance to the user.
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdc", math.NewInt(1_000*ONE)))

	// ACT: Attempt to query the simulation without the msg.
	_, err = queryServer.SimulateSwap(ctx, nil)
	assert.Error(t, err)

	// ARRANGE: Simulate again the Swap.
	swapRequest := &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(90*ONE)),
	}
	responseSimulation, err := queryServer.SimulateSwap(ctx, &types.QuerySimulateSwap{
		Amount: swapRequest.Amount,
		Routes: swapRequest.Routes,
		Min:    swapRequest.Min,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.MsgSwapResponse{
		Result: sdk.NewCoin("uusdn", math.NewInt(99999999)),
		Swaps: []*types.Swap{
			{
				In:   sdk.NewCoin("uusdc", math.NewInt(100000000)),
				Out:  sdk.NewCoin("uusdn", math.NewInt(99999999)),
				Fees: sdk.Coins{},
			},
		},
	}, responseSimulation)
}

func TestPausing(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	server := keeper.NewMsgServer(k)
	queryServer := keeper.NewQueryServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	user := utils.TestAccount()

	// ARRANGE: Create a Pool.
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdc", math.NewInt(100)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusde", math.NewInt(100)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdn", math.NewInt(100)))
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
		InitialA: 100,
		FutureA:  100,
	})
	assert.NoError(t, err)

	// ARRANGE: Pause pools by algorithm type.
	_, err = server.PauseByAlgorithm(ctx, &types.MsgPauseByAlgorithm{
		Signer:    "authority",
		Algorithm: types.STABLESWAP,
	})
	assert.NoError(t, err)

	// ACT: Query with invalid request.
	_, err = queryServer.Paused(ctx, nil)
	assert.Error(t, err)

	// ASSERT: Correct paused pools.
	res, err := queryServer.Paused(ctx, &types.QueryPaused{})
	assert.Nil(t, err)
	assert.Equal(t, res.PausedPools, []uint64{0})

	// ARRANGE: Add a new pool with the same algorithm.
	_, err = stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusde",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusde", math.NewInt(1000000000000000000)),
		),
		InitialA: 100,
		FutureA:  100,
	})
	assert.NoError(t, err)

	// ARRANGE: Add a new pool with the same algorithm.
	_, err = stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "ueure",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("ueure", math.NewInt(1000000000000000000)),
		),
		InitialA: 100,
		FutureA:  100,
	})
	assert.NoError(t, err)

	// ARRANGE: Unpause by algorithm.
	_, err = server.UnpauseByAlgorithm(ctx, &types.MsgUnpauseByAlgorithm{
		Signer:    "authority",
		Algorithm: types.STABLESWAP,
	})
	assert.NoError(t, err)

	// ASSERT: Correct paused pools.
	res, err = queryServer.Paused(ctx, &types.QueryPaused{})
	assert.Nil(t, err)
	assert.Equal(t, res.PausedPools, []uint64(nil))

	// ARRANGE: Pause pools by algorithm type.
	_, err = server.PauseByAlgorithm(ctx, &types.MsgPauseByAlgorithm{
		Signer:    "authority",
		Algorithm: types.STABLESWAP,
	})
	assert.NoError(t, err)

	// ASSERT: All pool are paused.
	res, err = queryServer.Paused(ctx, &types.QueryPaused{})
	assert.Nil(t, err)
	assert.Equal(t, res.PausedPools, []uint64{0, 1, 2})

	// ARRANGE: Unpause all the pools.
	_, err = server.UnpauseByAlgorithm(ctx, &types.MsgUnpauseByAlgorithm{
		Signer:    "authority",
		Algorithm: types.STABLESWAP,
	})
	assert.NoError(t, err)

	// ARRANGE: Pause a single pool by its id.
	_, err = server.PauseByPoolIds(ctx, &types.MsgPauseByPoolIds{
		Signer:  "authority",
		PoolIds: []uint64{1},
	})
	assert.NoError(t, err)

	// ASSERT: Pool 1 is paused.
	res, err = queryServer.Paused(ctx, &types.QueryPaused{})
	assert.Nil(t, err)
	assert.Equal(t, res.PausedPools, []uint64{1})

	// ARRANGE: Unpause pool by its id.
	_, err = server.UnpauseByPoolIds(ctx, &types.MsgUnpauseByPoolIds{
		Signer:  "authority",
		PoolIds: []uint64{1},
	})
	assert.NoError(t, err)

	// ASSERT: All pools are active.
	res, err = queryServer.Paused(ctx, &types.QueryPaused{})
	assert.Nil(t, err)
	assert.Equal(t, res.PausedPools, []uint64(nil))

	// ARRANGE: Pause multiple pools by its id.
	_, err = server.PauseByPoolIds(ctx, &types.MsgPauseByPoolIds{
		Signer:  "authority",
		PoolIds: []uint64{1, 0},
	})
	assert.NoError(t, err)

	// ASSERT: The paused pools are not active.
	res, err = queryServer.Paused(ctx, &types.QueryPaused{})
	assert.Nil(t, err)
	assert.Equal(t, res.PausedPools, []uint64{0, 1})

	// ARRANGE: Unpause multiple pools by their ids.
	_, err = server.UnpauseByPoolIds(ctx, &types.MsgUnpauseByPoolIds{
		Signer:  "authority",
		PoolIds: []uint64{1, 0},
	})
	assert.NoError(t, err)

	// ASSERT: The paused pools are active.
	res, err = queryServer.Paused(ctx, &types.QueryPaused{})
	assert.Nil(t, err)
	assert.Equal(t, res.PausedPools, []uint64(nil))
}

func TestPool(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	queryServer := keeper.NewQueryServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	user := utils.TestAccount()

	// ARRANGE: Create a Pool.
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdc", math.NewInt(100)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusde", math.NewInt(100)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdn", math.NewInt(100)))
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
		InitialA: 100,
		FutureA:  100,
	})
	assert.NoError(t, err)

	// ACT: Attempt to query a pool with an invalid message.
	_, err = queryServer.Pool(ctx, nil)
	assert.Error(t, err)

	// ACT: Attempt to query a non-existing pool.
	_, err = queryServer.Pool(ctx, &types.QueryPool{
		PoolId: 10,
	})
	assert.Error(t, err)

	// ACT: Query a valid pool.
	res, err := queryServer.Pool(ctx, &types.QueryPool{
		PoolId: 0,
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), res.Pool.Id)

	// ACT: Query a valid pool.
	res, err = queryServer.Pool(ctx, &types.QueryPool{
		PoolId: 0,
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), res.Pool.Id)
}

func TestPools(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	queryServer := keeper.NewQueryServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	user := utils.TestAccount()

	// ARRANGE: Create 2 pools.
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdc", math.NewInt(100)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusde", math.NewInt(100)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdn", math.NewInt(100)))
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
		InitialA: 100,
		FutureA:  100,
	})
	assert.NoError(t, err)
	_, err = stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusde",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusde", math.NewInt(1000000000000000000)),
		),
		InitialA: 100,
		FutureA:  100,
	})
	assert.NoError(t, err)

	// ACT: Attempt to query the pools with an invalid payload.
	_, err = queryServer.Pools(ctx, nil)
	assert.Error(t, err)

	// ARRANGE: Add a Pool with a different algorithm.
	err = k.SetPool(ctx, 2, types.Pool{
		Id:        2,
		Address:   "",
		Algorithm: types.UNSPECIFIED,
		Pair:      "",
	})
	assert.NoError(t, err)

	// ACT: Query Pools.
	res, err := queryServer.Pools(ctx, &types.QueryPools{})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Pools))
}

func TestRate(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	queryServer := keeper.NewQueryServer(k)
	server := keeper.NewMsgServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	bob, alice := utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Create a Pool.
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdc", math.NewInt(1000*ONE)))
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdn", math.NewInt(1000*ONE)))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], sdk.NewCoin("uusde", math.NewInt(1000*ONE)))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], sdk.NewCoin("uusdn", math.NewInt(1000*ONE)))
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
		InitialA:    1000,
		FutureA:     1000,
		FutureATime: 0,
	})
	assert.NoError(t, err)

	// ARRANGE: Add liquidity.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: bob.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
			sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		),
	})
	assert.NoError(t, err)

	// ACT: Query the rates with an invalid request.
	_, err = queryServer.Rate(ctx, nil)
	assert.Error(t, err)

	// ACT: Query the rates with a valid request.
	res, err := queryServer.Rate(ctx, &types.QueryRate{
		Denom:     "uusdn",
		Algorithm: types.STABLESWAP,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Rates))
	assert.Equal(t, &types.QueryRateResponse{
		Rates: []types.Rate{
			{
				Denom:     "uusdn",
				Vs:        "uusdc",
				Price:     math.LegacyMustNewDecFromStr("1.000001000001000001"),
				Algorithm: types.STABLESWAP,
			},
		},
	}, res)

	// ARRANGE: Add a second Pool.
	_, err = stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusde",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusde", math.NewInt(1000000000000000000)),
		),
		InitialA:    1000,
		FutureA:     1000,
		FutureATime: 0,
	})
	assert.NoError(t, err)

	// ARRANGE: Add liquidity.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: alice.Address,
		PoolId: 1,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(800*ONE)),
			sdk.NewCoin("uusde", math.NewInt(800*ONE)),
		),
	})
	assert.NoError(t, err)

	// ARRANGE: Perform a swap removing almost all the liquidity in the Pool 0.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: alice.Address,
		Amount: sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
		Routes: []types.Route{
			{
				PoolId:  0,
				DenomTo: "uusdc",
			},
		},
		Min: sdk.NewCoin("uusdc", math.NewInt(1*ONE)),
	})
	assert.NoError(t, err)

	// ACT: Query the new rates.
	res, err = queryServer.Rate(ctx, &types.QueryRate{
		Denom:     "uusdc",
		Algorithm: types.STABLESWAP,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Rates))
	assert.Equal(t, &types.QueryRateResponse{
		Rates: []types.Rate{
			{
				Denom:     "uusdc",
				Vs:        "uusdn",
				Price:     math.LegacyMustNewDecFromStr("0.002232000000000000"),
				Algorithm: types.STABLESWAP,
			},
		},
	}, res)

	// ARRANGE: Add a Pool with a different algorithm.
	err = k.SetPool(ctx, 2, types.Pool{
		Id:        2,
		Address:   "",
		Algorithm: types.UNSPECIFIED,
		Pair:      "",
	})
	assert.NoError(t, err)

	// ACT: Query the rate.
	res, err = queryServer.Rate(ctx, &types.QueryRate{
		Denom:     "uusdc",
		Algorithm: types.STABLESWAP,
	})
	// ASSERT: Expect the same response.
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.Rates))
	assert.Equal(t, &types.QueryRateResponse{
		Rates: []types.Rate{
			{
				Denom:     "uusdc",
				Vs:        "uusdn",
				Price:     math.LegacyMustNewDecFromStr("0.002232000000000000"),
				Algorithm: types.STABLESWAP,
			},
		},
	}, res)

	// ACT: Query the rate.
	res, err = queryServer.Rate(ctx, &types.QueryRate{
		Denom:     "uusdc",
		Algorithm: types.UNSPECIFIED,
	})
	assert.NoError(t, err)

	// ASSERT: Expect the same response.
	assert.Equal(t, 1, len(res.Rates))
	assert.Equal(t, &types.QueryRateResponse{
		Rates: []types.Rate{
			{
				Denom:     "uusdc",
				Vs:        "uusdn",
				Price:     math.LegacyMustNewDecFromStr("0.002232000000000000"),
				Algorithm: types.STABLESWAP,
			},
		},
	}, res)
}

func TestRates(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	queryServer := keeper.NewQueryServer(k)
	server := keeper.NewMsgServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	bob, alice := utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Create a Pool.
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusde", math.NewInt(1000*ONE)))
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdn", math.NewInt(1000*ONE)))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], sdk.NewCoin("uusdc", math.NewInt(1000*ONE)))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], sdk.NewCoin("uusdn", math.NewInt(1000*ONE)))
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusde",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusde", math.NewInt(1000000000000000000)),
		),
		InitialA: 100,
		FutureA:  100,
	})
	assert.NoError(t, err)

	// ARRANGE: Add liquidity.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: bob.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
			sdk.NewCoin("uusde", math.NewInt(100*ONE)),
		),
	})
	assert.NoError(t, err)

	// ACT: Query the rates with an invalid request.
	_, err = queryServer.Rates(ctx, nil)
	assert.Error(t, err)

	// ACT: Query the rates with a valid request.
	res, err := queryServer.Rates(ctx, &types.QueryRates{})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Rates))
	assert.Equal(t, &types.QueryRatesResponse{
		Rates: []types.Rate{
			{
				Denom:     "uusde",
				Vs:        "uusdn",
				Price:     math.LegacyMustNewDecFromStr("0.999999000000000000"),
				Algorithm: types.STABLESWAP,
			}, {
				Denom:     "uusdn",
				Vs:        "uusde",
				Price:     math.LegacyMustNewDecFromStr("1.000001000001000001"),
				Algorithm: types.STABLESWAP,
			},
		},
	}, res)

	// ARRANGE: Add a second Pool with liquidity.
	_, err = stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
		InitialA: 100,
		FutureA:  100,
	})
	assert.NoError(t, err)

	// ARRANGE: Add liquidity.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: alice.Address,
		PoolId: 1,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(800*ONE)),
			sdk.NewCoin("uusdc", math.NewInt(800*ONE)),
		),
	})
	assert.NoError(t, err)

	// ARRANGE: Perform a swap.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: alice.Address,
		Amount: sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
		Routes: []types.Route{
			{
				PoolId:  0,
				DenomTo: "uusde",
			},
		},
		Min: sdk.NewCoin("uusde", math.NewInt(1*ONE)),
	})
	assert.NoError(t, err)

	// ACT: Query the new rates.
	res, err = queryServer.Rates(ctx, &types.QueryRates{})
	assert.NoError(t, err)
	assert.Equal(t, 4, len(res.Rates))
	assert.Equal(t, &types.QueryRatesResponse{
		Rates: []types.Rate{
			{
				Denom:     "uusdc",
				Vs:        "uusdn",
				Price:     math.LegacyMustNewDecFromStr("0.999999000000000000"),
				Algorithm: types.STABLESWAP,
			},
			{
				Denom:     "uusde",
				Vs:        "uusdn",
				Price:     math.LegacyMustNewDecFromStr("0.007021000000000000"),
				Algorithm: types.STABLESWAP,
			},
			{
				Denom:     "uusdn",
				Vs:        "uusdc",
				Price:     math.LegacyMustNewDecFromStr("1.000001000001000001"),
				Algorithm: types.STABLESWAP,
			},
			{
				Denom:     "uusdn",
				Vs:        "uusde",
				Price:     math.LegacyMustNewDecFromStr("142.429853297251103831"),
				Algorithm: types.STABLESWAP,
			},
		},
	}, res)

	// ARRANGE: Add a Pool with a different algorithm.
	err = k.SetPool(ctx, 2, types.Pool{
		Id:        2,
		Address:   "",
		Algorithm: types.UNSPECIFIED,
		Pair:      "uusdx",
	})
	assert.NoError(t, err)

	// ACT: Expect the same response.
	res, err = queryServer.Rates(ctx, &types.QueryRates{})
	assert.NoError(t, err)
	assert.Equal(t, 4, len(res.Rates))
	assert.Equal(t, &types.QueryRatesResponse{
		Rates: []types.Rate{
			{
				Denom:     "uusdc",
				Vs:        "uusdn",
				Price:     math.LegacyMustNewDecFromStr("0.999999000000000000"),
				Algorithm: types.STABLESWAP,
			},
			{
				Denom:     "uusde",
				Vs:        "uusdn",
				Price:     math.LegacyMustNewDecFromStr("0.007021000000000000"),
				Algorithm: types.STABLESWAP,
			},
			{
				Denom:     "uusdn",
				Vs:        "uusdc",
				Price:     math.LegacyMustNewDecFromStr("1.000001000001000001"),
				Algorithm: types.STABLESWAP,
			},
			{
				Denom:     "uusdn",
				Vs:        "uusde",
				Price:     math.LegacyMustNewDecFromStr("142.429853297251103831"),
				Algorithm: types.STABLESWAP,
			},
		},
	}, res)

	// ACT: Query with a different requested algorithm and expect the same response.
	res, err = queryServer.Rates(ctx, &types.QueryRates{
		Algorithm: types.STABLESWAP,
	})
	assert.NoError(t, err)
	assert.Equal(t, 4, len(res.Rates))
	assert.Equal(t, &types.QueryRatesResponse{
		Rates: []types.Rate{
			{
				Denom:     "uusdc",
				Vs:        "uusdn",
				Price:     math.LegacyMustNewDecFromStr("0.999999000000000000"),
				Algorithm: types.STABLESWAP,
			},
			{
				Denom:     "uusde",
				Vs:        "uusdn",
				Price:     math.LegacyMustNewDecFromStr("0.007021000000000000"),
				Algorithm: types.STABLESWAP,
			},
			{
				Denom:     "uusdn",
				Vs:        "uusdc",
				Price:     math.LegacyMustNewDecFromStr("1.000001000001000001"),
				Algorithm: types.STABLESWAP,
			},
			{
				Denom:     "uusdn",
				Vs:        "uusde",
				Price:     math.LegacyMustNewDecFromStr("142.429853297251103831"),
				Algorithm: types.STABLESWAP,
			},
		},
	}, res)
}
