package enumify_test

import (
	"fmt"
	"go/types"
	"testing"

	"github.com/stretchr/testify/require"
	. "go.rtnl.ai/enumify"
	"golang.org/x/tools/go/packages"
)

type NameVarTypeCheck func(typ types.Type) bool

func TestNameVarTypeChecks(t *testing.T) {
	testNameVarTypeCheck := func(nvtc NameVarTypeCheck, valid string) func(t *testing.T) {
		return func(t *testing.T) {
			testCases, err := ParseTypes("testdata/namevars.go")
			require.NoError(t, err, "could not parse types from testdata/namevars.go")
			require.GreaterOrEqual(t, len(testCases), 1, "expected at least 1 test case")

			validated := false
			for _, obj := range testCases {
				if obj.Name() == valid {
					require.True(t, IsNamesType(obj.Type()), "expected %q to be a names type", obj.Name())
					require.True(t, nvtc(obj.Type()), "expected %q to pass type check", obj.Name())
					validated = true
				} else {
					require.False(t, nvtc(obj.Type()), "expected %q to not be a valid type", obj.Name())
				}
			}
			require.True(t, validated, "no valid example for type check found")
		}
	}

	t.Run("StringArray", testNameVarTypeCheck(IsStringArray, "namesArray"))
	t.Run("StringSlice", testNameVarTypeCheck(IsStringSlice, "namesSlice"))
	t.Run("StringTable", testNameVarTypeCheck(IsStringTable, "namesTable"))
	t.Run("Strings2DArray", testNameVarTypeCheck(IsStrings2DArray, "names2DArray"))
}

func ParseTypes(file string) (out []types.Object, err error) {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo,
	}

	var pkgs []*packages.Package
	if pkgs, err = packages.Load(cfg, file); err != nil {
		return nil, fmt.Errorf("could not load package from %q: %w", file, err)
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("expected only 1 package returned from packages.Load, got %d", len(pkgs))
	}

	if len(pkgs[0].Errors) > 0 {
		return nil, fmt.Errorf("package errors: %v", pkgs[0].Errors)
	}

	names := pkgs[0].Types.Scope().Names()
	out = make([]types.Object, 0, len(names))
	for _, name := range names {
		obj := pkgs[0].Types.Scope().Lookup(name)
		out = append(out, obj)
	}
	return out, nil
}
