package database

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	//"fmt"
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
	Body   BlockBody   `json:"body"`
}

type BlockHeader struct {
	Index 		uint64  	`json:"index"`
	Parent 		Hash    	`json:"parent"`
	Repository	string		`json:"repository"`
	Commit		[20]byte	`json:"commit"`
	PrevCommit 	[20]byte	`json:"prevCommit"`
	Time   		uint64  	`json:"time"`
	Miner  		Account		`json:"miner"`
}

type BlockBody struct {
	TX				Tx			`json:"tx"`
	ArtifactUrl		string		`json:"artifactUrl`
	ArtifactHash	string		`json:"artifactHash"`
	BuildLog		string		`json:"buildLog"`
}

type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

func NewBlock(index uint64, parent Hash, repository string, commit [20]byte, prevCommit [20]byte, time uint64, miner Account, tx Tx, artifactUrl string, artifactHash string, buildLog string) Block {
	return Block{BlockHeader{index, parent, repository, commit, prevCommit, time, miner}, BlockBody{tx, artifactUrl, artifactHash, buildLog}}
}

func (b Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(blockJson), nil
}