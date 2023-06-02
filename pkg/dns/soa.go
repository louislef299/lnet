package dns

import (
	"fmt"

	"github.com/miekg/dns"
)

type SoaResponse struct {
	Msg    *dns.Msg
	Server string
}

func GetSoa(nameserver string, servers []string) ([]*SoaResponse, error) {
	c := new(dns.Client)
	msg := new(dns.Msg)
	msg.SetEdns0(4096, true)
	fmt.Println("sending question using name server", nameserver)

	var msgs []*SoaResponse
	for _, e := range servers {
		msg.SetQuestion(dns.Fqdn(e), dns.TypeSOA)
		resp, _, err := c.Exchange(msg, nameserver+":53")
		if err != nil {
			return nil, err
		}

		s := SoaResponse{
			Msg:    resp,
			Server: e,
		}
		msgs = append(msgs, &s)
	}
	return msgs, nil
}
