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

package keeper

import (
	"context"
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	swapv1 "swap.noble.xyz/api/v1"
	"swap.noble.xyz/types"
	"swap.noble.xyz/types/stableswap"
)

var _ stableswap.MsgServer = &stableswapMsgServer{}

type stableswapMsgServer struct {
	*Keeper
}

func NewStableSwapMsgServer(keeper *Keeper) stableswap.MsgServer {
	return stableswapMsgServer{Keeper: keeper}
}

// CreatePool creates a new `StableSwap` Pool with the provided params.
func (s stableswapMsgServer) CreatePool(ctx context.Context, msg *stableswap.MsgCreatePool) (*stableswap.MsgCreatePoolResponse, error) {
	// Ensure that the signer has the required authority.
	if msg.Signer != s.authority {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", s.authority, msg.Signer)
	}

	// Check if the Pair denom is valid and exists on chain.
	if msg.Pair == "" {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "missing pair value")
	}

	// Check if the Pair denom is different from the base denom.
	if msg.Pair == s.baseDenom {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "pair denom must be different from %s", s.baseDenom)
	}

	// Check if the Pair denom exists on the bank module.
	if !s.bankKeeper.GetSupply(ctx, msg.Pair).IsPositive() {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "%s does not exists on chain", msg.Pair)
	}

	// Check if the Initial A value is valid.
	if msg.InitialA <= 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "invalid InitialA value")
	}

	// Manually sort and validate the rate multipliers.
	rateMultipliers := msg.RateMultipliers.Sort()
	if msg.RateMultipliers == nil || rateMultipliers.Len() != 2 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "RateMultipliers length must be 2, got %d", msg.RateMultipliers.Len())
	}
	if !rateMultipliers.AmountOf(msg.Pair).IsPositive() {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "%s rate multiplier must be positive, got %s",
			msg.Pair,
			rateMultipliers.AmountOf(msg.Pair).String(),
		)
	}
	if !rateMultipliers.AmountOf(s.baseDenom).IsPositive() {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "%s rate multiplier must be positive, got %s",
			s.baseDenom,
			rateMultipliers.AmountOf(s.baseDenom).String(),
		)
	}

	// If set ensure that the MaxFee is positive.
	if msg.MaxFee < 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "MaxFee cannot be negative")
	}

	// If set ensure that the MaxFee is positive.
	if msg.RewardsFee < 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "RewardsFee cannot be negative")
	}

	// If set ensure that the ProtocolFeePercentage is positive.
	if msg.ProtocolFeePercentage < 0 || msg.ProtocolFeePercentage > 100 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "invalid ProtocolFeePercentage value")
	}

	// Check if a Pool with the same Algorithm and Pair already exists.
	algorithm := types.Algorithm(swapv1.Algorithm_STABLESWAP)
	for _, pool := range s.GetPools(ctx) {
		if pool.Pair == msg.Pair && pool.Algorithm == algorithm {
			return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "pool with pair %s and %s algorithm already exists", msg.Pair, algorithm.String())
		}
	}

	// Increase and get the next Pool ID.
	poolId, err := s.IncreaseNextPoolID(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to set next pool id")
	}

	// Create the Pool address.
	prefix := fmt.Sprintf("%s/pool/%d", types.ModuleName, poolId)
	account := authtypes.NewEmptyModuleAccount(prefix)
	account = s.accountKeeper.NewAccount(ctx, account).(*authtypes.ModuleAccount)
	s.accountKeeper.SetModuleAccount(ctx, account)

	// Create the Protocol Fees Pool address.
	protocolFeesAccount := authtypes.NewEmptyModuleAccount(fmt.Sprintf("%s/protocol_fees", prefix))
	protocolFees := s.accountKeeper.NewAccount(ctx, protocolFeesAccount).(*authtypes.ModuleAccount)
	s.accountKeeper.SetModuleAccount(ctx, protocolFees)

	// Create the Rewards Fees Pool address.
	rewardFeesAccount := authtypes.NewEmptyModuleAccount(fmt.Sprintf("%s/reward_fees", prefix))
	rewardFees := s.accountKeeper.NewAccount(ctx, rewardFeesAccount).(*authtypes.ModuleAccount)
	s.accountKeeper.SetModuleAccount(ctx, rewardFees)

	// Set the new Pool on state.
	if err = s.SetPool(ctx, poolId, types.Pool{
		Id:        poolId,
		Address:   account.GetAddress().String(),
		Algorithm: algorithm,
		Pair:      msg.Pair,
	}); err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to set pool")
	}

	// Add the new Pool ID to the `Paused` state.
	if err = s.SetPaused(ctx, poolId, false); err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to create paused entry")
	}

	// Set the `StableSwap` data on state.
	if err = s.Stableswap.SetPool(ctx, poolId, stableswap.Pool{
		ProtocolFeePercentage: msg.ProtocolFeePercentage,
		RewardsFee:            msg.RewardsFee,
		MaxFee:                msg.MaxFee,
		InitialA:              msg.InitialA,
		FutureA:               msg.FutureA,
		InitialATime:          s.headerService.GetHeaderInfo(ctx).Time.Unix(),
		FutureATime:           msg.FutureATime,
		RateMultipliers:       msg.RateMultipliers,
		TotalShares:           math.LegacyZeroDec(),
	}); err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to set stableswap pool")
	}

	return &stableswap.MsgCreatePoolResponse{}, s.eventService.EventManager(ctx).Emit(ctx, &stableswap.PoolCreated{
		Algorithm:             algorithm.String(),
		ProtocolFeePercentage: msg.ProtocolFeePercentage,
		RewardsFee:            msg.RewardsFee,
		MaxFee:                msg.MaxFee,
		InitialA:              msg.InitialA,
		FutureA:               msg.FutureA,
		InitialATime:          s.headerService.GetHeaderInfo(ctx).Time.Unix(),
		FutureATime:           msg.FutureATime,
		RateMultipliers:       msg.RateMultipliers,
	})
}

