package checkpoint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "checkpoint" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgAddCheckpoint:
			return handleMsgAddCheckpoint(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized checkpoint Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgAddCheckpoint(ctx sdk.Context, keeper Keeper, msg MsgAddCheckpoint) sdk.Result {
	keeper.SetCheckpoint(ctx, msg.BlockNumber, msg.BlockHash)
	return sdk.Result{}
}
