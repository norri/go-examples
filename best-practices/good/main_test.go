package main

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescribeUser(t *testing.T) {
	t.Parallel()

	assert.Empty(t, describeUser(nil))
	assert.Equal(t, "ada", describeUser(&user{Name: "ada"}))
}

func TestCompareErr(t *testing.T) {
	t.Parallel()

	assert.True(t, compareErr(io.EOF), "plain io.EOF should match")
	assert.True(t, compareErr(wrap(io.EOF)), "wrapped io.EOF should match via errors.Is")
}

func wrap(err error) error { return errors.Join(errors.New("ctx"), err) }

func TestProcessModes(t *testing.T) {
	t.Parallel()

	items := []string{"a", "bb", "ccc", "dddd"}
	tests := []struct {
		name string
		mode int
		want string
	}{
		{name: "mode1", mode: 1, want: "?x?x"},
		{name: "mode2", mode: 2, want: "abbcccdddd"},
		{name: "default", mode: 99, want: "----"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, process(items, tc.mode))
		})
	}
}
