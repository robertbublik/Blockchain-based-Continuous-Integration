# Blockchain-based Continuous Integration

Table of Contents
=================

   * [TBB Training Ledger](#tbb-training-ledger)
      * [Sneak peek to Chapter 13](#sneak-peek-to-chapter-13)
         * [1/4 Check the current blockchain network status](#14-check-the-current-blockchain-network-status)
         * [2/4 Download the pre-compiled blockchain program](#24-download-the-pre-compiled-blockchain-program)
            * [Install](#install)
               * [Download](#download)
                  * [Linux](#linux)
                  * [MacOS](#macos)
               * [Verify the version](#verify-the-version)
         * [3/4 Connect to the training network](#34-connect-to-the-training-network)
         * [4/4 Check the current blockchain network status](#44-check-the-current-blockchain-network-status)
   * [Introduction](#introduction)
      * [How?](#how)
      * [What will you build?](#what-will-you-build)
         * [1) You will build a peer-to-peer system from scratch](#1-you-will-build-a-peer-to-peer-system-from-scratch)
         * [2) You will secure the system with a day-to-day practical cryptography](#2-you-will-secure-the-system-with-a-day-to-day-practical-cryptography)
         * [3) You will implement Bitcoin, Ethereum and XRP backend components](#3-you-will-implement-bitcoin-ethereum-and-xrp-backend-components)
         * [4) You will write unit tests and integration tests for all core components](#4-you-will-write-unit-tests-and-integration-tests-for-all-core-components)
      * [How to use this repository](#how-to-use-this-repository)
      * [Installation](#installation)
      * [Getting started](#getting-started)
   * [Usage](#usage)
      * [Install](#install-1)
      * [CLI](#cli)
         * [Show available commands and flags](#show-available-commands-and-flags)
            * [Show available run settings](#show-available-run-settings)
         * [Run a TBB node connected to the official book's test network](#run-a-tbb-node-connected-to-the-official-books-test-network)
         * [Run a TBB bootstrap node in isolation, on your localhost only](#run-a-tbb-bootstrap-node-in-isolation-on-your-localhost-only)
            * [Run a second TBB node connecting to your first one](#run-a-second-tbb-node-connecting-to-your-first-one)
         * [Create a new account](#create-a-new-account)
      * [HTTP](#http)
         * [List all balances](#list-all-balances)
         * [Send a signed TX](#send-a-signed-tx)
         * [Check node's status (latest block, known peers, pending TXs)](#check-nodes-status-latest-block-known-peers-pending-txs)
      * [Tests](#tests)
   * [Start](#start)
      * [Get the first 7 chapters for FREE](#get-the-first-7-chapters-for-free)
      * [Buy complete eBook](#buy-complete-ebook)
   * [Finish](#finish)
      * [Request 1000 TBB testing tokens](#request-1000-tbb-testing-tokens)

## Warning
This implementation currently only works on Linux.

## Installation

[Open instructions.](./Installation.md)

## Getting started

# Usage

## Install
```
cd $GOPATH/src/github.com/robertbublik
go install ./cmd/bci
```

#### Show available commands
```bash
bci help

Blockchain-based Continuous Integration CLI                                                                                                                                                                                                     Usage:                                     
	bci [flags]                                                                                                             
	bci [command]                                                                                                                                                                                                                                 Available Commands:                                                                                                       
	balances    Interact with balances (list...).                                                                           
	help        Help about any command                                                                                      
	run         Launches the BCI node and its HTTP API.                                                                     
	status      Displays status of BCI.                                                                                     
	tx          Interact with transactions (add, list...).                                                                                                                                                                                                                                                                              Flags:                                                                                                                    
	-h, --help   help for bci 
```

### Start Docker registry
```bash
docker run -d -p 5000:5000 --restart=always --name registry registry:2
```

### Start BCI nodes 
```bash
bci run --datadir=$HOME/bci-nodes/bootstrap --ip=127.0.0.1 --port=8080
bci run --datadir=$HOME/bci-nodes/miner-1 --account=miner-1 --ip=127.0.0.1 --port=8081
bci run --datadir=$HOME/bci-nodes/miner-2 --account=miner-2 --ip=127.0.0.1 --port=8082
bci run --datadir=$HOME/bci-nodes/miner-3 --account=miner-3 --ip=127.0.0.1 --port=8083
```

### Add a transaction
```
bci tx add --from=developer-1 --value=100 --language=docker --repository=https://github.com/robertbublik/BCI_docker
bci tx add --from=developer-2 --value=200 --language=docker --repository=https://github.com/robertbublik/BCI_docker
bci tx add --from=developer-3 --value=300 --language=docker --repository=https://github.com/robertbublik/BCI_docker
```

#### Run a second TBB node connecting to your first one
```
tbb run --datadir=$HOME/.tbb --ip=127.0.0.1 --port=8081 --bootstrap-ip=127.0.0.1 --bootstrap-port=8080
```

