#!/usr/bin/env bash

set -ex

# Default GOVERSION
[[ ! "$GOVERSION" ]] && GOVERSION=1.12
REPO=btcregtest

testrepo () {
  GO=go
  PROJECT=btcsuite
  NODE_REPO=btcd
  WALLET_REPO=btcwallet

  $GO version

  # binary needed for RPC tests
  env CC=gcc

  # run tests on all modules

  pushd ../../
  git clone --depth=50 --branch=bug_fix https://github.com/jfixby/dcrd.git ${PROJECT}/${NODE_REPO}
  git clone --depth=50 --branch=master https://github.com/${PROJECT}/${WALLET_REPO}.git ${PROJECT}/${WALLET_REPO}
  popd

  $GO fmt ./...
  $GO build ./...

  pushd ../../${PROJECT}/${NODE_REPO}
  $GO build ./...
  $GO install -v . ./cmd/...
  popd

  pushd ../../${PROJECT}/${WALLET_REPO}
  $GO build ./...
  $GO install -v . ./cmd/...
  popd

  GO111MODULE=on
  ${NODE_REPO} --version
  ${WALLET_REPO} --version
  $GO clean -testcache
  $GO build ./...
  $GO test -v ./...

  echo "------------------------------------------"
  echo "Tests completed successfully!"
}

testrepo
exit
