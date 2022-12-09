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

	"errors"
	"os"
	"os/exec"
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

func dump(f *os.File, format string, a ...any) {
	data := fmt.Sprintf(format, a...)
	_, err := f.WriteString(data)
	if err != nil {
		fmt.Println(err)
	}
}

func fileOpen(file string, flag int) *os.File {
	var f *os.File
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		f, err = os.Create(file)
		if err != nil {
			fmt.Printf("Can't create dump file\n")
			return nil
		}
	} else {
		f, err = os.OpenFile(file, flag, 0660)
		if err != nil {
			fmt.Printf("Can't open %s file\n", file)
			return nil
		}
	}
	return f
}

func dumpType(path string, itype string) {
	file := path + "/type"
	f := fileOpen(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE)
	if f == nil {
		return
	}
	dump(f, "%s", itype)
	f.Close()
}

func createDir(path string) int {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("Can't create intf config dir : %s", path)
			return -1
		}
	}
	return 0
}

func dumpMtu(path string, mtu int) {
	file := path + "/mtu"
	f := fileOpen(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE)
	if f == nil {
		return
	}
	dump(f, "%d", mtu)
	f.Close()
}

func dumpBondMode(path string, mode int) {
	file := path + "/mode"
	f := fileOpen(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE)
	if f == nil {
		return
	}
	dump(f, "%d", mode)
	f.Close()
}

func dumpMaster(path string, master string, mtype string) {
	file := path + "/master"
	f := fileOpen(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE)
	if f == nil {
		return
	}
	dump(f, "%s|%s", master, mtype)
	f.Close()
}

func dumpVxlan(name string, id int, local string, uif string) {
	file := path + name + "/info"
	f := fileOpen(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE)
	if f == nil {
		return
	}
	dump(f, "%d|%s|%s", id, local, uif)
	f.Close()
}

func dumpReal(path string, real string, vid string) {
	file := path + "/real"
	f := fileOpen(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE)
	if f == nil {
		return
	}
	dump(f, "%s type vlan id %s", real, vid)
	f.Close()
}

