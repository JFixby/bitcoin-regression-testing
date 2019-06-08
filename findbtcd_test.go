package btcregtest

import (
	"github.com/jfixby/pin"
	"github.com/picfight/pfcd_builder/fileops"
	"testing"
)

func TestFindDCR(t *testing.T) {
	path := fileops.Abs("../../btcsuite/btcd")
	pin.D("path", path)
}
