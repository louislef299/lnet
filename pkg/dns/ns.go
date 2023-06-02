package dns

import (
	"fmt"
	"runtime"

	"github.com/moby/moby/libnetwork/resolvconf"
)

// Returns local name servers
func GetLocalNS() ([]string, error) {
	switch runtime.GOOS {
	case "darwin", "linux":
		f, err := resolvconf.Get()
		if err != nil {
			return nil, err
		}
		return resolvconf.GetNameservers(f.Content, 0), nil
	default:
		return nil, fmt.Errorf("%s is not currently supported", runtime.GOOS)
	}
}

// Returns local resolv.conf options
func GetLocalOptions() ([]string, error) {
	switch runtime.GOOS {
	case "darwin", "linux":
		f, err := resolvconf.Get()
		if err != nil {
			return nil, err
		}
		return resolvconf.GetOptions(f.Content), nil
	default:
		return nil, fmt.Errorf("%s is not currently supported", runtime.GOOS)
	}
}

// Returns local resolv.conf search domains
func GetLocalSearchDomains() ([]string, error) {
	switch runtime.GOOS {
	case "darwin", "linux":
		f, err := resolvconf.Get()
		if err != nil {
			return nil, err
		}
		return resolvconf.GetSearchDomains(f.Content), nil
	default:
		return nil, fmt.Errorf("%s is not currently supported", runtime.GOOS)
	}
}

// Returns local resolv.conf path
func GetPath() (string, error) {
	switch runtime.GOOS {
	case "darwin", "linux":
		return resolvconf.Path(), nil
	default:
		return "", fmt.Errorf("%s is not currently supported", runtime.GOOS)
	}
}
