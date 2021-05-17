package node

import (
	"fmt"
	"context"
	"path/filepath"
	"github.com/robertbublik/bci/database"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"math/rand"
	"time"
	"strings"
	"os"
)

var registryUrl = "localhost:5000/"

type PendingBlock struct {
	index 		uint64
	parent 		database.Hash
	time   		uint64
	miner  		database.Account
	txs    		[]database.Tx
}

func NewPendingBlock(index uint64, parent database.Hash, miner database.Account, txs []database.Tx) PendingBlock {
	return PendingBlock{index, parent, uint64(time.Now().Unix()), miner, txs}
}

func Mine(ctx context.Context, pb PendingBlock, dataDir string) (database.Block, error) {
	if len(pb.txs) == 0 {
		return database.Block{}, fmt.Errorf("mining empty blocks is not allowed")
	}
	
	var block 	database.Block
	//var hash 	database.Hash
	var tx 		database.Tx
	var url 	string
	var done 	bool = false

	randomIndex := rand.Intn(len(pb.txs))
	tx = pb.txs[randomIndex]

	for !done {
		lastIndex := strings.LastIndex(tx.Repository, "/")
		repoName := strings.ToLower(tx.Repository[lastIndex + 1:])
		
		checkoutDir := filepath.Join(dataDir, repoName)

		// Clone the given repository to the given directory
		checkoutRepository(tx, checkoutDir)

		switch tx.Language {
		case "docker":
			fmt.Println("Docker build")
			dockerfilePath := filepath.Join(checkoutDir, "Dockerfile")
			url = registryUrl + repoName
			DockerBuildAndPush(ctx, tx, dockerfilePath, url)
		default:
			fmt.Println("Unknown language.")
		}
		block = database.NewBlock(pb.index, pb.parent, tx.Repository, tx.Commit, tx.PrevCommit, pb.time, pb.miner, tx, url)
		done = true
	}
	return block, nil
}

func (n *Node) removeMinedPendingTXs(block database.Block) {
	txHash, _ := block.Body.TX.Hash()
	if _, exists := n.pendingTXs[txHash.Hex()]; exists {
		fmt.Printf("\t-archiving mined TX: %s\n", txHash.Hex())

		n.archivedTXs[txHash.Hex()] = block.Body.TX
		delete(n.pendingTXs, txHash.Hex())
	}
}

func checkoutRepository(tx database.Tx, dir string) {
	
	Info("git clone %s\n %s", tx.Repository, dir)
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		fmt.Printf("Repository already cloned, pulling\n")
		r, err := git.PlainOpen(dir)
		CheckIfError(err)

		// Get the working directory for the repository
		w, err := r.Worktree()
		CheckIfError(err)

		// Pull the latest changes from the origin remote and merge into the current branch
		Info("git pull origin")
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		CheckIfError(err)

		// Print the latest commit that was just pulled
		ref, err := r.Head()
		CheckIfError(err)
		commit, err := r.CommitObject(ref.Hash())
		CheckIfError(err)

		fmt.Println(commit)
		return
	}
	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: tx.Repository,
	})

	CheckIfError(err)
	
	// ... retrieving the commit being pointed by HEAD
	Info("git show-ref --head HEAD")
	ref, err := r.Head()
	CheckIfError(err)
	fmt.Println(ref.Hash())

	w, err := r.Worktree()
	CheckIfError(err)
	if tx.Commit != "" {
		// ... checking out to commit
		Info("git checkout %s", tx.Commit)
		err = w.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(tx.Commit),})
		CheckIfError(err)

		// ... retrieving the commit being pointed by HEAD, it shows that the
		// repository is pointing to the giving commit in detached mode
		Info("git show-ref --head HEAD")
		ref, err = r.Head()
		CheckIfError(err)
		fmt.Println(ref.Hash())
	}
}