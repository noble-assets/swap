package keeper

import (
	"context"
)

// BeginBlocker is called at the beginning of each block. It performs the necessary
// state transitions and processing required for each Swap Submodule (ex. StableSwap).
func (k *Keeper) BeginBlocker(ctx context.Context) error {
	// Process the StableSwap BeginBlocker.
	k.StableSwapBeginBlocker(ctx)

	return nil
}
