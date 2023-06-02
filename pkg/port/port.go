package port

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	closed = "Closed"
	open   = "Open"
)

type ScanResult struct {
	Port     int
	State    string
	Protocol string
}

func IsOpen(state string) bool {
	if strings.Compare(state, open) == 0 {
		return true
	}
	return false
}

func ScanWithRoutines(ctx context.Context, hostname string, portrange int) (chan ScanResult, chan struct{}) {
	r := make(chan ScanResult)
	done := make(chan struct{})

	g, _ := errgroup.WithContext(ctx)
	protocol := "tcp"
	for port := 0; port <= portrange; port++ {
		port := port // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			result := ScanResult{Port: port, Protocol: protocol}
			address := hostname + ":" + strconv.Itoa(port)
			conn, err := net.DialTimeout(protocol, address, 10*time.Second)
			if err != nil {
				result.State = closed
				r <- result
				return err
			}
			defer conn.Close()
			result.State = open
			r <- result
			return nil
		})
	}

	go func() {
		g.Wait()
		close(r)
		done <- struct{}{}
	}()

	return r, done
}
