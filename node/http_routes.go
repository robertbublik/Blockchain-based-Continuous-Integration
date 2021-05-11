package node

import (
	"fmt"
	"github.com/robertbublik/bci/database"
	"net/http"
	"strconv"
)

type ErrRes struct {
	Error string `json:"error"`
}

type BalancesRes struct {
	Hash     database.Hash             `json:"blockHash"`
	Balances map[database.Account]uint `json:"balances"`
}

type PendingTxRes struct {
	PendingTXs 	[]database.Tx 		`json:"pendingTxs"`
}

type TxReq struct {
	From  		string 	`json:"from"`
	Value 		uint    			`json:"value"`
	Repository  string  			`json:"repository"`
	Commit 		[20]byte 			`json:"commit"`
	prevCommit 	[20]byte 			`json:"prevCommit"`
	Time  		uint64  			`json:"time"`
}

type TxAddRes struct {
	Success bool `json:"success"`
}

type TxMineRes struct {
	Success bool `json:"success"`
}

type StatusRes struct {
	Hash       	database.Hash       `json:"blockHash"`
	BlockIndex 	uint64              `json:"blockIndex"`
	KnownPeers 	map[string]PeerNode `json:"peersKnown"`
	PendingTXs 	[]database.Tx       `json:"pendingTxs"`
}

type SyncRes struct {
	Blocks []database.Block `json:"blocks"`
}

type AddPeerRes struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func listBalancesHandler(w http.ResponseWriter, r *http.Request, state *database.State) {
	writeRes(w, BalancesRes{state.LatestBlockHash(), state.Balances})
}

func listPendingTxHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	writeRes(w, PendingTxRes{node.getPendingTXsAsArray()})
}

// adds transaction to BCI
func txAddHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	req := TxReq{}
	err := readReq(r, &req)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	if req.Value > node.state.Balances[database.NewAccount(req.From)] {
		writeErrRes(w, fmt.Errorf("Balance too low. %s", err.Error()))
		return
	}

	tx := database.NewTx(database.NewAccount(req.From), req.Value, req.Repository, req.Commit, req.prevCommit)
	err = node.AddPendingTX(tx, node.info)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	writeRes(w, TxAddRes{Success: true})
}

// worker chooses transaction to mine
func txMineHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	req := TxReq{}
	err := readReq(r, &req)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	writeRes(w, TxAddRes{Success: true})
}

func statusHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	res := StatusRes{
		Hash:       node.state.LatestBlockHash(),
		BlockIndex: node.state.LatestBlock().Header.Index,
		KnownPeers: node.knownPeers,
		PendingTXs: node.getPendingTXsAsArray(),
	}

	writeRes(w, res)
}

func syncHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	reqHash := r.URL.Query().Get(endpointSyncQueryKeyFromBlock)

	hash := database.Hash{}
	err := hash.UnmarshalText([]byte(reqHash))
	if err != nil {
		writeErrRes(w, err)
		return
	}

	blocks, err := database.GetBlocksAfter(hash, node.dataDir)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	writeRes(w, SyncRes{Blocks: blocks})
}

func addPeerHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	peerIP := r.URL.Query().Get(endpointAddPeerQueryKeyIP)
	peerPortRaw := r.URL.Query().Get(endpointAddPeerQueryKeyPort)
	minerRaw := r.URL.Query().Get(endpointAddPeerQueryKeyMiner)

	peerPort, err := strconv.ParseUint(peerPortRaw, 10, 32)
	if err != nil {
		writeRes(w, AddPeerRes{false, err.Error()})
		return
	}

	peer := NewPeerNode(peerIP, peerPort, false, database.NewAccount(minerRaw), true)

	node.AddPeer(peer)

	fmt.Printf("Peer '%s' was added into KnownPeers\n", peer.TcpAddress())

	writeRes(w, AddPeerRes{true, ""})
}
