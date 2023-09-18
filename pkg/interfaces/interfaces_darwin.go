//go:build darwin
// +build darwin

package interfaces

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

func InterfaceMode(mode string, link netlink.Link) error {
	var err error
	switch mode {
	case Up:
		err = netlink.LinkSetUp(link)
	case Down:
		err = netlink.LinkSetDown(link)
	default:
		return fmt.Errorf("%v is not a supported mode", mode)
	}
	return err
}
