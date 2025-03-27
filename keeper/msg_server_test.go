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
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/header"
	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"swap.noble.xyz/keeper"
	stableswapkeeper "swap.noble.xyz/keeper/stableswap"
	"swap.noble.xyz/types"
	"swap.noble.xyz/types/stableswap"
	"swap.noble.xyz/utils"
	"swap.noble.xyz/utils/mocks"
)

const ONE = int64(1e6)

func TestUnbalanceAndRebalanceIsNotConservative(t *testing.T) {
	usdnLiquidity := sdk.NewCoin("uusdn", math.NewInt(1_000_000*ONE))
	usdcLiquidity := sdk.NewCoin("uusdc", math.NewInt(1_000_000*ONE))

	routeToUsdn := []types.Route{
		{
			PoolId:  0,
			DenomTo: "uusdn",
		},
	}

	routeToUsdc := []types.Route{
		{
			PoolId:  0,
			DenomTo: "uusdc",
		},
	}

	alice, bob := utils.TestAccount(), utils.TestAccount()

	type Servers struct {
		swap       types.MsgServer
		stableSwap stableswap.MsgServer
	}

	testCases := []struct {
		name      string
		initialA  int64
		swapInAmt math.Int
		// setup runs before executing the swap to unbalance the pool which will be rebalanced
		// via the swaps function below.
		// This is used to pre-unbalance the pool allowing to tests small amount rebalance swaps
		// in a short amount of time.
		setup func(context.Context, Servers)
		// swaps accepts the context, the total amount that has to be swapped in the function, and
		// returns the total output from swaps performed: amount out + fees.
		swaps func(context.Context, Servers, math.Int) math.Int
	}{
		{
			name:      "single rebalance",
			swapInAmt: usdcLiquidity.Amount.SubRaw(1),
			initialA:  800,
			setup:     func(_ context.Context, _ Servers) {},
			swaps: func(ctx context.Context, servers Servers, totIn math.Int) math.Int {
				resp, err := servers.swap.Swap(ctx, &types.MsgSwap{
					Signer: bob.Address,
					Amount: sdk.NewCoin("uusdn", totIn),
					Routes: routeToUsdc,
					Min:    sdk.NewCoin("uusdc", math.NewInt(0)),
				})
				require.Nil(t, err)

				fees := math.ZeroInt()
				if resp.Swaps[0].Fees != nil {
					fees = resp.Swaps[0].Fees[0].Amount
				}

				return fees.Add(resp.Swaps[0].Out.Amount)
			},
		},
		{
			name:      "multiple rebalances swaps of 1uusdn",
			swapInAmt: math.NewInt(100),
			initialA:  800,
			setup: func(ctx context.Context, servers Servers) {
				// Unbalance the pool towards uusdn
				_, err := servers.swap.Swap(ctx, &types.MsgSwap{
					Signer: bob.Address,
					Amount: sdk.NewCoin("uusdc", usdcLiquidity.Amount.SubRaw(101)),
					Routes: routeToUsdn,
					Min:    sdk.NewCoin("uusdn", math.NewInt(0)),
				})
				require.Nil(t, err)
			},
			swaps: func(ctx context.Context, servers Servers, totIn math.Int) math.Int {
				totOut := math.ZeroInt()
				for range totIn.Int64() {
					resp, err := servers.swap.Swap(ctx, &types.MsgSwap{
						Signer: bob.Address,
						Amount: sdk.NewCoin("uusdn", math.NewInt(1)),
						Routes: routeToUsdc,
						Min:    sdk.NewCoin("uusdc", math.NewInt(0)),
					})
					require.Nil(t, err)

					fees := math.ZeroInt()
					if len(resp.Swaps[0].Fees) != 0 {
						fees = resp.Swaps[0].Fees[0].Amount
					}
					totOut = totOut.Add(fees.Add(resp.Swaps[0].Out.Amount))
				}

				return totOut
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			// ARRANGE
			account := mocks.AccountKeeper{
				Accounts: make(map[string]sdk.AccountI),
			}
			bank := mocks.BankKeeper{
				Balances:    make(map[string]sdk.Coins),
				Restriction: mocks.NoOpSendRestrictionFn,
			}
			k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
			server := keeper.NewMsgServer(k)
			stableswapServer := keeper.NewStableSwapMsgServer(k)

			// ARRANGE: Create the Pool.
			_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
				Signer:                "authority",
				Pair:                  "uusdc",
				RewardsFee:            10_000_000,
				ProtocolFeePercentage: 100,
				InitialA:              tC.initialA,
				FutureA:               tC.initialA,
				FutureATime:           0,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1_000_000_000_000_000_000)),
					sdk.NewCoin("uusdc", math.NewInt(1_000_000_000_000_000_000)),
				),
			})
			require.Nil(t, err)

			// ARRANGE: Bob receives more funds to cover swap fees.
			bank.Balances[alice.Address] = append(bank.Balances[alice.Address], usdcLiquidity)
			bank.Balances[alice.Address] = append(bank.Balances[alice.Address], usdnLiquidity)
			bank.Balances[bob.Address] = append(bank.Balances[bob.Address], usdcLiquidity.AddAmount(math.NewInt(1_000)))
			bank.Balances[bob.Address] = append(bank.Balances[bob.Address], usdnLiquidity.AddAmount(math.NewInt(1_000)))

			_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
				Signer: alice.Address,
				PoolId: 0,
				Amount: sdk.NewCoins(usdcLiquidity, usdnLiquidity),
			})
			assert.NoError(t, err)

			servers := Servers{
				swap:       server,
				stableSwap: stableswapServer,
			}

			tC.setup(ctx, servers)

			// ACT: Unbalance the pool towards uusdn.
			resp, err := server.Swap(ctx, &types.MsgSwap{
				Signer: bob.Address,
				Amount: sdk.NewCoin("uusdc", tC.swapInAmt),
				Routes: routeToUsdn,
				Min:    sdk.NewCoin("uusdn", math.NewInt(0)),
			})
			require.Nil(t, err)

			totOut := tC.swaps(ctx, servers, resp.Swaps[0].Out.Amount)
			require.GreaterOrEqual(t, tC.swapInAmt.Int64(), totOut.Int64(), "expected unbalance and rebalance to be not profitable")
		})
	}
}

func TestLowAmountSwapBalancedPool(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	server := keeper.NewMsgServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	alice, bob := utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Create the Pool.
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		RewardsFee:            0,
		ProtocolFeePercentage: 100,
		InitialA:              800,
		FutureA:               800,
		FutureATime:           0,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})
	assert.Nil(t, err)

	// Add to pool the balance we have on mainnet
	usdnLiquidity := sdk.NewCoin("uusdn", math.NewInt(1_000_000*ONE))
	usdcLiquidity := sdk.NewCoin("uusdc", math.NewInt(1_000_000*ONE))

	// Add funds to user balances (doesn't matter how much)
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], usdcLiquidity)
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], usdnLiquidity)

	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], usdcLiquidity)
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], usdnLiquidity)

	routes := []types.Route{
		{
			PoolId:  0,
			DenomTo: "uusdn",
		},
	}

	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: alice.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(usdcLiquidity, usdnLiquidity),
	})
	assert.NoError(t, err)

	response, err := server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(int64(1))),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(0)),
	})
	require.Nil(t, err)
	assert.Equal(t, math.NewInt(1).String(), response.Result.Amount.String())

	response, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(int64(10))),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(0)),
	})
	require.Nil(t, err)
	assert.Equal(t, math.NewInt(10).String(), response.Result.Amount.String())

	response, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(int64(100))),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(0)),
	})
	require.Nil(t, err)
	assert.Equal(t, math.NewInt(99).String(), response.Result.Amount.String())

	response, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(int64(1_000))),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(0)),
	})
	require.Nil(t, err)
	assert.Equal(t, math.NewInt(999).String(), response.Result.Amount.String())
}

