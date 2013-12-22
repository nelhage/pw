package gpg

import (
	"bytes"
	. "launchpad.net/gocheck"
	"testing"
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
		"kumquat"}, []byte("Some test data"))
	c.Assert(err, IsNil)
	c.Assert(bytes.HasPrefix(out, []byte("-----BEGIN PGP MESSAGE-----\n")),
		Equals, true)
	c.Assert(bytes.HasSuffix(out, []byte("\n-----END PGP MESSAGE-----\n")),
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
		[]byte("Some test data"))
	c.Assert(err, NotNil)
	c.Assert(out, IsNil)
}
