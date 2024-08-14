package library

import "strings"

func cap(s string) string {
	result := ""
	for i, v := range s {
		if i != 0 && s[i-1] != ' ' {
			result = result + string(v)
			continue
		}
		if v >= 'a' && v <= 'z' {
			result = result + string(v-32)
			continue
		}
	}
	return result
}

func TextBeautifier(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", ": ")
	s = cap(s)
	return s
}