func TestLowAmountSwapUnbalancedPool(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	server := keeper.NewMsgServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	alice, bob := utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Create the Pool.
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		RewardsFee:            0,
		ProtocolFeePercentage: 100,
		InitialA:              800,
		FutureA:               800,
		FutureATime:           0,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})
	assert.Nil(t, err)

	// Add to pool the balance we have on mainnet
	usdnLiquidity := sdk.NewCoin("uusdn", math.NewInt(int64(1_000_000_000_000)))
	usdcLiquidity := sdk.NewCoin("uusdc", math.NewInt(int64(1_000_000_000_000)))

	// Add funds to user balances (doesn't matter how much)
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], usdcLiquidity)
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], usdnLiquidity)

	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], usdcLiquidity)
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], usdnLiquidity)

	routes := []types.Route{
		{
			PoolId:  0,
			DenomTo: "uusdc",
		},
	}

	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: alice.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(usdcLiquidity, usdnLiquidity),
	})
	assert.NoError(t, err)

	// Unbalance the pool towards uusdn
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdn", math.NewInt(int64(250_000_000_000))),
		Routes: routes,
		Min:    sdk.NewCoin("uusdc", math.NewInt(0)),
	})
	require.Nil(t, err)

	routes = []types.Route{
		{
			PoolId:  0,
			DenomTo: "uusdn",
		},
	}

	// ACT
	response, err := server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(int64(1))),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(0)),
	})
	require.Nil(t, err)
	assert.Equal(t, math.NewInt(1).String(), response.Result.Amount.String())

	response, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(int64(10))),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(0)),
	})
	require.Nil(t, err)
	assert.Equal(t, math.NewInt(10).String(), response.Result.Amount.String())

	response, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(int64(100))),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(0)),
	})
	require.Nil(t, err)
	assert.Equal(t, math.NewInt(100).String(), response.Result.Amount.String())

	response, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(int64(1_000))),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(0)),
	})
	require.Nil(t, err)
	assert.Equal(t, math.NewInt(1000).String(), response.Result.Amount.String())
}

func TestConformance(t *testing.T) {
	bob := utils.TestAccount()

	tests := []struct {
		name           string
		msgAdLiquidity *stableswap.MsgAddLiquidity
		msgSwap        *types.MsgSwap
		swapResponse   types.MsgSwapResponse
		error          error
	}{
		{
			"Swap 1",
			&stableswap.MsgAddLiquidity{
				Signer: bob.Address,
				PoolId: 0,
				Amount: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1000*ONE)),
					sdk.NewCoin("uusdc", math.NewInt(1000*ONE)),
				),
			},
			&types.MsgSwap{
				Signer: bob.Address,
				Amount: sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
				Routes: []types.Route{{PoolId: 0, DenomTo: "uusdc"}},
				Min:    sdk.NewCoin("uusdc", math.NewInt(ONE)),
			},
			types.MsgSwapResponse{
				Result: sdk.NewCoin("uusdc", math.NewInt(99974999)),
				Swaps: []*types.Swap{
					{
						PoolId: 0,
						In:     sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
						Out:    sdk.NewCoin("uusdc", math.NewInt(99974999)),
						Fees:   sdk.NewCoins(sdk.NewCoin("uusdn", math.NewInt(24998))),
					},
				},
			},
			nil,
		},
		{
			"Swap 2",
			&stableswap.MsgAddLiquidity{
				Signer: bob.Address,
				PoolId: 0,
				Amount: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
					sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
				),
			},
			&types.MsgSwap{
				Signer: bob.Address,
				Amount: sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
				Routes: []types.Route{{PoolId: 0, DenomTo: "uusdc"}},
				Min:    sdk.NewCoin("uusdc", math.NewInt(ONE)),
			},
			types.MsgSwapResponse{
				Result: sdk.NewCoin("uusdc", math.NewInt(99972764)),
				Swaps: []*types.Swap{
					{
						PoolId: 0,
						In:     sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
						Out:    sdk.NewCoin("uusdc", math.NewInt(99972764)),
						Fees:   sdk.NewCoins(sdk.NewCoin("uusdn", math.NewInt(24998))),
					},
				},
			},
			nil,
		},
		{
			"Swap 3",
			&stableswap.MsgAddLiquidity{
				Signer: bob.Address,
				PoolId: 0,
				Amount: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(10*ONE)),
					sdk.NewCoin("uusdc", math.NewInt(10*ONE)),
				),
			},
			&types.MsgSwap{
				Signer: bob.Address,
				Amount: sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
				Routes: []types.Route{{PoolId: 0, DenomTo: "uusdc"}},
				Min:    sdk.NewCoin("uusdc", math.NewInt(ONE)),
			},
			types.MsgSwapResponse{
				Result: sdk.NewCoin("uusdc", math.NewInt(9997499)),
				Swaps: []*types.Swap{
					{
						PoolId: 0,
						In:     sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
						Out:    sdk.NewCoin("uusdc", math.NewInt(9997499)),
						Fees:   sdk.NewCoins(sdk.NewCoin("uusdn", math.NewInt(2498))),
					},
				},
			},
			nil,
		},
	}
	// ASSERT: Execute each test case and expect the related error and swap results.
	for _, tt := range tests {
		account := mocks.AccountKeeper{
			Accounts: make(map[string]sdk.AccountI),
		}
		bank := mocks.BankKeeper{
			Balances:    make(map[string]sdk.Coins),
			Restriction: mocks.NoOpSendRestrictionFn,
		}
		bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdc", math.NewInt(10000*ONE)))
		bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdn", math.NewInt(10000*ONE)))
		k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
		server := keeper.NewMsgServer(k)
		stableswapServer := keeper.NewStableSwapMsgServer(k)

		t.Run(tt.name, func(t *testing.T) {
			_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
				Signer:                "authority",
				Pair:                  "uusdc",
				ProtocolFeePercentage: 50,
				RewardsFee:            2_500_000,
				InitialA:              1000,
				FutureA:               1000,
				FutureATime:           0,
				RateMultipliers: sdk.NewCoins(
					sdk.NewCoin("uusdn", math.NewInt(1e18)),
					sdk.NewCoin("uusdc", math.NewInt(1e18)),
				),
			})
			require.NoError(t, err)

			_, err = stableswapServer.AddLiquidity(ctx, tt.msgAdLiquidity)
			require.NoError(t, err)

			res, err := server.Swap(ctx, tt.msgSwap)
			require.NoError(t, err)

			assert.Nil(t, err)
			if err != nil {
				require.Equal(t, tt.error.Error(), err.Error())
			}

			assert.Equal(t, tt.swapResponse, *res)
		})
	}
}

