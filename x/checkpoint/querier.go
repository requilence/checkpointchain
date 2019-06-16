package checkpoint

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the nameservice Querier
const (
	QueryBeforeCheckpoint = "beforeblock"
	QueryLastCheckpoint   = "lastcheckpoint"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryBeforeCheckpoint:
			return queryCheckpointBeforeBLock(ctx, path[1:], req, keeper)
			/*	case QueryLastCheckpoint:
					return queryl(ctx, path[1:], req, keeper)
				case QueryNames:
					return queryNames(ctx, req, keeper)*/
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("unknown checkpoint query endpoint: %s", path[0]))
		}
	}
}

// nolint: unparam
func queryCheckpointBeforeBLock(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	number, err := strconv.Atoi(path[0])
	if err != nil {
		return []byte{}, sdk.ErrUnknownRequest("wrong number format")
	}

	if number > math.MaxUint32 {
		return []byte{}, sdk.ErrUnknownRequest("number is more than uint32")
	}

	blockNumber, value := keeper.GetCheckPointBeforeIndex(ctx, uint32(number))
	if len(value) == 0 {
		return []byte{}, sdk.ErrUnknownRequest("could not resolve find checkpoint")
	}

	blockHashHex := hex.EncodeToString(value)
	res, err := codec.MarshalJSONIndent(keeper.cdc, QueryResCheckpoint{uint32(blockNumber), blockHashHex})
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryCheckpointLatest(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	number, err := strconv.Atoi(path[0])
	if err != nil {
		return []byte{}, sdk.ErrUnknownRequest("wrong number format")
	}

	if number > math.MaxUint32 {
		return []byte{}, sdk.ErrUnknownRequest("number is more than uint32")
	}

	index := keeper.GetMaximumCheckpointIndex(ctx)
	if index == 0 {
		return []byte{}, sdk.ErrUnknownRequest("could not resolve find checkpoint")
	}

	checkpoint := keeper.GetCheckPointAtIndex(ctx, index)

	blockHashHex := hex.EncodeToString(checkpoint)
	res, err := codec.MarshalJSONIndent(keeper.cdc, QueryResCheckpoint{uint32(blockNumber), blockHashHex})
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}


/*
// nolint: unparam
func queryLastestCheckpoint(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {

	value := keeper.GetCheckPointAtIndex(ctx, uint32(number))
	if len(value) == 0 {
		return []byte{}, sdk.ErrUnknownRequest("could not resolve find checkpoint")
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, QueryResResolve{value})
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
*/
