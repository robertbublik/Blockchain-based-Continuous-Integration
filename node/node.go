package node

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/robertbublik/bci/database"
	"net/http"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

const DefaultAccount = ""
const DefaultIP = "127.0.0.1"
const DefaultHTTPort = 8080
const endpointStatus = "/node/status"

const endpointSync = "/node/sync"
const endpointSyncQueryKeyFromBlock = "fromBlock"

const endpointAddPeer = "/node/peer"
const endpointAddPeerQueryKeyIP = "ip"
const endpointAddPeerQueryKeyPort = "port"
const endpointAddPeerQueryKeyAccount = "account"

const miningIntervalSeconds = 10

type PeerNode struct {
	IP          string           `json:"ip"`
	Port        uint64           `json:"port"`
	IsBootstrap bool             `json:"is_bootstrap"`
	Account     database.Account `json:"account"`

	// Whenever my node already established connection, sync with this Peer
	connected bool
}

func (pn PeerNode) TcpAddress() string {
	return fmt.Sprintf("%s:%d", pn.IP, pn.Port)
}

type Node struct {
	dataDir string
	info    PeerNode

	state           *database.State
	knownPeers      map[string]PeerNode
	pendingTXs      map[string]database.Tx
	archivedTXs     map[string]database.Tx
	newSyncedBlocks chan database.Block
	newPendingTXs   chan database.Tx
	isMining        bool
}

func New(dataDir string, ip string, port uint64, acc database.Account, bootstrap PeerNode) *Node {
	knownPeers := make(map[string]PeerNode)
	knownPeers[bootstrap.TcpAddress()] = bootstrap

	return &Node{
		dataDir:         dataDir,
		info:            NewPeerNode(ip, port, false, acc, true),
		knownPeers:      knownPeers,
		pendingTXs:      make(map[string]database.Tx),
		archivedTXs:     make(map[string]database.Tx),
		newSyncedBlocks: make(chan database.Block),
		newPendingTXs:   make(chan database.Tx, 10000),
		isMining:        false,
	}
}

func NewPeerNode(ip string, port uint64, isBootstrap bool, acc database.Account, connected bool) PeerNode {
	return PeerNode{ip, port, isBootstrap, acc, connected}
}

