package markovmodel

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/OJOMB/trigram-markov-model/trigram"
)

// pseudo-random number generator from 0 to max
type prng func(max int) (randomInt int)

// Model is an object capable of accepting trigrams and generating novel text
// using a basic Markov Model
type Model struct {
	freqs map[string]map[string]uint
	prng  prng
}

// New creates a new Model
func New(prng func(n int) int) *Model {
	return &Model{
		freqs: make(map[string]map[string]uint),
		prng:  prng,
	}
}

// getPrefixByNumericalIndex takes an integer and returns a prefix (key) from the
// the internal frequency table
func (s *Model) getPrefixByNumericalIndex(i int) string {
	// s.freqs should always be populated by the time this method is called
	i %= len(s.freqs)

	// add map keys to slice
	keys := make([]string, len(s.freqs))
	j := 0
	for k := range s.freqs {
		keys[j] = k
		j++
	}
	// sort slice so that return is deterministic if we know s.prng(len(keys))
	sort.Strings(keys)

	return keys[i]
}

// chooseAdjoiningWord chooses a key from a map at random based on the values
// which are treated as weights
func (s *Model) chooseAdjoiningWord(options map[string]uint) string {
	var choices []string
	for k, v := range options {
		var i uint = 0
		for ; i < v; i++ {
			choices = append(choices, k)
		}
	}
	// sort for deterministic behaviour
	sort.Strings(choices)

	return choices[s.prng(len(choices))]
}

// handleDeadEnd is a playful solution to the problem of dead-ending when the
// constructed prefix has no precedent in the input corpus. This is only really an issue
// for relatively small input corpora
func (s *Model) handleDeadEnd(words []string, numWords int) []string {
	words[len(words)-1] = fmt.Sprintf("%s.", words[len(words)-1])

	var appendage []string
	if len(words) <= numWords-3 {
		randInt := s.prng(len(s.freqs))
		appendage = strings.Split(s.getPrefixByNumericalIndex(randInt), " ")
		appendage[0] = strings.Title(appendage[0])
	} else if len(words) == numWords-2 {
		appendage = []string{"You", "heard"}
	} else {
		appendage = []string{"amen"}
	}
	return append(words, appendage...)
}

// Add stores a new trigram into the model via updating the internal frequency table
func (s *Model) Add(t trigram.Trigram) {
	prefix := fmt.Sprintf("%s %s", t.Word1, t.Word2)
	freqs, ok := s.freqs[prefix]
	if !ok {
		s.freqs[prefix] = map[string]uint{t.Word3: 1}
		return
	}

	if _, ok = freqs[t.Word3]; !ok {
		freqs[t.Word3] = 1
	} else {
		freqs[t.Word3]++
	}

	s.freqs[prefix] = freqs
}

// Generate returns a novel text generated using a simple Markov Chain
func (s *Model) Generate(numWords int) (string, error) {
	if len(s.freqs) == 0 {
		return "", errors.New("Model is empty")
	}
	randInt := s.prng(len(s.freqs))
	randPrefix := s.getPrefixByNumericalIndex(randInt)
	var words []string = append([]string{}, strings.Split(randPrefix, " ")...)

	if numWords == 1 {
		return strings.Title(words[0]) + ".", nil
	}

	for {
		if len(words) == numWords {
			break
		}
		prefix := strings.Join(words[len(words)-2:], " ")
		freqs, ok := s.freqs[prefix]
		if !ok {
			// if the constructed prefix has no entry in the frequencies map
			// it's a deadend and should be handled as such
			words = s.handleDeadEnd(words, numWords)
		} else {
			words = append(words, s.chooseAdjoiningWord(freqs))
		}
	}

	// capitalise the first word
	words[0] = strings.Title(words[0])

	return strings.Join(words, " ") + ".", nil
}
