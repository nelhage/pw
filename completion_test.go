package main

import (
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type CompletionSuite struct{}

var _ = Suite(&CompletionSuite{})

func (s *CompletionSuite) TestParseLine(c *C) {
	testCases := []struct {
		beforePoint string
		afterPoint  string
		words       []string
		word        int
	}{
		{"hello wo", "rld", []string{"hello", "world"}, 1},
		{"hello   wo", "rld", []string{"hello", "world"}, 1},
		{"hello 'there' wo", "rld", []string{"hello", "'there'", "world"}, 2},
		{"hello 'a b' wo", "rld", []string{"hello", "'a b'", "world"}, 2},
		{"hello 'a\\ b' wo", "rld", []string{"hello", "'a\\ b'", "world"}, 2},
		{"hello 'a\\' b' wo", "rld", []string{"hello", "'a\\' b'", "world"}, 2},
		{`" a string  " with words`, " in it", []string{`" a string  "`, "with", "words", "in", "it"}, 2},
		{"a b ", "", []string{"a", "b", ""}, 2},
		{"pw ", "", []string{"pw", ""}, 1},
	}
	for _, tc := range testCases {
		line := tc.beforePoint + tc.afterPoint
		point := len(tc.beforePoint)
		st := parseLineForCompletion(line, point)
		c.Check(st.words, DeepEquals, tc.words)
		c.Check(st.word, Equals, tc.word)
	}
}
