/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"net"
	"net/netip"
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
		ifaces, err := net.Interfaces()
		if err != nil {
			log.Fatal(err)
		}

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

			prefix, err := netip.ParsePrefix(srcIP)
			if err != nil {
				log.Fatal(err)
			}
			if prefix.Addr().IsLoopback() {
				continue
			}

			for addr := prefix.Masked().Addr(); prefix.Contains(addr); addr = addr.Next() {
				// Send echo, but don't worry about errors
				sendEcho(addr)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(icmpCmd)
}

func sendEcho(addr netip.Addr) error {
	c, err := icmp.ListenPacket("udp4", addr.String())
	if err != nil {
		return fmt.Errorf("listen err, %s", err)
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
		return err
	}
	if _, err := c.WriteTo(wb, &net.UDPAddr{IP: net.ParseIP(addr.String())}); err != nil {
		return fmt.Errorf("WriteTo err, %s", err)
	}

	rb := make([]byte, 1500)
	n, peer, err := c.ReadFrom(rb)
	if err != nil {
		return err
	}
	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), rb[:n])
	if err != nil {
		return err
	}

	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		log.Printf("got reflection from %v", peer)
	default:
		return fmt.Errorf("got %+v; want echo reply", rm)
	}

	// Validate echo response
	_, ok := rm.Body.(*icmp.Echo)
	if !ok {
		switch b := rm.Body.(type) {
		case *icmp.DstUnreach:
			dest, err := licmp.IcmpDestUnreachableCode(rm.Code)
			if err != nil {
				return err
			}
			return fmt.Errorf("icmp %s-unreachable", dest)
		case *icmp.PacketTooBig:
			return fmt.Errorf("icmp packet-too-big (mtu %d)", b.MTU)
		default:
			return fmt.Errorf("icmp non-echo response")
		}
	}
	return nil
}
