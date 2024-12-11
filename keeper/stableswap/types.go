package stableswap

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
	stableswaptypes "swap.noble.xyz/types/stableswap"
)

type BondedPositionIndexes struct {
	ByProvider *indexes.Multi[string, collections.Triple[uint64, string, int64], stableswaptypes.BondedPosition]
}

func (i BondedPositionIndexes) IndexesList() []collections.Index[collections.Triple[uint64, string, int64], stableswaptypes.BondedPosition] {
	return []collections.Index[collections.Triple[uint64, string, int64], stableswaptypes.BondedPosition]{i.ByProvider}
}

func NewBondedPositionIndexes(builder *collections.SchemaBuilder) BondedPositionIndexes {
	return BondedPositionIndexes{
		ByProvider: indexes.NewMulti(
			builder, []byte("position_by_provider"), "position_by_provider",
			collections.StringKey,
			collections.TripleKeyCodec(collections.Uint64Key, collections.StringKey, collections.Int64Key),
			func(key collections.Triple[uint64, string, int64], shares stableswaptypes.BondedPosition) (string, error) {
				return key.K2(), nil
			},
		),
	}
}

// UnbondingPositions Indexes

type UnbondingPositionIndexes struct {
	ByProvider *indexes.Multi[string, collections.Triple[int64, string, uint64], stableswaptypes.UnbondingPosition]
}

func (i UnbondingPositionIndexes) IndexesList() []collections.Index[collections.Triple[int64, string, uint64], stableswaptypes.UnbondingPosition] {
	return []collections.Index[collections.Triple[int64, string, uint64], stableswaptypes.UnbondingPosition]{i.ByProvider}
}

func NewUnbondingPositionIndexes(builder *collections.SchemaBuilder) UnbondingPositionIndexes {
	return UnbondingPositionIndexes{
		ByProvider: indexes.NewMulti(
			builder, []byte("unbonding_by_provider"), "unbonding_by_provider",
			collections.StringKey,
			collections.TripleKeyCodec(collections.Int64Key, collections.StringKey, collections.Uint64Key),
			func(key collections.Triple[int64, string, uint64], position stableswaptypes.UnbondingPosition) (string, error) {
				return key.K2(), nil
			},
		),
	}
}
