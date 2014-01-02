package completion

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
	}{
		{"hello wo", "rld", []string{"hello", "wo"}},
		{"hello   wo", "rld", []string{"hello", "wo"}},
		{"hello 'there' wo", "rld", []string{"hello", "'there'", "wo"}},
		{"hello 'a b' wo", "rld", []string{"hello", "'a b'", "wo"}},
		{"hello 'a\\ b' wo", "rld", []string{"hello", "'a\\ b'", "wo"}},
		{"hello 'a\\' b' wo", "rld", []string{"hello", "'a\\' b'", "wo"}},
		{`" a string  " with words`, " in it", []string{`" a string  "`, "with", "words"}},
		{"a b ", "", []string{"a", "b", ""}},
		{"pw ", "", []string{"pw", ""}},
	}
	for _, tc := range testCases {
		line := tc.beforePoint + tc.afterPoint
		point := len(tc.beforePoint)
		st := parseLineForCompletion(line, point)
		c.Check(st.Words, DeepEquals, tc.words)
	}
}