func (n *Node) Run(ctx context.Context) error {
	fmt.Println(fmt.Sprintf("Listening on: %s:%d", n.info.IP, n.info.Port))

	state, err := database.NewStateFromDisk(n.dataDir)
	if err != nil {
		return err
	}
	defer state.Close()

	n.state = state

	fmt.Println("Blockchain state:")
	fmt.Printf("	- height: %d\n", n.state.LatestBlock().Header.Index)
	fmt.Printf("	- hash: %s\n", n.state.LatestBlockHash().Hex())

	go n.sync(ctx)
	go n.mine(ctx)
	
	
	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		listBalancesHandler(w, r, state)
	})

	http.HandleFunc("/tx/list", func(w http.ResponseWriter, r *http.Request) {
		listTxHandler(w, r, n)
	})

	http.HandleFunc("/tx/add", func(w http.ResponseWriter, r *http.Request) {
		txAddHandler(w, r, n)
	})

	http.HandleFunc("/tx/get", func(w http.ResponseWriter, r *http.Request) {
		txGetHandler(w, r, n)
	})

	http.HandleFunc(endpointStatus, func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, n)
	})

	http.HandleFunc(endpointSync, func(w http.ResponseWriter, r *http.Request) {
		syncHandler(w, r, n)
	})

	http.HandleFunc(endpointAddPeer, func(w http.ResponseWriter, r *http.Request) {
		addPeerHandler(w, r, n)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", n.info.Port)}

	go func() {
		<-ctx.Done()
		_ = server.Close()
	}()

	err = server.ListenAndServe()
	// This shouldn't be an error!
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (n *Node) mine(ctx context.Context) error {
	var miningCtx context.Context
	var stopCurrentMining context.CancelFunc

	ticker := time.NewTicker(time.Second * miningIntervalSeconds)

	for {
		select {
		case <-ticker.C:
			go func() {
				if len(n.pendingTXs) > 0 && !n.isMining {
					n.isMining = true
					if n.info.Port != 8080 {
						miningCtx, stopCurrentMining = context.WithCancel(ctx)
						err := n.minePendingTXs(miningCtx)
						if err != nil {
							fmt.Printf("ERROR: %s\n", err)
						}
					}
					n.isMining = false
				}
			}()

		case block, _ := <-n.newSyncedBlocks:
			if n.isMining {
				blockHash, _ := block.Hash()
				fmt.Printf("\nPeer mined next Block '%s' faster :(\n", blockHash.Hex())

				n.removeMinedPendingTXs(block)
				stopCurrentMining()
			}

		case <-ctx.Done():
			ticker.Stop()
			return nil
		}
	}
}

func (n *Node) minePendingTXs(ctx context.Context) error {
	blockToMine := NewPendingBlock(
		n.state.NextBlockIndex(),
		n.state.LatestBlockHash(),
		n.info.Account,
		n.getPendingTXsAsArray(),
	)

	minedBlock, err := Mine(ctx, blockToMine, n.dataDir)
	if err != nil {
		return err
	}

	n.removeMinedPendingTX(minedBlock)

	_, err = n.state.AddBlock(minedBlock)
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) removeMinedPendingTX(block database.Block) {
	if (block.Body.TX != database.Tx{}) && len(n.pendingTXs) > 0 {
		fmt.Println("Updating in-memory Pending TXs Pool:")
	}

	txHash, _ := block.Body.TX.Hash()
	if _, exists := n.pendingTXs[txHash.Hex()]; exists {
		fmt.Printf("\t-archiving mined TX: %s\n", txHash.Hex())

		n.archivedTXs[txHash.Hex()] = block.Body.TX
		delete(n.pendingTXs, txHash.Hex())
	}
	
}

func (n *Node) LatestBlockHash() database.Hash {
	return n.state.LatestBlockHash()
}

func (n *Node) AddPeer(peer PeerNode) {
	n.knownPeers[peer.TcpAddress()] = peer
}

func (n *Node) RemovePeer(peer PeerNode) {
	delete(n.knownPeers, peer.TcpAddress())
}

func (n *Node) IsKnownPeer(peer PeerNode) bool {
	if peer.IP == n.info.IP && peer.Port == n.info.Port {
		return true
	}

	_, isKnownPeer := n.knownPeers[peer.TcpAddress()]

	return isKnownPeer
}

func (n *Node) AddPendingTX(tx database.Tx, fromPeer PeerNode) error {
	txHash, err := tx.Hash()
	if err != nil {
		return err
	}

	txJson, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	_, isAlreadyPending := n.pendingTXs[txHash.Hex()]
	_, isArchived := n.archivedTXs[txHash.Hex()]

	if !isAlreadyPending && !isArchived {
		fmt.Printf("Added Pending TX %s from Peer %s\n", txJson, fromPeer.TcpAddress())
		n.pendingTXs[txHash.Hex()] = tx
		n.newPendingTXs <- tx
	}

	return nil
}


func (n *Node) GetTx(id string) database.Tx {
	for _, tx := range n.pendingTXs {
		if tx.Id == id {
			return tx
		}
	}

	return database.Tx{}
}


func (n *Node) getPendingTXsAsArray() []database.Tx {
	txs := make([]database.Tx, len(n.pendingTXs))

	i := 0
	for _, tx := range n.pendingTXs {
		txs[i] = tx
		i++
	}

	return txs
}

func (n *Node) IsAlreadyPending(id string) bool {
	for _, tx := range n.pendingTXs {
		if tx.Id == id {
			return true
		}
	}
	return false
}

func TxRequestToString(req TxReq) string {
	txHash := sha256.Sum256([]byte(req.From + strconv.FormatUint(req.Value, 10) + req.Repository + req.Language + req.Commit + req.PrevCommit))
	return hex.EncodeToString(txHash[:])[:5]
}