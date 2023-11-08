package bencode

import (
	// "fmt"
	"slices"
	"strconv"
	"testing"
)

func TestParseIntShouldFail(t *testing.T) {
	values := []string{"ie", "123", "45e", "i33"}

	for _, v := range values {
		_, _, err := ParseInt(v, false)
		if err == nil {
			t.Fatalf("Got err = <nil> for ParseInt(%s). Want err != <nil>", v)
		}
	}
}

func TestParseIntWithoutReturn(t *testing.T) {
	values := map[string]int64{
		"i3e":   3,
		"i-3e":  -3,
		"i304e": 304,
	}

	for key, value := range values {

		result, rest, err := ParseInt(key, false)

		if value != result || rest != "" || err != nil {
			t.Fatalf(`ParseInt("%s")=%d, %v. Want %d, <nil>.`, key, result, err, value)
		}
	}

}

func TestParseIntWithReturn(t *testing.T) {
	values := map[string][]string{
		"i3ei5e":      {"3", "i5e"},
		"i-3e5:hello": {"-3", "5:hello"},
		"i304eelf":    {"304", "elf"},
	}

	for key, value := range values {
		result, rest, err := ParseInt(key, true)
		ExpectedResult, ExpectedRest := value[0], value[1]

		if strconv.Itoa(int(result)) != ExpectedResult || rest != ExpectedRest || err != nil {
			t.Fatalf(`ParseInt("%s")=%d, %s, %v. Want %s, %s, <nil>.`, key, result, rest, err, ExpectedResult, ExpectedRest)
		}
	}

}

// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------

func TestParseStringShouldFail(t *testing.T) {
	values := []string{"5hello", ":hello", "hello", "&:a", "1a:hello world", "4:hi", "5:"}

	for _, v := range values {
		_, _, err := ParseString(v, false)
		if err == nil {
			t.Fatalf("Got err = <nil> for ParseString(%s). Want err != <nil>", v)
		}
	}
}

func TestParseStringWithoutReturn(t *testing.T) {
	values := map[string]string{
		"5:hello":         "hello",
		"0:":              "",
		"1:a":             "a",
		"12: hello world": " hello world",
	}

	for key, value := range values {

		result, rest, err := ParseString(key, false)
		if value != result || rest != "" || err != nil {
			t.Fatalf(`ParseString("%s")=%s, %s, %v. Want %s, , <nil>.`, key, result, rest, err, value)
		}
	}
}

func TestParseStringWithReturn(t *testing.T) {
	values := map[string][]string{
		"5:grape":              {"grape", ""},
		"0:plane":              {"", "plane"},
		"6:panic at the disco": {"panic ", "at the disco"},
	}

	for key, value := range values {

		result, rest, err := ParseString(key, true)
		ExpectedResult, ExpectedRest := value[0], value[1]

		if result != ExpectedResult || rest != ExpectedRest || err != nil {
			t.Fatalf(`ParseString("%s")=%s, %s, %v. Want %s, %s, <nil>.`, key, result, rest, err, ExpectedResult, ExpectedRest)
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
		_, _, err := ParseList(v, false)
		if err == nil {
			t.Fatalf("Got err=<nil> for %s. Want err != <nil>", v)
		}
	}
}

func TestParseListWithoutReturn(t *testing.T) {
	values := map[string][]string{
		"l5:hello4:it's2:an7:amazing5:worlde": {"hello", "it's", "an", "amazing", "world"},
		"l7:consolee":                         {"console"},
		"le":                                  {},
	}

	for k, v := range values {
		result, rest, err := ParseList(k, false)

		if slices.Compare(result, v) != 0 || rest != "" || err != nil {
			t.Fatalf("ParseList(%s)=%s, %s, %v. Want %s, , <nil>", k, result, rest, err, v)
		}
	}
}

func TestParseListWithReturn(t *testing.T) {
	values := map[string][]string{
		"l5:hello4:it's2:an7:amazing5:worldeblahblah": {"hello", "it's", "an", "amazing", "world"},
		"l7:consoleei32e": {"console"},
		"lekekw":          {},
	}

	rests := map[string]string{
		"l5:hello4:it's2:an7:amazing5:worldeblahblah": "blahblah",
		"l7:consoleei32e": "i32e",
		"lekekw":          "kekw",
	}

	for k, v := range values {
		result, rest, err := ParseList(k, true)

		if slices.Compare(result, v) != 0 || rest != rests[k] || err != nil {
			t.Fatalf("ParseList(%s)=%s, %s, %v. Want %s, %s, <nil>", k, result, rest, err, v, rests[k])
		}
	}
}
