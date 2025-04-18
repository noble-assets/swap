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

	"github.com/gogo/protobuf/sortkeys"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"swap.noble.xyz/types"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return msgServer{Keeper: keeper}
}

// Swap allows a user to swap one type of token for another, using multiple routes.
func (s msgServer) Swap(ctx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
	// Compute the Swap (date validation is performed internally).
	result, err := s.Keeper.Swap(ctx, msg)
	if err != nil {
		return nil, err
	}

	// Get the sum of all the fees.
	fees := sdk.Coins{}
	for _, swap := range result.Swaps {
		fees = fees.Add(swap.Fees...)
	}

	return result, s.eventService.EventManager(ctx).Emit(ctx, &types.Swapped{
		Signer: msg.Signer,
		Input:  msg.Amount,
		Output: result.Result,
		Routes: msg.Routes,
		Fees:   fees,
	})
}

// PauseByAlgorithm pauses all pools using a specific algorithm.
func (s msgServer) PauseByAlgorithm(ctx context.Context, msg *types.MsgPauseByAlgorithm) (*types.MsgPauseByAlgorithmResponse, error) {
	// Ensure that the signer has the required authority.
	if msg.Signer != s.authority {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", s.authority, msg.Signer)
	}

	// Iterate through the pools and collect the pools to pause.
	var poolsToPause []uint64 // store the pools to pause
	for poolId, pool := range s.GetPools(ctx) {
		if msg.Algorithm != types.UNSPECIFIED && pool.Algorithm == msg.Algorithm {
			poolsToPause = append(poolsToPause, poolId)
		}
	}

	// Iterate and pause all the wanted pools
	var pausedPools []uint64
	for _, poolId := range poolsToPause {
		if s.HasPool(ctx, poolId) {
			err := s.SetPaused(ctx, poolId, true) // set to true to pause
			if err != nil {
				return nil, sdkerrors.Wrap(err, "unable to pause pool")
			}
			pausedPools = append(pausedPools, poolId)
		}
	}

	sortkeys.Uint64s(pausedPools)
	return &types.MsgPauseByAlgorithmResponse{
			PausedPools: pausedPools,
		}, s.eventService.EventManager(ctx).Emit(ctx, &types.PoolsPaused{
			PoolIds: poolsToPause,
		})
}

// PauseByPoolIds pauses specific pools identified by their pool IDs.
func (s msgServer) PauseByPoolIds(ctx context.Context, msg *types.MsgPauseByPoolIds) (*types.MsgPauseByPoolIdsResponse, error) {
	// Ensure that the signer has the required authority.
	if msg.Signer != s.authority {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", s.authority, msg.Signer)
	}

	// Iterate and pause all the wanted pools.
	var pausedPools []uint64
	for _, poolId := range msg.PoolIds {
		if s.HasPool(ctx, poolId) {
			err := s.SetPaused(ctx, poolId, true) // set to true to pause
			if err != nil {
				return nil, sdkerrors.Wrap(err, "unable to pause pool")
			}
			pausedPools = append(pausedPools, poolId)
		}
	}

	sortkeys.Uint64s(pausedPools)
	return &types.MsgPauseByPoolIdsResponse{
			PausedPools: pausedPools,
		}, s.eventService.EventManager(ctx).Emit(ctx, &types.PoolsPaused{
			PoolIds: pausedPools,
		})
}

// UnpauseByAlgorithm unpauses all pools using a specific algorithm.
func (s msgServer) UnpauseByAlgorithm(ctx context.Context, msg *types.MsgUnpauseByAlgorithm) (*types.MsgUnpauseByAlgorithmResponse, error) {
	// Ensure that the signer has the required authority.
	if msg.Signer != s.authority {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", s.authority, msg.Signer)
	}

	// Iterate through the pools and collect the pools to unpause.
	var poolsToUnpause []uint64 // store the pools to unpause
	for poolId, pool := range s.GetPools(ctx) {
		if msg.Algorithm != types.UNSPECIFIED && pool.Algorithm == msg.Algorithm {
			poolsToUnpause = append(poolsToUnpause, poolId)
		}
	}

	// Iterate and unpause all the wanted pools.
	var unpausedPools []uint64
	for _, poolId := range poolsToUnpause {
		if s.HasPool(ctx, poolId) {
			err := s.SetPaused(ctx, poolId, false) // set false to unpause
			if err != nil {
				return nil, sdkerrors.Wrap(err, "unable to unpause pool")
			}
			unpausedPools = append(unpausedPools, poolId)
		}
	}

	sortkeys.Uint64s(unpausedPools)
	return &types.MsgUnpauseByAlgorithmResponse{
			UnpausedPools: unpausedPools,
		}, s.eventService.EventManager(ctx).Emit(ctx, &types.PoolsPaused{
			PoolIds: unpausedPools,
		})
}

