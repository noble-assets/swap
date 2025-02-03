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
	"time"

	"cosmossdk.io/core/header"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"swap.noble.xyz/keeper"
	"swap.noble.xyz/types"
	"swap.noble.xyz/types/stableswap"
	"swap.noble.xyz/utils"
	"swap.noble.xyz/utils/mocks"
)

func TestBondedPositionsByProvider(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	stableswapQueryServer := keeper.NewStableSwapQueryServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	user := utils.TestAccount()

	// ARRANGE: Create a Pool.
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

	// ACT: Attempt to query with an invalid request.
	_, err = stableswapQueryServer.BondedPositionsByProvider(ctx, &stableswap.QueryBondedPositionsByProvider{})
	assert.Error(t, err)

	// ACT: Attempt to query with a valid request but 0 positions.
	res, err := stableswapQueryServer.BondedPositionsByProvider(ctx, &stableswap.QueryBondedPositionsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(res.BondedPositions))

	// ARRANGE: Create a provider position by adding liquidity.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(10*ONE)),
			sdk.NewCoin("uusdn", math.NewInt(10*ONE)),
		),
	})
	assert.NoError(t, err)

	// ACT: Attempt to query with a valid request and 1 current position.
	res, err = stableswapQueryServer.BondedPositionsByProvider(ctx, &stableswap.QueryBondedPositionsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.BondedPositions))

	// ARRANGE: Create a second pool.
	_, err = stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusde",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusde", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
		),
		InitialA: 100,
		FutureA:  100,
	})
	assert.NoError(t, err)

	// ARRANGE: Create a provider position by adding liquidity.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 1,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusde", math.NewInt(10*ONE)),
			sdk.NewCoin("uusdn", math.NewInt(10*ONE)),
		),
	})
	assert.NoError(t, err)
	// ACT: Attempt to query with a multiple current position.
	res, err = stableswapQueryServer.BondedPositionsByProvider(ctx, &stableswap.QueryBondedPositionsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.BondedPositions))
}

func TestUnbondingBondedPositionsByProvider(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	stableswapQueryServer := keeper.NewStableSwapQueryServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)

	user := utils.TestAccount()

	// ARRANGE: Create a Pool.
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdc", math.NewInt(100*ONE)))
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

	// ACT: Attempt to query with an invalid request.
	_, err = stableswapQueryServer.UnbondingPositionsByProvider(ctx, &stableswap.QueryUnbondingPositionsByProvider{})
	assert.Error(t, err)

	// ACT: Attempt to query with a valid request but 0 positions.
	res, err := stableswapQueryServer.UnbondingPositionsByProvider(ctx, &stableswap.QueryUnbondingPositionsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(res.UnbondingPositions))

	// ARRANGE: Create a provider position by adding liquidity.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(10*ONE)),
			sdk.NewCoin("uusdn", math.NewInt(10*ONE)),
		),
	})
	assert.NoError(t, err)

	// ARRANGE: Remove 10% of the total user liquidity.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC)})
	_, err = stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     user.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(10),
	})
	assert.NoError(t, err)

	// ACT: Query the unbonding positions.
	res, err = stableswapQueryServer.UnbondingPositionsByProvider(ctx, &stableswap.QueryUnbondingPositionsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.UnbondingPositions))

	// ARRANGE: Remove another 20% of the total user liquidity.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC)})
	_, err = stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     user.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(10),
	})
	assert.NoError(t, err)

	// ACT: Query the unbonding positions.
	res, err = stableswapQueryServer.UnbondingPositionsByProvider(ctx, &stableswap.QueryUnbondingPositionsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.UnbondingPositions))
}

