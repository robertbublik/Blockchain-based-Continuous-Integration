package node

import (
	"fmt"
	"github.com/robertbublik/bci/database"
	"net/http"
	"strconv"
	"errors"
	"crypto/sha256"
	"encoding/hex"
)

type ErrRes struct {
	Error string `json:"error"`
}

type BalancesRes struct {
	Hash     database.Hash             		`json:"blockHash"`
	Balances map[database.Account]uint64 	`json:"balances"`
}

type TxsListRes struct {
	TXsList		[]database.Tx	`json:"txsList"`		
}

type TxReq struct {
	From  		string 	`json:"from"`
	Value 		uint64 	`json:"value"`
	Repository  string  `json:"repository"`
	Language	string	`json:"language`
	Commit 		string 	`json:"commit"`
	PrevCommit 	string 	`json:"prevCommit"`
	Time  		uint64  `json:"time"`
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
	WriteRes(w, BalancesRes{state.LatestBlockHash(), state.Balances})
}

func listTxHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	WriteRes(w, TxsListRes{node.getPendingTXsAsArray()})
}

// adds transaction to BCI
func txAddHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	req := TxReq{}
	err := ReadReq(r, &req)
	if err != nil {
		WriteErrRes(w, err)
		return
	}

	if req.Value > node.state.Balances[database.NewAccount(req.From)] {
		err := errors.New("Balance too low.")
		WriteErrRes(w, err)
		return
	}

	txHash := sha256.Sum256([]byte(req.From + strconv.FormatUint(req.Value, 10) + req.Repository + req.Language + req.Commit + req.PrevCommit))
	txId := hex.EncodeToString(txHash[:])[:5]
	tx := database.NewTx(txId, database.NewAccount(req.From), req.Value, req.Repository, req.Language, req.Commit, req.PrevCommit)
	
	err = node.AddPendingTX(tx, node.info)
	if err != nil {
		WriteErrRes(w, err)
		return
	}

	WriteRes(w, TxAddRes{Success: true})
	fmt.Printf("Tx added succesfully\n")
}

// worker chooses transaction to mine
func txMineHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	req := TxReq{}
	err := ReadReq(r, &req)
	if err != nil {
		WriteErrRes(w, err)
		return
	}

	WriteRes(w, TxAddRes{Success: true})
}

func statusHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	res := StatusRes{
		Hash:       node.state.LatestBlockHash(),
		BlockIndex: node.state.LatestBlock().Header.Index,
		KnownPeers: node.knownPeers,
		PendingTXs: node.getPendingTXsAsArray(),
	}

	WriteRes(w, res)
}

func syncHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	reqHash := r.URL.Query().Get(endpointSyncQueryKeyFromBlock)

	hash := database.Hash{}
	err := hash.UnmarshalText([]byte(reqHash))
	if err != nil {
		WriteErrRes(w, err)
		return
	}

	blocks, err := database.GetBlocksAfter(hash, node.dataDir)
	if err != nil {
		WriteErrRes(w, err)
		return
	}

	WriteRes(w, SyncRes{Blocks: blocks})
}

func addPeerHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	peerIP := r.URL.Query().Get(endpointAddPeerQueryKeyIP)
	peerPortRaw := r.URL.Query().Get(endpointAddPeerQueryKeyPort)
	accountRaw := r.URL.Query().Get(endpointAddPeerQueryKeyAccount)
	peerPort, err := strconv.ParseUint(peerPortRaw, 10, 32)
	if err != nil {
		WriteRes(w, AddPeerRes{false, err.Error()})
		return
	}

	peer := NewPeerNode(peerIP, peerPort, false, database.NewAccount(accountRaw), true)

	node.AddPeer(peer)

	fmt.Printf("Peer '%s' was added into KnownPeers\n", peer.TcpAddress())

	WriteRes(w, AddPeerRes{true, ""})
}