// UnpauseByPoolIds unpauses specific pools identified by their pool IDs.
func (s msgServer) UnpauseByPoolIds(ctx context.Context, msg *types.MsgUnpauseByPoolIds) (*types.MsgUnpauseByPoolIdsResponse, error) {
	// Ensure that the signer has the required authority.
	if msg.Signer != s.authority {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", s.authority, msg.Signer)
	}

	// Iterate and unpause all the wanted pools.
	var unpausedPools []uint64
	for _, poolId := range msg.PoolIds {
		if s.HasPool(ctx, poolId) {
			err := s.SetPaused(ctx, poolId, false) // set false to unpause
			if err != nil {
				return nil, sdkerrors.Wrap(err, "unable to unpause pool")
			}
			unpausedPools = append(unpausedPools, poolId)
		}
	}

	sortkeys.Uint64s(unpausedPools)
	return &types.MsgUnpauseByPoolIdsResponse{
			UnpausedPools: unpausedPools,
		}, s.eventService.EventManager(ctx).Emit(ctx, &types.PoolsPaused{
			PoolIds: unpausedPools,
		})
}

// WithdrawProtocolFees allows the protocol to withdraw accumulated fees and move them to another account.
func (s msgServer) WithdrawProtocolFees(ctx context.Context, msg *types.MsgWithdrawProtocolFees) (*types.MsgWithdrawProtocolFeesResponse, error) {
	// Ensure that the signer has the required authority.
	if msg.Signer != s.authority {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", s.authority, msg.Signer)
	}

	// Ensure that the receiver is a valid address.
	receiver, err := s.addressCodec.StringToBytes(msg.To)
	if err != nil {
		return nil, fmt.Errorf("unable to decode receiver address: %s", msg.To)
	}

	// Collect all the Pools protocol fees addresses and amounts.
	var poolsProtocolFeesAddresses []sdk.AccAddress
	for poolId := range s.GetPools(ctx) {
		controller, err := GetGenericController(ctx, s.Keeper, poolId)
		if err != nil {
			continue
		}

		// Skip processing if the pool is paused.
		if controller.IsPaused() {
			continue
		}

		poolsProtocolFeesAddresses = append(poolsProtocolFeesAddresses, controller.GetProtocolFeesAddresses()...)
	}

	// Send all the collected amounts to the provided address.
	rewards := sdk.Coins{}
	for _, poolSender := range poolsProtocolFeesAddresses {
		balances := s.bankKeeper.GetAllBalances(ctx, poolSender)
		if err = s.bankKeeper.SendCoins(ctx, poolSender, receiver, balances); err != nil {
			return nil, err
		}
		rewards = append(rewards, balances...)
	}

	return &types.MsgWithdrawProtocolFeesResponse{}, s.eventService.EventManager(ctx).Emit(ctx, &types.WithdrawnProtocolFees{
		To:      msg.To,
		Rewards: rewards,
	})
}

// WithdrawRewards allows a user to claim their accumulated rewards.
func (s msgServer) WithdrawRewards(ctx context.Context, msg *types.MsgWithdrawRewards) (*types.MsgWithdrawRewardsResponse, error) {
	// Ensure that the signer is a valid address.
	_, err := s.addressCodec.StringToBytes(msg.Signer)
	if err != nil {
		return nil, fmt.Errorf("unable to decode user address: %s", msg.Signer)
	}

	// Iterate and collects the user rewards from all the pools.
	rewards := sdk.Coins{}
	currentTime := s.headerService.GetHeaderInfo(ctx)
	for poolId := range s.GetPools(ctx) {
		controller, err := GetGenericController(ctx, s.Keeper, poolId)
		if err != nil {
			return nil, err
		}

		// Skip processing if the pool is paused.
		if controller.IsPaused() {
			continue
		}

		// Process the user rewards.
		poolRewards, err := controller.ProcessUserRewards(ctx, msg.Signer, currentTime.Time)
		if err != nil {
			return nil, err
		}
		rewards = rewards.Add(poolRewards...)
	}

	return &types.MsgWithdrawRewardsResponse{
			Rewards: rewards,
		}, s.eventService.EventManager(ctx).Emit(ctx, &types.WithdrawnRewards{
			Signer:  msg.Signer,
			Rewards: rewards,
		})
}
