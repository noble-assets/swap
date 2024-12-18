## Queries

---

### PositionsByProvider
`types.QueryPositionsByProvider`

Fetches all the rewards, bonded,unbonding positions for a specific provider.

```json
{
  "bonded_positions": [
    {
      "pool_id": 1,
      "shares": "100",
      "timestamp": "2024-11-18T00:00:00Z"
    }
  ],
  "unbonding_positions": [
    {
      "pool_id": 1,
      "unbonding_shares": "50",
      "end_time": "2024-12-01T00:00:00Z"
    }
  ],
  "rewards": [
    {
      "pool_id": 1,
      "amount": [
        {
          "amount": 1000,
          "denom": "uusdc"
        }
       ]
    }
  ]
}
```

**Arguments**
- `provider` — The provider's address.

**Requirements**
- Provider must exist.

---

### RewardsByProvider
`types.QueryRewardsByProvider`

Fetches all the rewards for a specific provider.

```json
{
  "rewards": [
    {
      "pool_id": 1,
      "amount": [
        {
          "amount": 1000,
          "denom": "uusdc"
        }
       ]
    }
  ]
}
```

**Arguments**
- `provider` — The provider's address.

**Requirements**
- Provider must exist.

---

### BondedPositionsByProvider
`types.QueryBondedPositionsByProvider`

Fetches all the bonded positions for a specific provider.

```json
{
  "bonded_positions": [
    {
      "pool_id": 1,
      "shares": "100",
      "timestamp": "2024-11-18T00:00:00Z"
    }
  ]
}
```

**Arguments**
- `provider` — The provider's address.

**Requirements**
- Provider must exist.

---

### UnbondingPositionsByProvider
`types.QueryUnbondingPositionsByProvider`

Fetches all unbonding positions for a specific provider.

```json
{
  "unbonding_positions": [
    {
      "pool_id": 1,
      "unbonding_shares": "50",
      "end_time": "2024-12-01T00:00:00Z"
    }
  ]
}
```

**Arguments**
- `provider` — The provider's address.

**Requirements**
- Provider must exist.

---