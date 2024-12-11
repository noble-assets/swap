package types

const ModuleName = "swap"

// Prefixes must be unique and cannot overlap, meaning one cannot be the start of another (ex. <pools> and <pools_stableswap> are invalid).
// Distinct prefixes with shared roots (ex. <pools_generic> and <pools_stableswap>) are allowed as long as neither fully contains the other.

var (
	NextPoolIDPrefix = []byte("next_pool_id")
	PausedPrefix     = []byte("paused")
	PoolsPrefix      = []byte("pools_generic")

	StableSwapPoolsPrefix                     = []byte("stableswap_pools")
	StableSwapUsersTotalBondedSharesPrefix    = []byte("stableswap_users_total_bonded_shares")
	StableSwapUsersTotalUnbondingSharesPrefix = []byte("stableswap_users_total_unbonding_shares")
	StableSwapPoolsTotalUnbondingSharesPrefix = []byte("stableswap_pools_total_unbonding_shares")
	StableSwapUnbondingPositionsPrefix        = []byte("stableswap_unbonding_positions")
	StableSwapBondedPositionPrefix            = []byte("stableswap_bonded_positions")
)