func dumpSubIntf(path string, subintf string, real string, vid string) {
	file := path + real + "/subintf"
	f := fileOpen(file, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	if f == nil {
		return
	}
	dump(f, "%s|%s|%s\n", subintf, real, vid)
	f.Close()
}

func dumpMember(name string, member string) {
	file := path + name + "/members"
	f := fileOpen(file, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	if f == nil {
		return
	}
	dump(f, "%s\n", member)
	f.Close()
}

func dumpIpv4Addr(path string, ip string) {
	file := path + "/ipv4addr"
	f := fileOpen(file, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	if f == nil {
		return
	}
	dump(f, "%s\n", ip)
	f.Close()
}

func dumpIpv4Neigh(path string, ip string, mac [6]byte) {
	file := path + "/ipv4neigh"
	f := fileOpen(file, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	if f == nil {
		return
	}
	dump(f, "%s lladdr %v\n", ip, net.HardwareAddr(mac[:]))
	f.Close()
}

func dumpVxlanFdb(path string, mac [6]byte, dst string) {
	file := path + "/vxfdbs"
	f := fileOpen(file, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	if f == nil {
		return
	}
	dump(f, "%v dst %s\n", net.HardwareAddr(mac[:]), dst)
	f.Close()
}

func dumpL2Fdb(path string, mac [6]byte) {
	file := path + "/l2fdbs"
	f := fileOpen(file, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	if f == nil {
		return
	}
	dump(f, "%v\n", net.HardwareAddr(mac[:]))
	f.Close()
}

func dumpIpv4Route(path string, ip string, nh string) {
	file := path + "/ipv4route"
	f := fileOpen(file, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	if f == nil {
		return
	}
	dump(f, "%s via %s\n", ip, nh)
	f.Close()
}

func AddLink(link nlp.Link) int {
	var ifMac [6]byte
	var ret int
	var intfpath string
	var upath string

	attrs := link.Attrs()
	name := attrs.Name

	if len(attrs.HardwareAddr) > 0 {
		copy(ifMac[:], attrs.HardwareAddr[:6])
	}

	mtu := attrs.MTU
	state := uint8(attrs.OperState) == IF_OPER_UP
	intfpath = path + name

	if _, ok := link.(*nlp.Bridge); ok {

		dump(f, "ip link add %v type bridge\n", name)
		if state {
			dump(f, "ip link set %v up\n", name)
		}

		createDir(intfpath)
		dumpType(intfpath, "bridge")
		return 0
	}

	/* Tagged Vlan port */
	if strings.Contains(name, ".") {
		pname := strings.Split(name, ".")
		dump(f, "ip link add %v link %v type vlan id %v\n", name, pname[0], pname[1])
		if state {
			dump(f, "ip link set %v up\n", pname[0])
		}

		createDir(intfpath)
		dumpType(intfpath, "subintf")
		dumpReal(intfpath, pname[0], pname[1])
		dumpSubIntf(path, name, pname[0], pname[1])
		//createDir(path + pname[0]+"/subintf/"+ name)
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
			dump(f, "ip link add %v type vxlan id %d local %s dev %v dstport 4789\n", name, tunId, vxlan.SrcAddr.String(), real)
			dump(f, "ip link set %v up\n", name)

			//intfpath = path + name
			createDir(intfpath)
			dumpType(intfpath, "vxlan")
			dumpVxlan(name, tunId, vxlan.SrcAddr.String(), real)
			upath = path + real
			createDir(upath)
			if _, ok := uif.(*nlp.Bridge); ok {
				dumpType(upath, "bridge")
			} else {
				dumpType(upath, "phy")
			}
			dumpMaster(upath, name, "vxlan")
		} else if bond, ok := link.(*nlp.Bond); ok {
			dump(f, "ip link add %v type bond\n", name)
			dump(f, "ip link set %v type bond mode %v\n", name, bond.Mode)
			//intfpath = path + name
			createDir(intfpath)
			dumpType(intfpath, "bond")
			dumpBondMode(intfpath, int(bond.Mode))
			if mtu != 1500 {
				dump(f, "ip link set %v mtu %d\n", name, mtu)
				dumpMtu(intfpath, mtu)
			}
		} else {
			//intfpath = path + name
			createDir(intfpath)
			dumpType(intfpath, "phy")
			if mtu != 1500 {
				dump(f, "ip link set %v mtu %d\n", name, mtu)
				dumpMtu(intfpath, mtu)
			}
		}
	}

	if state {
		dump(f, "ip link set %v up\n", name)
	}

	/* Untagged vlan ports */
	if attrs.MasterIndex > 0 {
		master, err := nlp.LinkByIndex(attrs.MasterIndex)
		if err != nil {
			fmt.Println(err)
			return -1
		}

		if _, ok := master.(*nlp.Bridge); ok {
			dumpMaster(intfpath, master.Attrs().Name, "bridge")
			dumpMember(master.Attrs().Name, name)
		} else if _, ok := master.(*nlp.Bond); ok {
			dumpMaster(intfpath, master.Attrs().Name, "bond")
			dumpMember(master.Attrs().Name, name)
			dump(f, "ip link set %s down\n", name)
		}
		dump(f, "ip link set %s master %s\n", name, master.Attrs().Name)
		dump(f, "ip link set %s up\n", name)
	}
	return ret
}

func AddAddr(addr nlp.Addr, link nlp.Link) int {
	var ret int
	var intfpath string

	attrs := link.Attrs()
	name := attrs.Name
	ipStr := (addr.IPNet).String()

	dump(f, "ip addr add %v dev %v\n", ipStr, name)

	intfpath = path + name

	dumpIpv4Addr(intfpath, ipStr)

	return ret
}

func AddNeigh(neigh nlp.Neigh, link nlp.Link) int {
	var ret int
	var mac [6]byte
	var brMac [6]byte
	var dst net.IP
	var intfpath string

	attrs := link.Attrs()
	name := attrs.Name

	if neigh.State&unix.NUD_PERMANENT == 0 {
		return -1
	}
	//fmt.Printf("%v\n", neigh)
	if len(neigh.HardwareAddr) == 0 {
		return -1
	}
	copy(mac[:], neigh.HardwareAddr[:6])

	intfpath = path + name

	if neigh.Family == unix.AF_INET {
		dump(f, "ip neigh add %v lladdr %s dev %v permanent\n", neigh.IP.String(), net.HardwareAddr(mac[:]), name)
		dumpIpv4Neigh(intfpath, neigh.IP.String(), mac)
	} else if neigh.Family == unix.AF_BRIDGE {

		if len(neigh.HardwareAddr) == 0 {
			return -1
		}
		copy(mac[:], neigh.HardwareAddr[:6])

		if neigh.Vlan == 1 {
			return 0
		}

		if mac[0]&0x01 == 1 {
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
				return 0
			}
		}

		if _, ok := link.(*nlp.Vxlan); ok {

			if len(neigh.IP) > 0 /*&& (neigh.MasterIndex == 0)*/ {
				dst = neigh.IP
			} else {
				return 0
			}

			dump(f, "bridge fdb append %v dst %v dev %v permanent\n", net.HardwareAddr(mac[:]), dst.String(), name)
			dumpVxlanFdb(intfpath, mac, dst.String())
		} else {
			dump(f, "bridge fdb add %v dev %v permanent\n", net.HardwareAddr(mac[:]), name)
			dumpL2Fdb(intfpath, mac)
		}
	}

	return ret

}

func AddRoute(route nlp.Route) int {
	var ipNet net.IPNet
	var intfpath string

	if route.Dst == nil {
		r := net.IPv4(0, 0, 0, 0)
		m := net.CIDRMask(0, 32)
		r = r.Mask(m)
		ipNet = net.IPNet{IP: r, Mask: m}
	} else {
		ipNet = *route.Dst
	}

	link, err := nlp.LinkByIndex(route.LinkIndex)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	name := link.Attrs().Name

	intfpath = path + name

	dump(f, "ip route add %s via %s proto static\n", ipNet.String(), route.Gw.String())
	dumpIpv4Route(intfpath, ipNet.String(), route.Gw.String())
	return 0
}

func GetIpAddrs(link nlp.Link) {
	addrs, err := nlp.AddrList(link, nlp.FAMILY_V4)
	if err == nil {
		if len(addrs) > 0 {
			for _, addr := range addrs {
				AddAddr(addr, link)
			}
		}
	}
}

func GetIpNeigh(link nlp.Link) {
	neighs, err := nlp.NeighList(link.Attrs().Index, nlp.FAMILY_ALL)
	if err == nil {
		if len(neighs) > 0 {
			for _, neigh := range neighs {
				AddNeigh(neigh, link)
			}
		}
	}
}

func GetFdbs(link nlp.Link) {
	/* Get FDBs */
	//if link.Attrs().MasterIndex > 0 {
	if true {
		neighs, err := nlp.NeighList(link.Attrs().Index, unix.AF_BRIDGE)
		if err == nil {
			if len(neighs) > 0 {
				for _, neigh := range neighs {
					AddNeigh(neigh, link)
				}
			}
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
				GetIpAddrs(link)
			}
		}
	}
}

func GetBonds() {
	links, err := nlp.LinkList()
	if err != nil {
		return
	}
	for _, link := range links {
		switch link.(type) {
		case *nlp.Bond:
			{
				AddLink(link)
				GetIpAddrs(link)
			}
		}
	}
}

var f *os.File
var path string

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

	path = "ipconfig_" + t.Local().Format("2006-01-02_15:04:05") + "/"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println("Can't create config dir")
		}
	}

	/*Get bridge info first */
	GetBridges()
	GetBonds()
	/*Get other link info */
	links, err := nlp.LinkList()
	if err != nil {
		ret = -1
		fmt.Println("Can't get device info")
		return file
	}

	for _, link := range links {
		_, br := link.(*nlp.Bridge)
		_, bo := link.(*nlp.Bond)
		if !br && !bo {
			ret = AddLink(link)

			if ret == -1 {
				continue
			}
			GetIpAddrs(link)
		}
		if !br {
			GetFdbs(link)
		}
		GetIpNeigh(link)
	}

	/* Get Static Routes only */
	rFilter := nlp.Route{Protocol: syscall.RTPROT_STATIC}
	routes, err := nlp.RouteListFiltered(nlp.FAMILY_V4, &rFilter, nlp.RT_FILTER_PROTOCOL)

	if err == nil {
		if len(routes) > 0 {
			for _, route := range routes {
				AddRoute(route)
			}
		}
	}
	if _, err := os.Stat("/opt/loxilb/ipconfig"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println("Can't create config dir")
		}
	} else {
		command := "mv /opt/loxilb/ipconfig /opt/loxilb/ipconfig.bk"
		cmd := exec.Command("bash", "-c", command)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("Can't backup /opt/loxilb/ipconfig")
			return file
		}
	}
	command := "cp -R " + path + " /opt/loxilb/ipconfig/"
	cmd := exec.Command("bash", "-c", command)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	return file
}
