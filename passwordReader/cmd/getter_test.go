package cmd

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("Received %v of type %v, but expected %v of type %v", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
	fmt.Println("test success")
}

type mockPWReader struct {
	Password    string
	ReturnError bool
}

func (mpw mockPWReader) ReadPassword() (string, error) {
	if mpw.ReturnError {
		return "", errors.New("stubbed error")
	}

	if mpw.Password == "" {
		return "", errors.New("no value was provided for password")
	}

	return mpw.Password, nil
}

func TestRunReturnsErrorWhenPasswordFails(t *testing.T) {
	_, err := Run(mockPWReader{ReturnError: true})

	err, ok := err.(error)

	if !ok {
		t.Errorf("expected to receive an error, but instead received value of type %v", reflect.TypeOf(err))
	}
}

func TestRunReturnsPasswordString(t *testing.T) {
	pw, err := Run(mockPWReader{Password: "cats"})

	if err != nil {
		t.Error("call to run failed without yielding a password string")
	}

	assertEqual(t, pw, "cats")
}
