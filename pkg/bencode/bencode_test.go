package bencode

import (
	// "fmt"
	"slices"
	"testing"
)

func TestParseIntShouldFail(t *testing.T) {
	values := []string{"ie", "123", "45e", "i33"}

	for _, v := range values {
		_, _, err := ParseInt(v)
		if err == nil {
			t.Fatalf("Got err = <nil> for ParseInt(%s). Want err != <nil>", v)
		}
	}
}

func TestParseInt(t *testing.T) {
	values := map[string][]int{
		"i3e":       {3, 2},
		"i-3e":      {-3, 3},
		"i304e":     {304, 4},
		"i101eblah": {101, 4},
	}

	for key, value := range values {

		result, stopIndex, err := ParseInt(key)

		if result != value[0] || stopIndex != value[1] || err != nil {
			t.Fatalf(`ParseInt("%s")=%d, %d, %v. Want %d, %d, <nil>.`, key, result, stopIndex, err, value[0], value[1])
		}
	}

}

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------

func TestParseStringShouldFail(t *testing.T) {
	values := []string{"5hello", ":hello", "hello", "&:a", "1a:hello world", "4:hi", "5:"}

	for _, v := range values {
		_, _, err := ParseString(v)
		if err == nil {
			t.Fatalf("Got err = <nil> for ParseString(%s). Want err != <nil>", v)
		}
	}
}

func TestParseString(t *testing.T) {
	values := map[string]string{
		"5:hello":         "hello",
		"0:":              "",
		"1:a":             "a",
		"12: hello world": " hello world",
		"5:helloworld":    "hello",
	}

	stopIndexes := map[string]int{
		"5:hello":         6,
		"0:":              1,
		"1:a":             2,
		"12: hello world": 14,
		"5:helloworld":    6,
	}

	for key, value := range values {

		result, stopIndex, err := ParseString(key)
		if result != value || stopIndex != stopIndexes[key] || err != nil {
			t.Fatalf(`ParseString("%s")=%s, %d, %v. Want %s, %d, <nil>.`, key, result, stopIndex, err, value, stopIndexes[key])
		}
	}
}

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------

func TestParseListShouldFail(t *testing.T) {
	values := []string{
		"l5:hello5:world",
		"7:console",
		"l1&:game of the tear",
		"l&:ae",
	}

	for _, v := range values {
		_, _, err := ParseList(v)
		if err == nil {
			t.Fatalf("Got err=<nil> for %s. Want err != <nil>", v)
		}
	}
}

func TestParseList(t *testing.T) {
	values := map[string][]string{
		"l5:hello4:it's2:an7:amazing5:worlde": {"hello", "it's", "an", "amazing", "world"},
		"l7:consolee":                         {"console"},
		"leKekW":                              {},
	}
	stopIndexes := map[string]int{
		"l5:hello4:it's2:an7:amazing5:worlde": 34,
		"l7:consolee":                         10,
		"leKekW":                              1,
		"l4:eggs2:in7:chick3:oute5:plate":     23,
	}

	for k, v := range values {
		result, stopIndex, err := ParseList(k)

		if slices.Compare(result, v) != 0 || stopIndex != stopIndexes[k] || err != nil {
			t.Fatalf("ParseList(%s)=%s, %d, %v. Want %s, %d, <nil>", k, result, stopIndex, err, v, stopIndexes[k])
		}
	}
}

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------

func TestParseDictionaryShouldFail(t *testing.T) {
	values := []string{
		"d5:hello5:world",
		"d7:consolee",
		"d7:consolesl3:ps43:ps5ee",
		"d1&:game of the year10:elden ringe",
	}

	for _, v := range values {
		_, _, err := ParseDictionary(v)
		if err == nil {
			t.Fatalf("Got err=<nil> for %s. Want err != <nil>", v)
		}
	}
}

func TestParseDictionary(t *testing.T){
	values := map[string]map[string]DictionaryElement{
		"d5:hello5:worlde": {
			"hello": DictionaryElement{Kind: STRING, String: "world"},
		},
		"d7:consoled4:name3:ps4ee": {
			"console": DictionaryElement{
				Kind: DICTIONARY,
				Dictionary: map[string]DictionaryElement{
					"name": {Kind: STRING, String: "ps4"},
				},
			},
		},
		"d8:consolesl3:ps43:ps5ee": {
			"consoles": DictionaryElement{Kind: LIST,List: []string{"ps4", "ps5"}},		},
		"d4:yeari2023ee": {
			"year": DictionaryElement{Kind: INTEGER, Integer: 2023},
		},
	}

	for k, v := range values {
		result, _, err := ParseDictionary(k)

		if !Equals(result, v) || err != nil  {
			t.Fatalf("ParseDictionary(%s)=%v %d, %v. Want %v, %d, <nil>", k, result, 0, err, v, 0)
		}
	}
}