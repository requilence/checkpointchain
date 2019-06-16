package checkpoint

import (
	"github.com/requilence/checkpointchain/x/checkpoint/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewMsgAddCheckpoint = types.NewMsgAddCheckpoint
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
)

type (
	MsgAddCheckpoint       = types.MsgAddCheckpoint
	QueryResCheckpoint     = types.QueryResCheckpoint
	QueryResLastCheckpoint = types.QueryResLastCheckpoint
	Checkpoint             = types.Checkpoint
)
