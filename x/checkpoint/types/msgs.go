package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgSetName defines a SetName message
type MsgAddCheckpoint struct {
	SectionIndex uint32         `json:"sectionIndex"`
	SectionHead  string         `json:"sectionHead"`
	CHTRoot      string         `json:"chtRoot"`
	BloomRoot    string         `json:"bloomRoot"`
	Creator      sdk.AccAddress `json:"creator"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgAddCheckpoint(sectionIndex uint32, sectionHead string, chtRoot string, bloomRoot string, creator sdk.AccAddress) MsgAddCheckpoint {
	return MsgAddCheckpoint{
		SectionIndex: sectionIndex,
		SectionHead:  sectionHead,
		CHTRoot:      chtRoot,
		BloomRoot:    bloomRoot,
		Creator:      creator,
	}
}

// Route should return the name of the module
func (msg MsgAddCheckpoint) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAddCheckpoint) Type() string { return "add_checkpoint" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAddCheckpoint) ValidateBasic() sdk.Error {
	if msg.SectionIndex == 0 {
		return sdk.ErrUnknownRequest("Incorrect values for block number and block hash")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAddCheckpoint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAddCheckpoint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}
