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
package dump

import (
	"errors"
	"fmt"
	"loxicmd/cmd/create"
	"loxicmd/pkg/api"
	"os"

	"gopkg.in/yaml.v2"
)

func ApplyFileConfig(file string, restOptions *api.RESTOptions) error {
	var comm api.TypeMeta
	var err error

	// open file
	byteBuf, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Unmashal to yaml
	if err := yaml.Unmarshal(byteBuf, &comm); err != nil {
		fmt.Printf("Error: Failed to unmarshal File: (%s)\n", err.Error())
		return err
	}
	switch comm.Kind {
	case "Loadbalancer", "lb", "LB":
		err = create.LoadBalancerCreateWithFile(restOptions, byteBuf)
	case "Endpoint", "ep", "endpoints":
		err = create.EndPointCreateWithFile(restOptions, byteBuf)
	case "FDB", "fdb":
		err = create.FDBCreateWithFile(restOptions, byteBuf)
	case "Firewall", "fw", "firewalls", "firewall":
		err = create.FirewallCreateWithFile(restOptions, byteBuf)
	case "ipv4address", "ipv4", "ipaddress", "ip", "IP", "IPaddress":
		err = create.IPv4AddressCreateWithFile(restOptions, byteBuf)
	case "mirror", "mirr", "mirrors", "Mirror":
		err = create.MirrorCreateWithFile(restOptions, byteBuf)
	case "nei", "neigh", "Neighbor", "Neigh", "neighbor":
		err = create.NeighborsCreateWithFile(restOptions, byteBuf)
	case "Policy", "pol", "policys", "pols", "polices":
		err = create.PolicyCreateWithFile(restOptions, byteBuf)
	case "Route", "route":
		err = create.RouteCreateWithFile(restOptions, byteBuf)
	case "Session", "session", "sessions":
		err = create.SessionCreateWithFile(restOptions, byteBuf)
	case "SessionULCL", "ulcl", "sessionulcls", "ulcls", "ULCL":
		err = create.SessionUlClCreateWithFile(restOptions, byteBuf)
	case "VlanMember", "vlanMember", "vlan-member", "vlan_member", "vlanmember":
		err = create.VlanMemberCreateWithFile(restOptions, byteBuf)
	case "Vlan", "vlan":
		err = create.VlanBridgeCreateWithFile(restOptions, byteBuf)
	case "VxlanPeer", "vxlanpeer", "vxlan-peer", "vxlan_peer":
		err = create.VxlanPeerCreateWithFile(restOptions, byteBuf)
	case "Vxlan", "vxlan":
		err = create.VxlanBridgeCreateWithFile(restOptions, byteBuf)
	case "BFD", "bfd":
		err = create.BFDCreateWithFile(restOptions, byteBuf)
	default:
		fmt.Printf("Not Supported\n")
		return errors.New("not supported")
	}
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
