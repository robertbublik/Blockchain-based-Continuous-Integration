package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/robertbublik/bci/database"
	"github.com/robertbublik/bci/node"
	"os"
)

func runCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Launches the BCI node and its HTTP API.",
		Run: func(cmd *cobra.Command, args []string) {
			account, _ := cmd.Flags().GetString(flagAccount)
			ip, _ := cmd.Flags().GetString(flagIP)
			port, _ := cmd.Flags().GetUint64(flagPort)

			fmt.Println("Launching BCI node and its HTTP API...")

			bootstrap := node.NewPeerNode(
				"127.0.0.1",
				8070,
				true,
				database.NewAccount("bootstrap"),
				false,
			)

			n := node.New(getDataDirFromCmd(cmd), ip, port, database.NewAccount(account), bootstrap)
			err := n.Run(context.Background())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	addDefaultRequiredFlags(runCmd)
	runCmd.Flags().String(flagMiner, node.DefaultMiner, "miner account of this node to receive block rewards")
	runCmd.Flags().String(flagIP, node.DefaultIP, "exposed IP for communication with peers")
	runCmd.Flags().Uint64(flagPort, node.DefaultHTTPort, "exposed HTTP port for communication with peers")

	return runCmd
}
