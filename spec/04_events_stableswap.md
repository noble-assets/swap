
## PoolCreated

This event is emitted whenever a new stable swap pool is created.

```json
{
  "type": "noble.swap.stableswap.v1.PoolCreated",
  "attributes": [
    {
      "key": "pool_id",
      "value": "12"
    },
    {
      "key": "algorithm",
      "value": "STABLESWAP"
    },
    {
      "key": "pair",
      "value": "uusdn"
    },
    {
      "key": "protocol_fee_percentage",
      "value": "0.1"
    },
    ...
  ]
}
```

This event is emitted by the following transactions:

- [`noble.swap.stableswap.v1.MsgCreatePool`](./02_messages.md#create-pool)

## PoolUpdated

This event is emitted whenever an existing pool's parameters are updated.

```json
{
  "type": "noble.swap.stableswap.v1.PoolUpdated",
  "attributes": [
    {
      "key": "pool_id",
      "value": "12"
    },
    {
      "key": "protocol_fee_percentage",
      "value": "0.2"
    },
    ...
  ]
}
```

This event is emitted by the following transactions:

- [`noble.swap.stableswap.v1.MsgUpdatePool`](./02_messages.md#update-pool)

## LiquidityAdded

This event is emitted whenever liquidity is added to a pool.

```json
{
  "type": "noble.swap.stableswap.v1.LiquidityAdded",
  "attributes": [
    {
      "key": "provider",
      "value": "noble1signer"
    },
    {
      "key": "pool_id",
      "value": "15"
    },
    {
      "key": "amount",
      "value": "1000uusdc, 1000uusdy"
    },
    {
      "key": "shares",
      "value": "200"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`noble.swap.stableswap.v1.MsgAddLiquidity`](./02_messages.md#add-liquidity)

## LiquidityRemoved

This event is emitted whenever liquidity is removed from a pool.

```json
{
  "type": "noble.swap.stableswap.v1.LiquidityRemoved",
  "attributes": [
    {
      "key": "provider",
      "value": "noble1signer"
    },
    {
      "key": "pool_id",
      "value": "15"
    },
    {
      "key": "amount",
      "value": "500uusdc, 500uusdn"
    },
    {
      "key": "shares",
      "value": "100"
    },
    {
      "key": "unlock_time",
      "value": "2024-12-01T00:00:00Z"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`noble.swap.stableswap.v1.MsgRemoveLiquidity`](./02_messages.md#remove-liquidity)

