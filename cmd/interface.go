/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
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

		return nil
	},
}

// interfaceDownCmd represents the down command
var interfaceDownCmd = &cobra.Command{
	Use:   "down",
	Short: "configure and find system network interfaces",
	Long:  `Used to configure and find system network interfaces.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) == 0 {
			return cmd.Usage()
		}

		for _, iface := range args {
			link, err := netlink.LinkByName(iface)
			if err != nil {
				log.Fatalf("could not find network interface %s: %v", iface, err)
			}

			err = netlink.LinkSetDown(link)
			if err != nil {
				log.Fatalf("could not disable interface %s: %v", iface, err)
			}
			fmt.Printf("successfully disabled interface %s\n", iface)
		}

		return nil
	},
}

// interfaceUpCmd represents the up command
var interfaceUpCmd = &cobra.Command{
	Use:   "up",
	Short: "configure and find system network interfaces",
	Long:  `Used to configure and find system network interfaces.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) == 0 {
			return cmd.Usage()
		}

		for _, iface := range args {
			link, err := netlink.LinkByName(iface)
			if err != nil {
				log.Fatalf("could not find network interface %s: %v", iface, err)
			}

			err = netlink.LinkSetUp(link)
			if err != nil {
				log.Fatalf("could not enable interface %s: %v", iface, err)
			}
			fmt.Printf("successfully enabled interface %s\n", iface)
		}

		return nil
	},
}

// interfacePromiscCmd represents the promiscuous command
var interfacePromiscCmd = &cobra.Command{
	Use:     "promiscuous",
	Aliases: []string{"prom"},
	Short:   "configure and find system network interfaces",
	Long:    `Used to configure and find system network interfaces.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) == 0 {
			return cmd.Usage()
		}

		for _, iface := range args {
			link, err := netlink.LinkByName(iface)
			if err != nil {
				log.Fatalf("could not find network interface %s: %v", iface, err)
			}

			err = netlink.SetPromiscOn(link)
			if err != nil {
				log.Fatalf("could not enable promiscuous mode on interface %s: %v", iface, err)
			}
			fmt.Printf("successfully enabled promiscuous mode on interface %s\n", iface)
		}

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
	interfaceCmd.AddCommand(interfaceUpCmd)
	interfaceCmd.AddCommand(interfaceDownCmd)
	interfaceCmd.AddCommand(interfacePromiscCmd)
}

func printInterface(i net.Interface) string {
	return fmt.Sprintf("(%d)%s:\n  flags: <%v>\n  mtu: %d\n  hardware address: %s",
		i.Index, i.Name, i.Flags, i.MTU, i.HardwareAddr.String())
}
