package ip

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

type ipCheck func(string) bool

type IpExpectedOutput struct {
	input    string
	expected bool
	test     ipCheck
}

func TestGetSoa(t *testing.T) {
	testCases := []IpExpectedOutput{
		{"192.168.1.1", true, IsIPv4},
		{"10000.4.3.4", false, IsIPv4},
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", true, IsIPv6},
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334:3445", false, IsIPv6},
	}

	for _, c := range testCases {
		t.Run(fmt.Sprintf("Validating %s", c.input), func(t *testing.T) {
			resp := c.test(c.input)
			assert.Equal(t, resp, c.expected)
		})
	}
}
