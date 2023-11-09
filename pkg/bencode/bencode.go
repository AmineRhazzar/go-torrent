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
	Kind       DictionaryElementType
	Integer    int
	String     string
	List       []string
	Dictionary map[string]DictionaryElement
}

// returns parsedInt, stopIndex, err
// stop index is the index of the "e" that indicates stop int parsing
// i3e => 3, 2, <nil>
// i-3e => -3, 3, <nil>
// i3eblah => 3, 2, <nil>
func ParseInt(s string) (int, int, error) {
	n := len(s)
	if n < 3 {
		return 0, 0, fmt.Errorf("cannot parse int from bencoded string %s", s)
	}
	if s[0] != 'i' {
		return 0, 0, fmt.Errorf("cannot parse int from bencoded string %s", s)
	}

	for i := 1; i < n; i++ {
		if s[i] == 'e' {
			p, err := strconv.ParseInt(s[1:i], 10, 64)
			return (int(p)), i, err
		}
	}
	return 0, 0, fmt.Errorf("cannot parse int from bencoded string %s. Reason: end character (e) not found", s)
}

func ParseString(s string) (string, int, error) {
	n := len(s)
	// get index of first occurence of : in string s
	var colonIndex int = -1
	var i int = 0
	for i = 0; i < n; i++ {
		if string(s[i]) == ":" {
			colonIndex = i
			break
		}
	}
	if colonIndex == -1 {
		return "", 0, fmt.Errorf("cannot parse bencoded string %s. Reason: can't find \":\"", s)
	}

	stringLength, err := strconv.ParseInt(s[0:colonIndex], 10, 32)
	if err != nil {
		return "", 0, fmt.Errorf("cannot parse bencoded string %s. Reason: can't parse length", s)
	}

	// check if there's enough length after colon (minimum length of stringLength)
	stop := colonIndex + int(stringLength)
	if stop > n {
		return "", 0, fmt.Errorf("cannot parse bencoded string %s. Reason: string not long enough (must be: %d)", s, stringLength)
	}

	return s[colonIndex+1 : stop+1], stop, nil

}

func ParseList(s string) ([](string), int, error) {
	n := len(s)

	if s[0] != 'l' {
		return make([]string, 0), 0, fmt.Errorf("cannot parse list from bencoded string %s", s)
	}

	var list []string
	var err error

	// stores the parsed string length. Example: "5:hello" -> currentStringLength=5
	var currentStringLength int64 = 0
	// index iterating through whole string s
	i := 1
	// indicates the index where the currentStringLength starts at in string s
	start := 1
	stopIndex := 0
	for i < n {
		if s[i] == ':' {
			currentStringLength, err = strconv.ParseInt(s[start:i], 10, 32)
			if err != nil {
				return make([]string, 0), 0, fmt.Errorf("cannot parse list from bencoded string %s. Reason: error parsing string length from interval [%d, %d[", s, start, i-1)
			}
			// so we parse from i+1 (directly after colon) until i+1+currentStringLength
			list = append(list, s[i+1:i+1+int(currentStringLength)])
			i = i + 1 + int(currentStringLength)
			start = i
			// if i is at the end of string s now (i==n) we return error
			// we parsed the string successfully but there's no 'e' at the end to the list, which must start with l and end with e
			if i == n {
				return make([]string, 0), 0, fmt.Errorf("cannot parse list from bencoded string %s. Reason: string has no character 'e' to mark list end", s)
			}
		} else if s[i] == 'e' {
			stopIndex = i
			break
		} else {
			i++
		}
	}

	return list, stopIndex, nil
}

/**
d3:cow3:moo4:spam4:eggs3:qtei3ee
      ^
- i=0
  push to stack: [d]
- i=1
  s[i] != 'l', 'd', 'e', 'i' => must be a string
  parsedString, stopIndex = parseString(s[i:]) = cow, 5
  if currKey == "" currKey = parsedString
  else dict[currKey] = DictionaryElement{..., parsedString,...}
  i+=stopIndex
- i=6


*/

func ParseDictionary(s string) (map[string]DictionaryElement, int, error) {
	emptyDict := make(map[string]DictionaryElement)
	dict := emptyDict

	if s[0] != 'd' {
		return dict, 0, fmt.Errorf("can't parse dictionary from string %s", s)
	}
	n := len(s)
	currKey := ""
	i, finalStopIndex := 1, -1
	for i < n {
		if s[i] == 'i' {
			// must be an int
			integer, j, err := ParseInt(s[i:])
			if err != nil {
				return emptyDict, 0, fmt.Errorf("can't parse dictionary from string %s. Reason: %v", s, err)
			}
			if currKey == "" {
				return emptyDict, 0, fmt.Errorf("can't parse dictionary from string %s. Reason: No key for int %d (position: %d)", s, integer, i)
			} else {
				dict[currKey] = DictionaryElement{
					Kind:    INTEGER,
					Integer: integer,
				}
				currKey = ""
			}
			i += j+1
		} else if s[i] == 'l' {
			list, j, err := ParseList(s[i:])
			if err != nil {
				return emptyDict, 0, fmt.Errorf("can't parse dictionary from string %s. Reason: %v", s, err)
			}
			if currKey == "" {
				return emptyDict, 0, fmt.Errorf("can't parse dictionary from string %s. Reason: No key for list %s (position: %d)", s, list, i)
			} else {
				dict[currKey] = DictionaryElement{
					Kind: LIST,
					List: list,
				}
				currKey = ""
			}
			i += j+1
		} else if s[i] == 'd' {
			innerDict, j, err := ParseDictionary(s[i:])
			if err != nil {
				return emptyDict, 0, fmt.Errorf("can't parse dictionary from string %s. Reason: %v", s, err)
			}
			if currKey == "" {
				return emptyDict, 0, fmt.Errorf("can't parse dictionary from string %s. Reason: No key for dictionary %+v (position: %d)", s, innerDict, i)
			} else {
				dict[currKey] = DictionaryElement{
					Kind: DICTIONARY,
					Dictionary: innerDict,
				}
				currKey = ""
			}
			i += j+1
		} else if s[i] == 'e' {
			finalStopIndex = i
			i++
		} else {
			// must be a string
			str, j, err := ParseString(s[i:])
			if err != nil {
				return emptyDict, 0, fmt.Errorf("can't parse dictionary from string %s. Reason: %v", s, err)
			}
			if currKey == "" {
				currKey = str
			} else {
				dict[currKey] = DictionaryElement{
					Kind:   STRING,
					String: str,
				}
				currKey = ""
			}
			i += j+1
		}
	}

	if finalStopIndex == -1 {
		return emptyDict, 0, fmt.Errorf("can't parse dictionary from string %s. Reason: can't find stop character 'e' for dictionary", s)
	}
	return dict, finalStopIndex, nil
}
