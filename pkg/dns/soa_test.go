package dns

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

type SoaExpected struct {
	input    []string
	expected error
}

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

			for i, r := range resp {
				assert.Equal(t, r.OwnerName, c.input[i], "dns message owner didn't match")
			}
		})
	}
}
