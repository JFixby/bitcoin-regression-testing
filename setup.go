
package btcregtest

import (
	"fmt"
	"github.com/jfixby/btcharness/memwallet"
	"github.com/jfixby/btcharness"
	"github.com/jfixby/btcharness/nodecls"
	"github.com/jfixby/coinharness"
	"github.com/jfixby/pin"
	"github.com/jfixby/pin/commandline"
	"github.com/jfixby/pin/fileops"
	"github.com/jfixby/pin/gobuilder"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg"
)

// Default harness name
const mainHarnessName = "main"

// SimpleTestSetup harbours:
// - rpctest setup
// - csf-fork test setup
// - and bip0009 test setup
type SimpleTestSetup struct {
	// harnessPool stores and manages harnesses
	// multiple harness instances may be run concurrently, to allow for testing
	// complex scenarios involving multiple nodes.
	harnessPool *pin.Pool

	// Regnet25 creates a regnet test harness
	// with 25 mature outputs.
	Regnet25 *coinharness.ChainWithMatureOutputsSpawner

	// Regnet5 creates a regnet test harness
	// with 5 mature outputs.
	Regnet5 *coinharness.ChainWithMatureOutputsSpawner

	// Regnet1 creates a regnet test harness
	// with 1 mature output.
	Regnet1 *coinharness.ChainWithMatureOutputsSpawner

	// Simnet1 creates a simnet test harness
	// with 1 mature output.
	Simnet1 *coinharness.ChainWithMatureOutputsSpawner

	// Regnet0 creates a regnet test harness
	// with only the genesis block.
	Regnet0 *coinharness.ChainWithMatureOutputsSpawner

	// Simnet0 creates a simnet test harness
	// with only the genesis block.
	Simnet0 *coinharness.ChainWithMatureOutputsSpawner

	// ConsoleNodeFactory produces a new TestNode instance upon request
	NodeFactory coinharness.TestNodeFactory

	// WalletFactory produces a new TestWallet instance upon request
	WalletFactory coinharness.TestWalletFactory

	// WorkingDir defines test setup working dir
	WorkingDir *pin.TempDirHandler
}

// TearDown all harnesses in test Pool.
// This includes removing all temporary directories,
// and shutting down any created processes.
func (setup *SimpleTestSetup) TearDown() {
	setup.harnessPool.DisposeAll()
	//setup.nodeGoBuilder.Dispose()
	setup.WorkingDir.Dispose()
}

// Setup deploys this test setup
func Setup() *SimpleTestSetup {
	setup := &SimpleTestSetup{
		WalletFactory: &memwallet.WalletFactory{},
		//Network:       &chaincfg.RegressionNetParams,
		WorkingDir: pin.NewTempDir(setupWorkingDir(), "simpleregtest").MakeDir(),
	}

	btcdEXE := &commandline.ExplicitExecutablePathString{
		PathString: "btcd",
	}
	setup.NodeFactory = &nodecls.ConsoleNodeFactory{
		NodeExecutablePathProvider: btcdEXE,
	}

	portManager := &coinharness.LazyPortManager{
		BasePort: 30000,
	}

	// Deploy harness spawner with generated
	// test chain of 25 mature outputs
	setup.Regnet25 = &coinharness.ChainWithMatureOutputsSpawner{
		WorkingDir:        setup.WorkingDir.Path(),
		DebugNodeOutput:   true,
		DebugWalletOutput: true,
		NumMatureOutputs:  25,
		NetPortManager:    portManager,
		WalletFactory:     setup.WalletFactory,
		NodeFactory:       setup.NodeFactory,
		ActiveNet:         &chaincfg.RegressionNetParams,
	}

	// Deploy harness spawner with generated
	// test chain of 5 mature outputs
	setup.Regnet5 = &coinharness.ChainWithMatureOutputsSpawner{
		WorkingDir:        setup.WorkingDir.Path(),
		DebugNodeOutput:   true,
		DebugWalletOutput: true,
		NumMatureOutputs:  5,
		NetPortManager:    portManager,
		WalletFactory:     setup.WalletFactory,
		NodeFactory:       setup.NodeFactory,
		ActiveNet:         &chaincfg.RegressionNetParams,
	}

	setup.Regnet1 = &coinharness.ChainWithMatureOutputsSpawner{
		WorkingDir:        setup.WorkingDir.Path(),
		DebugNodeOutput:   true,
		DebugWalletOutput: true,
		NumMatureOutputs:  1,
		NetPortManager:    portManager,
		WalletFactory:     setup.WalletFactory,
		NodeFactory:       setup.NodeFactory,
		ActiveNet:         &chaincfg.RegressionNetParams,
		NodeStartExtraArguments: map[string]interface{}{
			"rejectnonstd": commandline.NoArgumentValue,
		},
	}

	setup.Simnet1 = &coinharness.ChainWithMatureOutputsSpawner{
		WorkingDir:        setup.WorkingDir.Path(),
		DebugNodeOutput:   true,
		DebugWalletOutput: true,
		NumMatureOutputs:  1,
		NetPortManager:    portManager,
		WalletFactory:     setup.WalletFactory,
		NodeFactory:       setup.NodeFactory,
		ActiveNet:         &chaincfg.SimNetParams,
		NodeStartExtraArguments: map[string]interface{}{
			"rejectnonstd": commandline.NoArgumentValue,
		},
	}

	// Deploy harness spawner with empty test chain
	setup.Regnet0 = &coinharness.ChainWithMatureOutputsSpawner{
		WorkingDir:        setup.WorkingDir.Path(),
		DebugNodeOutput:   true,
		DebugWalletOutput: true,
		NumMatureOutputs:  0,
		NetPortManager:    portManager,
		WalletFactory:     setup.WalletFactory,
		NodeFactory:       setup.NodeFactory,
		ActiveNet:         &chaincfg.RegressionNetParams,
	}
	// Deploy harness spawner with empty test chain
	setup.Simnet0 = &coinharness.ChainWithMatureOutputsSpawner{
		WorkingDir:        setup.WorkingDir.Path(),
		DebugNodeOutput:   true,
		DebugWalletOutput: true,
		NumMatureOutputs:  0,
		NetPortManager:    portManager,
		WalletFactory:     setup.WalletFactory,
		NodeFactory:       setup.NodeFactory,
		ActiveNet:         &chaincfg.SimNetParams,
	}

	setup.harnessPool = pin.NewPool(setup.Regnet25)

	return setup
}

func setupWorkingDir() string {
	testWorkingDir, err := ioutil.TempDir("", "integrationtest")
	if err != nil {
		fmt.Println("Unable to create working dir: ", err)
		os.Exit(-1)
	}
	return testWorkingDir
}

func setupBuild(buildName string, workingDir string, nodeProjectGoPath string) *gobuilder.GoBuider {
	tempBinDir := filepath.Join(workingDir, "bin")
	pin.MakeDirs(tempBinDir)

	nodeGoBuilder := &gobuilder.GoBuider{
		GoProjectPath:    nodeProjectGoPath,
		OutputFolderPath: tempBinDir,
		BuildFileName:    buildName,
	}
	return nodeGoBuilder
}
