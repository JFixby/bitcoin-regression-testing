package btcregtest

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/jfixby/btcharness"
	"github.com/jfixby/coinharness"
	"github.com/jfixby/pin"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func genSpend(t *testing.T, r *coinharness.Harness, amt btcutil.Amount) *chainhash.Hash {
	// Grab a fresh address from the wallet.
	addr, err := r.Wallet.NewAddress(&coinharness.NewAddressArgs{"default"})
	if err != nil {
		t.Fatalf("unable to get new address: %v", err)
	}

	// Next, send amt to this address, spending from one of our
	// mature coinbase outputs.
	addrScript, err := txscript.PayToAddrScript(addr.Internal().(btcutil.Address))
	if err != nil {
		t.Fatalf("unable to generate pkscript to addr: %v", err)
	}
	output := wire.NewTxOut(int64(amt), addrScript)
	arg := &coinharness.SendOutputsArgs{
		Outputs: []coinharness.OutputTx{&btcharness.OutputTx{output}},
		FeeRate: 10,
	}
	txid, err := r.Wallet.SendOutputs(arg)
	if err != nil {
		t.Fatalf("coinbase spend failed: %v", err)
	}
	return txid.(*chainhash.Hash)
}

func assertTxMined(t *testing.T, r *coinharness.Harness, txid *chainhash.Hash, blockHash *chainhash.Hash) {
	block, err := r.NodeRPCClient().Internal().(*rpcclient.Client).GetBlock(blockHash)
	if err != nil {
		t.Fatalf("unable to get block: %v", err)
	}

	numBlockTxns := len(block.Transactions)
	if numBlockTxns < 2 {
		t.Fatalf("crafted transaction wasn't mined, block should have "+
			"at least %v transactions instead has %v", 2, numBlockTxns)
	}

	txHash1 := block.Transactions[1].TxHash()

	if txHash1 != *txid {
		t.Fatalf("txid's don't match, %v vs %v", txHash1, txid)
	}
}

func TestBallance(t *testing.T) {
	// Skip tests when running with -short
	//if testing.Short() {
	//	t.Skip("Skipping RPC harness tests in short mode")
	//}
	r := ObtainHarness(t.Name() + ".8")

	expectedBalance := btcutil.Amount(7200 * btcutil.SatoshiPerBitcoin)
	actualBalance := coinharness.GetBalance(t, r.Wallet).TotalSpendable.(btcutil.Amount)

	if actualBalance != expectedBalance {
		t.Fatalf("expected wallet balance of %v instead have %v",
			expectedBalance, actualBalance)
	}
}

func TestSendOutputs(t *testing.T) {
	// Skip tests when running with -short
	//if testing.Short() {
	//	t.Skip("Skipping RPC harness tests in short mode")
	//}
	r := ObtainHarness("TestSendOutputs")
	_, H, e := r.NodeRPCClient().GetBestBlock()
	pin.CheckTestSetupMalfunction(e)
	r.Wallet.Sync(H)
	// First, generate a small spend which will require only a single
	// input.
	txid := genSpend(t, r, btcutil.Amount(5*btcutil.SatoshiPerBitcoin))

	// Generate a single block, the transaction the wallet created should
	// be found in this block.
	blockHashes, err := r.NodeRPCClient().Generate(1)
	if err != nil {
		t.Fatalf("unable to generate single block: %v", err)
	}
	assertTxMined(t, r, txid, blockHashes[0].(*chainhash.Hash))

	// Next, generate a spend much greater than the block reward. This
	// transaction should also have been mined properly.
	txid = genSpend(t, r, btcutil.Amount(1000*btcutil.SatoshiPerBitcoin))
	blockHashes, err = r.NodeRPCClient().Generate(1)
	if err != nil {
		t.Fatalf("unable to generate single block: %v", err)
	}
	assertTxMined(t, r, txid, blockHashes[0].(*chainhash.Hash))

	// Generate another block to ensure the transaction is removed from the
	// mempool.
	if _, err := r.NodeRPCClient().Generate(1); err != nil {
		t.Fatalf("unable to generate block: %v", err)
	}
}
