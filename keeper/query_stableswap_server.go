package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"
	"time"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"swap.noble.xyz/types/stableswap"
)

var _ stableswap.QueryServer = &queryStableSwapServer{}

type queryStableSwapServer struct {
	keeper *Keeper
}

func NewStableSwapQueryServer(keeper *Keeper) stableswap.QueryServer {
	return queryStableSwapServer{keeper: keeper}
}

func (s queryStableSwapServer) PositionsByProvider(ctx context.Context, req *stableswap.QueryPositionsByProvider) (*stableswap.QueryPositionsByProviderResponse, error) {
	if req == nil || req.Provider == "" {
		return nil, errors.ErrInvalidRequest
	}

	var bondedPositions []stableswap.QueryBondedPositionResponseEntry
	if rawBondedPositions, err := s.BondedPositionsByProvider(ctx, &stableswap.QueryBondedPositionsByProvider{Provider: req.Provider}); err == nil && rawBondedPositions != nil {
		bondedPositions = append(bondedPositions, rawBondedPositions.BondedPositions...)
	}

	var unbondingPositions []stableswap.QueryUnbondingPositionResponseEntry
	if rawUnbondingPositions, err := s.UnbondingPositionsByProvider(ctx, &stableswap.QueryUnbondingPositionsByProvider{Provider: req.Provider}); err == nil && rawUnbondingPositions != nil {
		unbondingPositions = append(unbondingPositions, rawUnbondingPositions.UnbondingPositions...)
	}

	var rewards []stableswap.QueryRewardsResponseEntry
	if rawRewards, err := s.RewardsByProvider(ctx, &stableswap.QueryRewardsByProvider{Provider: req.Provider}); err == nil && rawRewards != nil {
		rewards = append(rewards, rawRewards.Rewards...)
	}

	return &stableswap.QueryPositionsByProviderResponse{
		BondedPositions:    bondedPositions,
		UnbondingPositions: unbondingPositions,
		Rewards:            rewards,
	}, nil
}

func (s queryStableSwapServer) RewardsByProvider(ctx context.Context, req *stableswap.QueryRewardsByProvider) (*stableswap.QueryRewardsByProviderResponse, error) {
	if req == nil || req.Provider == "" {
		return nil, errors.ErrInvalidRequest
	}

	var rewards []stableswap.QueryRewardsResponseEntry
	for _, controller := range GetStableSwapControllers(ctx, s.keeper) {
		poolRewards, err := controller.GetTotalPoolUserRewards(ctx, req.Provider, s.keeper.headerService.GetHeaderInfo(ctx).Time)
		if err != nil {
			continue
		}
		amount := sdk.Coins{}
		for _, reward := range poolRewards {
			amount = amount.Add(reward.Amount...)
		}
		rewards = append(rewards, stableswap.QueryRewardsResponseEntry{
			PoolId: controller.GetId(),
			Amount: amount,
		})
	}

	// Sort Rewards by Pool id.
	sort.Slice(rewards, func(i, j int) bool {
		return rewards[i].PoolId < rewards[j].PoolId
	})
	return &stableswap.QueryRewardsByProviderResponse{
		Rewards: rewards,
	}, nil
}

func (s queryStableSwapServer) BondedPositionsByProvider(ctx context.Context, req *stableswap.QueryBondedPositionsByProvider) (*stableswap.QueryBondedPositionsByProviderResponse, error) {
	if req == nil || req.Provider == "" {
		return nil, errors.ErrInvalidRequest
	}

	var positions []stableswap.QueryBondedPositionResponseEntry
	for _, position := range s.keeper.Stableswap.GetBondedPositionsByProvider(ctx, req.Provider) {
		positions = append(positions, stableswap.QueryBondedPositionResponseEntry{
			PoolId:    position.PoolId,
			Shares:    position.BondedPosition.Balance,
			Timestamp: time.Unix(position.Timestamp, 0),
		})
	}

	// Sort Rewards by Pool id.
	sort.Slice(positions, func(i, j int) bool {
		return positions[i].PoolId < positions[j].PoolId
	})
	return &stableswap.QueryBondedPositionsByProviderResponse{BondedPositions: positions}, nil
}

func (s queryStableSwapServer) UnbondingPositionsByProvider(ctx context.Context, req *stableswap.QueryUnbondingPositionsByProvider) (*stableswap.QueryUnbondingPositionsByProviderResponse, error) {
	if req == nil || req.Provider == "" {
		return nil, errors.ErrInvalidRequest
	}

	var positions []stableswap.QueryUnbondingPositionResponseEntry
	for _, position := range s.keeper.Stableswap.GetUnbondingPositionsByProvider(ctx, req.Provider) {
		positions = append(positions, stableswap.QueryUnbondingPositionResponseEntry{
			PoolId:          position.PoolId,
			UnbondingShares: position.UnbondingPosition.Shares,
			EndTime:         position.UnbondingPosition.EndTime,
		})
	}

	// Sort Rewards by Pool id.
	sort.Slice(positions, func(i, j int) bool {
		return positions[i].PoolId < positions[j].PoolId
	})
	return &stableswap.QueryUnbondingPositionsByProviderResponse{UnbondingPositions: positions}, nil
}
