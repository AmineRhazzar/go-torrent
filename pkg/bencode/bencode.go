package bencode

import (
	"fmt"
	"strconv"
)

func ParseInt(s string, returnsRest bool) (int64, string, error) {
	n := len(s)
	if n < 3 {
		return 0, "", fmt.Errorf("cannot parse int from bencoded string %s", s)
	}
	if s[0] != 'i' {
		return 0, "", fmt.Errorf("cannot parse int from bencoded string %s", s)
	}

	for i := 1; i < n; i++ {
		if s[i] == 'e' {
			p, err := strconv.ParseInt(s[1:i], 10, 64)
			if !returnsRest {
				return p, "", err
			}
			return p, s[i+1 : n], err
		}
	}
	return 0, "", fmt.Errorf("cannot parse int from bencoded string %s. Reason: end character (e) not found", s)
}

