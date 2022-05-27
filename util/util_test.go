package util

import (
	"strings"
	"testing"
)

func TestRandomString(t *testing.T) {
	t.Logf("Upper random 32: %s", strings.ToUpper(RandomString(32)))
	t.Logf("Lower random 37: %s", strings.ToLower(RandomString(37)))
	t.Logf("Lower random 43: %s", strings.ToLower(RandomString(43)))
}
