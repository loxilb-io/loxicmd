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
 package get

 import (
	 "fmt"
	 "net"
	 "strings"
	 "syscall"
 
	 "os"
	 "time"

	 nlp "github.com/vishvananda/netlink"
	 "golang.org/x/sys/unix"
 )
 
 const (
	 IF_OPER_UNKNOWN uint8 = iota
	 IF_OPER_NOTPRESENT
	 IF_OPER_DOWN
	 IF_OPER_LOWERLAYERDOWN
	 IF_OPER_TESTING
	 IF_OPER_DORMANT
	 IF_OPER_UP
 )
 
 func dump(format string, a ...any) {
	 data := fmt.Sprintf(format, a...)
	 _, err := f.WriteString(data)
	 if err != nil {
		 fmt.Println(err)
	 }
 }
 
 func AddLink(link nlp.Link) int {
	 var ifMac [6]byte
	 var ret int
 
	 attrs := link.Attrs()
	 name := attrs.Name
 
	 if len(attrs.HardwareAddr) > 0 {
		 copy(ifMac[:], attrs.HardwareAddr[:6])
	 }
 
	 mtu := attrs.MTU
	 state := uint8(attrs.OperState) == IF_OPER_UP
 
	 //dump("Device %v mac(%v) info recvd\n", name, ifMac)
 
	 if _, ok := link.(*nlp.Bridge); ok {
 
		 dump("ip link add %v type bridge\n", name)
		 if state {
			 dump("ip link set %v up\n", name)
		 }
		 return 0
	 }
 
	 /* Tagged Vlan port */
	 if strings.Contains(name, ".") {
		 pname := strings.Split(name, ".")
		 dump("ip link add %v link %v type vlan id %v\n", name, pname[0], pname[1])
		 if state {
			 dump("ip link set %v up\n", pname[0])
		 }
 
	 } else {
		 /* Physical port/ Bond/ VxLAN */
		 tunId := 0
		 if vxlan, ok := link.(*nlp.Vxlan); ok {
			 tunId = vxlan.VxlanId
			 uif, err := nlp.LinkByIndex(vxlan.VtepDevIndex)
			 if err != nil {
				 fmt.Println(err)
				 return -1
			 }
			 real := uif.Attrs().Name
			 dump("ip link add %v type vxlan id %d local %s dev %v dstport 4789\n", name, tunId, vxlan.SrcAddr.String(), real)
			 dump("ip link set %v up\n", name)
		 } else {
			 if mtu != 1500 {
				 dump("ip link set %v mtu %d\n", name, mtu)
			 }
		 }
	 }
 
	 if state {
		 dump("ip link set %v up\n", name)
	 }
 
	 /* Untagged vlan ports */
	 if attrs.MasterIndex > 0 {
		 br, err := nlp.LinkByIndex(attrs.MasterIndex)
		 if err != nil {
			 fmt.Println(err)
			 return -1
		 }
 
		 dump("ip link set %s master %s\n", name, br.Attrs().Name)
	 }
	 return ret
 }
 
 func AddAddr(addr nlp.Addr, link nlp.Link) int {
	 var ret int
 
	 attrs := link.Attrs()
	 name := attrs.Name
	 ipStr := (addr.IPNet).String()
 
	 dump("ip addr add %v dev %v\n", ipStr, name)
 
	 return ret
 }
 
 func AddNeigh(neigh nlp.Neigh, link nlp.Link) int {
	 var ret int
	 var mac [6]byte
	 var brMac [6]byte
	 var dst net.IP
 
	 attrs := link.Attrs()
	 name := attrs.Name
	 if neigh.Flags != unix.NUD_PERMANENT {
		 return -1
	 }
	 if len(neigh.HardwareAddr) == 0 {
		 return -1
	 }
	 copy(mac[:], neigh.HardwareAddr[:6])
 
	 if neigh.Family == unix.AF_INET {
		 dump("ip neigh add %v lladdr %v dev %v\n", neigh.IP.String(), mac, name)
 
	 } else if neigh.Family == unix.AF_BRIDGE {
 
		 if len(neigh.HardwareAddr) == 0 {
			 return -1
		 }
		 copy(mac[:], neigh.HardwareAddr[:6])
 
		 if neigh.Vlan == 1 {
			 /*FDB comes with vlan 1 also */
			 return 0
		 }
 
		 if mac[0]&0x01 == 1 {
			 /* Multicast MAC or ZERO address --- IGNORED */
			 return 0
		 }
 
		 if neigh.MasterIndex > 0 {
			 brLink, err := nlp.LinkByIndex(neigh.MasterIndex)
			 if err != nil {
				 fmt.Println(err)
				 return -1
			 }
 
			 copy(brMac[:], brLink.Attrs().HardwareAddr[:6])
			 if mac == brMac {
				 /*Same as bridge mac --- IGNORED */
				 return 0
			 }
		 }
 
		 if _, ok := link.(*nlp.Vxlan); ok {
			 /* Interested in only VxLAN FDB */
			 if len(neigh.IP) > 0 && (neigh.MasterIndex == 0) {
				 dst = neigh.IP
			 } else {
				 return 0
			 }
			 dump("bridge fdb append %v dst %v dev %v\n", mac[:], dst.String(), name)
		 } else {
			 dump("bridge fdb add %v dev %v\n", mac[:], name)
		 }
	 }
 
	 return ret
 
 }
 
 func AddRoute(route nlp.Route) int {
	 var ipNet net.IPNet
	 if route.Dst == nil {
		 r := net.IPv4(0, 0, 0, 0)
		 m := net.CIDRMask(0, 32)
		 r = r.Mask(m)
		 ipNet = net.IPNet{IP: r, Mask: m}
	 } else {
		 ipNet = *route.Dst
	 }
 
	 dump("ip route add %s via %s proto static\n", ipNet.String(), route.Gw.String())
	 return 0
 }
 
 func GetL3Config(link nlp.Link) {
 
	 addrs, err := nlp.AddrList(link, nlp.FAMILY_V4)
	 if err != nil {
		 dump("Error getting address list %v for intf %s\n",
			 err, link.Attrs().Name)
	 }
 
	 if len(addrs) > 0 {
		 for _, addr := range addrs {
			 AddAddr(addr, link)
		 }
	 }
 
	 neighs, err := nlp.NeighList(link.Attrs().Index, nlp.FAMILY_ALL)
	 if err != nil {
		 dump("Error getting neighbors list %v for intf %s\n",
			 err, link.Attrs().Name)
	 }
 
	 if len(neighs) > 0 {
		 for _, neigh := range neighs {
			 AddNeigh(neigh, link)
		 }
	 }
 }
 
 func GetFdbs(link nlp.Link) {
	 /* Get FDBs */
	 if link.Attrs().MasterIndex > 0 {
		 neighs, err := nlp.NeighList(link.Attrs().Index, unix.AF_BRIDGE)
		 if err != nil {
			 dump("Error getting neighbors list %v for intf %s\n",
				 err, link.Attrs().Name)
		 }
 
		 if len(neighs) > 0 {
			 for _, neigh := range neighs {
				 AddNeigh(neigh, link)
			 }
 
		 } else {
			 dump("No FDBs found for intf %s\n", link.Attrs().Name)
		 }
	 }
 }
 
 func GetBridges() {
	 links, err := nlp.LinkList()
	 if err != nil {
		 return
	 }
	 for _, link := range links {
		 switch link.(type) {
		 case *nlp.Bridge:
			 {
				 AddLink(link)
				 GetL3Config(link)
			 }
		 }
	 }
 }
 
 var f *os.File
 
 func Nlpdump() string {
	 var ret int
	 var err error
	 fileP := []string{"ipconfig_", ".txt"}
	 t := time.Now()
	 file := strings.Join(fileP, t.Local().Format("2006-01-02_15:04:05"))
	 f, err = os.Create(file)
	 if err != nil {
		 fmt.Printf("Can't create dump file\n")
		 os.Exit(1)
	 }
 
	 defer f.Close()
 
	 /*Get bridge info first */
	 GetBridges()
 
	 /*Get other link info */
	 links, err := nlp.LinkList()
	 if err != nil {
		 //dump("Error in getting device info(%v)\n", err)
		 ret = -1
	 }
 
	 for _, link := range links {
		 if _, ok := link.(*nlp.Bridge); !ok {
			 ret = AddLink(link)
 
			 if ret == -1 {
				 continue
			 }
			 GetFdbs(link)
			 GetL3Config(link)
		 }
	 }
 
	 /* Get Static Routes only */
	 rFilter := nlp.Route{Protocol: syscall.RTPROT_STATIC}
	 routes, err := nlp.RouteListFiltered(nlp.FAMILY_V4, &rFilter, nlp.RT_FILTER_PROTOCOL)
 
	 if err != nil {
		 dump("Error getting route list %v\n", err)
	 }
 
	 if len(routes) > 0 {
		 for _, route := range routes {
			 AddRoute(route)
		 }
	 }
	 return file
 }
 