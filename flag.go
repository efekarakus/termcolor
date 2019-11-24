package termcolor

import (
	"os"
	"strings"
)

// Point to os.Args for testing.
var args = os.Args

// See https://github.com/sindresorhus/has-flag/blob/ecd4cb75870f5d49eef1e0faee328b2019960de3/index.js#L1-L8
func hasFlag(flag string) bool {
	// Prefix the flag with the necessary dashes.
	var prefix string
	if !strings.HasPrefix(flag, "-") {
		if len(flag) == 1 {
			// Short flag.
			prefix = "-"
		} else {
			prefix = "--"
		}
	}
	pos := indexOf(args, prefix + flag)
	if pos == -1 {
		return false
	}
	// Flag parsing stops after the "--" flag.
	terminatorPos := indexOf(args, "--")
	if terminatorPos == -1 {
		// The flag exists and there is no terminator
		return true
	}
	return pos < terminatorPos
}

func indexOf(ss []string, s string) int {
	for i, el := range ss {
		if el == s {
			return i
		}
	}
	return -1
}
