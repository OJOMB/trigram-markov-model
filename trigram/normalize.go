package trigram

import "strings"

// Normalise makes the string lower case and removes any special characters and
// numbers.
func Normalise(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(
		func(r rune) rune {
			if r < 'a' || r > 'z' {
				return -1
			}
			return r
		},
		s,
	)

	return s
}

// NormaliseSlice normalises a slice of strings. Elements that are normalised to
// empty strings are discarded.
func NormaliseSlice(w []string) []string {
	var newWords []string
	for i := range w {
		word := Normalise(w[i])
		if word != "" {
			newWords = append(newWords, word)
		}
	}
	return newWords
}
