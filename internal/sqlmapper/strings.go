package sqlmaper

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

func camelCase(str string) string {
	return camelCaseToLowerCase(str, '_')
}
func camelCaseToLowerCase(str string, connector rune) string {
	if len(str) == 0 {
		return ""
	}

	buf := &bytes.Buffer{}
	var prev, r0, r1 rune
	var size int

	r0 = connector

	for len(str) > 0 {
		prev = r0
		r0, size = utf8.DecodeRuneInString(str)
		str = str[size:]

		switch {
		case r0 == utf8.RuneError:
			buf.WriteRune(r0)

		case unicode.IsUpper(r0):
			if prev != connector && !unicode.IsNumber(prev) {
				buf.WriteRune(connector)
			}

			buf.WriteRune(unicode.ToLower(r0))

			if len(str) == 0 {
				break
			}

			r0, size = utf8.DecodeRuneInString(str)
			str = str[size:]

			if !unicode.IsUpper(r0) {
				buf.WriteRune(r0)
				break
			}

			// find next non-upper-case character and insert connector properly.
			// it's designed to convert `HTTPServer` to `http_server`.
			// if there are more than 2 adjacent upper case characters in a word,
			// treat them as an abbreviation plus a normal word.
			for len(str) > 0 {
				r1 = r0
				r0, size = utf8.DecodeRuneInString(str)
				str = str[size:]

				if r0 == utf8.RuneError {
					buf.WriteRune(unicode.ToLower(r1))
					buf.WriteRune(r0)
					break
				}

				if !unicode.IsUpper(r0) {
					if r0 == '_' || r0 == ' ' || r0 == '-' {
						r0 = connector

						buf.WriteRune(unicode.ToLower(r1))
					} else if unicode.IsNumber(r0) {
						// treat a number as an upper case rune
						// so that both `http2xx` and `HTTP2XX` can be converted to `http_2xx`.
						buf.WriteRune(unicode.ToLower(r1))
						buf.WriteRune(connector)
						buf.WriteRune(r0)
					} else {
						buf.WriteRune(connector)
						buf.WriteRune(unicode.ToLower(r1))
						buf.WriteRune(r0)
					}

					break
				}

				buf.WriteRune(unicode.ToLower(r1))
			}

			if len(str) == 0 || r0 == connector {
				buf.WriteRune(unicode.ToLower(r0))
			}

		case unicode.IsNumber(r0):
			if prev != connector && !unicode.IsNumber(prev) {
				buf.WriteRune(connector)
			}

			buf.WriteRune(r0)

		default:
			if r0 == ' ' || r0 == '-' || r0 == '_' {
				r0 = connector
			}

			buf.WriteRune(r0)
		}
	}

	return buf.String()
}
