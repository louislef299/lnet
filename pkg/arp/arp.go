//go:build linux
// +build linux

package arp

// This file contains ARP related functions for the Seesaw Network Control
// component.

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	opARPRequest = 1
	opARPReply   = 2
	hwLen        = 6
)

var (
	ethernetBroadcast = net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

func htons(p uint16) uint16 {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], p)
	return *(*uint16)(unsafe.Pointer(&b))
}

// arpHeader specifies the header for an ARP message.
type arpHeader struct {
	hardwareType          uint16
	protocolType          uint16
	hardwareAddressLength uint8
	protocolAddressLength uint8
	opcode                uint16
}

// arpMessage represents an ARP message.
type arpMessage struct {
	arpHeader
	senderHardwareAddress []byte
	senderProtocolAddress []byte
	targetHardwareAddress []byte
	targetProtocolAddress []byte
}

// bytes returns the wire representation of the ARP message.
func (m *arpMessage) bytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, m.arpHeader); err != nil {
		return nil, fmt.Errorf("binary write failed: %v", err)
	}
	buf.Write(m.senderHardwareAddress)
	buf.Write(m.senderProtocolAddress)
	buf.Write(m.targetHardwareAddress)
	buf.Write(m.targetProtocolAddress)

	return buf.Bytes(), nil
}

// gratuitousARPReply returns an ARP message that contains a gratuitous ARP
// reply from the specified sender.
func GratuitousARPReply(ip net.IP, mac net.HardwareAddr) (*arpMessage, error) {
	if ip.To4() == nil {
		return nil, fmt.Errorf("%q is not an IPv4 address", ip)
	}
	if len(mac) != hwLen {
		return nil, fmt.Errorf("%q is not an Ethernet MAC address", mac)
	}

	m := &arpMessage{
		arpHeader{
			1,           // Ethernet
			0x0800,      // IPv4
			hwLen,       // 48-bit MAC Address
			net.IPv4len, // 32-bit IPv4 Address
			opARPReply,  // ARP Reply
		},
		mac,
		ip.To4(),
		ethernetBroadcast,
		net.IPv4bcast,
	}

	return m, nil
}

// sendARP sends the given ARP message via the specified interface.
func SendARP(iface *net.Interface, m *arpMessage) error {
	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_DGRAM, int(htons(unix.ETH_P_ARP)))
	if err != nil {
		return fmt.Errorf("failed to get raw socket: %v", err)
	}
	defer unix.Close(fd)

	if err := unix.BindToDevice(fd, iface.Name); err != nil {
		return fmt.Errorf("failed to bind to device: %v", err)
	}

	ll := unix.SockaddrLinklayer{
		Protocol: htons(unix.ETH_P_ARP),
		Ifindex:  iface.Index,
		Pkttype:  0, // unix.PACKET_HOST
		Hatype:   m.hardwareType,
		Halen:    m.hardwareAddressLength,
	}
	target := ethernetBroadcast
	if m.opcode == opARPReply {
		target = m.targetHardwareAddress
	}
	copy(ll.Addr[:], target)

	b, err := m.bytes()
	if err != nil {
		return fmt.Errorf("failed to convert ARP message: %v", err)
	}

	if err := unix.Bind(fd, &ll); err != nil {
		return fmt.Errorf("failed to bind: %v", err)
	}
	if err := unix.Sendto(fd, b, 0, &ll); err != nil {
		return fmt.Errorf("failed to send: %v", err)
	}

	return nil
}

// ARPSendGratuitous sends a gratuitous ARP message via the specified interface.
func ARPSendGratuitous(arpMap map[string][]net.IP, out *int) error {
	for ifName, ips := range arpMap {
		iface, err := net.InterfaceByName(ifName)
		if err != nil {
			log.Printf("failed to get interface %q: %v", ifName, err)
			continue
		}
		for _, ip := range ips {
			log.Printf("Sending gratuitous ARP for %s (%s) via %s", ip, iface.HardwareAddr, iface.Name)
			m, err := GratuitousARPReply(ip, iface.HardwareAddr)
			if err != nil {
				log.Printf("failed to build gratuitous arp: %v", err)
				continue
			}
			if err := SendARP(iface, m); err != nil {
				log.Printf("failed to send gratuitous arp: %v", err)
				continue
			}
		}
	}
	return nil
}
