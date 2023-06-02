package dns

import (
	"fmt"
	"runtime"

	"github.com/moby/moby/libnetwork/resolvconf"
)

func GetLocalNS() ([]string, error) {
	switch runtime.GOOS {
	case "darwin", "linux":
		f, err := resolvconf.Get()
		if err != nil {
			return nil, err
		}
		return resolvconf.GetNameservers(f.Content, 0), nil
	default:
		fmt.Printf("%s is not currently supported", runtime.GOOS)
		return nil, nil
	}
}
