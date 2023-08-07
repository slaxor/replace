package replace

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// Subvert converts a string to the case of the match
// i = input string
// o = old string
// n = new string
// one restriction is that the case of the first word of the match is
// used for the replacement, so if you have a match of "foo Foo" it will
// be replaced with "bar bar" and not "bar Bar"
func Subvert(i, o, n string) string {
	o = normalize(o)
	n = normalize(n)
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
		isPrivate := unicode.IsLower(rune(wholeMatch[0]))
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
		if isPrivate {
			replacement = strings.ToLower(replacement[0:1]) + replacement[1:]
		}
		i = strings.Replace(i, wholeMatch, replacement, m[1])
	}
	return i
}

// SubvertBytes is a convenience function that uses []byte instead of
// string as input and output
// i = input []byte
// o = old string
// n = new string
// one restriction is that the case of the first word of the match is
// used for the replacement, so if you have a match of "foo Foo" it will
// be replaced with "bar bar" and not "bar Bar"
func SubvertBytes(i []byte, o, n string) []byte {
	return []byte(Subvert(string(i), o, n))
}

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

func isAllCaps(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

// SubvertFileContent replaces all occurrences of the old string with the new string in the file content.
func SubvertFileContent(filePath string, o string, n string) error {
	stat, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	perm := stat.Mode().Perm()
	file, err := os.OpenFile(filePath, os.O_RDWR, perm)
	if err != nil {
		return err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	newContent := Subvert(string(content), o, n)
	if string(content) == newContent {
		return nil
	}
	err = file.Truncate(0) // maybe a  bit dangerous
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = file.WriteString(newContent)
	if err != nil {
		return err
	}

	return nil
}
