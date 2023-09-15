package icmp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/netip"
	"os"
	"regexp"
	"sync/atomic"

	"github.com/jpillora/icmpscan"
	"golang.org/x/net/icmp"
	"golang.org/x/sync/errgroup"
)

var (
	icmpCode           = []string{"network", "host", "protocol", "port", "must-fragment", "dest"}
	sequenceNum uint32 = 1

	ErrInvalidICMPCode = errors.New("the provided code is invalid")
	ErrInvalidAddress  = errors.New("the IP address provided is invalid")
)

type ICMP struct {
	// Address range to send ICMP Requests
	Prefix *netip.Prefix

	// Socket to send the ICMP request on
	Conn     *icmp.PacketConn
	Response chan *icmp.Message
	Done     chan struct{}
}

// When ICMP returned message is of type "Destination Unreachable", can
// call the code to get the hardware error.
func IcmpDestUnreachableCode(code int) (string, error) {
	if code > len(icmpCode) {
		return "", ErrInvalidICMPCode
	} else {
		return icmpCode[code], nil
	}
}

func (i *ICMP) Scan(ctx context.Context) error {
	group, _ := errgroup.WithContext(ctx)
	group.Go(func() error {
		return ReadAndInterpretEcho(i.Conn, i.Response)
	})

	count := 0
	for addr := i.Prefix.Masked().Addr(); i.Prefix.Contains(addr); addr = addr.Next() {
		// addrInLoop := addr
		group.Go(func() error {
			return SendEcho(i.Conn, addr, int(sequenceNum))
		})
		atomic.AddUint32(&sequenceNum, 1)
		count++
	}

	log.Printf("waiting on %d scans", count)
	group.Wait()
	i.Done <- struct{}{}
	return nil
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

func NewICMP(conn *icmp.PacketConn, prefix *netip.Prefix) *ICMP {
	return &ICMP{
		Conn:     conn,
		Prefix:   prefix,
		Response: make(chan *icmp.Message),
		Done:     make(chan struct{}),
	}
}

func OldIcmpScan() {
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
