package trigram

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

// ParseFileToNormalisedTrigrams parses all trigrams from a given file.
func ParseFileToNormalisedTrigrams(path string) ([]Trigram, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	words, err := Parse(f)
	if err != nil {
		return nil, err
	}
	return trigramsFromSlice(words), nil
}

// Parse uses buffered I/O to collate a slice of words from an io.Reader
func Parse(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	var words []string
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		words = append(words, word)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(words) < 3 {
		return nil, errors.New("input has less than 3 words")
	}
	return words, nil
}

// trigramsFromSlice takes a slice of strings
// normalises that slice and returns a slice of trigrams
func trigramsFromSlice(words []string) []Trigram {
	words = NormaliseSlice(words)

	var trigrams []Trigram
	for i := 0; i < len(words)-2; i++ {
		trigram := Trigram{words[i], words[i+1], words[i+2]}
		trigrams = append(trigrams, trigram)
	}
	return trigrams
}
