// Author: Inho gog <inhogog2@netlox.io>
package get

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"loxicmd/pkg/api"
	"net/http"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

func NewGetPortCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetPortCmd = &cobra.Command{
		Use:   "port",
		Short: "Get a Port dump",
		Long:  `It shows port dump Information`,
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
			resp, err := client.Port().Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetPortResult(resp, *restOptions)
				return
			}

		},
	}

	return GetPortCmd
}

func PrintGetPortResult(resp *http.Response, o api.RESTOptions) {
	portresp := api.PortGet{}
	var data [][]string
	resultByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &portresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(portresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	// Table Init
	table := TableInit()

	// Sort port Data
	sort.Slice(portresp.Ports, func(i, j int) bool {
		return portresp.Ports[i].PortNo < portresp.Ports[j].PortNo
	})

	// Making Port data
	for _, port := range portresp.Ports {
		if o.PrintOption == "wide" {
			table.SetHeader([]string{"index", "portname", "MAC", "link/state", "mtu", "isActive/bpf", "Statistics", "L3Info", "L2Info", "Sync"})
			data = append(data, []string{fmt.Sprintf("%d", port.PortNo), port.Name, // Default Info
				port.HInfo.MacAddrStr, fmt.Sprintf("%v/%v", port.HInfo.Link, port.HInfo.State), fmt.Sprintf("%d", port.HInfo.Mtu), // HW info
				fmt.Sprintf("%v/%v", port.SInfo.PortActive, port.SInfo.BpfLoaded), // SW info
				fmt.Sprintf("rx/tx byte : %d/%d \nrx/tx packets : %d/%d \nrx/tx error : %d/%d ", // Statistic infor
					port.Stats.RxBytes, port.Stats.TxBytes,
					port.Stats.RxPackets, port.Stats.TxPackets,
					port.Stats.RxError, port.Stats.TxError),
				MakeL3InfoRoString(port.L3), // L3 info
				MakeL2InfoRoString(port.L2), // L2 info
				fmt.Sprintf("%v", port.Sync),
			})
		} else {
			table.SetHeader([]string{"index", "portname", "MAC", "link/state", "L3Info", "L2Info"})
			data = append(data, []string{fmt.Sprintf("%d", port.PortNo), port.Name,
				port.HInfo.MacAddrStr, fmt.Sprintf("%v/%v", port.HInfo.Link, port.HInfo.State),
				MakeL3InfoRoString(port.L3), MakeL2InfoRoString(port.L2)})
		}
	}

	// Rendering the load balance data to table
	TableShow(data, table)
}

func MakeL3InfoRoString(l3 api.PortLayer3Info) (ret string) {
	ret = fmt.Sprintf("Routed: %v\nIPv4 : %s \nIPv6 : %s", l3.Routed, l3.Ipv4_addrs, l3.Ipv6_addrs)
	return ret
}
func MakeL2InfoRoString(l2 api.PortLayer2Info) (ret string) {
	ret = fmt.Sprintf("IsPVID: %v\nVID : %d", l2.IsPvid, l2.Vid)
	return ret
}
