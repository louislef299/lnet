/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"net/netip"

	"github.com/spf13/cobra"
)

// ipcalcCmd represents the ipcalc command
var ipcalcCmd = &cobra.Command{
	Use:     "ipcalc",
	Example: "  lnet ipcalc 192.168.11.5/24",
	Short:   "Returns information for an IPv4 CIDR range.",
	Long: `Returns information for an IPv4 CIDR range. An 
IP address to operate on must always be specified 
along with a netmask or a CIDR prefix as well. `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Println("requires an ip address")
		}
		prefix, err := netip.ParsePrefix(args[0])
		if err != nil {
			log.Fatal("could not parse provided ip address", err)
		}

		fmt.Println("Printing information for IP", prefix.String())
		prettyPrintNetwork(os.Stdout, "Address", prefix.Addr())

		if prefixMask(prefix) == 32 {
			fmt.Println("Network IP Space: 1")
			return
		}

		netPrefix := prefix.Masked()

		count := 0
		var addr netip.Addr
		for addr = netPrefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
			var header string
			if count == 0 {
				header = "Network"
			} else if count == 1 && prefixMask(prefix) < 31 {
				header = "HostMin"
			}

			if header != "" {
				prettyPrintNetwork(os.Stdout, header, addr)
			}
			count++
		}

		if prefixMask(prefix) < 31 {
			prettyPrintNetwork(os.Stdout, "HostMax", addr.Prev().Prev())
		}
		prettyPrintNetwork(os.Stdout, "Broadcast", addr.Prev())
		fmt.Println("Network IP Space:", count)
	},
}

func init() {
	rootCmd.AddCommand(ipcalcCmd)
}

// Pretty prints the network IP and binary address
func prettyPrintNetwork(out io.Writer, prefix string, addr netip.Addr) {
	addrBin, err := addr.MarshalBinary()
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(out, "%s:\t%s\t", prefix, addr)
	for i, b := range addrBin {
		fmt.Fprintf(out, "%08b", b)
		if i < len(addrBin)-1 {
			fmt.Fprintf(out, ".")
		}
	}
	fmt.Fprintf(out, "\n")
}

func prefixMask(networkPrefix netip.Prefix) int {
	n, err := networkPrefix.MarshalBinary()
	if err != nil {
		panic(err)
	}

	return int(n[len(n)-1])
}
