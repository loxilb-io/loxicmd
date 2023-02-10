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
	"errors"
	"fmt"
	"loxicmd/pkg/api"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type CreateMirrorOptions struct {
	MirrID    string
	MirrInfo  []string
	TargerObj []string
}

func NewCreateMirrorCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := CreateMirrorOptions{}

	var createmirrorCmd = &cobra.Command{
		Use:   "mirror <mirrorIdent> --mirrorInfo=<InfoOption>:<InfoValue>,... --targetObject=attachement:<port1,rule2>,mirrObjName:<ObjectName>",
		Short: "Create a Mirror",
		Long: `Create a Mirror using LoxiLB
--<infoOption>s of mirrorInfo
type(int) : Mirroring type as like 0 == SPAN, 1 == RSPAN, 2 == ERSPAN 
port(string) : The port where mirrored traffic needs to be sent
vlan(int) : for RSPAN we may need to send tagged mirror traffic
remoteIP(string) : For ERSPAN we may need to send tunnelled mirror traffic
sourceIP(string): For ERSPAN we may need to send tunnelled mirror traffic
tunnelID(int): For ERSPAN we may need to send tunnelled mirror traffic


ex) loxicmd create mirror mirr-1 --mirrorInfo="type:0,port:hs0" --targetObject="attachement:1,mirrObjName:hs1"

`,
		Aliases: []string{"mirror", "mirr", "mirrors"},
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var mirrorMods api.MirrMod
			// Make mirrorMod
			if err := ReadCreateMirrorOptions(&mirrorMods, args); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if err := GetMirrorInfoPairList(&mirrorMods, o.MirrInfo); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if err := GetTargetObjPairList(&mirrorMods, o.TargerObj); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			resp, err := MirrorAPICall(restOptions, mirrorMods)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			defer resp.Body.Close()

			fmt.Printf("Debug: response.StatusCode: %d\n", resp.StatusCode)
			if resp.StatusCode == http.StatusOK {
				PrintCreateResult(resp, *restOptions)
				return
			}

		},
	}

	createmirrorCmd.Flags().StringSliceVar(&o.MirrInfo, "mirrorInfo", o.MirrInfo, "Information about the mirror")
	createmirrorCmd.Flags().StringSliceVar(&o.TargerObj, "targetObject", o.TargerObj, "Information about object to which mirror needs to be attached")

	return createmirrorCmd
}

func ReadCreateMirrorOptions(o *api.MirrMod, args []string) error {
	if len(args) > 2 {
		return errors.New("create Mirror command get so many args")
	} else if len(args) == 0 {
		return errors.New("create Mirror need <MirrID> arg")
	}
	o.Ident = args[0]

	return nil
}

func GetMirrorInfoPairList(o *api.MirrMod, MirrInfo []string) error {
	for _, mirrorArg := range MirrInfo {
		mirrorArgsPair := strings.Split(mirrorArg, ":")
		if len(mirrorArgsPair) != 2 {
			return fmt.Errorf("mirrorArgs '%s' is invalid format", MirrInfo)
		} else if mirrorArgsPair[0] == "type" {
			mirrorType, err := strconv.Atoi(mirrorArgsPair[1])
			if err != nil {
				return err
			}
			o.Info.MirrType = mirrorType
		} else if mirrorArgsPair[0] == "port" {
			o.Info.MirrPort = mirrorArgsPair[1]
		} else if mirrorArgsPair[0] == "vlan" {
			mirrorType, err := strconv.Atoi(mirrorArgsPair[1])
			if err != nil {
				return err
			}
			o.Info.MirrVlan = mirrorType
		} else if mirrorArgsPair[0] == "remoteIP" {
			o.Info.MirrRip = mirrorArgsPair[1]
		} else if mirrorArgsPair[0] == "sourceIP" {
			o.Info.MirrSip = mirrorArgsPair[1]
		} else if mirrorArgsPair[0] == "tunnelID" {
			tunnelID, err := strconv.Atoi(mirrorArgsPair[1])
			if err != nil {
				return err
			}
			o.Info.MirrTid = tunnelID
		}

	}
	return nil
}

func GetTargetObjPairList(o *api.MirrMod, TargerObj []string) error {
	for _, mirrorArg := range TargerObj {
		mirrorArgsPair := strings.Split(mirrorArg, ":")
		if len(mirrorArgsPair) != 2 {
			return fmt.Errorf("TargerObj '%s' is invalid format", TargerObj)
		} else if mirrorArgsPair[0] == "mirrObjName" {
			o.Target.MirrObjName = mirrorArgsPair[1]
		} else if mirrorArgsPair[0] == "attachement" {
			attachement, err := strconv.Atoi(mirrorArgsPair[1])
			if err != nil {
				return err
			}
			o.Target.AttachMent = api.MirrObjType(attachement)
		}

	}
	return nil
}

func MirrorAPICall(restOptions *api.RESTOptions, mirrorModel api.MirrMod) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.Mirror().Create(ctx, mirrorModel)
}
