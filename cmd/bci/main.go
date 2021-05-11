package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/robertbublik/bci/fs"
	"os"
)

const flagDataDir = "datadir"
const flagAccount = "account"
const flagIP = "ip"
const flagPort = "port"

func main() {
	var bciCmd = &cobra.Command{
		Use:   "bci",
		Short: "Blockchain-based Continuous Integration CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	bciCmd.AddCommand(versionCmd)
	bciCmd.AddCommand(runCmd())
	bciCmd.AddCommand(balancesCmd())

	err := bciCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagDataDir, "", "Absolute path to the node data dir where the DB will/is stored")
	cmd.MarkFlagRequired(flagDataDir)
}

func getDataDirFromCmd(cmd *cobra.Command) string {
	dataDir, _ := cmd.Flags().GetString(flagDataDir)

	return fs.ExpandPath(dataDir)
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
