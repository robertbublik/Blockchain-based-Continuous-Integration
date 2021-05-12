package main

import (
/* 	"fmt"
	"net/http" */
	"github.com/spf13/cobra"
	"github.com/robertbublik/bci/node"
/* 	"os" */
	//"encoding/json"
)

func statusCmd() *cobra.Command {
	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Displays status of BCI.",
		Run: func(cmd *cobra.Command, args []string) {
			/* res, err := http.Get("http://localhost:8080/node/status")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			} */
			tools.Test()
			//txListRes := node.TxsListRes{}
			
			/* err = utils.ReadRes(res, &txListRes)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			json.MarshalIndent(txListRes, "", "\t") */
			
		},
	}

	return statusCmd
}

/* func HandleRequest(w http.ResponseWriter, req *http.Request) {
    body := App.MustReadBody(req, w)
    if body == nil {
        return
    }

    var prettyJSON bytes.Buffer
    error := json.Indent(&prettyJSON, body, "", "\t")
    if error != nil {
        log.Println("JSON parse error: ", error)
        App.BadRequest(w)
        return
    }

    log.Println(string(prettyJSON.Bytes()))
} */