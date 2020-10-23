package markovmodel

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/OJOMB/trigram-markov-model/trigram"
)

var testTrigrams []trigram.Trigram = []trigram.Trigram{
	{Word1: "to", Word2: "be", Word3: "or"},
	{Word1: "be", Word2: "or", Word3: "not"},
	{Word1: "or", Word2: "not", Word3: "to"},
	{Word1: "not", Word2: "to", Word3: "be"},
	{Word1: "to", Word2: "be", Word3: "that"},
	{Word1: "be", Word2: "that", Word3: "is"},
	{Word1: "that", Word2: "is", Word3: "the"},
	{Word1: "is", Word2: "the", Word3: "question"},
}

func mockPRNG() func(int) int {
	counter := 0
	return func(n int) (deterministicInt int) {
		result := counter % n
		counter++
		return result
	}
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func TestNewSolution(t *testing.T) {
	s := New(rand.Intn)
	want := "*solution.Solution"
	got := fmt.Sprintf("%T", s)
	if want != got {
		t.Errorf("Expected s type to be %s, got: %s", want, got)
	}
	if len(s.freqs) != 0 {
		t.Errorf("Expected s.freqs to be empty. Instead found a populated map: %v", s.freqs)
	}
	want = "solution.prng"
	got = fmt.Sprintf("%T", s.prng)
	if want != got {
		t.Errorf("Expected s.prng type to be %s, got: %s", want, got)
	}
}

func TestAddTrigramToModel(t *testing.T) {
	s := New(mockPRNG())
	s.Add(trigram.Trigram{Word1: "to", Word2: "be", Word3: "or"})

	want := map[string]map[string]uint{
		"to be": {"or": 1},
	}

	eq := reflect.DeepEqual(s.freqs, want)
	if !eq {
		t.Errorf("Expected: %v, got: %v", want, s.freqs)
	}
}

func TestAddTrigramsToModel(t *testing.T) {
	s := New(mockPRNG())

	for _, tgram := range testTrigrams {
		s.Add(tgram)
	}

	want := map[string]map[string]uint{
		"to be":   {"or": 1, "that": 1},
		"be or":   {"not": 1},
		"or not":  {"to": 1},
		"not to":  {"be": 1},
		"be that": {"is": 1},
		"that is": {"the": 1},
		"is the":  {"question": 1},
	}

	eq := reflect.DeepEqual(s.freqs, want)

	if !eq {
		t.Errorf("Expected: %v, got: %v", want, s.freqs)
	}
}

func TestChoosAdjoiningWord(t *testing.T) {
	s := New(mockPRNG())
	m := map[string]uint{
		"or":  1,
		"not": 5,
		"to":  3,
		"be":  1,
	}
	expectedArray := []string{
		"be", "not", "not", "not", "not",
		"not", "or", "to", "to", "to",
		"be", "not", "not", "not", "not",
		"not", "or", "to", "to", "to",
	}
	for _, expected := range expectedArray {
		got := s.chooseAdjoiningWord(m)
		if got != expected {
			t.Errorf("Expected: %s, got: %s", expected, got)
		}
	}
}

func TestTableGetPrefixByNumericalIndex(t *testing.T) {
	// prefixes in expected alphabetical order
	// ["be or", "be that", "is the", "not to", "or not", "that is", "to be"]

	testTable := []struct {
		Input    int
		Expected string
	}{
		{0, "be or"},
		{1, "be that"},
		{2, "is the"},
		{3, "not to"},
		{4, "or not"},
		{5, "that is"},
		{6, "to be"},
	}

	s := New(mockPRNG())
	for _, tgram := range testTrigrams {
		s.Add(tgram)
	}

	for _, test := range testTable {
		prefix := s.getPrefixByNumericalIndex(test.Input)
		if prefix != test.Expected {
			t.Errorf("\nexpected: %s\ngot: %s", test.Expected, prefix)
		}
	}

}

func TestTableHandleDeadEnd(t *testing.T) {
	var testTable = []struct {
		numWords int
		input    []string
		expected []string
	}{
		{
			3,
			[]string{"test", "test"},
			[]string{"test", "test.", "amen"},
		},
		{
			4,
			[]string{"test", "test"},
			[]string{"test", "test.", "You", "heard"},
		},
		{
			10,
			[]string{"test", "test"},
			[]string{"test", "test.", "New", "prefix"},
		},
	}

	s := New(mockPRNG())
	s.Add(trigram.Trigram{Word1: "new", Word2: "prefix", Word3: "word3"})
	var got []string
	for _, test := range testTable {
		got = s.handleDeadEnd(test.input, test.numWords)
		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, got)
		}
	}
}

func TestGenerate(t *testing.T) {
	// the counter in the mockPRNG closure increments on each call
	// this means that each time we run s.Generate we start from a
	// different initial prefix and accordingly should receive different
	// output. this test runs s.Generate a few times, checking that the
	// resulting output is what it should be given the starting prefix
	// and the numWords.

	testTable := []struct {
		numWords int
		expected string
	}{
		{
			10,
			"Be or not to be or not to be or.",
		},
		{
			9,
			"Is the question. Or not. That is. You heard.",
		},
		{
			7,
			"To be or not to be or.",
		},
		{
			8,
			"That is the question. Be that. You heard.",
		},
	}

	s := New(mockPRNG())
	// check s.Generate throws error if model is empty
	got, err := s.Generate(5)
	if err.Error() != "Model is empty" {
		t.Errorf("Expected error message: 'Model is empty'. Got error: %v", err)
	}

	for _, tgram := range testTrigrams {
		s.Add(tgram)
	}

	for _, test := range testTable {
		got, _ = s.Generate(test.numWords)
		if test.expected != got {
			t.Errorf("Expected: %s, Got: %s", test.expected, got)
		}
	}
}