func TestRewardsSingleUser(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	server := keeper.NewMsgServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)
	bob, alice := utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Create a Pool.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)})
	_, _ = stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		RewardsFee:            1_000_000,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           0,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})

	// ASSERT: Empty liquidity.
	pool, _ := k.Pools.Get(ctx, 0)
	poolAddress, _ := account.AddressCodec().StringToBytes(pool.Address)
	poolLiquidity := bank.GetAllBalances(ctx, poolAddress)
	assert.Equal(t, len(poolLiquidity), 0)

	// ACT: Attempt to remove liquidity without a position.
	_, err := stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     bob.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(10),
	})
	assert.Error(t, err)

	// ARRANGE: Create a liquidity position for bob.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)})
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdc", math.NewInt(1_000_000_000)))
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdn", math.NewInt(1_000_000_000)))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], sdk.NewCoin("uusdc", math.NewInt(100_000_000)))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], sdk.NewCoin("uusdn", math.NewInt(100_000_000)))
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: bob.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000_000_000)),
			sdk.NewCoin("uusdc", math.NewInt(1000_000_000)),
		),
	})
	assert.NoError(t, err)

	// ARRANGE: Perform swaps and track the total fees amount (minus the protocol fees).
	cumulativeFees := sdk.Coins{}
	for i := 0; i < 1; i++ {
		res, err := server.Swap(ctx, &types.MsgSwap{
			Signer: alice.Address,
			Amount: sdk.NewCoin("uusdc", math.NewInt(10_000_000)),
			Routes: []types.Route{{PoolId: 0, DenomTo: "uusdn"}},
			Min:    sdk.NewCoin("uusdn", math.NewInt(1_999_000)),
		})
		assert.NoError(t, err)
		cumulativeFees = cumulativeFees.Add(res.Swaps[0].Fees...)

		res, err = server.Swap(ctx, &types.MsgSwap{
			Signer: alice.Address,
			Amount: sdk.NewCoin("uusdn", math.NewInt(10_000_000)),
			Routes: []types.Route{{PoolId: 0, DenomTo: "uusdc"}},
			Min:    sdk.NewCoin("uusdc", math.NewInt(1_999_000)),
		})
		assert.NoError(t, err)
		cumulativeFees = cumulativeFees.Add(res.Swaps[0].Fees...)
	}
	expectedRewards := sdk.Coins{}
	for _, coin := range cumulativeFees {
		expectedRewards = append(expectedRewards, sdk.NewCoin(coin.Denom, coin.Amount.Sub(coin.Amount.Mul(math.NewInt(1)).Quo(math.NewInt(100)))))
	}

	// ARRANGE: Increase time.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 2, 1, 1, 1, 1, time.UTC)})

	// ARRANGE: Set failing collections for BondingPositions.
	tmpBondedPositions := k.Stableswap.BondedPositions
	builder := collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName)))
	k.Stableswap.BondedPositions = collections.NewIndexedMap(
		builder, types.StableSwapBondedPositionsPrefix, "stableswap_bonded_positions",
		collections.TripleKeyCodec(collections.Uint64Key, collections.StringKey, collections.Int64Key),
		codec.CollValue[stableswap.BondedPosition](mocks.MakeTestEncodingConfig("noble").Codec),
		stableswapkeeper.NewBondedPositionIndexes(builder),
	)

	// ACT: Attempt to withdraw rewards.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 3, 1, 1, 1, 1, time.UTC)})
	_, err = server.WithdrawRewards(ctx, &types.MsgWithdrawRewards{
		Signer: bob.Address,
	})
	assert.Error(t, err)
	assert.Equal(t, "error accessing store", err.Error())

	// ARRANGE: Restore BondingPositions collection.
	k.Stableswap.BondedPositions = tmpBondedPositions

	// ACT: Withdraw rewards.
	res, err := server.WithdrawRewards(ctx, &types.MsgWithdrawRewards{
		Signer: bob.Address,
	})
	assert.NoError(t, err)

	// ASSERT: Expect the correct amount of withdrawn rewards (minus the protocol fees).
	assert.Equal(t, res.Rewards.Len(), bank.Balances[bob.Address].Len())
	totalRewards := sdk.Coins{}
	for _, reward := range res.Rewards {
		totalRewards = append(totalRewards, reward)
		assert.Equal(t, reward.Amount, bank.Balances[bob.Address].AmountOf(reward.Denom))
	}
	assert.Equal(t, expectedRewards, totalRewards)

	// ARRANGE: Increase the time.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 2, 1, 1, 2, 1, time.UTC)})

	// ACT: Attempt to withdraw rewards again.
	res, err = server.WithdrawRewards(ctx, &types.MsgWithdrawRewards{
		Signer: bob.Address,
	})
	assert.NoError(t, err)

	// ASSERT: Expect 0 rewards.
	assert.True(t, res.Rewards.IsZero())
}

