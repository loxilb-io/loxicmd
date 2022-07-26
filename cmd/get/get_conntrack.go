// Author: Inho gog <inhogog2@netlox.io>
package get

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"loxicmd/pkg/api"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

func NewGetConntrackCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetctCmd = &cobra.Command{
		Use:     "conntrack",
		Aliases: []string{"ct", "conntracks", "cts"},
		Short:   "Get a Conntrack",
		Long:    `It shows connection track Information`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd
			_ = args
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			resp, err := client.Conntrack().Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetCTResult(resp, *restOptions)
				return
			}

		},
	}

	return GetctCmd
}

func PrintGetCTResult(resp *http.Response, o api.RESTOptions) {
	ctresp := api.CtInformationGet{}
	var data [][]string
	resultByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &ctresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(ctresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	// Table Init
	table := TableInit()
	table.SetHeader([]string{"destinationIP", "sourceIP", "destinationPort", "sourcePort", "protocol", "state", "act", "packets", "bytes"})
	// Making load balance data
	data = makeConntrackData(ctresp)

	// Rendering the load balance data to table
	TableShow(data, table)
}

func makeConntrackData(ctresp api.CtInformationGet) (data [][]string) {
	for _, conntrack := range ctresp.CtInfo {
		data = append(data, []string{
			conntrack.Dip,
			conntrack.Sip,
			fmt.Sprintf("%d", conntrack.Dport),
			fmt.Sprintf("%d", conntrack.Sport),
			conntrack.Proto,
			conntrack.CState,
			conntrack.CAct,
			fmt.Sprintf("%v", conntrack.Pkts),
			fmt.Sprintf("%v", conntrack.Bytes),
		})
	}
	return data
}
