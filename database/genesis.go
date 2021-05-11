package database

import (
	"io/ioutil"
	"encoding/json"
)

var genesisJson = `
{
  "genesis_time": "2019-03-18T00:00:00.000000000Z",
  "chain_id": "BCI-ledger",
  "balances": {
    "developer-1": 1000,
	"developer-2": 1000,
	"developer-3": 1000,
	"miner-1": 0,
	"miner-2": 0,
	"miner-3": 0
  }
}`

type genesis struct {
	Balances map[Account]uint `json:"balances"`
}

func loadGenesis(path string) (genesis, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return genesis{}, err
	}

	var loadedGenesis genesis
	err = json.Unmarshal(content, &loadedGenesis)
	if err != nil {
		return genesis{}, err
	}

	return loadedGenesis, nil
}

func writeGenesisToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(genesisJson), 0644)
}