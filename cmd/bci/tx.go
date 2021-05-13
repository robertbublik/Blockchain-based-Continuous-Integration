package main

import (
	"fmt"
	"net/http"
	"github.com/spf13/cobra"
	"github.com/robertbublik/bci/node"
	"os"
	"encoding/json"
	"bytes"
	"time"
)

func txCmd() *cobra.Command {
	var txCmd = &cobra.Command{
		Use:   "tx",
		Short: "Interact with transactions (add, list...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	txCmd.AddCommand(txListCmd())
	txCmd.AddCommand(txAddCmd())
	txCmd.AddCommand(txMineCmd())

	return txCmd
}

func txListCmd() *cobra.Command {
	var txListCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists all transactions in the BCI.",
		Run: func(cmd *cobra.Command, args []string) {
			res, err := http.Get("http://localhost:8080/tx/list")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			txsListRes := node.TxsListRes{}
			err = node.ReadRes(res, &txsListRes)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			b, err := json.MarshalIndent(txsListRes, "", "\t")
			fmt.Println(string(b))
			
		},
	}

	return txListCmd
}

func txAddCmd() *cobra.Command {
	var txAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a transaction to the BCI.",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(flagFrom)
			value, _ := cmd.Flags().GetUint64(flagValue)
			repository, _ := cmd.Flags().GetString(flagRepository)
			language, _ := cmd.Flags().GetString(flagLanguage)
			commit, _ := cmd.Flags().GetString(flagCommit)
			prevCommit, _ := cmd.Flags().GetString(flagPrevCommit)

			tx := node.TxReq{from, value, repository, language, commit, prevCommit, uint64(time.Now().Unix())}
			payloadBytes, err := json.Marshal(tx)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			body := bytes.NewReader(payloadBytes)
			req, err := http.NewRequest("POST", "http://127.0.0.1:8080/tx/add", body)
			if err != nil {
				fmt.Printf("1")
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("2")
				fmt.Println(err)
				os.Exit(1)
			}

			printResponse(resp)

			defer resp.Body.Close() */
		},
	}
	addDefaultStringRequiredFlags(txAddCmd, flagFrom, "", "Account name of developer submitting a transaction")
	addDefaultUint64RequiredFlags(txAddCmd, flagValue, 0, "Value of the transaction's reward")
	addDefaultStringRequiredFlags(txAddCmd, flagRepository, "", "URL of the repository")
	addDefaultStringRequiredFlags(txAddCmd, flagLanguage, "", "Java, Docker")
	addDefaultStringRequiredFlags(txAddCmd, flagCommit, "", "Requested commit hash of repository to checkout")
	txAddCmd.Flags().String(flagPrevCommit, "", "Link transaction to an earlier build in the BCI through the previously used commit hash")
	return txAddCmd
}

func txMineCmd() *cobra.Command {
	var txMineCmd = &cobra.Command{
		Use:   "mine",
		Short: "Mine a transaction in the BCI.",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString(flagId)

			txReq := node.TxMineReq{id}
			payloadBytes, err := json.Marshal(txReq)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			body := bytes.NewReader(payloadBytes)
			req, err := http.NewRequest("POST", "http://127.0.0.1:8080/tx/mine", body)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer resp.Body.Close()
		},
	}
	addDefaultStringRequiredFlags(txMineCmd, flagId, "", "Id of transaction to be mined")
	return txMineCmd
}