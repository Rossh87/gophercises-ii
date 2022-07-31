package cmd

import (
	"os"

	"golang.org/x/term"
)

type PasswordReader interface {
	ReadPassword() (string, error)
}

type StdinPasswordReader struct {
}

func (pr StdinPasswordReader) ReadPassword() (string, error) {
	pw, err := term.ReadPassword(int(os.Stdin.Fd()))

	if err != nil {
		return "", err
	}

	return string(pw), err
}

func Run(pr PasswordReader) (string, error) {
	pwd, err := pr.ReadPassword()

	if err != nil {
		return "", err
	}

	return string(pwd), nil
}
