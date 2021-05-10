package database

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

//const BlockReward = 100

type Hash [32]byte

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.Hex()), nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) IsEmpty() bool {
	emptyHash := Hash{}

	return bytes.Equal(emptyHash[:], h[:])
}

type Block struct {
	Header BlockHeader `json:"header"`
	TXs    []Tx        `json:"payload"`
}

type BlockHeader struct {
	Parent 		Hash    	`json:"parent"`
	Repository	String		`json:"repository"`
	Commit		20[byte]	`json:"commit"`
	PrevCommit 	20[byte]	`json:"prevCommit"`
	Number 		uint64  	`json:"number"`
	Time   		uint64  	`json:"time"`
	Miner  		Account		`json:"miner"`
}

type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

func NewBlock(parent Hash, repository string, commit 20[byte], prevCommit 20[byte], number uint64, time uint64, miner Account, txs []Tx) Block {
	return Block{BlockHeader{parent, repository, commit, prevCommit, number, nonce, time, miner}, txs}
}

func (b Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(blockJson), nil
}