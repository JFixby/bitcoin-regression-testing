Btcsuite (Bitcoin) regression testing
=======
[![Build Status](http://img.shields.io/travis/jfixby/btcregtest.svg)](https://travis-ci.org/jfixby/btcregtest)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)

Harbours a pre-configured test setup and unit tests to run RPC-driven node tests.

Builds a btcd-specific RPC testing harness crafting and executing integration
tests by driving a `btcd` instance via the `RPC` interface.

Each instance of an active harness comes equipped with a simple in-memory
HD wallet capable of properly syncing to the generated chain, creating new
addresses, and crafting fully signed transactions paying to an arbitrary
set of outputs. 

## Build 

```
set GO111MODULE=on
go build ./...
go test ./...
 
 ```
 
## License
This code is licensed under the [copyfree](http://copyfree.org) ISC License.