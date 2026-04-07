package enumify_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/enumify"
)

func TestLowerFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},
		{input: "foo", expected: "foo"},
		{input: "Foo", expected: "foo"},
		{input: "FOO", expected: "fOO"},
		{input: "FooBar", expected: "fooBar"},
		{input: "ΩTime", expected: "ωTime"},
		{input: "ωTime", expected: "ωTime"},
	}

	for _, test := range tests {
		actual := enumify.LowerFirst(test.input)
		require.Equal(t, test.expected, actual)
	}
}

func TestUpperFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},
		{input: "foo", expected: "Foo"},
		{input: "Foo", expected: "Foo"},
		{input: "FOO", expected: "FOO"},
		{input: "FooBar", expected: "FooBar"},
		{input: "fooBar", expected: "FooBar"},
		{input: "ΩTime", expected: "ΩTime"},
		{input: "ωTime", expected: "ΩTime"},
	}

	for _, test := range tests {
		actual := enumify.UpperFirst(test.input)
		require.Equal(t, test.expected, actual)
	}
}
