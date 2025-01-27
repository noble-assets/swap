# StableSwap State

## StableSwapPools

The `StableSwapPools` field is a collection (`collections.Map`) that maps the unique [pool id](01_state.md#nextpoolid) to their corresponding `stableswap.Pool` object. Each `stableswap.Pool` object represents the specific state of a StableSwap pool.
```go
const StableSwapPools = []byte("stableswap_pools")
```

It is updated by the following messages:
- [`noble.swap.stableswap.v1.MsgCreatePool`](./02_messages_stableswap.md#create-pool)
- [`noble.swap.stableswap.v1.MsgUpdatePool`](./02_messages_stableswap.md#update-pool)

---

## BondedPositions

The `BondedPositions` field is a collection (`collections.IndexedMap`) that maps the bonded liquidity of each user by a triple key `<pool_id, address, timestamp>` to the [`Position`](01_types.md#position) object, which stores the necessary state for a bonded liquidity position.
```go
const Positions = []byte("stableswap_bonded_positions")
```

It is updated by the following messages:
- [`noble.swap.stableswap.v1.MsgAddLiquidity`](./02_messages.md#addliquidity)
- [`noble.swap.stableswap.v1.MsgRemoveLiquidity`](./02_messages.md#removeliquidity)

## UnbondingPositions

The `UnbondingPositions` field is a collection (`collections.IndexedMap`) that maps the current liquidity unbonding positions of each user by a triple key `<end_time, address, pool_id>` to the [`UnbondingPosition`](01_types.md#unbondingposition) object, which stores the necessary state for an unbonding position.
```go
const UnbondingPositions = []byte("stableswap_unbonding_positions")
```

It is updated by the following messages:
- [`noble.swap.stableswap.v1.MsgRemoveLiquidity`](./02_messages.md#removeliquidity)


## StableSwapPoolTotalUnbondingShares

The `StableSwapPoolTotalUnbondingShares` field is a collection (`collections.Map`) that maps each unique [pool id](01_state.md#nextpoolid) to a `math.LegacyDec` value representing the total unbonding liquidity shares for that specific StableSwap pool.
```go
const StableSwapPoolTotalUnbondingShares = []byte("stableswap_pool_total_unbonding_shares")
```

It is updated by the following messages:
- [`noble.swap.v1.MsgRemoveLiquidity`](./02_messages.md#remove-liquidity)

---

## StableSwapUsersTotalBondedShares

The `StableSwapUsersTotalBondedShares` field is a collection (`collections.Map`) that maps a pair `<pool_id, address>` to a `math.LegacyDec` value representing the total bonded liquidity shares held by a specific user in a specific StableSwap pool.
```go
const StableSwapUsersTotalBondedShares = []byte("stableswap_users_total_bonded_shares")
```

It is updated by the following messages:
- [`noble.swap.v1.MsgAddLiquidity`](./02_messages.md#add-liquidity)
- [`noble.swap.v1.MsgRemoveLiquidity`](./02_messages.md#remove-liquidity)

---

## StableSwapUsersTotalUnbondingShares

The `StableSwapUsersTotalUnbondingShares` field is a collection (`collections.Map`) that maps a pair `<pool_id, address>` to a `math.LegacyDec` value representing the total unbonding liquidity shares held by a specific user in a specific StableSwap pool.
```go
const StableSwapUsersTotalUnbondingShares = []byte("stableswap_users_total_unbonding_shares")
```

It is updated by the following messages:
- [`noble.swap.v1.MsgRemoveLiquidity`](./02_messages.md#remove-liquidity)

