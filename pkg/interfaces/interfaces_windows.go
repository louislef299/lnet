//go:build default
// +build default

package interfaces

import (
	"fmt"
	"runtime"

	"github.com/vishvananda/netlink"
)

func InterfaceMode(mode string, link netlink.Link) error {
	return fmt.Errorf("%v not supported", runtime.GOOS)
}
