## Messages

### Create Pool
`noble.swap.stableswap.v1.MsgCreatePool`

Creates a new StableSwap liquidity pool with specific AMM parameters.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.swap.stableswap.v1.MsgCreatePool",
        "signer": "noble1signer",
        "pair": "uusdc",
        "protocol_fee_percentage": 10,
        "rewards_fee": 5,
        "max_fee": 20,
        "initial_a": 1000,
        "future_a": 2000,
        "future_a_time": 1750000000,
        "rate_multipliers": [
          {
            "denom": "uusdc",
            "amount": 1000000000
          },
          {
            "denom": "uusdc",
            "amount": 1000000000
          }
        ]
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

**Arguments**
- `signer` — Address of the account creating the pool.
- `pair` — The token pair for the pool (e.g., `uusd`). The default main pair is `uusdn`
- `rewards_fee` — Rewards fee as a percentage.
- `protocol_fee_percentage` — Protocol fee as a percentage off from the `rewards_fee`.
- `max_fee` — Maximum allowable fee.
- `initial_a` — Initial amplification coefficient.
- `future_a` — Future amplification coefficient.
- `future_a_time` — Timestamp for the future amplification coefficient to take effect.
- `rate_multipliers` — Rate multipliers for the tokens in the pool.

**State Changes**
- Creates a new StableSwap liquidity pool.
- Initializes the Pool with specified parameters.

### Update Pool
`noble.swap.stableswap.v1.MsgUpdatePool`

Updates an existing StableSwap liquidity pool with specific AMM parameters.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.swap.stableswap.v1.MsgUpdatePool",
        "signer": "noble1signer",
        "id": 1,
        "protocol_fee_percentage": 10,
        "rewards_fee": 5,
        "max_fee": 20,
        "future_a": 3000,
        "future_a_time": 1760000000,
        "rate_multipliers": [
          {
            "denom": "uusdc",
            "amount": 1000000000
          },
          {
            "denom": "uusdc",
            "amount": 1000000000
          }
        ]
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

**Arguments**
- `signer` — Address of the account creating the pool.
- `pair` — The token pair for the pool (e.g., `uusd`). The default main pair is `uusdn`
- `rewards_fee` — Rewards fee as a percentage.
- `protocol_fee_percentage` — Protocol fee as a percentage off from the `rewards_fee`.
- `max_fee` — Maximum allowable fee.
- `initial_a` — Initial amplification coefficient.
- `future_a` — Future amplification coefficient.
- `future_a_time` — Timestamp for the future amplification coefficient to take effect.
- `rate_multipliers` — Rate multipliers for the tokens in the pool.

**State Changes**
- Creates a new StableSwap liquidity pool.
- Initializes the Pool with specified parameters.

### Add Liquidity
`noble.swap.v1.MsgAddLiquidity`

Provides tokens to a specified liquidity pool, increasing its reserves and earning liquidity shares for the provider.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.swap.v1.MsgAddLiquidity",
        "signer": "noble1signer",
        "pool_id": "1",
        "amount": [
          { "denom": "uusdc", "amount": "1000000" },
          { "denom": "uusdn", "amount": "1000000" }
        ]
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

**Arguments**
- `signer` — The Noble address of the account providing liquidity.
- `pool_id` — ID of the pool.
- `amount` — Amount of tokens to add.

**Requirements**
- `amount` — The base token (USDN) amount must be at least 1 unit (1000000).

**State Changes**
- Updates `Pool` reserves.
- Mints liquidity shares to the signer.
- Creates a new user `Position`.
- Updates `StableSwapUsersTotalBondedShares`

---

### Remove Liquidity
`noble.swap.v1.MsgRemoveLiquidity`

Removes a user’s share in a pool by withdrawing a specified percentage of the provided liquidity and adjusting the pool's reserves accordingly. This message initiate automatically also a [MsgWithdrawRewards](02_messages.md#withdraw-rewards)

```json
{
  "body": {
    "messages": [
      {
        "@type": "/noble.swap.v1.MsgRemoveLiquidity",
        "signer": "noble1signer",
        "pool_id": "1",
        "percentage": "0.5"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

**Arguments**
- `signer` — Address of the account removing liquidity.
- `pool_id` — ID of the pool.
- `percentage` — Percentage of liquidity to remove.

**State Changes**
- Adjusts pool reserves.
- Burns liquidity shares from the signer.
- Updates the user `Position`.
- Updates `StableSwapUsersTotalBondedShares`, `StableSwapUsersTotalUnbondingShares`, `StableSwapPoolTotalUnbondingShares`

---