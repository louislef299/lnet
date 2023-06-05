package dns

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

type SoaResponse struct {
	OwnerName    string
	TTL          int64
	MName        string
	RName        string
	SerialNumber int64
	RefreshTime  int64
	RetryTime    int64
	ExpireTime   int64
	MinimumTTL   int64

	Msg *dns.Msg
}

func (s *SoaResponse) String() string {
	return fmt.Sprintf("ANSWER SECTION:\n  Owner Name: %s\n  TTL: %d\n  Primary NS: %s\n  Responsible Person: %s\n  Serial Number: %d\n  Refresh time in seconds: %d\n  Retry time in seconds: %d\n  Expire time in seconds: %d\n  Minimum TTL: %d",
		s.OwnerName, s.TTL, s.MName, s.RName, s.SerialNumber, s.RefreshTime, s.RetryTime, s.ExpireTime, s.MinimumTTL)
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

		s := &SoaResponse{
			Msg: resp,
		}
		err = s.ConvertMsgToSoa()
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, s)
	}
	return msgs, nil
}

// Converts a plain dns.Msg into a SoaResponse
func (s *SoaResponse) ConvertMsgToSoa() error {
	if len(s.Msg.Answer) < 1 {
		return ErrInvalidMsgAnswer
	}
	a := strings.FieldsFunc(s.Msg.Answer[0].String(), SoaSplit)
	if len(a) < 11 {
		return ErrInvalidMsgAnswer
	}

	m := map[string]int64{
		"TTL":          1,
		"SerialNumber": 6,
		"RefreshTime":  7,
		"RetryTime":    8,
		"ExpireTime":   9,
		"MinimumTTL":   10,
	}
	for k, v := range m {
		i, err := strconv.ParseInt(a[v], 10, 64)
		if err != nil {
			return err
		}
		m[k] = i
	}

	s.OwnerName = a[0]
	s.TTL = m["TTL"]
	s.MName = a[4]
	s.RName = a[5]
	s.SerialNumber = m["SerialNumber"]
	s.RefreshTime = m["RefreshTime"]
	s.RetryTime = m["RetryTime"]
	s.ExpireTime = m["ExpireTime"]
	s.MinimumTTL = m["MinimumTTL"]

	return nil
}

func SoaSplit(r rune) bool {
	return r == '\t' || r == ' '
}
