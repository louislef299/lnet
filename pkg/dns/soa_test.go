package dns

import (
	"fmt"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

type SoaExpected struct {
	input    []string
	expected error
}

const (
	ownerLoc = 0
)

func TestGetSoa(t *testing.T) {
	testCases := []SoaExpected{
		{[]string{"google.com."}, nil},
		{[]string{"google.com.", "apple.com."}, nil},
	}

	ns, err := GetLocalNS()
	if err != nil {
		t.Fatal("could not gather name server:", err)
	}

	for _, c := range testCases {
		t.Run(fmt.Sprintf("Validating %s", c.input), func(t *testing.T) {
			resp, err := GetSoa(ns[0], c.input)
			assert.Equal(t, err, c.expected, "errors don't match")

			for _, r := range resp {
				ans := strings.FieldsFunc(r.Msg.Answer[0].String(), Split)
				assert.Equal(t, ans[ownerLoc], r.Server, "dns message owner didn't match")
			}
		})
	}
}

func Split(r rune) bool {
	return r == '\t' || r == ' '
}
