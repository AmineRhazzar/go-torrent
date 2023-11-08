package bencode

import (
	"fmt"
	"strconv"
)

type DictionaryElementType int

const (
	INTEGER DictionaryElementType = iota
	STRING
	LIST
	DICTIONARY
)

type DictionaryElement struct {
	kind       DictionaryElementType
	Integer    int
	String     string
	List       []string
	Dictionary map[string]DictionaryElement
}

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

func ParseList(s string, returnsRest bool) ([](string), string, error) {
	n := len(s)

	if s[0] != 'l' {
		return make([]string, 0), "", fmt.Errorf("cannot parse list from bencoded string %s", s)
	}

	var list []string
	var err error

	// stores the parsed string length. Example: "5:hello" -> currentStringLength=5
	var currentStringLength int64 = 0
	// index iterating through whole string s
	i := 1
	// indicates the index where the currentStringLength starts at in string s
	start := 1

	for i < n {
		if s[i] == ':' {
			currentStringLength, err = strconv.ParseInt(s[start:i], 10, 32)
			if err != nil {
				return make([]string, 0), "", fmt.Errorf("cannot parse list from bencoded string %s. Reason: error parsing string length from interval [%d, %d[", s, start, i-1)
			}
			// so we parse from i+1 (directly after colon) until i+1+currentStringLength
			list = append(list, s[i+1:i+1+int(currentStringLength)])
			i = i + 1 + int(currentStringLength)
			start = i
			// if i is at the end of string s now (i==n) we return error
			// we parsed the string successfully but there's no 'e' at the end to the list, which must start with l and end with e
			if i == n {
				return make([]string, 0), "", fmt.Errorf("cannot parse list from bencoded string %s. Reason: string has no character 'e' to mark list end", s)
			}
		} else if i == 'e' {
			break
		} else {
			i++
		}
	}

	if returnsRest {
		return list, s[start+1:], nil
	}
	return list, "", nil
}

func ParseDictionary(s string) (map[string]DictionaryElement, error) {
	//d3:cow3:moo4:spam4:eggse
	//     ^
	//d4:spaml1:a1:be

	dict := make(map[string]DictionaryElement)
	n := len(s)
	if s[0] != 'd' || s[n-1] != 'e' {
		return dict, fmt.Errorf("cannot parse bencoded string %s. Reason: Doesn't start with \"d\" and end with \"e\"", s)
	}

	// i := 1
	// currKey := ""
	// start := 1
	// var currentStringLength int64 = 0
	// for i < n-1 {
	// 	if s[i] == ':' {
	// 		// first we check from 'start' until colon ':', that's the currentStringLength
	// 		// and then we extract that string [i+1, i+1+currentStringLength]
	// 		// we check if it's a

	// 	}
	// }
	return dict, nil
}
