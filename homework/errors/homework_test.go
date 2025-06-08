package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	if len(e.errors) == 0 {
		return ""
	}

	errMsg := strings.Builder{}

	errMsg.WriteString(fmt.Sprintf("%d errors occured:\n", len(e.errors)))

	for _, err := range e.errors {
		errMsg.WriteString("\t* ")
		errMsg.WriteString(err.Error())
	}

	errMsg.WriteString("\n")

	return errMsg.String()
}

func Append(err error, errs ...error) *MultiError {
	var multiErr *MultiError

	if err == nil {
		multiErr = &MultiError{errors: make([]error, 0)}
	} else {
		if me, ok := err.(*MultiError); ok {
			multiErr = me
		} else {
			multiErr = &MultiError{errors: []error{err}}
		}
	}

	for _, e := range errs {
		if e != nil {
			multiErr.errors = append(multiErr.errors, e)
		}
	}

	return multiErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
