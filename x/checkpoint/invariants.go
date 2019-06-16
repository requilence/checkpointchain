package checkpoint

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/requilence/checkpointchain/x/checkpoint/types"
	"github.com/requilence/checkpointchain/x/checkpoint/watchers/ethereum"
)

// register all staking invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "tx-proof",
		CheckpointsProofInvariants(k))
}

// AllInvariants runs all invariants of the checkpoint proof module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) error {
		err := CheckpointsProofInvariants(k)(ctx)
		if err != nil {
			return err
		}

		return nil
	}
}

func CheckpointsProofInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) error {
		iterator := k.GetCheckpointsIterator(ctx)
		for ; iterator.Valid(); iterator.Next() {
			blockNumber := binaryToNumber(iterator.Key())
			blockHash := k.GetCheckPointAtIndex(ctx, blockNumber)

			actualHash, err := ethereum.GetBlockHash(int64(blockNumber))
			if err != nil {
				return err
			}

			txBlockHash := "0x" + hex.EncodeToString(blockHash)
			if actualHash != txBlockHash {
				return fmt.Errorf("wrong hash for ethereum block %d: real %s, in tx %s", blockNumber, actualHash, txBlockHash)
			}
		}
		return nil
	}
}
