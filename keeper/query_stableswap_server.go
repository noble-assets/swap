package keeper

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/types/errors"
	stableswapkeeper "swap.noble.xyz/keeper/stableswap"
	"swap.noble.xyz/types/stableswap"
)

var _ stableswap.QueryServer = &queryStableSwapServer{}

type queryStableSwapServer struct {
	*stableswapkeeper.Keeper
}

func NewStableSwapQueryServer(keeper *stableswapkeeper.Keeper) stableswap.QueryServer {
	return queryStableSwapServer{Keeper: keeper}
}

func (s queryStableSwapServer) PositionsByProvider(ctx context.Context, req *stableswap.QueryPositionsByProvider) (*stableswap.QueryPositionsByProviderResponse, error) {
	if req == nil || req.Provider == "" {
		return nil, errors.ErrInvalidRequest
	}

	var positions []stableswap.QueryPositionsByProviderResponse_Position
	for _, position := range s.GetBondedPositionsByProvider(ctx, req.Provider) {
		positions = append(positions, stableswap.QueryPositionsByProviderResponse_Position{
			PoolId:    position.PoolId,
			Shares:    position.BondedPosition.Balance,
			Timestamp: time.Unix(position.Timestamp, 0),
		})
	}

	return &stableswap.QueryPositionsByProviderResponse{Positions: positions}, nil
}

func (s queryStableSwapServer) UnbondingPositionsByProvider(ctx context.Context, req *stableswap.QueryUnbondingPositionsByProvider) (*stableswap.QueryUnbondingPositionsByProviderResponse, error) {
	if req == nil || req.Provider == "" {
		return nil, errors.ErrInvalidRequest
	}

	var positions []stableswap.QueryUnbondingPositionsByProviderResponse_UnbondingPosition
	for _, position := range s.GetUnbondingPositionsByProvider(ctx, req.Provider) {
		positions = append(positions, stableswap.QueryUnbondingPositionsByProviderResponse_UnbondingPosition{
			PoolId:          position.PoolId,
			UnbondingShares: position.UnbondingPosition.Shares,
			EndTime:         position.UnbondingPosition.EndTime,
		})
	}

	return &stableswap.QueryUnbondingPositionsByProviderResponse{UnbondingPositions: positions}, nil
}
