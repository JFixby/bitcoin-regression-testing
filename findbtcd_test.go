package btcregtest

import (
	"github.com/jfixby/pin"
	"github.com/jfixby/pin/fileops"
	"testing"
)

func TestFindDCR(t *testing.T) {
	path := fileops.Abs("../../btcsuite/btcd")
	pin.D("path", path)
}
