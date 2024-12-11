package keeper

import (
	"context"
)

func (k *Keeper) BeginBlocker(ctx context.Context) error {
	// Process the StableSwap BeginBlocker.
	k.StableSwapBeginBlocker(ctx)
	return nil
}