func TestRewardsMultiUser(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	server := keeper.NewMsgServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)
	bob, alice, tom := utils.TestAccount(), utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Create a Pool.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)})
	_, _ = stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		RewardsFee:            1_000_000,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           0,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})

	// ACT: Attempt to remove liquidity without any.
	_, err := stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     bob.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(10),
	})
	assert.Error(t, err)

	// ARRANGE: Give users different amount of balances.
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdc", math.NewInt(500_000_000)))
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdn", math.NewInt(500_000_000)))
	bank.Balances[tom.Address] = append(bank.Balances[tom.Address], sdk.NewCoin("uusdc", math.NewInt(500_000_000)))
	bank.Balances[tom.Address] = append(bank.Balances[tom.Address], sdk.NewCoin("uusdn", math.NewInt(500_000_000)))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], sdk.NewCoin("uusdc", math.NewInt(100_000_000)))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], sdk.NewCoin("uusdn", math.NewInt(100_000_000)))

	// ARRANGE: Provide liquidity from bob & tom.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)})
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: bob.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(500000000)),
			sdk.NewCoin("uusdc", math.NewInt(500000000)),
		),
	})
	assert.NoError(t, err)
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: tom.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(500000000)),
			sdk.NewCoin("uusdc", math.NewInt(500000000)),
		),
	})
	assert.NoError(t, err)

	// ARRANGE: Execute 100 swaps to create rewards.
	cumulativeFees := sdk.Coins{}
	for i := 0; i < 100; i++ {
		res, err := server.Swap(ctx, &types.MsgSwap{
			Signer: alice.Address,
			Amount: sdk.NewCoin("uusdc", math.NewInt(1_000_000)),
			Routes: []types.Route{{PoolId: 0, DenomTo: "uusdn"}},
			Min:    sdk.NewCoin("uusdn", math.NewInt(999_000)),
		})
		assert.NoError(t, err)
		cumulativeFees = cumulativeFees.Add(res.Swaps[0].Fees[0])

		res, err = server.Swap(ctx, &types.MsgSwap{
			Signer: alice.Address,
			Amount: sdk.NewCoin("uusdn", math.NewInt(1_000_000)),
			Routes: []types.Route{{PoolId: 0, DenomTo: "uusdc"}},
			Min:    sdk.NewCoin("uusdc", math.NewInt(999_000)),
		})
		assert.NoError(t, err)
		cumulativeFees = cumulativeFees.Add(res.Swaps[0].Fees[0])
	}

	// ACT: Withdraw Bob's rewards.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 2, 1, 1, 1, 1, time.UTC)})
	res, err := server.WithdrawRewards(ctx, &types.MsgWithdrawRewards{
		Signer: bob.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, res.Rewards.Len(), bank.Balances[bob.Address].Len())
	for _, reward := range res.Rewards {
		assert.Equal(t, reward.Amount, bank.Balances[bob.Address].AmountOf(reward.Denom))
	}

	// ACT: Withdraw Tom's rewards.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 2, 1, 1, 1, 1, time.UTC)})
	res, err = server.WithdrawRewards(ctx, &types.MsgWithdrawRewards{
		Signer: tom.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, res.Rewards.Len(), bank.Balances[tom.Address].Len())
	for _, reward := range res.Rewards {
		assert.Equal(t, reward.Amount, bank.Balances[tom.Address].AmountOf(reward.Denom))
	}
}

func TestWithdrawRewards(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	server := keeper.NewMsgServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)
	bob := utils.TestAccount()

	// ARRANGE: Create a Pool.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)})
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		RewardsFee:            1_000_000,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           0,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})
	assert.NoError(t, err)

	// ACT: Withdraw rewards with an invalid address.
	_, err = server.WithdrawRewards(ctx, &types.MsgWithdrawRewards{
		Signer: bob.Invalid,
	})
	assert.Error(t, err)

	// ARRANGE: Setup up failing collections on Pools.
	tmpPools := k.Pools
	k.Pools = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Get, utils.GetKVStore(ctx, types.ModuleName))),
		types.PoolsPrefix, "pools_generic", collections.Uint64Key, codec.CollValue[types.Pool](mocks.MakeTestEncodingConfig("noble").Codec),
	)

	// ACT: Attempt to withdraw the rewards.
	_, err = server.WithdrawRewards(ctx, &types.MsgWithdrawRewards{
		Signer: bob.Address,
	})
	assert.Equal(t, mocks.ErrorStoreAccess.Error(), err.Error())
	k.Pools = tmpPools

	// ARRANGE: Setup up failing collections on BondedPositions.
	tmpBondedPositions := k.Stableswap.BondedPositions
	builder := collections.NewSchemaBuilder(mocks.FailingStore(mocks.Iterator, utils.GetKVStore(ctx, types.ModuleName)))
	k.Stableswap.BondedPositions = collections.NewIndexedMap(
		builder, types.StableSwapBondedPositionsPrefix, "stableswap_bonded_positions",
		collections.TripleKeyCodec(collections.Uint64Key, collections.StringKey, collections.Int64Key),
		codec.CollValue[stableswap.BondedPosition](mocks.MakeTestEncodingConfig("noble").Codec),
		stableswapkeeper.NewBondedPositionIndexes(builder),
	)

	// ACT: Withdraw the rewards.
	_, err = server.WithdrawRewards(ctx, &types.MsgWithdrawRewards{
		Signer: bob.Address,
	})
	assert.NoError(t, err)
	k.Stableswap.BondedPositions = tmpBondedPositions

	// ARRANGE: Add liquidity.
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdc", math.NewInt(10_000_000)))
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdn", math.NewInt(10_000_000)))
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: bob.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(9_000_000)),
			sdk.NewCoin("uusdn", math.NewInt(9_000_000)),
		),
	})
	assert.NoError(t, err)
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(1_000_000)),
		Routes: []types.Route{{PoolId: 0, DenomTo: "uusdn"}},
		Min:    sdk.NewCoin("uusdn", math.NewInt(100)),
	})
	assert.NoError(t, err)

	// ACT: Withdraw rewards with a too short period.
	_, err = server.WithdrawRewards(ctx, &types.MsgWithdrawRewards{
		Signer: bob.Address,
	})
	assert.Equal(t, "period is too short", err.Error())

	// ARRANGE: Increase the time.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 2, 1, 1, 1, 1, time.UTC)})

	// ACT: Withdraw rewards.
	res, err := server.WithdrawRewards(ctx, &types.MsgWithdrawRewards{
		Signer: bob.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, res.Rewards.Len())

	// ARRANGE: Pause pools.
	_, err = server.PauseByAlgorithm(ctx, &types.MsgPauseByAlgorithm{
		Signer:    "authority",
		Algorithm: types.STABLESWAP,
	})
	assert.NoError(t, err)

	// ACT: Withdraw rewards with paused pools.
	_, err = server.WithdrawRewards(ctx, &types.MsgWithdrawRewards{
		Signer: bob.Address,
	})
	assert.NoError(t, err)

	// ARRANGE: Unpause pools.
	_, err = server.UnpauseByAlgorithm(ctx, &types.MsgUnpauseByAlgorithm{
		Signer:    "authority",
		Algorithm: types.STABLESWAP,
	})
	assert.NoError(t, err)
}

func TestWithdrawProtocolFees(t *testing.T) {
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
	bob, receiver := utils.TestAccount(), utils.TestAccount()
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdc", math.NewInt(10_000*ONE)))
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdn", math.NewInt(10_000*ONE)))

	// ARRANGE: Create a Pool.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)})
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 50,
		RewardsFee:            1_000_000,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           0,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})
	assert.NoError(t, err)

	// ARRANGE: Add liquidity.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: bob.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1_000*ONE)),
			sdk.NewCoin("uusdc", math.NewInt(1_000*ONE)),
		),
	})
	assert.NoError(t, err)

	// ARRANGE: Generate fees with a swap.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
		Routes: []types.Route{
			{
				PoolId:  0,
				DenomTo: "uusdc",
			},
		},
		Min: sdk.NewCoin("uusdc", math.NewInt(90*ONE)),
	})
	assert.NoError(t, err)

	// ARRANGE: Add an invalid Pool.
	err = k.SetPool(ctx, 2, types.Pool{
		Id:        2,
		Address:   "",
		Algorithm: 10,
		Pair:      "",
	})
	assert.NoError(t, err)

	// ASSERT: Correct ProtocolFees amount.
	poolInfo, err := queryServer.Pool(ctx, &types.QueryPool{
		PoolId: 0,
	})
	assert.NoError(t, err)
	assert.Equal(t, math.NewInt(4999), poolInfo.Pool.ProtocolFees.AmountOf("uusdn"))

	// ACT: Withdraw rewards with an invalid authority address.
	_, err = server.WithdrawProtocolFees(ctx, &types.MsgWithdrawProtocolFees{
		Signer: bob.Address,
	})
	assert.Error(t, err)

	// ACT: Withdraw rewards with paused pools.
	_, err = server.PauseByAlgorithm(ctx, &types.MsgPauseByAlgorithm{Signer: "authority", Algorithm: types.STABLESWAP})
	assert.NoError(t, err)

	// ASSERT: No rewards withdrawn.
	_, err = server.WithdrawProtocolFees(ctx, &types.MsgWithdrawProtocolFees{
		Signer: "authority",
		To:     receiver.Address,
	})
	assert.NoError(t, err)
	poolInfo, err = queryServer.Pool(ctx, &types.QueryPool{
		PoolId: 0,
	})
	assert.NoError(t, err)
	assert.Equal(t, math.NewInt(4999), poolInfo.Pool.ProtocolFees.AmountOf("uusdn"))
	assert.Equal(t, math.NewInt(0), bank.Balances[receiver.Address].AmountOf("uusdn"))

	// ARRANGE: Unpause the pools.
	_, err = server.UnpauseByAlgorithm(ctx, &types.MsgUnpauseByAlgorithm{Signer: "authority", Algorithm: types.STABLESWAP})
	assert.NoError(t, err)

	// ACT: Withdraw rewards to an invalid address.
	_, err = server.WithdrawProtocolFees(ctx, &types.MsgWithdrawProtocolFees{
		Signer: "authority",
		To:     receiver.Invalid,
	})
	assert.Equal(t, fmt.Sprintf("unable to decode receiver address: %s", receiver.Invalid), err.Error())

	// ACT: Withdraw rewards to an invalid address.
	_, err = server.WithdrawProtocolFees(ctx, &types.MsgWithdrawProtocolFees{
		Signer: "authority",
		To:     receiver.Invalid,
	})
	assert.Equal(t, fmt.Sprintf("unable to decode receiver address: %s", receiver.Invalid), err.Error())

	// ACT: Withdraw rewards with a valid message.
	_, err = server.WithdrawProtocolFees(ctx, &types.MsgWithdrawProtocolFees{
		Signer: "authority",
		To:     receiver.Address,
	})
	assert.NoError(t, err)

	// ASSERT: Correct ProtocolFees amount.
	poolInfo, err = queryServer.Pool(ctx, &types.QueryPool{
		PoolId: 0,
	})
	assert.NoError(t, err)
	assert.Equal(t, math.NewInt(0), poolInfo.Pool.ProtocolFees.AmountOf("uusdn"))
	assert.Equal(t, math.NewInt(4999), bank.Balances[receiver.Address].AmountOf("uusdn"))
}

