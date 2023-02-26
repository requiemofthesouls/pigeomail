package shortener

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type (
	LoggedFields map[string][]string

	Shortener interface {
		Shorten(handler string, text []byte) []byte
	}

	shortener struct {
		regexps map[string][]*regexp.Regexp
	}
)

func New(loggedFields LoggedFields) Shortener {
	var short = shortener{
		regexps: make(map[string][]*regexp.Regexp, len(loggedFields)),
	}
	for handler, fields := range loggedFields {
		var regexps = make([]*regexp.Regexp, 0, len(fields))
		for _, field := range fields {
			regexps = append(regexps, regexp.MustCompile(fmt.Sprintf(
				`\\*"%s\\*":\s*\\*([^[]"{0,1}(?:[^"\\,}]|\\.)*\\*"{0,1}|\[[^]]*\])`,
				field,
			)))
		}

		short.regexps[strings.ToLower(handler)] = regexps
	}

	return short
}

func (s shortener) Shorten(handler string, text []byte) []byte {
	if len(text) == 0 {
		return text
	}

	var (
		regexps []*regexp.Regexp
		ok      bool
	)
	if regexps, ok = s.regexps[strings.ToLower(handler)]; !ok {
		return text
	}

	if len(regexps) == 0 {
		return make([]byte, 0)
	}

	var shortenedParts [][]byte
	for _, re := range regexps {
		shortenedParts = append(shortenedParts, re.FindAll(text, -1)...)
	}

	if len(shortenedParts) == 0 {
		return text
	}

	return composeJSON(shortenedParts)
}

func composeJSON(fields [][]byte) []byte {
	var shortenedText = []byte("{")
	shortenedText = append(shortenedText, bytes.Join(fields, []byte(", "))...)

	return append(shortenedText, []byte("}")...)
}
