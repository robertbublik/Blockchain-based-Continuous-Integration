package database

import (
	"crypto/sha256"
	"encoding/json"
	"time"
)

type Account string

func NewAccount(value string) Account {
	return Account(value)
}

type Tx struct {
	From  		Account 	`json:"from"`
	Value 		uint64    	`json:"value"`
	Repository  string  	`json:"repository"`
	Commit 		string 		`json:"commit"`
	prevCommit 	string 		`json:"prevCommit"`
	Time  		uint64  	`json:"time"`
}

func NewTx(from Account, value uint64, repository string, commit string, prevCommit string) Tx {
	return Tx{from, value, repository, commit, prevCommit, uint64(time.Now().Unix())}
}

func (t Tx) Hash() (Hash, error) {
	txJson, err := json.Marshal(t)
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(txJson), nil
}
