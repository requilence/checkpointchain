package types

import (
	"fmt"
	"strings"
)

type Checkpoint struct {
	Index uint32 `json:"index"`
	Data   []byte `json:"data"`
}

// implement fmt.Stringer
func (w Checkpoint) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Index: %s
Data: %s`, w.Index, w.Data))
}

func NewCheckpoint() Checkpoint {
	return Checkpoint{}
}
