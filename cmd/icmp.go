/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const targetIP = "8.8.8.8"

// icmpCmd represents the icmp command
var icmpCmd = &cobra.Command{
	Use:   "icmp",
	Short: "Runs an ICMP scan on your local device",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		reply, ok := rm.Body.(*icmp.Echo)
		if !ok {
			switch b := rm.Body.(type) {
			case *icmp.DstUnreach:
				dest := ""
				switch rm.Code {
				case 0:
					dest = "network"
				case 1:
					dest = "host"
				case 2:
					dest = "protocol"
				case 3:
					dest = "port"
				case 4:
					dest = "must-fragment"
				default:
					dest = "dest"
				}
				log.Fatalf("icmp %s-unreachable", dest)
			case *icmp.PacketTooBig:
				log.Fatalf("icmp packet-too-big (mtu %d)", b.MTU)
			default:
				log.Fatal("icmp non-echo response")
			}
		}
		ip, err := icmp.ParseIPv4Header(reply.Data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ip)
	},
}

func init() {
	rootCmd.AddCommand(icmpCmd)
}
