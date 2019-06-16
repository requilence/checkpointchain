package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/requilence/checkpointchain/x/checkpoint/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	nameserviceQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the checkpoint module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       utils.ValidateCmd,
	}
	nameserviceQueryCmd.AddCommand(client.GetCommands(
		GetCmdBeforeBlock(storeKey, cdc),
		GetCmdLastCheckpoint(storeKey, cdc),
	)...)
	return nameserviceQueryCmd
}

// GetCmdResolveName queries information about a name
func GetCmdBeforeBlock(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "before [block_id]",
		Short: "get the checkpoint for block <= block_id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			block := args[0]

			_, err := strconv.Atoi(block)
			if err != nil {
				return fmt.Errorf("block number should be int")
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/beforeblock/%s", queryRoute, block), nil)
			if err != nil {
				return err
			}

			var out types.QueryResCheckpoint
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdResolveName queries information about a name
func GetCmdLastCheckpoint(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "last",
		Short: "get the last known checkpoint",
		//Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			//for blockNumber:= start; blockNumber>=0; blockNumber-- {
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/lastcheckpoint", queryRoute), nil)
			if err != nil {
				//fmt.Printf("could not resolve name - %s \n", block)
				//continue
				return err
			}

			var out types.QueryResLastCheckpoint
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
			//	}

			return nil
		},
	}
}
