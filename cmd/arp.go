/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"log"
	"net"
	"unsafe"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

type SockaddrInet4 struct {
	Port int
	Addr [4]byte
	raw  unix.RawSockaddrInet4
}

type arpreq struct {
	arp_pa    SockaddrInet4 /* protocol address */
	arp_ha    SockaddrInet4 /* hardware address */
	arp_flags int           /* flags */
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

		fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
		if err != nil {
			log.Fatal(err)
		}

		var areq *arpreq
		err = ioctl(uintptr(fd), unix.SIOCGARP, unsafe.Pointer(areq))
		if err != nil {
			log.Fatal(err)
		}
		log.Println(areq)

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
