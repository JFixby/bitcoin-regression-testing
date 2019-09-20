package btcregtest

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"testing"
)

func TestSetupValidity(t *testing.T) {
	coins50 := btcutil.Amount(50 /*BTC*/ * 1e8)
	stringVal := fmt.Sprintf("%v", coins50)
	expectedStringVal := "50 BTC"
	//pin.D("stringVal", stringVal)
	if expectedStringVal != stringVal {
		t.Fatalf("Incorrect coin: "+
			"expected %v, got %v", expectedStringVal, stringVal)
	}
}
