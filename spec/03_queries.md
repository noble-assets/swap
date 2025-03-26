## Queries

### Pool
`types.QueryPool`

Retrieves detailed information about a specific liquidity pool by its ID.

```json
{
  "pool":{
    "id":"1",
    "address":"noble1pooladdress",
    "algorithm":"STABLESWAP",
    "pair": "uusdc",
    "details":{
      "@type":"/noble.swap.stableswap.v1.Pool",
      "protocol_fee_percentage":"2",
      "rewards_fee":"400000",
      "initial_a":"100",
      "future_a":"100",
      "initial_a_time":"1730000000",
      "future_a_time":"1900000000",
      "rate_multipliers":[
        {
          "denom": "uusdc",
          "amount": "1000000000000000000"
        },
        {
          "denom":"uusdn",
          "amount":"1000000000000000000"
        }
      ]
    },
    "liquidity":[
      {
        "denom": "uusdc",
        "amount": "999000040011"
      },
      {
        "denom": "uusdn",
        "amount": "1000999960002"
      }
    ],
    "protocol_fees":[
      {
        "denom": "uusdn",
        "amount": "799"
      }
    ],
    "reward_fees":[
      {
        "denom": "uusdn",
        "amount": "39199"
      }
    ]
  }
}
```

**Arguments**
- `id` — The unique identifier of the liquidity pool.

**Requirements**
- Pool ID must exist.

---

### Pools
`types.QueryPools`

Fetches information about all existing liquidity pools.

```json
{
  "pools": [
    {
      "id":"1",
      "address":"noble1pooladdress",
      "algorithm":"STABLESWAP",
      "pair": "uusdc",
      "details":{
        "@type":"/noble.swap.stableswap.v1.Pool",
        "protocol_fee_percentage":"2",
        "rewards_fee":"400000",
        "initial_a":"100",
        "future_a":"100",
        "initial_a_time":"1730000000",
        "future_a_time":"1900000000",
        "rate_multipliers":[
          {
            "denom": "uusdc",
            "amount": "1000000000000000000"
          },
          {
            "denom":"uusdn",
            "amount":"1000000000000000000"
          }
        ]
      },
      "liquidity":[
        {
          "denom": "uusdc",
          "amount": "999000040011"
        },
        {
          "denom": "uusdn",
          "amount": "1000999960002"
        }
      ],
      "protocol_fees":[
        {
          "denom": "uusdn",
          "amount": "799"
        }
      ],
      "reward_fees":[
        {
          "denom": "uusdn",
          "amount": "39199"
        }
      ]
    },
    {
      ...
    }
  ]
}
```

**Arguments**
- None.

---

### Paused
`types.QueryPaused`

Returns the ids of the paused pools.

```json
{
  "paused": ["0", "3"]
}
```

**Arguments**
- None.

**Requirements**
- None.

---

### Rates
`types.QueryRates`

Fetches all rates for pools that match a specific algorithm.

```json
{
  "rates": [
    {
      "denom": "uusdc",
      "vs": "uusdn",
      "price": "1.002001921832733839",
      "algorithm": "STABLESWAP"
    },
    {
      "denom": "uusdn",
      "vs": "uusdc",
      "price": "0.998002077851235874",
      "algorithm": "STABLESWAP"
    }
  ]
}
```

**Arguments**
- `algorithm` — (Optional) The algorithm filter.

**Requirements**
- None.

---

### Rate
`types.QueryRate`

Fetches rates for a specific token.

```json
{
  "rates": [
    {
      "denom": "uusdc",
      "vs": "uusdn",
      "price": "1.002001921832733839",
      "algorithm": "STABLESWAP"
    }
  ]
}
```

**Arguments**
- `denom` — The token denomination.

**Requirements**
- Token denomination must be valid.

---

### Simulate Swap
`types.QuerySimulateSwap`

Simulate the expected output and associated fees for a token swap and [route](01_types.md#route) without executing the transaction or requiring a valid account or balance.

```json
{
  "result": {
    "denom": "uusdc",
    "amount": "999714654"
  },
  "swaps": [
    {
      "in": {
        "denom": "uusdn",
        "amount": "1000000000"
      },
      "out": {
        "denom": "uusdc",
        "amount": "999714654"
      },
      "fees": [
        {
          "denom": "uusdn",
          "amount": "249991"
        }
      ]
    }
  ]
}
```

**Arguments**
- `amount` — Input token.
- `routes` — Path of pools for the swap.
- `min` — Minimum output token wanted.

**Requirements**
- Signer must have sufficient input tokens.
- Token denominations must be correctly specified, and routing paths must be valid.
- The minimum output token denom must match the output token denom in the final route.

---