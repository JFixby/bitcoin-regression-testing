package btcregtest

import (
	"bytes"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"testing"
)

func TestGetBlockHash(t *testing.T) {
	//if testing.Short() {
	//	t.Skip("Skipping RPC harness tests in short mode")
	//}
	r := ObtainHarness(mainHarnessName)
	// Create a new block connecting to the current tip.
	generatedBlockHashes, err := r.NodeRPCClient().Generate(1)
	if err != nil {
		t.Fatalf("Unable to generate block: %v", err)
	}

	info, err := r.NodeRPCClient().Internal().(*rpcclient.Client).GetInfo()
	if err != nil {
		t.Fatalf("call to getinfo cailed: %v", err)
	}

	blockHash, err := r.NodeRPCClient().Internal().(*rpcclient.Client).GetBlockHash(int64(info.Blocks))
	if err != nil {
		t.Fatalf("Call to `getblockhash` failed: %v", err)
	}

	// Block hashes should match newly created block.
	if !bytes.Equal(generatedBlockHashes[0].(*chainhash.Hash)[:], blockHash[:]) {
		t.Fatalf("Block hashes do not match. Returned hash %v, wanted "+
			"hash %v", blockHash, generatedBlockHashes[0].(*chainhash.Hash)[:])
	}
}
