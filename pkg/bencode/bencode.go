package bencode

import (
	"fmt"
	"slices"
	"strconv"
	// "github.com/davecgh/go-spew/spew"
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

func ParseDictionary(s string) (map[string]DictionaryElement, int, error) {
	emptyDict := make(map[string]DictionaryElement)
	dict := emptyDict

	if s[0] != 'd' {
		return emptyDict, 0, fmt.Errorf("can't parse dictionary from string %s", s)
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
			i += j + 1
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
			i += j + 1
		} else if s[i] == 'd' {
			innerDict, j, err := ParseDictionary(s[i:])
			if err != nil {
				return emptyDict, 0, fmt.Errorf("can't parse dictionary from string %s. Reason: %v", s, err)
			}
			if currKey == "" {
				return emptyDict, 0, fmt.Errorf("can't parse dictionary from string %s. Reason: No key for dictionary %+v (position: %d)", s, innerDict, i)
			} else {
				dict[currKey] = DictionaryElement{
					Kind:       DICTIONARY,
					Dictionary: innerDict,
				}
				currKey = ""
			}
			i += j + 1
		} else if s[i] == 'e' {
			finalStopIndex = i
			if currKey != "" {
				// there's is a key with no corresponding value
				return emptyDict, 0, fmt.Errorf("can't parse dictionary from string %s. Reason: found key without corresponding value", s)
			}
			break
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
			i += j + 1
		}
	}

	if finalStopIndex == -1 {
		return emptyDict, 0, fmt.Errorf("can't parse dictionary from string %s. Reason: can't find stop character 'e' for dictionary", s)
	}
	return dict, finalStopIndex, nil
}

type jsonElement struct {
	k     string
	v     string
	depth int
}

func indent(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += " "
	}
	return s
}

func jsonifyDictionary(dict map[string]DictionaryElement, depth int) []jsonElement {
	var elements []jsonElement
	for key, val := range dict {
		switch val.Kind {
		case INTEGER:
			elements = append(elements, jsonElement{
				k: key, v: fmt.Sprintf("%d", val.Integer), depth: depth,
			})
		case STRING:
			elements = append(elements, jsonElement{
				k: key, v: val.String, depth: depth,
			})
		case LIST:
			elements = append(elements, jsonElement{
				k: key, v: fmt.Sprintf("%s", val.List), depth: depth,
			})
		case DICTIONARY:
			elements = append(elements, jsonElement{
				k: key, v: stringifyJson(jsonifyDictionary(val.Dictionary, depth+1)), depth: depth,
			})
		}
	}
	return elements
}

func stringifyJson(elements []jsonElement) string {
	numSpacesIn := elements[0].depth * 4
	numSpacesOut := (elements[0].depth - 1) * 4

	s := "\n" + indent(numSpacesOut) + "{\n"
	for _, el := range elements {
		s += indent(numSpacesIn) + el.k + ": " + el.v + "\n"
	}
	s += indent(numSpacesOut) + "}"
	return s
}

func PrintDictionary(dict map[string]DictionaryElement) {
	fmt.Println(stringifyJson(jsonifyDictionary(dict, 1)))
}

func Equals(dict1 map[string]DictionaryElement, dict2 map[string]DictionaryElement) bool {
	for k1, v1 := range dict1 {
		v2, ok := dict2[k1]
		if !ok || v1.Kind != v2.Kind {
			return false
		}

		switch v1.Kind {
		case INTEGER:
			if v1.Integer != v2.Integer {
				return false
			}
		case STRING:
			if v1.String != v2.String {
				return false
			}
		case LIST:
			if slices.Compare(v1.List, v2.List) != 0 {
				return false
			}
		case DICTIONARY:
			if !Equals(v1.Dictionary, v2.Dictionary) {
				return false
			}
		}
	}
	return true
}
