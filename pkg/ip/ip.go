package ip

import (
	"crypto/sha1"
	"net"
	"net/netip"
	"strings"
)

func IsIPv4(address string) bool {
	if net.ParseIP(address) == nil {
		return false
	}
	return strings.Count(address, ":") < 2
}

func IsIPv6(address string) bool {
	if net.ParseIP(address) == nil {
		return false
	}
	return strings.Count(address, ":") >= 2
}

// Hash an IP with SHA1
func Hash(ip netip.Addr) []byte {
	input := []byte(ip.String())
	h := sha1.New()
	h.Write(input)
	output := h.Sum(nil)
	return output
}
