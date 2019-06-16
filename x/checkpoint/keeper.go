package checkpoint

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/requilence/checkpointchain/x/checkpoint/watchers/ethereum"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

func numberToBinary(block uint32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, block)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func binaryToNumber(key []byte) uint32 {
	var k uint32
	buf := new(bytes.Buffer)
	buf.Write(key)
	err := binary.Read(buf, binary.LittleEndian, &k)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return k
}

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetCheckPointBeforeIndex(ctx sdk.Context, startBlock uint32) (blockNumber uint32, blockHash []byte) {
	store := ctx.KVStore(k.storeKey)
	for i := startBlock; i > 0; i-- {
		key := numberToBinary(i)
		if store.Has(key) {
			return i, store.Get(key)
		}
	}
	return
}

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetCheckPointAtIndex(ctx sdk.Context, blockNumber uint32)  *ethereum.SectionInfo {
	store := ctx.KVStore(k.storeKey)
	key := numberToBinary(blockNumber)
	if !store.Has(key) {
		return nil
	}

	var info ethereum.SectionInfo
	encoded := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(encoded, &info)
	return &info
}

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetMaximumCheckpointIndex(ctx sdk.Context) uint32 {
	store := ctx.KVStore(k.storeKey)
	key := numberToBinary(0)
	if !store.Has(key) {
		return 0
	}

	return binaryToNumber(store.Get(key))
}

// Sets the entire Whois metadata struct for a name
func (k Keeper) SetCheckpoint(ctx sdk.Context, info ethereum.SectionInfo) {
	store := ctx.KVStore(k.storeKey)

	// set maximum
	store.Set(numberToBinary(0), numberToBinary(uint32(info.SectionIndex)))
	store.Set(numberToBinary(uint32(info.SectionIndex)), k.cdc.MustMarshalBinaryBare(info))
}

// Sets the entire Whois metadata struct for a name
func (k Keeper) SetMaximum(ctx sdk.Context, index uint32) {
	store := ctx.KVStore(k.storeKey)

	// set maximum
	store.Set(numberToBinary(0), numberToBinary(index))
}

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k Keeper) GetCheckpointsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.TransientStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
