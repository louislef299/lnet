package icmp

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"regexp"
	"syscall"

	"github.com/jpillora/icmpscan"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

var (
	icmpCode           = []string{"network", "host", "protocol", "port", "must-fragment", "dest"}
	ErrInvalidICMPCode = errors.New("the provided code is invalid")
	ErrInvalidAddress  = errors.New("the IP address provided is invalid")
)

// When ICMP returned message is of type "Destination Unreachable", can
// call the code to get the hardware error.
func IcmpDestUnreachableCode(code int) (string, error) {
	if code > len(icmpCode) {
		return "", ErrInvalidICMPCode
	} else {
		return icmpCode[code], nil
	}
}

// Open an ICMP socket
func Listen(addr netip.Addr) (*icmp.PacketConn, error) {
	var network string
	priv := os.Getuid() == 0
	if priv && addr.Is4() {
		network = "ip4:icmp"
	} else if !priv && addr.Is4() {
		network = "udp4" // Use udp if not root user
	} else if priv && addr.Is6() {
		network = "ip6:ipv6-icmp"
	} else if !priv && addr.Is6() {
		network = "udp6"
	} else {
		return nil, ErrInvalidAddress
	}

	return icmp.ListenPacket(network, addr.String())
}

// Takes in an existing ICMP connection and returns the message
func ReadEcho(conn *icmp.PacketConn) (*icmp.Message, net.Addr, error) {
	rb := make([]byte, 1500)
	n, peer, err := conn.ReadFrom(rb)
	if err != nil {
		return nil, nil, err
	}
	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), rb[:n])
	if err != nil {
		return nil, nil, err
	}
	return rm, peer, err
}

// Send an ICMP echo to the provided IP address given an existing connection
func SendEcho(conn *icmp.PacketConn, addr netip.Addr, sequenceNum int) error {
	log.Println("pinging", addr.String())
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  sequenceNum,
			Data: hash(addr),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		return err
	}

	_, err = conn.WriteTo(wb, &net.IPAddr{IP: net.ParseIP(addr.String())})
	if neterr, ok := err.(*net.OpError); ok {
		if neterr.Err == syscall.ENOBUFS {
			return nil
		}
	}
	return err
}

// Hash an IP with SHA1
func hash(ip netip.Addr) []byte {
	input := []byte(ip.String())
	h := sha1.New()
	h.Write(input)
	output := h.Sum(nil)
	return output
}

func IcmpScan() {
	hosts, err := icmpscan.Run(icmpscan.Spec{
		Hostnames: true,
		MACs:      true,
		Log:       true,
	})
	if err != nil {
		log.Fatal("could not run local scan:", err)
	}

	decimals := regexp.MustCompile(`\.\d+`)
	for i, host := range hosts {
		if host.Active {
			if host.MAC == "" {
				host.MAC = "-"
			}
			if host.Hostname == "" {
				host.Hostname = "-"
			}
			rtt := decimals.ReplaceAllString(host.RTT.String(), "")
			fmt.Printf("[%03d] %15s, %6s, %17s, %s\n", i+1, host.IP, rtt, host.MAC, host.Hostname)
		}
	}
}
