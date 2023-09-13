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

func BottomOfIt() {
	intfs, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, n := range intfs {
		addrs, err := n.Addrs()
		if err != nil {
			log.Fatal(err)
		}

		// check if addr is ipv4
		valid := false
		var ipv4 net.IP
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				log.Fatal(err)
			}
			if i := ip.To4(); i != nil {
				valid = true
				ipv4 = i
			}
		}
		if valid {
			fmt.Println(n, "has an ip", ipv4)
		}
	}
}

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

// Hash an IP with SHA1
func hash(ip netip.Addr) []byte {
	input := []byte(ip.String())
	h := sha1.New()
	h.Write(input)
	output := h.Sum(nil)
	return output
}
