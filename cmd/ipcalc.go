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
	Use:   "ipcalc",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if len(args) < 1 {
		// 	log.Println("requires an ip address")
		// }
		prefix, err := netip.ParsePrefix("192.168.11.5/24")
		if err != nil {
			log.Fatal("could not parse provided ip address", err)
		}

		fmt.Println("Printing information for IP", prefix.String())
		prettyPrint(os.Stdout, "Address", prefix.Addr())
		netPrefix := prefix.Masked()

		count := 0
		var addr netip.Addr
		for addr = netPrefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
			var header string
			print := true
			switch count {
			case 0:
				header = "Network"
			case 1:
				header = "HostMin"
			default:
				print = false
			}

			if print {
				prettyPrint(os.Stdout, header, addr)
			}
			count++
		}

		prettyPrint(os.Stdout, "HostMax", addr.Prev().Prev())
		prettyPrint(os.Stdout, "Broadcast", addr.Prev())
	},
}

func init() {
	rootCmd.AddCommand(ipcalcCmd)
}

func prettyPrint(out io.Writer, prefix string, addr netip.Addr) {
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