func TestPausingAndUnpausingByPoolIds(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	server := keeper.NewMsgServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	user := utils.TestAccount()

	// ARRANGE: Create the initial pools and provide to the user the necessary liquidity.
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdc", math.NewInt(100*ONE)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusde", math.NewInt(100*ONE)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdn", math.NewInt(100*ONE)))
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		InitialA:              100,
		FutureA:               100,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
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

	// ARRANGE: Pause with invalid authority.
	_, err = server.PauseByPoolIds(ctx, &types.MsgPauseByPoolIds{
		Signer:  "user",
		PoolIds: []uint64{0},
	})
	assert.Equal(t, sdkerrors.Wrapf(types.ErrInvalidAuthority, "expected authority, got user").Error(), err.Error())

	// ARRANGE: Simulate Paused failing collection.
	tmpPaused := k.Paused
	k.Paused = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.PausedPrefix, "paused", collections.Uint64Key, codec.BoolValue,
	)
	_, err = server.PauseByPoolIds(ctx, &types.MsgPauseByPoolIds{
		Signer:  "authority",
		PoolIds: []uint64{0},
	})
	assert.Equal(t, sdkerrors.Wrapf(mocks.ErrorStoreAccess, "unable to pause pool").Error(), err.Error())
	k.Paused = tmpPaused

	// ARRANGE: Pause non existing pool.
	res, err := server.PauseByPoolIds(ctx, &types.MsgPauseByPoolIds{
		Signer:  "authority",
		PoolIds: []uint64{100},
	})
	assert.NoError(t, err)
	assert.Equal(t, res.PausedPools, []uint64(nil))

	// ARRANGE: Pause pool by its id.
	res, err = server.PauseByPoolIds(ctx, &types.MsgPauseByPoolIds{
		Signer:  "authority",
		PoolIds: []uint64{0},
	})
	assert.NoError(t, err)
	assert.Equal(t, res.PausedPools, []uint64{0})

	// ASSERT: It is impossible to operate on the pool.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(10)),
			sdk.NewCoin("uusdc", math.NewInt(10)),
		),
	})
	assert.Error(t, err, types.ErrPoolActivityPaused)

	// ARRANGE: pause multiple pools by their ids.
	_, err = server.PauseByPoolIds(ctx, &types.MsgPauseByPoolIds{
		Signer:  "authority",
		PoolIds: []uint64{0, 1},
	})
	assert.NoError(t, err)

	// ASSERT: It is impossible to operate on the pools.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(10)),
			sdk.NewCoin("uusdc", math.NewInt(10)),
		),
	})
	assert.Error(t, err, types.ErrPoolActivityPaused)
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 1,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(500000000)),
			sdk.NewCoin("uusde", math.NewInt(500000000)),
		),
	})
	assert.Error(t, err, types.ErrPoolActivityPaused)

	// ARRANGE: Unpause with invalid authority.
	_, err = server.UnpauseByPoolIds(ctx, &types.MsgUnpauseByPoolIds{
		Signer:  "user",
		PoolIds: []uint64{0},
	})
	assert.Equal(t, sdkerrors.Wrapf(types.ErrInvalidAuthority, "expected authority, got user").Error(), err.Error())

	// ARRANGE: Simulate Paused failing collection.
	tmpPaused = k.Paused
	k.Paused = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.PausedPrefix, "paused", collections.Uint64Key, codec.BoolValue,
	)
	_, err = server.UnpauseByPoolIds(ctx, &types.MsgUnpauseByPoolIds{
		Signer:  "authority",
		PoolIds: []uint64{0},
	})
	assert.Equal(t, sdkerrors.Wrapf(mocks.ErrorStoreAccess, "unable to unpause pool").Error(), err.Error())
	k.Paused = tmpPaused

	// ARRANGE: Unpause non-existing pool.
	res2, err := server.UnpauseByPoolIds(ctx, &types.MsgUnpauseByPoolIds{
		Signer:  "authority",
		PoolIds: []uint64{100},
	})
	assert.NoError(t, err)
	assert.Equal(t, res2.UnpausedPools, []uint64(nil))

	// ARRANGE: Unpause a pools by its id.
	resUnpause, err := server.UnpauseByPoolIds(ctx, &types.MsgUnpauseByPoolIds{
		Signer:  "authority",
		PoolIds: []uint64{0},
	})
	assert.NoError(t, err)
	assert.Equal(t, resUnpause.UnpausedPools, []uint64{0})

	// ASSERT: It is possible to operate on the unpaused pool.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(10*ONE)),
			sdk.NewCoin("uusdc", math.NewInt(10*ONE)),
		),
	})
	assert.NoError(t, err)
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 1,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(10*ONE)),
			sdk.NewCoin("uusde", math.NewInt(10*ONE)),
		),
	})
	assert.Error(t, err, types.ErrPoolActivityPaused)

	// ARRANGE: Re-pause the pools.
	_, err = server.PauseByPoolIds(ctx, &types.MsgPauseByPoolIds{
		Signer:  "authority",
		PoolIds: []uint64{0, 1},
	})
	assert.NoError(t, err)
	// ARRANGE: Unpause multiple pools by their ids.
	resUnpause, err = server.UnpauseByPoolIds(ctx, &types.MsgUnpauseByPoolIds{
		Signer:  "authority",
		PoolIds: []uint64{0, 1},
	})
	assert.NoError(t, err)
	assert.Equal(t, resUnpause.UnpausedPools, []uint64{0, 1})
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 3, 20, 1, 1, 1, 1, time.UTC)})

	// ASSERT: It is now possible to operate on both pools.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(10*ONE)),
			sdk.NewCoin("uusdc", math.NewInt(10*ONE)),
		),
	})
	assert.NoError(t, err)
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 1,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(10*ONE)),
			sdk.NewCoin("uusde", math.NewInt(10*ONE)),
		),
	})
	assert.NoError(t, err)
}

