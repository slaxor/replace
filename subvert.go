package replace

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

func Subvert(oldName, newName string) func(string) string {
	// oldSnake := ToSnakeCase(oldName)
	// newSnake := ToSnakeCase(newName)
	// oldCamel := ToCamelCase(oldName)
	// newCamel := ToCamelCase(newName)
	return func(s string) string {
		// s = strings.ReplaceAll(s, oldName, newName)
		// s = strings.ReplaceAll(s, strings.Title(oldName), strings.Title(newName))
		// s = strings.ReplaceAll(s, strings.ToUpper(oldName), strings.ToUpper(newName))
		// s = strings.ReplaceAll(s, strings.ToLower(oldName), strings.ToLower(newName))
		// s = strings.ReplaceAll(s, oldSnake, newSnake)
		// s = strings.ReplaceAll(s, oldCamel, newCamel)
		s = ToMulti(s, oldName, newName)
		return s
	}
}

// func init() {
// log.SetFlags(log.Lshortfile)
// }

func normalize(s string) string {
	s = strings.ReplaceAll(s, " ", "_")
	if strings.Contains(s, "_") {
		return s
	}
	var words []string
	var lastPos int
	rs := []rune(s)

	for i := 0; i < len(rs); i++ {
		if i > 0 && unicode.IsUpper(rs[i]) {
			words = append(words, s[lastPos:i])
			lastPos = i
		}
	}

	words = append(words, s[lastPos:])
	return strings.Join(words, "_")
}

// ToSnakeCase converts a camel case string to snake case
// func ToSnakeCase(s string) string {
// return strings.ToLower(normalize(s))
// }

// func _ToSnakeCase(s string) string {
// var result string
// var words []string
// var lastPos int
// rs := []rune(s)

// for i := 0; i < len(rs); i++ {
// if i > 0 && unicode.IsUpper(rs[i]) {
// words = append(words, s[lastPos:i])
// lastPos = i
// }
// }

// words = append(words, s[lastPos:])

// for k, word := range words {
// if k > 0 {
// result += "_"
// }
// result += strings.ToLower(word)
// }

// return result
// }

// ToCamelCase converts a snake case string to camel case
// func ToCamelCase(s string) string {
// words := strings.FieldsFunc(s, func(r rune) bool { return r == '_' })

// for i := 0; i < len(words); i++ {
// words[i] = strings.Title(words[i])
// }

// return strings.Join(words, "")
// }

func isAllCaps(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

// ToMulti converts a string to the case of the match
// i = input string
// o = old string
// n = new string
// one restriction is that the case of the first word of the match is
// used for the replacement, so if you have a match of "foo Foo" it will
// be replaced with "bar bar" and not "bar Bar"
func ToMulti(i, o, n string) string {
	o = normalize(o)
	n = normalize(n)
	o = regexp.QuoteMeta(o) // I ask myself if it would be a nice feature to be able to pass a regexp as old string
	ows := strings.Split(o, `_`)
	e := strings.Join(ows, `)([[:space:]_]*)(`)
	e = fmt.Sprintf("(?ims)(%s)", e)
	re := regexp.MustCompile(e)

	nws := strings.Split(n, `_`)

	ib := []byte(i)
	ms := re.FindAllSubmatchIndex(ib, -1)

	for _, m := range ms {
		var v []byte
		var joiner string
		wholeMatch := string(ib[m[0]:m[1]])
		var replacement string
		for j := 2; j < len(m); j += 2 {
			v = ib[m[j]:m[j+1]]
			if len(v) == 0 { // empty joiner
				continue
			}
			if unicode.IsSpace(rune(v[0])) || v[0] == '_' {
				joiner = string(v)
				continue
			}
			isTitle := unicode.IsUpper(rune(v[0]))
			isUpper := true
			for _, c := range v {
				isUpper = isUpper && unicode.IsUpper(rune(c))
			}
			for nwi := range nws {
				switch {
				case isUpper:
					nws[nwi] = strings.ToUpper(nws[nwi])
				case isTitle:
					nws[nwi] = strings.Title(nws[nwi])
				default:
					nws[nwi] = strings.ToLower(nws[nwi])
				}
			}
		}
		replacement = strings.Join(nws, joiner)

		i = strings.Replace(i, wholeMatch, replacement, m[1])
	}

	return i
}
