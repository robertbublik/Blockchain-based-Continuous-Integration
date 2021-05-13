package database

import (
	"path/filepath"
	"os"
	"io/ioutil"
	//"github.com/robertbublik/bci/node"
)

func initDataDirIfNotExists(dataDir string) error {
	if fileExist(getGenesisJsonFilePath(dataDir)) {
		return nil
	}

	if err := os.MkdirAll(getDatabaseDirPath(dataDir), os.ModePerm); err != nil {
		return err
	}

	if err := writeGenesisToDisk(getGenesisJsonFilePath(dataDir)); err != nil {
		return err
	}

	if err := writeEmptyBlocksDbToDisk(getBlocksDbFilePath(dataDir)); err != nil {
		return err
	}

/* 	if err := writeEmptyTxDbToDisk(getBlocksDbFilePath(dataDir)); err != nil {
		return err
	} */

	return nil
}

func getDatabaseDirPath(dataDir string) string {
	return filepath.Join(dataDir, "database")
}

func getGenesisJsonFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "genesis.json")
}

func getBlocksDbFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "block.db")
}

/* func getTxDbFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "tx.db")
} */

func fileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func writeEmptyBlocksDbToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(""), os.ModePerm)
}

/* func loadTransactions(path string) (txs, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return node.TxsListRes{}, err
	}

	var loadedTxs txs
	err = json.Unmarshal(content, &loadedGenesis)
	if err != nil {
		return genesis{}, err
	}

	return loadedGenesis, nil
}

func writeEmptyTxDbToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(""), os.ModePerm)
} */