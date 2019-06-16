package cli

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/requilence/checkpointchain/x/checkpoint/watchers/ethereum"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/requilence/checkpointchain/x/checkpoint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	checkpointTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Checkpoint add subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       utils.ValidateCmd,
	}

	checkpointTxCmd.AddCommand(client.PostCommands(
		GetCmdSetCheckpoint(cdc),
		GetCmdWatchEth(cdc),
	)...)

	return checkpointTxCmd
}

// GetCmdSetCheckpoint is the CLI command for sending a BuyName transaction
func GetCmdSetCheckpoint(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set [block_id] [block_hash]",
		Short: "add the new checkpoint",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			blockNumber, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("block_id is not a number: %s", err.Error())
			}

			if blockNumber > math.MaxUint32 {
				return fmt.Errorf("block_id is more than uint32")
			}

			if !strings.HasPrefix(args[1], "0x") {
				return fmt.Errorf("block_hash should start with 0x")
			}

			blockHashHex := args[1][2:]
			if len(blockHashHex) != 64 {
				return fmt.Errorf("block_hash should be 64 symbols after 0x")
			}

			blockHash, err := hex.DecodeString(blockHashHex)
			if err != nil {
				return fmt.Errorf("error decoding block_hex: %s", err.Error())
			}

			msg := types.NewMsgAddCheckpoint(uint32(blockNumber), blockHash, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdWatchEth(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "watch_eth [divider]",
		Short: "watch ethereum",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			cliCtx.PrintResponse = true

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			fromName := cliCtx.GetFromName()
			passphrase, err := keys.GetPassphrase(fromName)
			if err != nil {
				return err
			}

			var ch = make(chan ethereum.SectionInfo)

			divider, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			err = ethereum.Watcher(divider, ch)
			if err != nil {
				return err
			}

			for {
				select {
				case checkpoint := <-ch:
					fmt.Printf("Got checkpoint section %d: %s\n", checkpoint.SectionIndex, checkpoint.SectionHead)

					msg := types.NewMsgAddCheckpoint(uint32(checkpoint.Number), b, cliCtx.GetFromAddress())
					err = msg.ValidateBasic()
					if err != nil {
						return err
					}

					err = SignAndBroadcastTxCLI(txBldr, passphrase, cliCtx, []sdk.Msg{msg})
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}
