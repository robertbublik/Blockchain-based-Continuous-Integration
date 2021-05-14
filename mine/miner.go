package mine

import (
	"fmt"
	"path/filepath"
	"github.com/robertbublik/bci/database"
	//"github.com/robertbublik/bci/node"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	//"math/rand"
	"time"
)

const rootDir = "/tmp/gitdir"

type PendingBlock struct {
	index 		uint64
	parent 		database.Hash
	repository	string
	commit 		string
	prevCommit 	string
	time   		uint64
	miner  		database.Account
	tx    		database.Tx
}

func NewPendingBlock(index uint64, parent database.Hash, repository string, commit string, prevCommit string, miner database.Account, tx database.Tx) PendingBlock {
	return PendingBlock{index, parent, repository, commit, prevCommit, uint64(time.Now().Unix()), miner, tx}
}
/*
func Mine2(ctx context.Context, pb PendingBlock) (database.Block, error) {
	
	var block database.Block
	block = database.NewBlock(pb.index, pb.parent, pb.repository, pb.commit, pb.prevCommit, pb.time, pb.miner, pb.tx, "", "", "")
	return block, nil

	 if len(pb.txs) == 0 {
		return database.Block{}, fmt.Errorf("mining empty blocks is not allowed")
	}

	start := time.Now()
	attempt := 0
	var block database.Block
	var hash database.Hash
	var nonce uint32

	for !database.IsBlockHashValid(hash) {
		select {
		case <-ctx.Done():
			fmt.Println("Mining cancelled!")

			return database.Block{}, fmt.Errorf("mining cancelled. %s", ctx.Err())
		default:
		}

		attempt++
		nonce = generateNonce()

		if attempt%1000000 == 0 || attempt == 1 {
			fmt.Printf("Mining %d Pending TXs. Attempt: %d\n", len(pb.txs), attempt)
		}

		block = database.NewBlock(pb.parent, pb.index, nonce, pb.time, pb.miner, pb.txs)
		blockHash, err := block.Hash()
		if err != nil {
			return database.Block{}, fmt.Errorf("couldn't mine block. %s", err.Error())
		}

		hash = blockHash
	}

	fmt.Printf("\nMined new Block '%x' using PoW🎉🎉🎉%s:\n", hash, fs.Unicode("\\U1F389"))
	fmt.Printf("\tHeight: '%v'\n", block.Header.Index)
	fmt.Printf("\tNonce: '%v'\n", block.Header.Nonce)
	fmt.Printf("\tCreated: '%v'\n", block.Header.Time)
	fmt.Printf("\tMiner: '%v'\n", block.Header.Miner)
	fmt.Printf("\tParent: '%v'\n\n", block.Header.Parent.Hex())

	fmt.Printf("\tAttempt: '%v'\n", attempt)
	fmt.Printf("\tTime: %s\n\n", time.Since(start))

	return block, nil 
}*/


func Mine(tx database.Tx) {
	// Clone the given repository to the given directory
	dir := checkoutRepository(tx.Repository, tx.Commit)

	switch tx.Language {
	case "docker":
		fmt.Println("Docker build")
		DockerBuild(tx, dir)
	default:
		fmt.Println("Unknown language.")
	}

	/* var miningCtx context.Context
	var stopCurrentMining context.CancelFunc

	ticker := time.NewTicker(time.Second * miningIntervalSeconds)

	for {
		select {
		case <-ticker.C:
			go func() {
				if len(n.pendingTXs) > 0 && !n.isMining {
					n.isMining = true

					miningCtx, stopCurrentMining = context.WithCancel(ctx)
					err := n.minePendingTXs(miningCtx)
					if err != nil {
						fmt.Printf("ERROR: %s\n", err)
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
	} */
}
/*
func (n *Node) minePendingTXs(ctx context.Context) error {
	blockToMine := NewPendingBlock(
		n.state.LatestBlockHash(),
		n.state.NextBlockNumber(),
		n.info.Account,
		n.getPendingTXsAsArray(),
	)

	minedBlock, err := Mine(ctx, blockToMine)
	if err != nil {
		return err
	}

	n.removeMinedPendingTXs(minedBlock)

	_, err = n.state.AddBlock(minedBlock)
	if err != nil {
		return err
	}

	return nil
}


func (n *Node) removeMinedPendingTXs(block database.Block) {
	txHash, _ := block.Body.TX.Hash()
	if _, exists := n.pendingTXs[txHash.Hex()]; exists {
		fmt.Printf("\t-archiving mined TX: %s\n", txHash.Hex())

		n.archivedTXs[txHash.Hex()] = block.Body.TX
		delete(n.pendingTXs, txHash.Hex())
	}
} */

func checkoutRepository(repository string, commit string) string {
	dir := filepath.Join(rootDir, repository)
	Info("git clone %s %s", repository, dir)
	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: repository,
	})

	CheckIfError(err)
	
	// ... retrieving the commit being pointed by HEAD
	Info("git show-ref --head HEAD")
	ref, err := r.Head()
	CheckIfError(err)
	fmt.Println(ref.Hash())

	w, err := r.Worktree()
	CheckIfError(err)
	if commit != "" {
		// ... checking out to commit
		Info("git checkout %s", commit)
		err = w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(commit),
		})
		CheckIfError(err)

		// ... retrieving the commit being pointed by HEAD, it shows that the
		// repository is pointing to the giving commit in detached mode
		Info("git show-ref --head HEAD")
		ref, err = r.Head()
		CheckIfError(err)
		fmt.Println(ref.Hash())
	}
	return dir
}