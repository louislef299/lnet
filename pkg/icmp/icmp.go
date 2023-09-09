package icmp

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"

	"github.com/jpillora/icmpscan"
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
