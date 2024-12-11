package keeper

import (
	"context"
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gogo/protobuf/sortkeys"
	"swap.noble.xyz/types"
)

var _ types.QueryServer = &queryServer{}

type queryServer struct {
	*Keeper
}

func NewQueryServer(keeper *Keeper) types.QueryServer {
	return queryServer{Keeper: keeper}
}

func (s queryServer) Pool(ctx context.Context, req *types.QueryPool) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	controller, err := GetGenericController(ctx, s.Keeper, req.PoolId)
	if err != nil {
		return nil, err
	}

	return &types.QueryPoolResponse{Pool: &types.PoolDetails{
		Id:           controller.GetId(),
		Address:      controller.GetAddress(),
		Algorithm:    controller.GetAlgorithm(),
		Pairs:        []string{s.baseDenom, controller.GetPair()},
		Details:      controller.PoolDetails(),
		Liquidity:    s.bankKeeper.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(controller.GetAddress())),
		ProtocolFees: s.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(fmt.Sprintf("%s/pool/%d/protocol_fees", types.ModuleName, controller.GetId()))),
		RewardFees:   s.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(fmt.Sprintf("%s/pool/%d/rewards_fees", types.ModuleName, controller.GetId()))),
	}}, err
}

func (s queryServer) Pools(ctx context.Context, req *types.QueryPools) (*types.QueryPoolsResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	var pools []*types.PoolDetails
	for _, pool := range s.GetPools(ctx) {
		controller, err := GetGenericController(ctx, s.Keeper, pool.Id)
		if err != nil {
			continue
		}

		pools = append(pools, &types.PoolDetails{
			Id:           pool.Id,
			Address:      pool.Address,
			Algorithm:    pool.Algorithm,
			Pairs:        []string{s.baseDenom, pool.Pair},
			Details:      controller.PoolDetails(),
			Liquidity:    s.bankKeeper.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(pool.Address)),
			ProtocolFees: s.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(fmt.Sprintf("%s/pool/%d/protocol_fees", types.ModuleName, pool.Id))),
			RewardFees:   s.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(fmt.Sprintf("%s/pool/%d/rewards_fees", types.ModuleName, pool.Id))),
		})
	}
	return &types.QueryPoolsResponse{Pools: pools}, nil
}

func (s queryServer) SimulateSwap(ctx context.Context, req *types.MsgSwap) (*types.MsgSwapResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	// Create a cached context from the current context. This cached context will
	// simulate the swap operation without committing any state changes to the main
	// store. By doing so, we can observe the potential effects of the swap without
	// permanently altering the real state.
	cacheCtx, _ := sdk.UnwrapSDKContext(ctx).CacheContext()
	return s.Keeper.Swap(cacheCtx, req)
}

func (s queryServer) Paused(ctx context.Context, req *types.QueryPaused) (*types.QueryPausedResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	var pausedPools []uint64
	for poolId, isPaused := range s.GetPaused(ctx) {
		if isPaused {
			pausedPools = append(pausedPools, poolId)
		}
	}
	sortkeys.Uint64s(pausedPools)

	return &types.QueryPausedResponse{PausedPools: pausedPools}, nil
}

func (s queryServer) Rates(ctx context.Context, req *types.QueryRates) (*types.QueryRatesResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	var rates []types.Rate

	for _, pool := range s.GetPools(ctx) {
		if req.Algorithm != types.UNSPECIFIED && req.Algorithm != pool.Algorithm {
			continue
		}

		controller, err := GetGenericController(ctx, s.Keeper, pool.Id)
		if err != nil {
			continue
		}
		rates = append(rates, controller.GetRates(ctx)...)
	}

	// Sort rates by Denom first, then by Vs
	sort.Slice(rates, func(i, j int) bool {
		if rates[i].Denom == rates[j].Denom {
			return rates[i].Vs < rates[j].Vs // Secondary sort by Vs
		}
		return rates[i].Denom < rates[j].Denom // Primary sort by Denom
	})

	return &types.QueryRatesResponse{Rates: rates}, nil
}

func (s queryServer) Rate(ctx context.Context, req *types.QueryRate) (*types.QueryRateResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	var rates []types.Rate
	for _, pool := range s.GetPools(ctx) {
		if req.Algorithm != types.UNSPECIFIED && req.Algorithm != pool.Algorithm {
			continue
		}

		controller, err := GetGenericController(ctx, s.Keeper, pool.Id)
		if err != nil {
			continue
		}

		for _, rate := range controller.GetRates(ctx) {
			if rate.Denom == req.Denom {
				rates = append(rates, rate)
			}
		}
	}

	return &types.QueryRateResponse{Rates: rates}, nil
}
