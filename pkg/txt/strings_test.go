package txt

import (
	"fmt"
	"sort"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestBool(t *testing.T) {
	t.Run("utf8", func(t *testing.T) {
		cn := "中文"
		en := "en"
		v, w := utf8.DecodeRuneInString(cn)
		fmt.Printf("%#U %d\n", v, w)
		v, w = utf8.DecodeRuneInString(cn[w:])
		fmt.Printf("%#U %d\n", v, w)

		v, w = utf8.DecodeRuneInString(en)
		fmt.Printf("%#U %d\n", v, w)
		v, w = utf8.DecodeRuneInString(en[w:])
		fmt.Printf("%#U %d\n", v, w)
		fmt.Println(strings.Split(cn, ""))
	})
	t.Run("substrings", func(t *testing.T) {
		substrings := SplitWords("abc")
		sort.Strings(substrings)
		assert.True(t, equal(substrings, []string{"a", "ab", "abc", "b", "bc", "c"}))

		substrings = SplitWords("abcd")
		sort.Strings(substrings)
		assert.True(t, equal(substrings, []string{"a", "ab", "abc", "abcd", "b", "bc", "bcd", "c", "cd", "d"}))

		substrings = SplitWords("abcb")
		sort.Strings(substrings)
		assert.True(t, equal(substrings, []string{"a", "ab", "abc", "abcb", "b", "bc", "bcb", "c", "cb"}))

		substrings = SplitWords("北京")
		sort.Strings(substrings)
		assert.True(t, equal(substrings, []string{"京", "北", "北京"}))
	})
}
