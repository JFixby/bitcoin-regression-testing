package btcregtest

import (
	"bytes"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"testing"
)

func TestGetBestBlock(t *testing.T) {
	// Skip tests when running with -short
	//if testing.Short() {
	//	t.Skip("Skipping RPC harness tests in short mode")
	//}
	r := ObtainHarness(mainHarnessName)

	_, prevbestHeight, err := r.NodeRPCClient().GetBestBlock()
	if err != nil {
		t.Fatalf("Call to `getbestblock` failed: %v", err)
	}

	// Create a new block connecting to the current tip.
	generatedBlockHashes, err := r.NodeRPCClient().Generate(1)
	if err != nil {
		t.Fatalf("Unable to generate block: %v", err)
	}

	bestHash, bestHeight, err := r.NodeRPCClient().GetBestBlock()
	if err != nil {
		t.Fatalf("Call to `getbestblock` failed: %v", err)
	}

	// Hash should be the same as the newly submitted block.
	if !bytes.Equal(bestHash.(*chainhash.Hash)[:], generatedBlockHashes[0].(*chainhash.Hash)[:]) {
		t.Fatalf("Block hashes do not match. Returned hash %v, wanted "+
			"hash %v", bestHash, generatedBlockHashes[0].(*chainhash.Hash)[:])
	}

	// Block height should now reflect newest height.
	if bestHeight != prevbestHeight+1 {
		t.Fatalf("Block heights do not match. Got %v, wanted %v",
			bestHeight, prevbestHeight+1)
	}
}
