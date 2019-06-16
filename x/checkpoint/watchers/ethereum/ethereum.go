package ethereum

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var client *ethclient.Client

type SectionInfo struct {
	SectionIndex int64
	SectionHead  string
	CHTRoot      string
	BloomRoot    string
}

const defaultRpcAddress = "ws://127.0.0.1:8546"
const defaultGethBinary = "geth"

var ErrorFailedToGetABlock = fmt.Errorf("failed to get a block")
var ErrorFailedToConnect = fmt.Errorf("failed to connect to ethereum RPC")

var gethBinary string

func getClient() (*ethclient.Client, error) {
	rpcAddress := defaultRpcAddress
	gethBinary = defaultGethBinary
	if v := os.Getenv("CHECKPOINTCHAIN_ETH_RPC"); v != "" {
		rpcAddress = v
	}

	if v := os.Getenv("CHECKPOINTCHAIN_GETH_BINARY"); v != "" {
		gethBinary = v
	}

	var err error
	if client == nil {
		client, err = ethclient.Dial(rpcAddress)
		if err != nil {
			return nil, ErrorFailedToConnect
		}
	}

	return client, nil
}

func Watcher(divider int, ch chan SectionInfo) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %s", err.Error())
	}

	for {
		select {
		case err := <-sub.Err():
			fmt.Printf("Got error on subscription: %s", err.Error())
		case header := <-headers:
			s := []string{}
			fmt.Println(header.Number.Int64())
			b, err := exec.Command(gethBinary, "--verbosity", "0", "attach", "--exec", "[admin.nodeInfo.protocols.les.cht.bloomRoot, admin.nodeInfo.protocols.les.cht.chtRoot, admin.nodeInfo.protocols.les.cht.sectionHead, admin.nodeInfo.protocols.les.cht.sectionIndex]").Output()
			if err != nil {
				fmt.Printf("geth exec error: %s")
				continue
			}

			err = json.Unmarshal(b, &s)
			if err != nil {
				fmt.Printf("geth res unmarshal error: %s")
				continue
			}

			sectionIndex, err := strconv.Atoi(s[3])
			if err != nil {
				fmt.Printf("geth res unmarshal error: %s")
				continue
			}

			ch <- SectionInfo{SectionIndex: int64(sectionIndex), SectionHead: s[2], CHTRoot: s[1], BloomRoot: s[0]}

		}
	}
}

func GetBlockHash(number int64) (string, error) {
	client, err := getClient()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	block, err := client.BlockByNumber(ctx, big.NewInt(number))
	if err != nil {
		return "", ErrorFailedToGetABlock
	}

	return block.Hash().Hex(), nil
}
