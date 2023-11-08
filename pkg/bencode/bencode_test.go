package bencode

import (
	// "fmt"
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

