/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
)

const (
	// See linux/if_arp.h.
	// Note that Linux doesn't support IPv4 over IPv6 tunneling.
	sysARPHardwareIPv4IPv4 = 768 // IPv4 over IPv4 tunneling
	sysARPHardwareIPv6IPv6 = 769 // IPv6 over IPv6 tunneling
	sysARPHardwareIPv6IPv4 = 776 // IPv6 over IPv4 tunneling
	sysARPHardwareGREIPv4  = 778 // any over GRE over IPv4 tunneling
	sysARPHardwareGREIPv6  = 823 // any over GRE over IPv6 tunneling
)

// interfaceCmd represents the interface command. Utilizes
// [RFC3549](https://datatracker.ietf.org/doc/html/rfc3549)
var interfaceCmd = &cobra.Command{
	Use:     "interface",
	Aliases: []string{"if", "inter"},
	Short:   "configure and find system network interfaces",
	Long: `Used to configure and find system network interfaces. Controls
the kernel space interfaces and routes using the NETLINK
address family.(RFC3549)`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var interfaces []net.Interface
		if len(args) == 0 {
			interfaces, err = net.Interfaces()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			for _, n := range args {
				iname, err := net.InterfaceByName(n)
				if err != nil {
					log.Fatal(err)
				}
				interfaces = append(interfaces, *iname)
			}
		}
		for _, i := range interfaces {
			fmt.Println(printInterface(i))
			addrs, err := i.Addrs()
			if err != nil {
				return err
			}
			for i, a := range addrs {
				if i == 0 {
					fmt.Printf("  addrs: ")
				} else {
					fmt.Printf("   - ")
				}
				fmt.Println(a)
			}
		}
		t, err := interfaceTable(0)
		if err != nil {
			return err
		}
		fmt.Println("table:", t)
		return nil
	},
}

// interfaceSockCmd represents the interface command
var interfaceSockCmd = &cobra.Command{
	Use:     "socket",
	Aliases: []string{"sock", "mtu"},
	Short:   "configure and find system network interfaces",
	Long:    `Used to configure and find system network interfaces.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		link, err := netlink.LinkByName(args[0])
		if err != nil {
			panic(err)
		}

		mtu, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}

		err = netlink.LinkSetMTU(link, mtu)
		if err != nil {
			panic(err)
		}
		fmt.Println(link)

		// Communication directly with NETLINK in the kernel uses a socket
		// to communicate. A socket must first be created along with a send/
		// recv request to gather information:

		// fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, unix.NETLINK_GENERIC)
		// if err != nil {
		// 	panic(err)
		// }
		// To dive further into the subject, follow the Linux kernel introduction:
		// https://docs.kernel.org/userspace-api/netlink/intro.html

		return nil
	},
}

func init() {
	rootCmd.AddCommand(interfaceCmd)
	interfaceCmd.AddCommand(interfaceSockCmd)
}

func printInterface(i net.Interface) string {
	return fmt.Sprintf("(%d)%s:\n  flags: <%v>\n  mtu: %d\n  hardware address: %s",
		i.Index, i.Name, i.Flags, i.MTU, i.HardwareAddr.String())
}

// If the ifindex is zero, interfaceTable returns mappings of all
// network interfaces. Otherwise it returns a mapping of a specific
// interface.
func interfaceTable(ifindex int) ([]net.Interface, error) {
	tab, err := syscall.NetlinkRIB(syscall.RTM_GETLINK, syscall.AF_UNSPEC)
	if err != nil {
		return nil, os.NewSyscallError("netlinkrib", err)
	}
	msgs, err := syscall.ParseNetlinkMessage(tab)
	if err != nil {
		return nil, os.NewSyscallError("parsenetlinkmessage", err)
	}
	var ift []net.Interface
loop:
	for _, m := range msgs {
		switch m.Header.Type {
		case syscall.NLMSG_DONE:
			break loop
		case syscall.RTM_NEWLINK:
			ifim := (*syscall.IfInfomsg)(unsafe.Pointer(&m.Data[0]))
			if ifindex == 0 || ifindex == int(ifim.Index) {
				attrs, err := syscall.ParseNetlinkRouteAttr(&m)
				if err != nil {
					return nil, os.NewSyscallError("parsenetlinkrouteattr", err)
				}
				ift = append(ift, *newLink(ifim, attrs))
				if ifindex == int(ifim.Index) {
					break loop
				}
			}
		}
	}
	return ift, nil
}

func newLink(ifim *syscall.IfInfomsg, attrs []syscall.NetlinkRouteAttr) *net.Interface {
	ifi := &net.Interface{Index: int(ifim.Index), Flags: linkFlags(ifim.Flags)}
	for _, a := range attrs {
		switch a.Attr.Type {
		case syscall.IFLA_ADDRESS:
			// We never return any /32 or /128 IP address
			// prefix on any IP tunnel interface as the
			// hardware address.
			switch len(a.Value) {
			case net.IPv4len:
				switch ifim.Type {
				case sysARPHardwareIPv4IPv4, sysARPHardwareGREIPv4, sysARPHardwareIPv6IPv4:
					continue
				}
			case net.IPv6len:
				switch ifim.Type {
				case sysARPHardwareIPv6IPv6, sysARPHardwareGREIPv6:
					continue
				}
			}
			var nonzero bool
			for _, b := range a.Value {
				if b != 0 {
					nonzero = true
					break
				}
			}
			if nonzero {
				ifi.HardwareAddr = a.Value[:]
			}
		case syscall.IFLA_IFNAME:
			ifi.Name = string(a.Value[:len(a.Value)-1])
		case syscall.IFLA_MTU:
			ifi.MTU = int(*(*uint32)(unsafe.Pointer(&a.Value[:4][0])))
		}
	}
	return ifi
}

func linkFlags(rawFlags uint32) net.Flags {
	var f net.Flags
	if rawFlags&syscall.IFF_UP != 0 {
		f |= net.FlagUp
	}
	if rawFlags&syscall.IFF_RUNNING != 0 {
		f |= net.FlagRunning
	}
	if rawFlags&syscall.IFF_BROADCAST != 0 {
		f |= net.FlagBroadcast
	}
	if rawFlags&syscall.IFF_LOOPBACK != 0 {
		f |= net.FlagLoopback
	}
	if rawFlags&syscall.IFF_POINTOPOINT != 0 {
		f |= net.FlagPointToPoint
	}
	if rawFlags&syscall.IFF_MULTICAST != 0 {
		f |= net.FlagMulticast
	}
	return f
}
