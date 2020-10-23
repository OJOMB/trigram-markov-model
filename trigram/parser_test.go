package trigram

import (
	"bytes"
	"reflect"
	"testing"
)

var expectedTrigrams []Trigram = []Trigram{
	{Word1: "this", Word2: "is", Word3: "a"},
	{Word1: "is", Word2: "a", Word3: "test"},
	{Word1: "a", Word2: "test", Word3: "i"},
	{Word1: "test", Word2: "I", Word3: "repeat"},
	{Word1: "I", Word2: "repeat", Word3: "this"},
	{Word1: "repeat", Word2: "this", Word3: "is"},
	{Word1: "this", Word2: "is", Word3: "a"},
	{Word1: "is", Word2: "a", Word3: "test"},
}

func TestParse(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("This is a test, I repeat! This is a test !!")
	content, err := Parse(&buffer)
	if err != nil {
		t.Errorf("Failed to read from test buffer: %s", err.Error())
	}
	expected := []string{
		"This", "is", "a", "test,", "I", "repeat!",
		"This", "is", "a", "test", "!!",
	}

	if !reflect.DeepEqual(expected, content) {
		t.Errorf("Expected: %v, got: %v", expected, content)
	}

	buffer.Reset()
	buffer.WriteString("A Test")
	_, err = Parse(&buffer)
	expectedErr := "input has less than 3 words"
	if err.Error() != expectedErr {
		t.Errorf("Expected error: %s, got: %v", expectedErr, err)
	}
}

func TestTrigramsFromSlice(t *testing.T) {
	input := []string{
		"This", "is", "a", "test,", "I", "repeat!",
		"This", "is", "a", "test", "!!",
	}
	got := trigramsFromSlice(input)

	if reflect.DeepEqual(got, expectedTrigrams) {
		t.Errorf("Expected %v, got %v", expectedTrigrams, got)
	}
}

func TestParseFileToNormalisedTrigrams(t *testing.T) {
	got, err := ParseFileToNormalisedTrigrams("testdata/test.txt")
	if err != nil {
		t.Errorf("Received unexpected error: %v", err)
	}
	if reflect.DeepEqual(got, expectedTrigrams) {
		t.Errorf("Expected %v, got %v", expectedTrigrams, got)
	}
}
