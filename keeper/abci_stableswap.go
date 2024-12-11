package keeper

import "context"

func (k *Keeper) StableSwapBeginBlocker(ctx context.Context) (executed bool) {
	headerInfo := k.headerService.GetHeaderInfo(ctx)
	if headerInfo.Height%k.stableswapConfig.UnbondingBlockDelta != 0 {
		return false
	}
	k.Stableswap.Logger().Info("processing unbondings epoch")

	// iterate over the stableswap pools and process the unbonding positions whose unbonding period has ended.
	for poolId := range k.Stableswap.GetPools(ctx) {
		// Retrieve the Pool Controller.
		controller, err := GetStableSwapController(ctx, k, poolId)
		if err != nil {
			k.Stableswap.Logger().Error("processing unbondings epoch:", err.Error())
			continue
		}

		if controller.IsPaused() {
			continue
		}

		// Process the pending Unbondings and return the updated Pool.
		err = controller.ProcessUnbondings(ctx, headerInfo.Time)
		if err != nil {
			k.Stableswap.Logger().Error("processing unbondings epoch:", err.Error())
			continue
		}
	}
	return true
}