// UpdatePool updates the params of the `StableSwap` Pool.
func (s stableswapMsgServer) UpdatePool(ctx context.Context, msg *stableswap.MsgUpdatePool) (*stableswap.MsgUpdatePoolResponse, error) {
	// Ensure that the signer has the required authority.
	if msg.Signer != s.authority {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", s.authority, msg.Signer)
	}

	// Get the Pool Controller.
	controller, err := GetStableSwapController(ctx, s.Keeper, msg.PoolId)
	if err != nil {
		return nil, err
	}

	// Ensure that the requested Pool is a StableSwap pool.
	if controller.GetAlgorithm() != types.STABLESWAP {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPool, "invalid pool algorithm")
	}

	// Check if the A values are valid.
	if msg.InitialA <= 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "invalid InitialA value")
	}

	// Manually sort and validate the rate multipliers.
	rateMultipliers := msg.RateMultipliers.Sort()
	if rateMultipliers == nil || rateMultipliers.Len() != 2 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "RateMultipliers length must be 2, got %d", msg.RateMultipliers.Len())
	}
	if !rateMultipliers.AmountOf(controller.GetPair()).IsPositive() {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "%s rate multiplier must be positive, got %s",
			controller.GetPair(),
			rateMultipliers.AmountOf(controller.GetPair()).String(),
		)
	}
	if !rateMultipliers.AmountOf(s.baseDenom).IsPositive() {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "%s rate multiplier must be positive, got %s",
			s.baseDenom,
			rateMultipliers.AmountOf(s.baseDenom).String(),
		)
	}

	// If set ensure that the MaxFee is positive.
	if msg.MaxFee < 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "MaxFee cannot be negative")
	}

	// If set ensure that the MaxFee is positive.
	if msg.RewardsFee < 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "RewardsFee cannot be negative")
	}

	// If set ensure that the ProtocolFeePercentage is positive.
	if msg.ProtocolFeePercentage < 0 || msg.ProtocolFeePercentage > 100 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolParams, "invalid ProtocolFeePercentage value")
	}

	if err = controller.UpdatePool(
		ctx,
		msg.ProtocolFeePercentage,
		msg.RewardsFee,
		msg.MaxFee,
		msg.InitialA,
		s.headerService.GetHeaderInfo(ctx).Time.Unix(),
		msg.FutureA,
		msg.FutureATime,
		rateMultipliers,
	); err != nil {
		return nil, err
	}

	return &stableswap.MsgUpdatePoolResponse{}, s.eventService.EventManager(ctx).Emit(ctx, &stableswap.PoolUpdated{
		ProtocolFeePercentage: msg.ProtocolFeePercentage,
		RewardsFee:            msg.RewardsFee,
		MaxFee:                msg.MaxFee,
		FutureA:               msg.FutureA,
		FutureATime:           msg.FutureATime,
		RateMultipliers:       msg.RateMultipliers,
	})
}

