package keeper

import (
	"context"
	"fmt"
)

// StableSwapBeginBlocker handles the periodic updates for the StableSwap submodule,
// specifically managing and processing unbonding positions that have reached the end
// of their unbonding period. This function is invoked at the beginning of each block.
//
// This function checks whether the current block height aligns with the unbonding block
// interval configured for StableSwap. If not, the function exits early. When the condition
// is met, it processes all eligible unbonding positions across active stableswap pools.
func (k *Keeper) StableSwapBeginBlocker(ctx context.Context) (executed bool) {
	headerInfo := k.headerService.GetHeaderInfo(ctx)
	if headerInfo.Height%k.stableswapConfig.UnbondingBlockDelta != 0 {
		return false
	}
	k.Stableswap.Logger().Info("processing unbondings epoch")

	// Get all the UnbondingPositions expiring at the current block time.
	unbondingPositions := k.Stableswap.GetUnbondingPositionsUntil(ctx, headerInfo.Time.Unix())

	// If no unbonding entries, early exit.
	if len(unbondingPositions) == 0 {
		return true
	}

	// Get all the StableSwap Controllers.
	controllers := GetStableSwapControllers(ctx, k)

	// Iterate through the UnbondingPositions and process eligible unbondings.
	for _, position := range unbondingPositions {
		// Get the specific Pool Controller
		controller, exists := controllers[position.PoolId]
		if !exists {
			k.Stableswap.Logger().Error(fmt.Sprintf("Pool %d Controller does not exists", position.PoolId))
			continue
		}

		// Skip processing if the pool is paused.
		if controller.IsPaused() {
			continue
		}

		// Process the pending unbondings.
		err := controller.ProcessUnbondings(ctx, headerInfo.Time)
		if err != nil {
			k.Stableswap.Logger().Error(fmt.Sprintf("failed to process Pool %d unbondings: %s", position.PoolId, err.Error()))
			continue
		}
	}
	return true
}
