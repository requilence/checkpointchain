package types

import (
	"encoding/json"
	"fmt"
)

// Query Result Payload for a resolve query
type QueryResCheckpoint struct {
	BlockNumber uint32 `json:"blockNumber"`
	BlockHash   string `json:"blockHash"`
}

// implement fmt.Stringer
func (r QueryResCheckpoint) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// Query Result Payload for a resolve query
type QueryResLastCheckpoint struct {
	BlockHash string `json:"blockNumber"`
}

// implement fmt.Stringer
func (r QueryResLastCheckpoint) String() string {
	return fmt.Sprintf("%x", r.BlockHash)
}
