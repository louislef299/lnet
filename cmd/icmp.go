/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"log"
	"net"
	"net/netip"
	"os"
	"sync/atomic"

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
		ifaces, err := net.Interfaces()
		if err != nil {
			log.Fatal(err)
		}

		var sequenceNum uint32 = 1
		for _, iface := range ifaces {
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
				log.Println("moving along")
				continue
			}

			// Use udp if not root user
			var network string
			if os.Getuid() == 0 {
				network = "ip4:icmp"
			} else {
				network = "udp4"
			}

			c, err := icmp.ListenPacket(network, prefix.Addr().String())
			if err != nil {
				log.Fatalf("listen err, %s", err)
			}
			defer c.Close()

			count := 0
			for addr := prefix.Masked().Addr(); prefix.Contains(addr); addr = addr.Next() {
				if count < 3 {
					count++
					continue
				}
				// Send echo, but don't worry about errors
				err := licmp.SendEcho(c, addr, int(sequenceNum))
				if err != nil {
					log.Fatal(err)
				}
				atomic.AddUint32(&sequenceNum, 1)

				rm, peer, err := licmp.ReadEcho(c)
				if err != nil {
					log.Fatal(err)
				}

				switch rm.Type {
				case ipv4.ICMPTypeEchoReply:
					log.Printf("got reflection from %v", peer)
				case ipv4.ICMPTypeDestinationUnreachable:
					log.Fatal("destination is unreachable!")
				default:
					log.Fatalf("got %+v; want echo reply", rm)
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

				count++
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(icmpCmd)
}