func TestRewardsByProvider(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	stableswapQueryServer := keeper.NewStableSwapQueryServer(k)
	stableswapServer := keeper.NewStableSwapMsgServer(k)
	server := keeper.NewMsgServer(k)

	user := utils.TestAccount()

	// ARRANGE: Create a Pool.
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdc", math.NewInt(10_000*ONE)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusde", math.NewInt(10_000*ONE)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdn", math.NewInt(10_000*ONE)))
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		RewardsFee:            100000,
		MaxFee:                100000,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           0,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})
	assert.NoError(t, err)

	// ACT: Query the rewards with an invalid request.
	_, err = stableswapQueryServer.RewardsByProvider(ctx, nil)
	assert.Error(t, err)

	// ACT: Attempt to query with a valid request but 0 positions.
	rewards, err := stableswapQueryServer.RewardsByProvider(ctx, &stableswap.QueryRewardsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(rewards.Rewards))
	assert.Equal(t, 0, len(rewards.Rewards[0].Amount))

	// ARRANGE: Create a provider position by adding liquidity.
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(1000*ONE)),
			sdk.NewCoin("uusdn", math.NewInt(1000*ONE)),
		),
	})
	assert.NoError(t, err)

	// ARRANGE: Remove 10% of the total user liquidity.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC)})
	_, err = stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     user.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(10),
	})
	assert.NoError(t, err)

	// ARRANGE: Remove another 20% of the total user liquidity.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC)})
	_, err = stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     user.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(10),
	})
	assert.NoError(t, err)

	// ACT: Generate rewards.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: user.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(51*ONE)),
		Routes: []types.Route{{PoolId: 0, DenomTo: "uusdn"}},
		Min:    sdk.NewCoin("uusdn", math.NewInt(1*ONE)),
	})

	// ASSERT: The action should've succeeded.
	assert.NoError(t, err)

	// ACT: Attempt to query the rewards.
	rewards, err = stableswapQueryServer.RewardsByProvider(ctx, &stableswap.QueryRewardsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(rewards.Rewards))
	assert.Equal(t, 1, len(rewards.Rewards[0].Amount))

	// ARRANGE: Create a second pool.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC)})
	_, err = stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusde",
		ProtocolFeePercentage: 1,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusde", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
		),
		InitialA: 100,
		FutureA:  100,
	})
	assert.NoError(t, err)

	// ARRANGE: Create a provider position by adding liquidity.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC)})
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 1,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusde", math.NewInt(100*ONE)),
			sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
		),
	})
	assert.NoError(t, err)

	// ACT: Query the rewards.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 5, 0, 0, 0, time.UTC)})
	rewards, err = stableswapQueryServer.RewardsByProvider(ctx, &stableswap.QueryRewardsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(rewards.Rewards))
	assert.Equal(t, 1, len(rewards.Rewards[0].Amount))
	assert.Equal(t, 0, len(rewards.Rewards[1].Amount))

	// ACT: Generate rewards.
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: user.Address,
		Amount: sdk.NewCoin("uusde", math.NewInt(51*ONE)),
		Routes: []types.Route{{PoolId: 1, DenomTo: "uusdn"}},
		Min:    sdk.NewCoin("uusdn", math.NewInt(1*ONE)),
	})
	assert.NoError(t, err)

	// ACT: Query the rewards.
	rewards, err = stableswapQueryServer.RewardsByProvider(ctx, &stableswap.QueryRewardsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(rewards.Rewards))
	assert.Equal(t, 1, len(rewards.Rewards[0].Amount))

	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)})
	_, err = stableswapQueryServer.RewardsByProvider(ctx, &stableswap.QueryRewardsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
}

func TestPositions(t *testing.T) {
	account := mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.SwapKeeperWithKeepers(t, account, bank)
	stableswapServer := keeper.NewStableSwapMsgServer(k)
	stableswapQueryServer := keeper.NewStableSwapQueryServer(k)
	server := keeper.NewMsgServer(k)

	user := utils.TestAccount()

	// ARRANGE: Create a Pool.
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdc", math.NewInt(10_000*ONE)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusde", math.NewInt(10_000*ONE)))
	bank.Balances[user.Address] = append(bank.Balances[user.Address], sdk.NewCoin("uusdn", math.NewInt(10_000*ONE)))
	_, err := stableswapServer.CreatePool(ctx, &stableswap.MsgCreatePool{
		Signer:                "authority",
		Pair:                  "uusdc",
		ProtocolFeePercentage: 1,
		RewardsFee:            100000,
		MaxFee:                100000,
		InitialA:              100,
		FutureA:               100,
		FutureATime:           0,
		RateMultipliers: sdk.NewCoins(
			sdk.NewCoin("uusdn", math.NewInt(1000000000000000000)),
			sdk.NewCoin("uusdc", math.NewInt(1000000000000000000)),
		),
	})
	assert.NoError(t, err)

	// ARRANGE: Create a provider position by adding liquidity.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC)})
	_, err = stableswapServer.AddLiquidity(ctx, &stableswap.MsgAddLiquidity{
		Signer: user.Address,
		PoolId: 0,
		Amount: sdk.NewCoins(
			sdk.NewCoin("uusdc", math.NewInt(100*ONE)),
			sdk.NewCoin("uusdn", math.NewInt(100*ONE)),
		),
	})
	assert.NoError(t, err)

	// ACT: Query the rewards.
	_, err = stableswapQueryServer.RewardsByProvider(ctx, &stableswap.QueryRewardsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)

	// ACT: Generate rewards.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 5, 0, 0, 0, time.UTC)})
	_, err = server.Swap(ctx, &types.MsgSwap{
		Signer: user.Address,
		Amount: sdk.NewCoin("uusdc", math.NewInt(51*ONE)),
		Routes: []types.Route{{PoolId: 0, DenomTo: "uusdn"}},
		Min:    sdk.NewCoin("uusdn", math.NewInt(1*ONE)),
	})
	assert.NoError(t, err)

	// ACT: Add an unbonding position.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 6, 0, 0, 0, time.UTC)})
	_, err = stableswapServer.RemoveLiquidity(ctx, &stableswap.MsgRemoveLiquidity{
		Signer:     user.Address,
		PoolId:     0,
		Percentage: math.LegacyNewDec(50),
	})
	assert.NoError(t, err)

	// ACT: Attempt to query with an invalid request.
	ctx = ctx.WithHeaderInfo(header.Info{Time: time.Date(2020, 1, 1, 7, 0, 0, 0, time.UTC)})
	_, err = stableswapQueryServer.PositionsByProvider(ctx, nil)
	assert.Error(t, err)

	// ACT: Query the positions.
	_, err = stableswapQueryServer.PositionsByProvider(ctx, &stableswap.QueryPositionsByProvider{
		Provider: user.Address,
	})
	assert.NoError(t, err)
}
