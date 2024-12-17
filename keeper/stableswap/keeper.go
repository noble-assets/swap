package stableswap

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/header"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	swapv1 "swap.noble.xyz/api/v1"
	"swap.noble.xyz/types"
	stableswaptypes "swap.noble.xyz/types/stableswap"
)

type Keeper struct {
	cdc codec.BinaryCodec

	eventService  event.Service
	headerService header.Service
	logger        log.Logger

	Pools                     collections.Map[uint64, stableswaptypes.Pool]
	PoolsTotalUnbondingShares collections.Map[uint64, math.LegacyDec]
	UsersTotalBondedShares    collections.Map[collections.Pair[uint64, string], math.LegacyDec]
	UsersTotalUnbondingShares collections.Map[collections.Pair[uint64, string], math.LegacyDec]
	BondedPositions           *collections.IndexedMap[collections.Triple[uint64, string, int64], stableswaptypes.BondedPosition, BondedPositionIndexes]
	UnbondingPositions        *collections.IndexedMap[collections.Triple[int64, string, uint64], stableswaptypes.UnbondingPosition, UnbondingPositionIndexes]
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	eventService event.Service,
	headerService header.Service,
	logger log.Logger,
) *Keeper {
	builder := collections.NewSchemaBuilder(storeService)

	return &Keeper{
		cdc:           cdc,
		eventService:  eventService,
		headerService: headerService,
		logger:        logger,

		Pools:                     collections.NewMap(builder, types.StableSwapPoolsPrefix, "stableswap_pools", collections.Uint64Key, codec.CollValue[stableswaptypes.Pool](cdc)),
		UsersTotalBondedShares:    collections.NewMap(builder, types.StableSwapUsersTotalBondedSharesPrefix, "stableswap_users_total_shares", collections.PairKeyCodec(collections.Uint64Key, collections.StringKey), sdk.LegacyDecValue),
		UsersTotalUnbondingShares: collections.NewMap(builder, types.StableSwapUsersTotalUnbondingSharesPrefix, "stableswap_users_total_unbonding_shares", collections.PairKeyCodec(collections.Uint64Key, collections.StringKey), sdk.LegacyDecValue),
		PoolsTotalUnbondingShares: collections.NewMap(builder, types.StableSwapPoolsTotalUnbondingSharesPrefix, "stableswap_pools_total_unbonding_shares", collections.Uint64Key, sdk.LegacyDecValue),

		BondedPositions: collections.NewIndexedMap(
			builder, types.StableSwapBondedPositionsPrefix, "stableswap_bonded_positions",
			collections.TripleKeyCodec(collections.Uint64Key, collections.StringKey, collections.Int64Key),
			codec.CollValue[stableswaptypes.BondedPosition](cdc),
			NewBondedPositionIndexes(builder),
		),
		UnbondingPositions: collections.NewIndexedMap(
			builder, types.StableSwapUnbondingPositionsPrefix, "stableswap_unbonding_positions",
			collections.TripleKeyCodec(collections.Int64Key, collections.StringKey, collections.Uint64Key),
			codec.CollValue[stableswaptypes.UnbondingPosition](cdc),
			NewUnbondingPositionIndexes(builder),
		),
	}
}

func (k *Keeper) Logger() log.Logger {
	return k.logger.With("module", types.ModuleName, "algorithm", swapv1.Algorithm_STABLESWAP.String())
}
