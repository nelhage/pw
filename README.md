# pw -- A simple GPG-based secret store

`pw` is a simple GPG-based secret store. It's designed for storing
passwords or other secrets, using GPG to encrypt them for security. It
relies on GPG for all the crypto, including relying on gpg-agent to do
any key storage or UI.

## Synopsis

    # Installation
    go get github.com/nelhage/pw

    # Prompt for a password to save as "google.com"
    pw add google.com

    # Retrieve that password and print to stdout
    pw get google.com

    # Copy that password to the clipboard
    pw copy google.com

    # Generate a new password, save it as "reddit.com", and copy it to the clipboad
    pw generate reddit.com

    # List passwords
    pw ls

## Configuration

`pw` can be configured via command-line options, or by specifying the
same options in `~/.pw`, in the form `<option> = <value>`. Supported
configuration options are:

- `root` -- The root directory of the password store. Defaults to
  `~/pw`.
- `gpgkey` -- The GPG key ID to encrypt passwords to. Can be in any
  form understood by `gpg -r`, but a full key fingerprint is suggested
  for maximum explicitness.
- `cmd.copy` -- A command to use for copying text. This command will
  be invoked with the text to copy on STDIN. Defaults to `xclip -i` on
  Linux and `pbcopy` on OS X.
- `cmd.generate` -- A command that should return a freshly-generated
  password when invoked. Defaults to `pwgen 12 1 -s`.

## Tab-completion

To enable tab-completion for `pw`, simply add the following line to
your `.bashrc`:

    complete -C 'pw -do-completion' pw

Completion support is currently only implemented/tested for
bash. Patches or bug reports in other shells are welcome.
