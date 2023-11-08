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

func ParseString(s string, returnsRest bool) (string, string, error) {
	n := int32(len(s))
	// get index of first occurence of : in string s
	var colonIndex int32 = -1
	var i int32 = 0
	for i = 0; i < n; i++ {
		if string(s[i]) == ":" {
			colonIndex = i
			break
		}
	}
	if colonIndex == -1 {
		return "", "", fmt.Errorf("cannot parse bencoded string %s. Reason: can't find \":\"", s)
	}
	stringLength, err := strconv.ParseInt(s[0:colonIndex], 10, 32)
	if err != nil {
		return "", "", fmt.Errorf("cannot parse bencoded string %s. Reason: can't parse length", s)
	}
	// check if there's enough length after colon (minimum length of stringLength)
	stop := colonIndex + 1 + int32(stringLength)
	if stop > n {
		return "", "", fmt.Errorf("cannot parse bencoded string %s. Reason: string not long enough (must be: %d)", s, stringLength)
	}

	if returnsRest {
		return s[colonIndex+1 : stop], s[stop:n], nil
	} else {
		return s[colonIndex+1 : stop], "", nil
	}

}