func TestPausingAndUnpausingByAlgorithm(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	server := keeper.NewMsgServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	user := utils.TestAccount()

	// ARRANGE: Create 2 Pools.
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdc", math.NewInt(100*ONE)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusde", math.NewInt(100*ONE)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdn", math.NewInt(100*ONE)))
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

	// ARRANGE: Pause with invalid authority.
	_, err = server.PauseByAlgorithm(ctx, &types.MsgPauseByAlgorithm{
		Signer:    "user",
		Algorithm: types.STABLESWAP,
	})
	assert.Equal(t, sdkerrors.Wrapf(types.ErrInvalidAuthority, "expected authority, got user").Error(), err.Error())

	// ARRANGE: Simulate Paused failing collection.
	tmpPaused := k.Paused
	k.Paused = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.PausedPrefix, "paused", collections.Uint64Key, codec.BoolValue,
	)
	_, err = server.PauseByAlgorithm(ctx, &types.MsgPauseByAlgorithm{
		Signer:    "authority",
		Algorithm: types.STABLESWAP,
	})
	assert.Equal(t, sdkerrors.Wrapf(mocks.ErrorStoreAccess, "unable to pause pool").Error(), err.Error())
	k.Paused = tmpPaused

	// ARRANGE: Pause non-existing algorithm.
	res, err := server.PauseByAlgorithm(ctx, &types.MsgPauseByAlgorithm{
		Signer:    "authority",
		Algorithm: types.UNSPECIFIED,
	})
	assert.NoError(t, err)
	assert.Equal(t, res.PausedPools, []uint64(nil))

	// ARRANGE: Pause pools by algorithm type.
	res, err = server.PauseByAlgorithm(ctx, &types.MsgPauseByAlgorithm{
		Signer:    "authority",
		Algorithm: types.STABLESWAP,
	})
	assert.NoError(t, err)
	assert.Equal(t, res.PausedPools, []uint64{0, 1})

	// ASSERT: pools are paused, activities are blocked.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(10)),
			sdk.NewCoin("uusdc", math.NewInt(10)),
		),
	})
	assert.Error(t, err, types.ErrPoolActivityPaused)
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 1,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(10)),
			sdk.NewCoin("uusde", math.NewInt(10)),
		),
	})
	assert.Error(t, err, types.ErrPoolActivityPaused)

	// ARRANGE: Unpause with invalid authority.
	_, err = server.UnpauseByAlgorithm(ctx, &types.MsgUnpauseByAlgorithm{
		Signer:    "user",
		Algorithm: types.STABLESWAP,
	})
	assert.Equal(t, sdkerrors.Wrapf(types.ErrInvalidAuthority, "expected authority, got user").Error(), err.Error())

	// ARRANGE: Simulate Paused failing collection.
	tmpPaused = k.Paused
	k.Paused = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Set, utils.GetKVStore(ctx, types.ModuleName))),
		types.PausedPrefix, "paused", collections.Uint64Key, codec.BoolValue,
	)
	_, err = server.UnpauseByAlgorithm(ctx, &types.MsgUnpauseByAlgorithm{
		Signer:    "authority",
		Algorithm: types.STABLESWAP,
	})
	assert.Equal(t, sdkerrors.Wrapf(mocks.ErrorStoreAccess, "unable to unpause pool").Error(), err.Error())
	k.Paused = tmpPaused

	// ARRANGE: Unpause non-existing algorithm.
	res2, err := server.UnpauseByAlgorithm(ctx, &types.MsgUnpauseByAlgorithm{
		Signer:    "authority",
		Algorithm: types.UNSPECIFIED,
	})
	assert.NoError(t, err)
	assert.Equal(t, res2.UnpausedPools, []uint64(nil))

	// ARRANGE: Unpause pools by algorithm type.
	resUnpause, err := server.UnpauseByAlgorithm(ctx, &types.MsgUnpauseByAlgorithm{
		Signer:    "authority",
		Algorithm: types.STABLESWAP,
	})
	assert.NoError(t, err)
	assert.Equal(t, resUnpause.UnpausedPools, []uint64{0, 1})

	// ASSERT: Pools are unpaused, normal activity.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(10*ONE)),
			sdk.NewCoin("uusdc", math.NewInt(10*ONE)),
		),
	})
	assert.NoError(t, err)
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 1,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(10*ONE)),
			sdk.NewCoin("uusde", math.NewInt(10*ONE)),
		),
	})
	assert.NoError(t, err)
}

func TestSwapAgainstBondedLiquidity(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)

	server := keeper.NewMsgServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)
	provider, user := utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Create a Pool.
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
		RewardsFee: 4e3,
		InitialA:   100,
		FutureA:    100,
	})
	assert.NoError(t, err)
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)})
	bank.Balances[provider.Address] = append(bank.Balances[provider.Address], sdk.NewCoin("uusdc", math.NewInt(1100*ONE)))
	bank.Balances[provider.Address] = append(bank.Balances[provider.Address], sdk.NewCoin("uusdn", math.NewInt(1100*ONE)))

	// ARRANGE: Add liquidity.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: provider.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000*ONE)),
			sdk.NewCoin("uusdc", math.NewInt(1000*ONE)),
		),
	})
	assert.NoError(t, err)

	// ARRANGE: Create a liquidity position for the user.
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdc", math.NewInt(100*ONE)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdn", math.NewInt(100*ONE)))
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(50*ONE)),
			sdk.NewCoin("uusdc", math.NewInt(50*ONE)),
		),
	})
	assert.NoError(t, err)
	assert.Equal(t, math.NewInt(50*ONE), bank.Balances[user.Address].AmountOf("uusdc"))

	// ACT: Attempt to swap with not enough balance since bonded.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: user.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(51*ONE)),
		Routes: []types.Route{{PoolId: 0, DenomTo: "uusdn"}},
		Min:    sdk.NewCoin("uusdn", math.NewInt(30*ONE)),
	})
	// ASSERT: The action should've failed due to bonded balance.
	assert.ErrorIs(t, types.ErrInsufficientBalance, err)

	// ACT: Begin unbonding the position
	_, err = stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     user.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(100),
	})
	assert.NoError(t, err)

	// ACT: Attempt to swap without enough balance since bonding time.
	err = k.BeginBlocker(ctx)
	assert.NoError(t, err)
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: user.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(51*ONE)),
		Routes: []types.Route{{PoolId: 0, DenomTo: "uusdn"}},
		Min:    sdk.NewCoin("uusdn", math.NewInt(30*ONE)),
	})
	// ASSERT: The action should've failed due to low balance since unbonding time is not yet elapsed.
	assert.ErrorIs(t, types.ErrInsufficientBalance, err)

	// ACT: Attempt to swap with not enough balance since bonding time.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 10, 1, 1, 1, 1, time.UTC)})
	err = k.BeginBlocker(ctx)
	assert.NoError(t, err)
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: user.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(51)),
		Routes: []types.Route{{PoolId: 0, DenomTo: "uusdn"}},
		Min:    sdk.NewCoin("uusdn", math.NewInt(30)),
	})
	// ASSERT: The action should've succeeded.
	assert.NoError(t, err)
}

