/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"log"
	"net"
	"net/netip"
	"os"
	"time"

	licmp "github.com/louislef299/lnet/pkg/icmp"
	"github.com/spf13/cobra"
)

// Packet represents a received and processed ICMP echo packet.
type Packet struct {
	// Rtt is the round-trip time it took to ping.
	Rtt time.Duration

	// IPAddr is the address of the host being pinged.
	IPAddr *net.IPAddr

	// Addr is the string address of the host being pinged.
	Addr string

	// NBytes is the number of bytes in the message.
	Nbytes int

	// Seq is the ICMP sequence number.
	Seq int

	// TTL is the Time To Live on the packet.
	Ttl int

	// ID is the ICMP identifier.
	ID int
}

var (
	icmpCode = []string{"network", "host", "protocol", "port", "must-fragment", "dest"}
	timeout  string
)

// icmpCmd represents the icmp command
var icmpCmd = &cobra.Command{
	Use:   "icmp",
	Short: "Runs an ICMP scan on your local device",
	Long:  `ref: rfc-editor.org/rfc/rfc792`,
	Run: func(cmd *cobra.Command, args []string) {
		t, err := time.ParseDuration(timeout)
		if err != nil {
			log.Fatal("couldn't parse timeout duration:", err)
		}

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

		c, err := licmp.Listen(prefix.Addr(), time.Now().Add(t))
		if err != nil {
			log.Fatalf("listen err, %s", err)
		}
		defer c.Close()

		i := licmp.NewICMP(c, &prefix)
		go i.Scan(cmd.Context())

		for {
			select {
			case r := <-i.Response:
				log.Println("got a response:", r)
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

	icmpCmd.Flags().StringVarP(&timeout, "timeout", "t", "2m", "timeout for the entire icmp scan")
}
