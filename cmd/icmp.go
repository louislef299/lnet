/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"net"
	"os"

	licmp "github.com/louislef299/lnet/pkg/icmp"
	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
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
		fmt.Println(addrs)

		var srcIP string
		for _, addr := range addrs {
			src, _, _ := net.ParseCIDR(addr.String())
			if src.To4() != nil {
				//first ipv4 address on interface
				srcIP = src.String()
				break
			}
		}

		c, err := icmp.ListenPacket("udp4", srcIP)
		if err != nil {
			log.Fatalf("listen err, %s", err)
		}
		defer c.Close()

		wm := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  1,
				Data: []byte("HELLO-R-U-THERE"),
			},
		}
		wb, err := wm.Marshal(nil)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := c.WriteTo(wb, &net.UDPAddr{IP: net.ParseIP("192.168.4.1")}); err != nil {
			log.Fatalf("WriteTo err, %s", err)
		}

		rb := make([]byte, 1500)
		n, peer, err := c.ReadFrom(rb)
		if err != nil {
			log.Fatal(err)
		}
		rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), rb[:n])
		if err != nil {
			log.Fatal(err)
		}
		switch rm.Type {
		case ipv4.ICMPTypeEchoReply:
			log.Printf("got reflection from %v", peer)
		default:
			log.Printf("got %+v; want echo reply", rm)
			os.Exit(1)
		}

		// Validate echo response
		_, ok := rm.Body.(*icmp.Echo)
		if !ok {
			switch b := rm.Body.(type) {
			case *icmp.DstUnreach:
				dest, err := licmp.IcmpDestUnreachableCode(rm.Code)
				if err != nil {
					log.Fatal(err)
				}
				log.Fatalf("icmp %s-unreachable", dest)
			case *icmp.PacketTooBig:
				log.Fatalf("icmp packet-too-big (mtu %d)", b.MTU)
			default:
				log.Fatal("icmp non-echo response")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(icmpCmd)
}