func TestSwap(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	server := keeper.NewMsgServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	alice, bob := utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Create the Pool.
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		RewardsFee:            4e3,
		ProtocolFeePercentage: 1,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           1893452400,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})
	assert.Nil(t, err)
	nLiquidity := sdk.NewCoin("uusdn", math.NewInt(1_000_000*ONE))
	usdcLiquidity := sdk.NewCoin("uusdc", math.NewInt(1_000_000*ONE))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], usdcLiquidity)
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], nLiquidity)

	routes := []types.Route{
		{
			PoolId:  0,
			DenomTo: "uusdn",
		},
	}

	// ACT: Attempt to perform a Swap with an empty balance.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(0)),
	})
	assert.ErrorIs(t, types.ErrInsufficientBalance, err)

	// ARRANGE: Add funds to Bob.
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdc", math.NewInt(1_000*ONE)))
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdx", math.NewInt(10*ONE)))

	// ACT: Attempt to perform a Swap without Pool liquidity.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(10)),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(10)),
	})
	assert.Equal(t, "error computing swap routes plan: pool liquidity must be positive", err.Error())

	// ARRANGE: Add liquidity to the Pool.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: alice.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(usdcLiquidity, nLiquidity),
	})
	assert.NoError(t, err)

	// ACT: Attempt to perform a Swap with an invalid address.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Invalid,
		Amount: sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(0)),
	})
	assert.Errorf(t, err, "unable to decode signer address: %s", bob.Invalid)

	// ACT: Attempt to perform a Swap with an invalid route.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		Routes: []types.Route{
			{
				PoolId:  0,
				DenomTo: "uusdx",
			},
		},
		Min: sdk.NewCoin("uusdn", math.NewInt(0)),
	})
	assert.ErrorIs(t, types.ErrInvalidSwapRoutingPlan, err)

	// ACT: Attempt to perform a Swap with an invalid route.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		Routes: []types.Route{
			{
				PoolId:  0,
				DenomTo: "uusdc",
			},
		},
		Min: sdk.NewCoin("uusdc", math.NewInt(0)),
	})
	assert.Equal(t, "error computing swap routes plan: cannot swap for the same denom uusdc: invalid swap routing plan", err.Error())

	// ACT: Attempt to perform a Swap with an invalid route.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdx", math.NewInt(1*ONE)),
		Routes: []types.Route{
			{
				PoolId:  0,
				DenomTo: "uusdc",
			},
		},
		Min: sdk.NewCoin("uusdc", math.NewInt(1)),
	})
	assert.Equal(t, "error computing swap routes plan: uusdx is not a paired asset in pool 0: invalid swap routing plan", err.Error())

	// ARRANGE: Create a second pool without liquidity.
	_, err = stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusde",
		RewardsFee:            4e3,
		ProtocolFeePercentage: 1,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           1893452400,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusde", math.NewInt(1000000000000000000)),
		),
	})
	assert.Nil(t, err)

	// ACT: Attempt to perform a Swap with an invalid pool id.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		Routes: []types.Route{
			{
				PoolId:  1,
				DenomTo: "uusde",
			},
		},
		Min: sdk.NewCoin("uusdn", math.NewInt(0)),
	})
	assert.ErrorIs(t, types.ErrInvalidSwapRoutingPlan, err)

	// ARRANGE: Pause the routed Pool.
	_, err = server.PauseByAlgorithm(ctx, &types.MsgPauseByAlgorithm{Signer: "authority", Algorithm: types.STABLESWAP})
	assert.NoError(t, err)

	// ACT: Attempt to perform a Swap within a paused Pool.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(10)),
	})
	// ASSERT: Routing Plan validation failed.
	assert.Error(t, err)

	// ARRANGE: Unpause the routed Pool.
	_, err = server.UnpauseByAlgorithm(ctx, &types.MsgUnpauseByAlgorithm{Signer: "authority", Algorithm: types.STABLESWAP})
	assert.NoError(t, err)

	// ACT: Attempt to perform a Swap with a failing slippage condition.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
	})
	// ASSERT: Routing Plan validation failed.
	assert.Error(t, err)

	// ARRANGE: Set up failing collections for the StableSwap Pool.
	tmpPools := k.Pools
	k.Pools = collections.NewMap(
		collections.NewSchemaBuilder(mocks.FailingStore(mocks.Get, utils.GetKVStore(ctx, types.ModuleName))),
		types.PoolsPrefix, "pools_generic", collections.Uint64Key, codec.CollValue[types.Pool](mocks.MakeTestEncodingConfig("noble").Codec),
	)

	// ACT: Attempt to perform a Swap with a failing Pool collection.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
	})
	// ASSERT: Routing Plan validation failed.
	assert.Error(t, err)

	// ARRANGE: Restore collection.
	k.Pools = tmpPools

	// ACT: Attempt to perform a Swap with an invalid DenomTo and Min.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		Routes: []types.Route{
			{
				PoolId:  0,
				DenomTo: "uusdx",
			},
		},
		Min: sdk.NewCoin("uusdx", math.NewInt(99*ONE)),
	})
	assert.Equal(t, "error computing swap routes plan: pool 0 doesn't contain denom uusdx: invalid swap routing plan", err.Error())

	// ACT: Perform a valid Swap.
	response, err := server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(99999900)),
	})
	// ACT: Expect a successful swap and validate all the resulting amounts.
	assert.Nil(t, err)
	assert.Equal(t, sdk.NewCoin("uusdn", math.NewInt(99999959)), response.Result)
	pool, _ := k.Pools.Get(ctx, 0)
	assert.Equal(t, bank.Balances[pool.Address].AmountOf("uusdc"), math.NewInt(nLiquidity.Amount.Int64()+(100*ONE)-response.Swaps[0].Fees.AmountOf("uusdc").Int64()))
	assert.Equal(t, bank.Balances[pool.Address].AmountOf("uusdn"), math.NewInt(usdcLiquidity.Amount.Int64()-response.Result.Amount.Int64()))
	assert.Equal(t, bank.Balances[bob.Address].AmountOf("uusdc"), math.NewInt((1_000-100)*ONE))
	assert.Equal(t, bank.Balances[bob.Address].AmountOf("uusdn"), response.Result.Amount)

	// ARRANGE: Set a different InitialA value.
	_, err = stableswapServer.UpdatePool(ctx, &stableswap.MsgUpdatePool{
		Signer:                "authority",
		PoolId:                0,
		RewardsFee:            4e3,
		ProtocolFeePercentage: 1,
		InitialA:              1,
		FutureA:               100,
		FutureATime:           1893452400,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})
	assert.NoError(t, err)

	// ACT: Perform a new Swap expecting higher fees.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		Routes: routes,
		Min:    sdk.NewCoin("uusdn", math.NewInt(99990000)),
	})
	assert.NoError(t, err)
}

