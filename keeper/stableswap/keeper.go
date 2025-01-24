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

	// Pools stores the StableSwap pools in the state, mapped by their unique uint64 ID.
	Pools collections.Map[uint64, stableswaptypes.Pool]

	// UsersTotalBondedShares keeps track of the total bonded shares for each user in a pool,
	// keyed by a pair of pool ID (uint64) and user address (string).
	UsersTotalBondedShares collections.Map[collections.Pair[uint64, string], math.LegacyDec]

	// UsersTotalUnbondingShares keeps track of the total unbonding shares for each user in a pool,
	// keyed by a pair of pool ID (uint64) and user address (string).
	UsersTotalUnbondingShares collections.Map[collections.Pair[uint64, string], math.LegacyDec]

	// PoolsTotalUnbondingShares tracks the total unbonding shares for each pool, keyed by the pool ID (uint64).
	PoolsTotalUnbondingShares collections.Map[uint64, math.LegacyDec]

	// BondedPositions stores the bonded positions of users in a pool,
	// keyed by a triple of pool ID (uint64), user address (string), and current_time (int64).
	BondedPositions *collections.IndexedMap[collections.Triple[uint64, string, int64], stableswaptypes.BondedPosition, BondedPositionIndexes]

	// UnbondingPositions stores the unbonding positions of users in a pool,
	// keyed by a triple of unbonding_end_time (int64), user address (string), and pool ID (uint64).
	UnbondingPositions *collections.IndexedMap[collections.Triple[int64, string, uint64], stableswaptypes.UnbondingPosition, UnbondingPositionIndexes]
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

		Pools:                     collections.NewMap(builder, types.StableSwapPoolsPrefix, "pools_stableswap", collections.Uint64Key, codec.CollValue[stableswaptypes.Pool](cdc)),
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