// RemoveLiquidity allows a user to remove liquidity from a `StableSwap` liquidity pool.
func (s stableswapMsgServer) RemoveLiquidity(ctx context.Context, msg *stableswap.MsgRemoveLiquidity) (*stableswap.MsgRemoveLiquidityResponse, error) {
	// Check if the provider address is valid.
	_, err := s.addressCodec.StringToBytes(msg.Signer)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to decode provider address %s", msg.Signer)
	}

	// Check if the unbonding percentage is a valid number.
	if msg.Percentage.LT(math.LegacyZeroDec()) || msg.Percentage.GT(math.LegacyNewDec(100)) {
		return nil, types.ErrInvalidUnbondPercentage
	}

	// Get the StableswapController associated to the Pool.
	stableswapController, err := GetStableSwapController(ctx, s.Keeper, msg.PoolId)
	if err != nil {
		return nil, err
	}

	// Allow interactions only if the Pool is not paused.
	if stableswapController.IsPaused() {
		return nil, sdkerrors.Wrapf(types.ErrPoolActivityPaused, "pool %d is paused", msg.PoolId)
	}

	// Calculate the new user Unbonding BondedPosition to apply.
	unbondingCommitment, err := stableswapController.RemoveLiquidity(ctx, s.headerService.GetHeaderInfo(ctx).Time, msg)
	if err != nil {
		return nil, err
	}

	return &stableswap.MsgRemoveLiquidityResponse{
			UnbondingShares: unbondingCommitment.UnbondingPosition.Shares,
		}, s.eventService.EventManager(ctx).Emit(ctx, &stableswap.LiquidityRemoved{
			Provider:   msg.Signer,
			PoolId:     msg.PoolId,
			Amount:     unbondingCommitment.UnbondingPosition.Amount,
			Shares:     unbondingCommitment.UnbondingPosition.Shares,
			UnlockTime: unbondingCommitment.UnbondingPosition.EndTime,
		})
}

// AddLiquidity allows a user to add liquidity to a `StableSwap` liquidity pool.
func (s stableswapMsgServer) AddLiquidity(ctx context.Context, msg *stableswap.MsgAddLiquidity) (*stableswap.MsgAddLiquidityResponse, error) {
	// Check if the provider address is valid.
	provider, err := s.addressCodec.StringToBytes(msg.Signer)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to decode provider address %s", msg.Signer)
	}

	// Get the `StableSwap` Controller.
	stableswapController, err := GetStableSwapController(ctx, s.Keeper, msg.PoolId)
	if err != nil {
		return nil, err
	}

	// Allow interactions only if the Pool is not paused.
	if stableswapController.IsPaused() {
		return nil, sdkerrors.Wrapf(types.ErrPoolActivityPaused, "pool %d is paused", msg.PoolId)
	}

	// Sort and validate the amount.
	amount := msg.Amount.Sort()

	// Check if the pairs are provided correctly.
	if !amount.AmountOf(s.baseDenom).IsPositive() {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAmount, "must provide positive amount of %s", s.baseDenom)
	}
	if !amount.AmountOf(stableswapController.GetPair()).IsPositive() {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAmount, "must provide positive amount of %s", stableswapController.GetPair())
	}
	// Check if the input coins to add are valid coins.
	if msg.Amount.Len() != 2 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAmount, "coins should be 2, got %d", msg.Amount.Len())
	}

	// Get the Pool Address.
	poolAddress, err := s.addressCodec.StringToBytes(stableswapController.GetAddress())
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to decode pool address, got %s", stableswapController.GetAddress())
	}

	// Validate the liquidity amount.
	baseRate := stableswapController.GetRate(ctx)
	baseAmount := amount.AmountOf(s.baseDenom).ToLegacyDec()
	pairAmount := amount.AmountOf(stableswapController.GetPair()).ToLegacyDec()
	if !pairAmount.TruncateInt().Equal(baseAmount.Mul(baseRate).TruncateInt()) {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidAmount,
			"must provide balanced amount of %s%s and %s%s",
			baseAmount.TruncateInt().String(),
			s.baseDenom,
			baseRate.Mul(pairAmount).TruncateInt().String(),
			stableswapController.GetPair(),
		)
	}

	// Ensure that deposit amount of the base token is at least 1 unit (1e6).
	if baseAmount.LT(math.LegacyNewDec(1e6)) {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidAmount,
			"must provide a minimum amount of 1000000%s but got: %s%s",
			s.baseDenom,
			baseAmount.TruncateInt().String(),
			s.baseDenom,
		)
	}

	// Create the new user BondedPosition.
	newPosition, err := stableswapController.AddLiquidity(ctx, s.headerService.GetHeaderInfo(ctx).Time, msg)
	if err != nil {
		return nil, err
	}

	// Transfer the tokens to the Pool.
	if err = s.bankKeeper.SendCoins(ctx, provider, poolAddress, amount); err != nil {
		return nil, sdkerrors.Wrap(err, "unable to transfer from provider to pool")
	}

	return &stableswap.MsgAddLiquidityResponse{
			MintedShares: newPosition.BondedPosition.Balance.TruncateInt64(),
		}, s.eventService.EventManager(ctx).Emit(ctx, &stableswap.LiquidityAdded{
			Provider: msg.Signer,
			PoolId:   msg.PoolId,
			Amount:   amount,
			Shares:   newPosition.BondedPosition.Balance,
		})
}
