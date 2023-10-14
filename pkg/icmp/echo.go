package icmp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"syscall"

	"github.com/louislef299/lnet/pkg/ip"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/sync/errgroup"
)

var (
	ErrDestinationUnreachable = errors.New("destination is unreachable")
	ErrNonEchoResponse        = errors.New("icmp non-echo response")
	ErrNotEchoReply           = errors.New("did not receive echo reply")
)

type Packet struct {
	bytes  []byte
	nbytes int
	ttl    int
}

// Takes in an existing ICMP connection and returns the message
func ReadEcho(conn *icmp.PacketConn) (*icmp.Message, net.Addr, error) {
	rb := make([]byte, 1500)
	n, peer, err := conn.ReadFrom(rb)
	if err != nil {
		return nil, nil, err
	}
	msg, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), rb[:n])
	if err != nil {
		return nil, nil, err
	}
	return msg, peer, err
}

func ReadAndInterpretEcho(ctx context.Context, conn *icmp.PacketConn) (chan *icmp.Message, error) {
	resp := make(chan *icmp.Message)
	group, _ := errgroup.WithContext(ctx)

	group.Go(func() error {
		for {
			msg, peer, err := ReadEcho(conn)
			if err != nil {
				return err
			}
			if err := ValidateEcho(msg, peer); err != nil {
				return err
			}
			resp <- msg
		}
	})
	return resp, nil
}

// Send an ICMP echo to the provided IP address given an existing connection
func SendEcho(conn *icmp.PacketConn, addr netip.Addr, sequenceNum int) error {
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  sequenceNum,
			Data: ip.Hash(addr),
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
	log.Println("scan sent to", addr.String())
	return err
}

// Validates that the ICMP message is an echo response, will return nil if it is
func ValidateEcho(msg *icmp.Message, addr net.Addr) error {
	switch msg.Type {
	case ipv4.ICMPTypeEchoReply:
		return validateEchoTargetReachable(msg)
	case ipv4.ICMPTypeDestinationUnreachable:
		return ErrDestinationUnreachable
	default:
		return ErrNotEchoReply
	}
}

// Validates that the target is reachable based on the message given
func validateEchoTargetReachable(msg *icmp.Message) error {
	// Validate echo response
	_, ok := msg.Body.(*icmp.Echo)
	if !ok {
		switch b := msg.Body.(type) {
		case *icmp.DstUnreach:
			dest, err := IcmpDestUnreachableCode(msg.Code)
			if err != nil {
				return err
			}
			return fmt.Errorf("icmp %s-unreachable", dest)
		case *icmp.PacketTooBig:
			return fmt.Errorf("icmp packet-too-big (mtu %d)", b.MTU)
		default:
			return ErrNonEchoResponse
		}
	}
	return nil
}
