package trimmer

import (
	"fmt"
	"regexp"
	"strings"
)

type (
	TrimmedFields map[string][]string

	Trimmer interface {
		Trim(handler string, text []byte) []byte
	}

	trimmer struct {
		regexps map[string]map[string]*regexp.Regexp
	}
)

func New(trimmedFields TrimmedFields) Trimmer {
	var trim = &trimmer{regexps: make(map[string]map[string]*regexp.Regexp)}
	for handler, fields := range trimmedFields {
		var regexps = make(map[string]*regexp.Regexp, len(fields))
		for _, field := range fields {
			regexps[field] = regexp.MustCompile(fmt.Sprintf(
				`\\*"%s\\*":\s*\\*([^[]"{0,1}(?:[^"\\,}]|\\.)*\\*"{0,1}|\[[^]]*\])`,
				field,
			))
		}

		trim.regexps[strings.ToLower(handler)] = regexps
	}

	return trim
}

func (t *trimmer) Trim(handler string, text []byte) []byte {
	if len(text) == 0 {
		return text
	}

	var (
		regexps map[string]*regexp.Regexp
		ok      bool
	)
	if regexps, ok = t.regexps[strings.ToLower(handler)]; !ok {
		return text
	}

	for field, re := range regexps {
		text = re.ReplaceAll(text, []byte(fmt.Sprintf(`"%s": "TRIMMED_CONTENT"`, field)))
	}

	return text
}
