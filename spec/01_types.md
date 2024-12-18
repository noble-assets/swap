## State Types

### Algorithm
`swap.v1.Algorithm`

Defines the algorithm used by a pool for token exchange.

```json
{
  "UNSPECIFIED": 0,
  "STABLESWAP": 1
}
```

**Description**
- `UNSPECIFIED` — Default value, indicates an undefined algorithm.
- `STABLESWAP` — Algorithm for stable swaps with minimal slippage.

---

### Pool
`swap.v1.Pool`

Represents a liquidity pool within the state.

```json
{
  "id": 1,
  "address": "noble1pooladdress",
  "algorithm": "STABLESWAP",
  "pair": "uusdc",
  "total_shares": "1000.0",
  "timestamp": "2024-11-18T00:00:00Z"
}
```

**Fields**
- `id` — Unique identifier for the pool.
- `address` — Cosmos address associated with the pool.
- `algorithm` — Algorithm used by the pool for swaps (`Algorithm` enum).
- `pair` — Token pair associated with the pool.
- `total_shares` — Total shares representing liquidity in the pool.
- `timestamp` — Creation or last update time of the pool.

---

### BondedPosition
`swap.v1.BondedPosition`

Represents a liquidity bonded position in a pool.

```json
{
  "balance": "500.0",
  "timestamp": "2024-11-18T00:00:00Z",
  "rewards_period_start": "2024-11-01T00:00:00Z"
}
```

**Fields**
- `balance` — Number of shares in the position.
- `timestamp` — Creation or last update time of the position.
- `rewards_period_start` — Start time for calculating rewards for the position.

---

### UnbondingPosition
`swap.v1.UnbondingPosition`

Represents an unbonding position in a pool.

```json
{
  "shares": "50.0",
  "amount": [
    { "denom": "uusdc", "amount": "50000" }
  ],
  "end_time": "2024-12-01T00:00:00Z"
}
```

**Fields**
- `shares` — Number of unbonding shares.
- `amount` — Tokens being unbonded.
- `end_time` — Time when unbonding completes.

---

### Route
`swap.v1.Route`

Represents a step in a swap route.

```json
{
  "pool_id": 1,
  "denom_to": "uusdn"
}
```

**Fields**
- `pool_id` — Identifier of the pool used for the swap step.
- `denom_to` — Target token denomination after the swap.

---