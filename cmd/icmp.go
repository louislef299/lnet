/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"log"
	"net"
	"net/netip"
	"os"

	licmp "github.com/louislef299/lnet/pkg/icmp"
	"github.com/spf13/cobra"
)

var (
	icmpCode = []string{"network", "host", "protocol", "port", "must-fragment", "dest"}
)

// icmpCmd represents the icmp command
var icmpCmd = &cobra.Command{
	Use:   "icmp",
	Short: "Runs an ICMP scan on your local device",
	Long:  `ref: rfc-editor.org/rfc/rfc792`,
	Run: func(cmd *cobra.Command, args []string) {
		iface, err := net.InterfaceByName("wlp1s0")
		if err != nil {
			log.Fatal(err)
		}

		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}

		var srcIP string
		for _, addr := range addrs {
			src, _, _ := net.ParseCIDR(addr.String())
			if src.To4() != nil {
				//first ipv4 address on interface
				srcIP = addr.String()
				break
			}
		}

		log.Println("pinging off connection", srcIP)
		prefix, err := netip.ParsePrefix(srcIP)
		if err != nil {
			log.Fatal(err)
		}
		if prefix.Addr().IsLoopback() {
			log.Fatal("loopback address")
		}

		c, err := licmp.Listen(prefix.Addr())
		if err != nil {
			log.Fatalf("listen err, %s", err)
		}
		defer c.Close()

		i := licmp.NewICMP(c, &prefix)
		go i.Scan(cmd.Context())

		for {
			select {
			case r := <-i.Response:
				log.Println("got a response:", r.Body)
			case <-i.Done:
				log.Println("scan finished")
				return
			case <-cmd.Context().Done():
				log.Println("context cancelled")
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(icmpCmd)
}
