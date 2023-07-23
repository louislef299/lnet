/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"encoding/binary"
	"log"
	"net"
	"unsafe"

	"github.com/spf13/cobra"
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

// arpCmd represents the arp command
var arpCmd = &cobra.Command{
	Use:   "arp",
	Short: "A brief description of your command",
	Long: `Address Resolution Protocol(ARP) is a protocol for mapping
IP addresses to fixed hardware(MAC) addresses. When a new 
computer joins a local area network (LAN), it will receive 
a unique IP address to use for identification and 
communication. 

Packets of data arrive at a gateway, destined for a 
particular host machine. The gateway, or the piece of 
hardware on a network that allows data to flow from one 
network to another, asks the ARP program to find a MAC address 
that matches the IP address. The ARP cache keeps a list of 
each IP address and its matching MAC address. The ARP cache is 
dynamic, but users on a network can also configure a static ARP 
table containing IP addresses and MAC addresses.

Read more in RFC 826:
datatracker.ietf.org/doc/html/rfc826`,
	Run: func(cmd *cobra.Command, args []string) {
		iface, err := getMainInterface()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(iface)

		// arp.ARPSendGratuitous(map[string][]net.IP{"en0"})

		// fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// var areq *arpHeader
		// err = ioctl(uintptr(fd), unix.SIOCGARP, unsafe.Pointer(areq))
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// log.Println(areq)

		// mdlayher didn't implement working arp protocol :(
		// decode the packet layer with:
		// https://pkg.go.dev/github.com/google/gopacket
	},
}

func init() {
	rootCmd.AddCommand(arpCmd)
}

// could use some work, going on a lot of assumptions
func getMainInterface() (*net.Interface, error) {
	i, err := net.InterfaceByName("en0")
	if err != nil {
		// if no en0, just grab first interface
		return net.InterfaceByIndex(1)
	}
	return i, nil
}

func ioctl(fd uintptr, name int, data unsafe.Pointer) unix.Errno {
	_, _, err := unix.RawSyscall(unix.SYS_IOCTL, fd, uintptr(name), uintptr(data))
	return err
}
