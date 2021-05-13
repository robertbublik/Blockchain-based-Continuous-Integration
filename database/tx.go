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
	Id			string		`json:"id"`
	From  		Account 	`json:"from"`
	Value 		uint64    	`json:"value"`
	Repository  string  	`json:"repository"`
	Language 	string		`json:"language`
	Commit 		string 		`json:"commit"`
	prevCommit 	string 		`json:"prevCommit"`
	Time  		uint64  	`json:"time"`
}

func NewTx(id string, from Account, value uint64, repository string, language string, commit string, prevCommit string) Tx {
	return Tx{id, from, value, repository, language, commit, prevCommit, uint64(time.Now().Unix())}
}

func (t Tx) Hash() (Hash, error) {
	txJson, err := json.Marshal(t)
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(txJson), nil
}

