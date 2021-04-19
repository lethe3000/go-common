package txt

import "strings"

// SplitWords 获取数组的所有子字符串
func SplitWords(s string) []string {
	if len(s) == 0 {
		return make([]string, 0)
	}
	sliced := strings.Split(s, "")
	slicedLen := len(sliced)

	substrings := make(map[string]string)
	for substringLen := 1; substringLen <= len(sliced); substringLen++ {
		for idx := 0; idx+substringLen <= slicedLen; idx++ {
			sb := strings.Join(sliced[idx:idx+substringLen], "")
			substrings[sb] = sb
		}
	}

	splitted := make([]string, 0, len(substrings))
	for k := range substrings {
		splitted = append(splitted, k)
	}
	return splitted
}
