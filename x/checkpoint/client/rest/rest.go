package rest

import (
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/requilence/checkpointchain/x/checkpoint/types"

	"github.com/gorilla/mux"
)

const (
	restNumber = "number"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	fmt.Printf("RegisterRoutes storeName = %s\n", storeName)
	r.HandleFunc(fmt.Sprintf("/%s", storeName), addCheckpointHandler(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/before/{%s}", storeName, restNumber), checkpointBeforeHandler(cdc, cliCtx, storeName)).Methods("GET")
}

type addCheckpointReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	BlockNumber string `json:"name"`
	BlockHash   string `json:"amount"`
}

func addCheckpointHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req addCheckpointReq

		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		blockNumber, err := strconv.Atoi(req.BlockNumber)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("block_id is not a number: %s", err.Error()))
			return
		}

		if blockNumber > math.MaxUint32 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "block_id is more than uint32")
			return
		}

		if !strings.HasPrefix(req.BlockHash, "0x") {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "block_hash should start with 0x")
			return
		}

		blockHashHex := req.BlockHash[2:]
		if len(blockHashHex) != 64 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "block_hash should be 64 symbols after 0x")
			return
		}

		blockHash, err := hex.DecodeString(blockHashHex)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("error decoding block_hex: %s", err.Error()))
			return
		}

		msg := types.NewMsgAddCheckpoint(uint32(blockNumber), blockHash, cliCtx.GetFromAddress())
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("error validating checkpoint: %s", err.Error()))
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func checkpointBeforeHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restNumber]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/beforeblock/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
