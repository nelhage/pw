package gpg

import (
	"strings"
	"testing"

	. "launchpad.net/gocheck"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type GPGSuite struct{}

var _ = Suite(&GPGSuite{})

func (s *GPGSuite) TestPipe(c *C) {
	out, err := runGPGPipe([]string{"--encrypt",
		"--symmetric",
		"--batch",
		"--no-use-agent",
		"--armor",
		"--passphrase",
		"kumquat"}, "Some test data")
	c.Assert(err, IsNil)
	c.Assert(strings.HasPrefix(out, "-----BEGIN PGP MESSAGE-----\n"),
		Equals, true)
	c.Assert(strings.HasSuffix(out, "\n-----END PGP MESSAGE-----\n"),
		Equals, true)
}

func (s *GPGSuite) TestError(c *C) {
	out, err := runGPGPipe([]string{"--encrypt",
		"--symmetric",
		"--batch",
		"--no-use-agent",
		"--armor",
		"--passphrase",
		"kumquat",
		"--no-such-argument"},
		"Some test data")
	c.Assert(err, NotNil)
	c.Assert(out, Equals, "")
}
