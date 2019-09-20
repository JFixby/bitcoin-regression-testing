
package btcregtest

import (
	"github.com/btcsuite/btcutil"
	"github.com/jfixby/btcharness"
	"github.com/jfixby/coinharness"
	"testing"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func TestMemWalletLockedOutputs(t *testing.T) {
	// Skip tests when running with -short
	//if testing.Short() {
	//	t.Skip("Skipping RPC h tests in short mode")
	//}
	r := ObtainHarness(mainHarnessName)
	// Obtain the initial balance of the wallet at this point.
	startingBalance, err := r.Wallet.GetBalance("")
	if err != nil {
		t.Fatalf("unable to get balance: %v", err)
	}

	// First, create a signed transaction spending some outputs.
	addr, err := r.Wallet.NewAddress(nil)
	if err != nil {
		t.Fatalf("unable to generate new address: %v", err)
	}
	pkScript, err := txscript.PayToAddrScript(addr.(btcutil.Address))
	if err != nil {
		t.Fatalf("unable to create script: %v", err)
	}
	outputAmt := btcutil.Amount(50 * btcutil.SatoshiPerBitcoin)
	output := wire.NewTxOut(int64(outputAmt), pkScript)
	ctargs := &coinharness.CreateTransactionArgs{
		Outputs: []coinharness.OutputTx{&btcharness.OutputTx{output}},
		FeeRate: 10,
		Change:  true,
	}
	tx, err := r.Wallet.CreateTransaction(ctargs)
	if err != nil {
		t.Fatalf("unable to create transaction: %v", err)
	}

	// The current wallet balance should now be at least 50 BTC less
	// (accounting for fees) than the period balance
	currentBalance, err := r.Wallet.GetBalance("")
	if err != nil {
		t.Fatalf("unable to get balance: %v", err)
	}
	if !(currentBalance.TotalSpendable <= startingBalance.TotalSpendable-outputAmt) {
		t.Fatalf("spent outputs not locked: previous balance %v, "+
			"current balance %v", startingBalance, currentBalance)
	}

	// Now unlocked all the spent inputs within the unbroadcast signed
	// transaction. The current balance should now be exactly that of the
	// starting balance.
	txin := tx.TxIn()
	inpts := make([]coinharness.InputTx, len(txin))
	for i, j := range txin {
		inpts[i] = j
	}
	r.Wallet.UnlockOutputs(inpts)
	currentBalance, err = r.Wallet.GetBalance("")
	if err != nil {
		t.Fatalf("unable to get balance: %v", err)
	}
	if currentBalance.TotalSpendable != startingBalance.TotalSpendable {
		t.Fatalf("current and starting balance should now match: "+
			"expected %v, got %v", startingBalance, currentBalance)
	}
}
