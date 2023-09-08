package ip

import (
	"encoding/binary"
	"math"
	"net"
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

// Converts an IPv4 address into a uint32
func ToUInt32(ip net.IP) uint32 {
	return binary.BigEndian.Uint32([]byte(ip.To4()))
}

// Converts a uint32 into an IPv4 address
func FromUInt32(u uint32) net.IP {
	buff := make([]byte, 4)
	binary.BigEndian.PutUint32(buff, u)
	return net.IP(buff)
}

// DeltaIP returns the IPv4 delta-many places away
func DeltaIP(ip net.IP, delta int) net.IP {
	if delta == 0 {
		return ip
	}
	i := ToUInt32(ip)
	if delta < 0 {
		i -= uint32(delta * -1)
	} else if delta > 0 {
		i += uint32(delta)
	}
	if i == math.MaxUint32 {
		return ip //cant increment past broadcast
	}
	return FromUInt32(i)
}

// NextIP returns the next IPv4 in sequence
func NextIP(ip net.IP) net.IP {
	return DeltaIP(ip, 1)
}

// PrevIP returns the previous IPv4 in sequence
func PrevIP(ip net.IP) net.IP {
	return DeltaIP(ip, -1)
}
