## Messages

### Swap
`swap.v1.MsgSwap`

Executes a token exchange between pools based on a defined [route](01_types.md#route) and ensures a minimum output token amount is received.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/swap.v1.MsgSwap",
        "signer": "noble1signer",
        "amount": {
          "denom": "uusdc",
          "amount": "1000000"
        },
        "routes": [
          { "pool_id": "1", "denom_to": "uusdn" },
          { "pool_id": "2", "denom_to": "uusde" }
        ],
        "min": {
          "denom": "uusde",
          "amount": "950000"
        }
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
- `signer` — The Noble address of the account performing the swap.
- `Amount` — Input token.
- `routes` — Path of pools for the swap.
- `min` — Minimum output token wanted.

**Requirements**
- Signer must have sufficient input tokens.
- Pools 

**State Changes**
- Updates the pools liquidity and user balances.

---

### Withdraw Protocol Fees
`swap.v1.MsgWithdrawProtocolFees`

Transfers accumulated protocol fees from the system to a specified destination address initiated by an authorized account.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/swap.v1.MsgWithdrawProtocolFees",
        "signer": "noble1signer",
        "to": "noble1destination"
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
- `signer` — Noble address of the authority account initiating the withdrawal.
- `to` — Destination Noble address.

**State Changes**
- Transfers protocol fees to the specified address.

---

### Withdraw Rewards
`swap.v1.MsgWithdrawRewards`

Collects accumulated rewards earned by a user, transferring them to the user’s account balance.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/swap.v1.MsgWithdrawRewards",
        "signer": "noble1signer"
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
- `signer` — Address of the Noble account withdrawing rewards.

**State Changes**
- Transfers rewards to the signer’s account.

---

### Pause By Algorithm
`swap.v1.MsgPauseByAlgorithm`

Temporarily halts all operations within the pools with the given algorithm, preventing further transactions and interactions (swaps, liquidity, unbonding, etc.).

```json
{
  "body": {
    "messages": [
      {
        "@type": "/swap.v1.MsgPauseByAlgorithm",
        "signer": "noble1signer",
        "algorithm": "STABLESWAP"
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
- `signer` — Noble address of the authority account initiating the pause.
- `algorithm` - Pool [algorithm](01_types.md#algorithm) to pause.

**State Changes**
- Sets the pools with the given algorithm state to paused.

---

### Pause By Pool Ids
`swap.v1.MsgPauseByPoolIds`

Temporarily halts all operations of the provided pool ids, preventing further transactions and interactions (swaps, liquidity, unbonding, etc.).

```json
{
  "body": {
    "messages": [
      {
        "@type": "/swap.v1.MsgPauseByPoolIds",
        "signer": "noble1signer",
        "pool_ids": ["0", "2"]
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
- `signer` — Noble address of the authority account initiating the pause.
- `pool_ids` - Ids of the pools to pause.

**State Changes**
- Sets the pools with the given pool ids state to paused.

---

### Unpause By Algorithm
`swap.v1.MsgUnpauseByAlgorithm`

Unpauses all operations within the pools with the given algorithm, restoring transactions and interactions (swaps, liquidity, unbonding, etc.).

```json
{
  "body": {
    "messages": [
      {
        "@type": "/swap.v1.MsgUnpauseByAlgorithm",
        "signer": "noble1signer",
        "algorithm": "STABLESWAP"
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
- `signer` — Noble address of the authority account initiating to unpause.
- `algorithm` - Pool [algorithm](01_types.md#algorithm) to unpause.

**State Changes**
- Sets the pools with the given algorithm state to unpaused.

---

### Unpause By Pool Ids
`swap.v1.MsgUnpauseByPoolIds`
Unpauses all operations of the provided pool ids, restoring transactions and interactions (swaps, liquidity, unbonding, etc.).

```json
{
  "body": {
    "messages": [
      {
        "@type": "/swap.v1.MsgUnpauseByPoolIds",
        "signer": "noble1signer",
        "pool_ids": ["0", "2"]
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
- `signer` — Noble address of the authority account initiating the unpause.
- `pool_ids` - Ids of the pools to unpause.

**State Changes**
- Sets the pools with the given pool ids state to unpaused.

---
