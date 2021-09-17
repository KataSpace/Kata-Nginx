package KataGin

import (
	"fmt"
	"regexp"
	"strings"
)

func defaultConvert(s string) string {

	return s
}

func SlashConvertWithOne(s string) string {
	return slashConvert(s, 1)
}

func SlashConvertWithTwo(s string) string {
	return slashConvert(s, 2)
}

func SlashConvertWithThree(s string) string {
	return slashConvert(s, 3)
}

func SlashConvertWithFive(s string) string {
	return slashConvert(s, 5)
}

func SlashConvertWithTen(s string) string {
	return slashConvert(s, 10)
}

// slashConvert  add slash before every upper char
// n  Capital letters with less than N consecutive characters are ignored
// e.g.
// n = 5 GetAPIAllName will get Get、All、Name. API is three consecutive capital letters, which is less than the requirement of 5, so it will be ignored
// n = 2 GetAPIAllName will get Get、AP、IA、Name.
// more examples please referee convert_test.go
func slashConvert(s string, n int) string {
	r := "[A-Z][^A-Z]+"
	if n >= 1 {
		r = fmt.Sprintf("[A-Z][a-z]+|([A-Z]|[0-9]){%d}", n)
	}

	reg := regexp.MustCompile(r)

	sub := reg.FindAllString(s, -1)

	return strings.Join(sub, "/")
}

func defaultGetMethods(s string) (method string, name string) {

	var b strings.Builder
	b.Grow(len(s))

	b.WriteByte(s[0])
	for i := 1; i < len(s); i++ {
		c := s[i]
		if c == '/' || 'A' <= c && c <= 'Z' {
			name = s[i:]
			break
		}
		b.WriteByte(c)
	}

	return b.String(), name
}
