package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// (20) Failing test: describeUser panics on nil — no guard in bad/main.go.
func TestDescribeUser_NilPanics(t *testing.T) {
	assert.Empty(t, describeUser(nil)) // panics before assert is reached
}

// (21) Failing test: process(mode=1) — wrong output from the deeply nested logic.
func TestProcessMode1_ExpectedFragments(t *testing.T) {
	items := []string{"a", "bb", "ccc", "dddd"}
	assert.Equal(t, "?x?x", process(items, 1))
}

// (22) Failing test: compareErr uses == instead of errors.Is, so wrapped
// errors don't match.
func TestCompareErr_WrappedEOF(t *testing.T) {
	assert.True(t, compareErr(wrap(io.EOF)),
		"expected wrapped io.EOF to be matched (it isn't — bad/ uses ==)")
}

type wrappedErr struct{ inner error }

func (w wrappedErr) Error() string { return "wrapped: " + w.inner.Error() }
func (w wrappedErr) Unwrap() error { return w.inner }

func wrap(err error) error { return wrappedErr{inner: err} }