func TestMultiPoolSwap(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	server := keeper.NewMsgServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)
	alice, tom, bob := utils.TestAccount(), utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Create 2 Pools.
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		RewardsFee:            4e3,
		ProtocolFeePercentage: 1,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           1893452400,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})
	assert.Nil(t, err)
	_, err = stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusde",
		RewardsFee:            4e3,
		ProtocolFeePercentage: 1,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           1893452400,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusde", math.NewInt(1000000000000000000)),
		),
	})
	assert.Nil(t, err)

	// ARRANGE: Provide liquidity in both Pools.
	nLiquidity := sdk.NewCoin("uusdn", math.NewInt(1_000_000*ONE))
	usdcLiquidity := sdk.NewCoin("uusdc", math.NewInt(1_000_000*ONE))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], usdcLiquidity)
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], nLiquidity)
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: alice.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(usdcLiquidity, nLiquidity),
	})
	require.NoError(t, err)
	nLiquidity2 := sdk.NewCoin("uusdn", math.NewInt(1_000_000*ONE))
	usdeLiquidity := sdk.NewCoin("uusde", math.NewInt(1_000_000*ONE))
	bank.Balances[tom.Address] = append(bank.Balances[tom.Address], usdeLiquidity)
	bank.Balances[tom.Address] = append(bank.Balances[tom.Address], nLiquidity2)
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: tom.Address,
		PoolId: 1,
		Amount: sdk.NewCoins(nLiquidity, usdeLiquidity),
	})
	require.NoError(t, err)

	// ACT: Perform a multi-route swap.
	routes := []types.Route{
		{
			PoolId:  0,
			DenomTo: "uusdn",
		},
		{
			PoolId:  1,
			DenomTo: "uusde",
		},
	}
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdc", math.NewInt(1_000*ONE)))
	response, err := server.Swap(ctx, &types.MsgSwap{
		Signer: bob.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
		Routes: routes,
		Min:    sdk.NewCoin("uusde", math.NewInt(99999900)),
	})
	assert.Nil(t, err)

	// ASSERT: Expect matching values in state.
	assert.Equal(t, sdk.NewCoin("uusde", math.NewInt(99999918)), response.Result)
	pool0, _ := k.Pools.Get(ctx, 0)
	assert.Equal(t, bank.Balances[pool0.Address].AmountOf("uusdc"), math.NewInt(nLiquidity.Amount.Int64()+(100*ONE)-response.Swaps[0].Fees.AmountOf("uusdc").Int64()))
	assert.Equal(t, bank.Balances[pool0.Address].AmountOf("uusdn"), math.NewInt(usdcLiquidity.Amount.Int64()-response.Swaps[0].Out.Amount.Int64()))
	assert.Equal(t, bank.Balances[bob.Address].AmountOf("uusdc"), math.NewInt((1_000-100)*ONE))
	assert.Equal(t, bank.Balances[bob.Address].AmountOf("uusdn"), math.ZeroInt())
	assert.Equal(t, bank.Balances[bob.Address].AmountOf("uusde"), response.Swaps[len(response.Swaps)-1].Out.Amount)
}

func BenchmarkSwap(b *testing.B) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(b, account, bank)
	server := keeper.NewMsgServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	alice, bob := utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Create a Pool.
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		RewardsFee:            4e3,
		ProtocolFeePercentage: 1,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           1893452400,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
		),
	})
	assert.Nil(b, err)

	// ARRANGE: Provide liquidity.
	nLiquidity := sdk.NewCoin("uusdn", math.NewInt(1_000_000_000*ONE))
	usdcLiquidity := sdk.NewCoin("uusdc", math.NewInt(1_000_000_000*ONE))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], usdcLiquidity)
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], nLiquidity)
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: alice.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(usdcLiquidity, nLiquidity),
	})
	assert.NoError(b, err)

	// ACT: Perform swaps.
	routes := []types.Route{
		{
			PoolId:  0,
			DenomTo: "uusdn",
		},
	}
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdc", math.NewInt(1_000_000_000*ONE)))
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdn", math.NewInt(1_000_000_000*ONE)))
	for i := 0; i < b.N; i++ {
		_, err = server.Swap(ctx, &types.MsgSwap{
			Signer: bob.Address,
			Amount: sdk.NewCoin("uusdc", math.NewInt(1*ONE)),
			Routes: routes,
			Min:    sdk.NewCoin("uusdn", math.NewInt(1_000)),
		})
		assert.NoError(b, err)
	}
}

func BenchmarkMultiPoolSwap(b *testing.B) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(b, account, bank)
	server := keeper.NewMsgServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	alice, tom, bob := utils.TestAccount(), utils.TestAccount(), utils.TestAccount()

	// ARRANGE: Create 2 Pools.
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		RewardsFee:            4e3,
		ProtocolFeePercentage: 1,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           1893452400,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})
	assert.Nil(b, err)
	_, err = stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusde",
		RewardsFee:            4e3,
		ProtocolFeePercentage: 1,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           1893452400,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusde", math.NewInt(1000000000000000000)),
		),
	})
	assert.Nil(b, err)

	// ARRANGE: Provide liquidity in both Pools.
	nLiquidity := sdk.NewCoin("uusdn", math.NewInt(1_000_000*ONE))
	usdcLiquidity := sdk.NewCoin("uusdc", math.NewInt(1_000_000*ONE))
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], usdcLiquidity)
	bank.Balances[alice.Address] = append(bank.Balances[alice.Address], nLiquidity)
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: alice.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(usdcLiquidity, nLiquidity),
	})
	require.NoError(b, err)
	nLiquidity2 := sdk.NewCoin("uusdn", math.NewInt(1_000_000*ONE))
	usdeLiquidity := sdk.NewCoin("uusde", math.NewInt(1_000_000*ONE))
	bank.Balances[tom.Address] = append(bank.Balances[tom.Address], usdeLiquidity)
	bank.Balances[tom.Address] = append(bank.Balances[tom.Address], nLiquidity2)
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: tom.Address,
		PoolId: 1,
		Amount: sdk.NewCoins(nLiquidity, usdeLiquidity),
	})
	require.NoError(b, err)

	// ACT: Perform multi-pool swaps.
	routes := []types.Route{
		{
			PoolId:  0,
			DenomTo: "uusdn",
		},
		{
			PoolId:  1,
			DenomTo: "uusde",
		},
	}
	bank.Balances[bob.Address] = append(bank.Balances[bob.Address], sdk.NewCoin("uusdc", math.NewInt(1_000_000_000*ONE)))
	for i := 0; i < b.N; i++ {
		_, err = server.Swap(ctx, &types.MsgSwap{
			Signer: bob.Address,
			Amount: sdk.NewCoin("uusdc", math.NewInt(1*ONE)),
			Routes: routes,
			Min:    sdk.NewCoin("uusde", math.NewInt(1_000)),
		})
		assert.NoError(b, err)
	}
}
