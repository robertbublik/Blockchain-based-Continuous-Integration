package main

import (
	"fmt"
	"net/http"
	"github.com/spf13/cobra"
	"github.com/robertbublik/bci/node"
	"os"
	"encoding/json"
	"io/ioutil"
)

func statusCmd() *cobra.Command {
	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Displays status of BCI.",
		Run: func(cmd *cobra.Command, args []string) {
			res, err := http.Get("http://localhost:8080/node/status")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			statusRes := node.StatusRes{}
			err = node.ReadRes(res, &statusRes)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			b, err := json.MarshalIndent(statusRes, "", "\t")
			fmt.Println(string(b))
			
		},
	}

	return statusCmd
}
