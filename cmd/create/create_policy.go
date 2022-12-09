/*
 * Copyright (c) 2022 NetLOX Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package create

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"loxicmd/pkg/api"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type CreatePolicyOptions struct {
	Ident   string
	Rate    string
	Block   string
	Target  string
	Color   bool
	PolType int
}

type CreatePolicyResult struct {
	Result string `json:"result"`
}

func ReadCreatePolicyOptions(o *CreatePolicyOptions, args []string) error {
	if len(args) > 1 {
		fmt.Println("create Pol command get so many args")
		fmt.Println(args)
	} else if len(args) <= 0 {
		return errors.New("create Pol need Ident args")
	}
	o.Ident = args[0]
	return nil
}

func GetRatePair(body *api.PolMod, RateBlock string) error {
	RatePair := strings.Split(RateBlock, ":")
	if len(RatePair) != 2 {
		return errors.New("Lots ")
	}
	Peak, err := strconv.Atoi(RatePair[0])
	if err != nil {
		return fmt.Errorf("Peak '%s' is not integer", RatePair[0])
	}

	Commited, err := strconv.Atoi(RatePair[1])
	if err != nil {
		return fmt.Errorf("Commited '%s' is not integer", RatePair[1])
	}
	body.Info.CommittedInfoRate = uint64(Commited)
	body.Info.PeakInfoRate = uint64(Peak)
	return nil
}

func GetBlockPair(body *api.PolMod, Block string) error {
	BlockPair := strings.Split(Block, ":")
	if len(BlockPair) != 2 {
		return errors.New("error")
	}
	Excess, err := strconv.Atoi(BlockPair[0])
	if err != nil {
		return fmt.Errorf("Excess '%s' is not integer", BlockPair[0])
	}

	Commited, err := strconv.Atoi(BlockPair[1])
	if err != nil {
		return fmt.Errorf("Commited '%s' is not integer", BlockPair[1])
	}
	body.Info.ExcessBlkSize = uint64(Excess)
	body.Info.CommittedBlkSize = uint64(Commited)
	return nil
}

func GetTargetPair(body *api.PolMod, Block string) error {
	BlockPair := strings.Split(Block, ":")
	if len(BlockPair) != 2 {
		return errors.New("error")
	}

	AttachMent, err := strconv.Atoi(BlockPair[1])
	if err != nil {
		return fmt.Errorf("AttachMent '%s' is not integer", BlockPair[1])
	}
	body.Target.PolObjName = BlockPair[0]
	body.Target.AttachMent = api.PolObjType(AttachMent)
	return nil
}

func NewCreatePolicyCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := CreatePolicyOptions{}

	var createPolCmd = &cobra.Command{
		Use:   "policy IDENT --rate=<Peak>:<Commited> --target=<ObjectName>:<Attachment> [--block-size=<Excess>:<Committed>] [--color] [--pol-type=<policy type>]",
		Short: "Create a Policy",
		Long: `Create a Policy 
Ex) loxicmd create policy pol-hs0 --rate=100:100 --target=hs0:1
    loxicmd create policy pol-hs1 --rate=100:100 --target=hs0:1 --block-size=12000:6000
    loxicmd create policy pol-hs1 --rate=100:100 --target=hs0:1 --color
    loxicmd create policy pol-hs1 --rate=100:100 --target=hs0:1 --color --pol-type 0

rate unit : Mbps
block-size unit : bps
Policy type(pol-type) 0 : TrTCM,  1 : SrTCM

	`,
		Aliases: []string{"pol", "policys", "pols", "polices"},

		Run: func(cmd *cobra.Command, args []string) {
			if err := ReadCreatePolicyOptions(&o, args); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			// Make body
			body := api.PolMod{}

			body.Ident = o.Ident
			if err := GetRatePair(&body, o.Rate); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if err := GetBlockPair(&body, o.Block); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}

			if err := GetTargetPair(&body, o.Target); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			body.Info.ColorAware = o.Color
			body.Info.PolType = o.PolType

			resp, err := PolicyAPICall(restOptions, body)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			defer resp.Body.Close()

			fmt.Printf("Debug: response.StatusCode: %d\n", resp.StatusCode)
			if resp.StatusCode != http.StatusOK {
				PrintCreatePolResult(resp, *restOptions)
				return
			}

		},
	}

	createPolCmd.Flags().StringVar(&o.Rate, "rate", o.Rate, "Rate pairs can be specified as '<Peak>:<Commited>'")
	createPolCmd.Flags().StringVar(&o.Block, "block-size", o.Block, "Block Size pairs can be specified as '<Excess>:<Committed>'")
	createPolCmd.Flags().StringVar(&o.Target, "target", o.Target, "Target Interface pairs can be specified as '<ObjectName>:<Attachment>'")
	createPolCmd.Flags().BoolVarP(&o.Color, "color", "", false, "Policy color enbale or not")
	createPolCmd.Flags().IntVar(&o.PolType, "pol-type", o.PolType, "Target Interface pairs can be specified as '<ObjectName>:<Attachment>'")

	return createPolCmd
}

func PrintCreatePolResult(resp *http.Response, o api.RESTOptions) {
	result := CreatePolicyResult{}
	resultByte, err := io.ReadAll(resp.Body)
	fmt.Printf("Debug: response.Body: %s\n", string(resultByte))

	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}
	if err := json.Unmarshal(resultByte, &result); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(resp.Body, "", "\t")
		fmt.Println(string(resultIndent))
		return
	}

	fmt.Printf("%s\n", result.Result)
}

func PolicyAPICall(restOptions *api.RESTOptions, PolModel api.PolMod) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.Policy().Create(ctx, PolModel)
}
