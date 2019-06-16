package checkpoint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/requilence/checkpointchain/x/checkpoint/watchers/ethereum"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	Checkpoints []ethereum.SectionInfo `json:"checkpoints"`
}

func NewGenesisState(checkpoints []Checkpoint) GenesisState {
	return GenesisState{Checkpoints: checkpoints}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.Checkpoints {
		if record.BlockHash == nil {
			return fmt.Errorf("Invalid WhoisRecord: BlockHash: %s. Error: Missing", record.BlockHash)
		}
		if record.BlockNumber == 0 {
			return fmt.Errorf("Invalid WhoisRecord: BlockNumber: %d. Error: Missing", record.BlockNumber)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Checkpoints: []Checkpoint{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.Checkpoints {
		keeper.SetCheckpoint(ctx, record)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []Checkpoint
	iterator := k.GetCheckpointsIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		blockNumber := binaryToNumber(iterator.Key())
		blockHash := k.GetCheckPointAtIndex(ctx, blockNumber)
		records = append(records, Checkpoint{BlockNumber: blockNumber, BlockHash: blockHash})
	}
	return GenesisState{Checkpoints: records}
}
