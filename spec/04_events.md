## PoolsPaused

This event is emitted whenever a set of pools are paused in the Swap module.

```json
{
  "type": "noble.swap.v1.PoolsPaused",
  "attributes": [
    {
      "key": "pool_ids",
      "value": ["1", "2", "3"]
    }
  ]
}
```

This event is emitted by the following transactions:

- [`noble.swap.v1.MsgPauseByPoolIds`](./02_messages.md#pause-by-pool-ids)
- [`noble.swap.v1.MsgPauseByAlgorithm`](./02_messages.md#pause-by-algorithm)

## PoolsUnpaused

This event is emitted whenever a set of pools are unpaused in the Swap module.

```json
{
  "type": "noble.swap.v1.PoolsUnpaused",
  "attributes": [
    {
      "key": "pool_ids",
      "value": ["4", "5"]
    }
  ]
}
```

This event is emitted by the following transactions:

- [`noble.swap.v1.MsgUnpauseByPoolIds`](./02_messages.md#unpause-by-pool-ids)
- [`noble.swap.v1.MsgUnpauseByAlgorithm`](./02_messages.md#unpause-by-algorithm)

## Swap

This event is emitted whenever a swap operation is completed successfully.

```json
{
  "type": "noble.swap.v1.Swap",
  "attributes": [
    {
      "key": "signer",
      "value": "noble1signer"
    },
    {
      "key": "input",
      "value": "1000uusdc"
    },
    {
      "key": "output",
      "value": "998uusdn"
    },
    { 
      "key": "routes",
      "value": "[{'pool_id':'0','denom_to':'uusdn'}]"
    },
    {
      "key": "fees",
      "value": "2uusdc"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`noble.swap.v1.MsgSwap`](./02_messages.md#swap)

## WithdrawnProtocolFees

This event is emitted whenever protocol fees are withdrawn from a pool.

```json
{
  "type": "noble.swap.v1.WithdrawnProtocolFees",
  "attributes": [
    {
      "key": "to",
      "value": "noble1signer"
    },
    {
      "key": "rewards",
      "value": "500uusdn, 300uusdc"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`noble.swap.v1.MsgWithdrawProtocolFees`](./02_messages.md#withdraw-protocol-fees)

## WithdrawnRewards

This event is emitted whenever a user withdraws their rewards from a pool.

```json
{
  "type": "noble.swap.v1.WithdrawnRewards",
  "attributes": [
    {
      "key": "signer",
      "value": "noble1signer"
    },
    {
      "key": "rewards",
      "value": "200uusdn"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`noble.swap.v1.MsgWithdrawRewards`](./02_messages.md#withdraw-rewards)
