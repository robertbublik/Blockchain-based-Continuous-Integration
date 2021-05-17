package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/robertbublik/bci/fs"
	"github.com/robertbublik/bci/node"
	"os"
	"net/http"
)

const flagDataDir = "datadir"
const flagAccount = "account"
const flagIP = "ip"
const flagPort = "port"
const flagFrom = "from"
const flagTo = "to"
const flagValue = "value"
const flagRepository = "repository"
const flagLanguage = "language"
const flagCommit = "commit"
const flagPrevCommit = "prevCommit"
const flagId = "id"

func main() {
	var bciCmd = &cobra.Command{
		Use:   "bci",
		Short: "Blockchain-based Continuous Integration CLI",
		Run: func(cmd *cobra.Command, args []string) {
			
		},
	}

	bciCmd.AddCommand(runCmd())
	bciCmd.AddCommand(balancesCmd())
	bciCmd.AddCommand(statusCmd())
	bciCmd.AddCommand(txCmd())

	err := bciCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func addDefaultStringRequiredFlags(cmd *cobra.Command, flag string, defaultValue string, help string) {
	cmd.Flags().String(flag, defaultValue, help)
	cmd.MarkFlagRequired(flag)
}

func addDefaultUint64RequiredFlags(cmd *cobra.Command, flag string, defaultValue uint64, help string) {
	cmd.Flags().Uint64(flag, defaultValue, help)
	cmd.MarkFlagRequired(flag)
}

func getDataDirFromCmd(cmd *cobra.Command) string {
	dataDir, _ := cmd.Flags().GetString(flagDataDir)

	return fs.ExpandPath(dataDir)
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}

func connectionErr() error {
	return fmt.Errorf("connection error")
}

func printResponse(r *http.Response) {
	responseGeneric := map[string]interface{}{}
	node.ReadRes(r, &responseGeneric)
	for key,value := range responseGeneric {
		switch key {
		case "response":
			fmt.Println(value)
		case "error":
			fmt.Println(value)
		default:
			fmt.Printf("Unknown response")
		}
	}
}