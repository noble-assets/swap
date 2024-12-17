# State

## Paused

The `paused` field is a (`collections.Map`) that maps the unique [pool_id](01_state.md#nextpoolid) to a`bool` used to represent the pausing state of the pool. When set to `true`, the pool is in a paused state, indicating that no actions (swap/liquidity/rewards/etc.) can be initiated or processed.
```go
const PausedPrefix = []byte("paused")
```

It is updated by the following messages:

- [`swap.v1.MsgPauseByAlgorithm`](./02_messages.md#pause-by-algorithm)
- [`swap.v1.MsgPauseByPoolIds`](./02_messages.md#pause-by-pool-ids)
- [`swap.v1.MsgUnpauseByAlgorithm`](./02_messages.md#unpause-by-algorithm)
- [`swap.v1.MsgUnpauseByPoolIds`](./02_messages.md#unpause-by-pool-ids)


## NextPoolId

The `next_pool_id` field is a collection sequence (`collections.Sequence`) of an `uint64` integer used to track the identifier that will be assigned to the next pool created in the Swaps module.
```go
const NextPoolIdPrefix = []byte("next_pool_id")
```

It is an auto-incrementing value that cannot be updated manually or through external messages.


## Pools

The `Pools` field is a collection (`collections.Map`) that maps the unique [pool_id](01_state.md#nextpoolid) to their corresponding generic [`Pool`](01_types.md#pool) object, which stores the shared state across all pools and algorithm types.
```go
const Pools = []byte("pools_generic")
```

It is updated by the following messages:
- [`swap.stableswap.v1.MsgCreatePool`](./02_messages_stableswap.md#create-pool)
