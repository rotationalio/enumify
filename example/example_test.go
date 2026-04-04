package example_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/enumify/example"
)

// This is an unrelated test that should be ignored by the enumify generator.
func TestExample(t *testing.T) {
	example, err := example.New()
	require.NoError(t, err)
	require.NotNil(t, example)
	require.NotZero(t, example.Name)
	require.False(t, example.Date.IsZero())
}
