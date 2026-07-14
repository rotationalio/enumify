package enumify

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
	"unicode"

	g "github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/require"
)

var (
	DefaultInvalid = []any{"foo", "123", "INVALID", 257, -1, 314.314, struct{}{}, true, false}
	TestingT       = g.Params(g.Id("t").Op("*").Qual("testing", "T"))
)

const (
	DBCharLimit = 32
)

type TestSuite[T ~uint8, Names []string | [][]string] struct {
	Values      []T       // the ordered enum values
	Names       Names     // the names of the enum values
	Invalid     []any     // any invalid values to test
	ICase       bool      // whether the enum is case insensitive (adds multi-case tests) defaults to true
	ISpace      bool      // whether the enum is space insensitive (adds space repr tests) defaults to true
	DBCharLimit int       // the maximum length of a string representation for the database (defaults to 32)
	parser      Parser[T] // parsing function created by ParseFactory
	unknowns    string    // the string representation of the unknown value
}

//============================================================================
// Testing Utilities
//============================================================================

func (s *TestSuite[T, Names]) Run(t *testing.T) {
	t.Run("Interface", s.TestInterface)
	t.Run("Stringer", s.TestStringer)
	t.Run("StringBounds", s.TestStringBounds)
	t.Run("Parse", s.TestParse)
	t.Run("Database", s.TestDatabase)
}

func (s *TestSuite[T, Names]) TestInterface(t *testing.T) {
	enum := T(0)
	require.Implements(t, (*Enum)(nil), &enum, "%T must implement the Enum interface", enum)
}

func (s *TestSuite[T, Names]) TestStringer(t *testing.T) {
	names := s.Strings()
	for i, enum := range s.Values {
		e, ok := any(enum).(fmt.Stringer)
		require.True(t, ok, "expected %T to be a fmt.Stringer", enum)
		require.Equal(t, names[i], e.String(), "expected %T to have string representation %q, got %q", e, names[i], e.String())
	}

	// Test Zero Values
	zero := T(0)
	e, ok := any(zero).(fmt.Stringer)
	require.True(t, ok, "expected %T to be a fmt.Stringer", zero)
	require.Equal(t, s.Unknowns(), e.String(), "expected %T to have string representation %q, got %q", e, s.Unknowns(), e.String())
}

func (s *TestSuite[T, Names]) TestStringBounds(t *testing.T) {
	max := uint8(0)   // max starts at 0 and anything greater in Values is set to max
	min := uint8(255) // min starts at 255 and anything less in Values is set to min

	// Discover the maximum and minimum values of the enum.
	for _, e := range s.Values {
		if uint8(e) > max {
			max = uint8(e)
		}
		if uint8(e) < min {
			min = uint8(e)
		}
	}

	// Create a value above the maximum value.
	above := T(max + 1)
	aboves, ok := any(above).(fmt.Stringer)
	require.True(t, ok, "expected %T to be a fmt.Stringer", above)
	require.Equal(t, s.Unknowns(), aboves.String(), "expected %T to have string representation %q for unknown value above maximum enum value %d", above, s.Unknowns(), max)

	// Test zero value
	if min > 0 {
		zero := T(0)
		zeros, ok := any(zero).(fmt.Stringer)
		require.True(t, ok, "expected %T to be a fmt.Stringer", zero)
		require.Equal(t, s.Unknowns(), zeros.String(), "expected %T to have string representation %q for unknown value at zero", zero, s.Unknowns())
	}
}

func (s *TestSuite[T, Names]) TestParse(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		testCases := s.ValidCases()
		for _, val := range testCases {
			actual, err := s.Parser()(val)
			require.NoError(t, err, "expected parsing valid %T value %q to not error", T(0), val)
			require.Contains(t, s.Values, actual, "expected parsing valid %T value %q to return valid enum", T(0), val)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		testCases := s.InvalidCases()
		for _, val := range testCases {
			actual, err := s.Parser()(val)
			require.Error(t, err, "expected parsing invalid %T value %q to error", T(0), val)
			require.Equal(t, T(0), actual, "expected parsing invalid %T value %q to return unknown value %T", T(0), val, actual)
		}
	})
}

func (s *TestSuite[T, Names]) TestDatabase(t *testing.T) {
	// TODO: implement scan and value tests
	t.Run("VARCHAR", func(t *testing.T) {
		// Ensure that all string representations are less than or equal to the db VARCHAR limit
		var limit int = s.DBCharLimit
		if limit == 0 {
			limit = DBCharLimit
		}

		for _, enum := range s.Values {
			s, ok := any(enum).(fmt.Stringer)
			require.True(t, ok, "expected %T to be a fmt.Stringer", enum)
			require.LessOrEqual(t, len(s.String()), limit, "expected %T value %q to be less than or equal to %d characters", enum, s.String(), limit)
		}
	})
}

//============================================================================
// Helper Functions
//============================================================================

func (s *TestSuite[T, Names]) Parser() Parser[T] {
	if s.parser == nil {
		s.parser = ParseFactory[T](s.Names)
	}
	return s.parser
}

func (s *TestSuite[T, Names]) Unknowns() string {
	if s.unknowns == "" {
		s.unknowns = s.Strings()[0]
	}
	return s.unknowns
}

func (s *TestSuite[T, Names]) Strings() []string {
	var names []string
	switch col := any(s.Names).(type) {
	case []string:
		names = col
	case [][]string:
		names = col[0]
	}
	return names
}

func (s *TestSuite[T, Names]) ValidCases() []any {
	cases := make([]any, 0, len(s.Values)*18)
	names := s.Strings()

	// Add the numeric representations
	// Set 1: the enum values themselves
	// Set 2-6: uint, uint8, uint16, uint32, uint64
	// Set 7-11: int, int8, int16, int32, int64
	// Set 12-13: float32, float64
	for _, val := range s.Values {
		cases = append(cases, val)
		cases = append(cases, uint(val), uint8(val), uint16(val), uint32(val), uint64(val))
		cases = append(cases, int(val), int8(val), int16(val), int32(val), int64(val))
		cases = append(cases, float32(val), float64(val))
	}

	// Add the string representations
	// Set 14: the enum values as strings
	// Set 15: lowercase if case-insensitive
	// Set 16: mixed case if case-insensitive
	// Set 17: uppercase if case-insensitive
	// Set 18: spaces added if space-insensitive
	for _, name := range names {
		cases = append(cases, name)
		if s.ICase {
			cases = append(cases, strings.ToLower(name))
			cases = append(cases, mixedCase(name))
			cases = append(cases, strings.ToUpper(name))
		}
		if s.ISpace {
			cases = append(cases, addSpaces(name))
		}
	}

	return cases
}

func (s *TestSuite[T, Names]) InvalidCases() []any {
	if len(s.Invalid) > 0 {
		return s.Invalid
	}

	cases := make([]any, 0, len(DefaultInvalid)+len(s.Values)+16)
	cases = append(cases, DefaultInvalid...)
	cases = append(cases, math.MaxInt, math.MinInt, math.MaxInt8, math.MinInt8, math.MaxInt16, math.MinInt16, math.MaxInt32, math.MinInt32, math.MaxInt64, math.MinInt64)
	cases = append(cases, math.MaxUint16, math.MaxUint32)
	cases = append(cases, float32(-1.0), float64(-1.0), math.MaxFloat32, math.MaxFloat64)

	for _, name := range s.Strings() {
		cases = append(cases, mangle(name))
		if !s.ICase {
			cases = append(cases, mixedCase(name))
		}
		if !s.ISpace {
			cases = append(cases, addSpaces(name))
		}
	}
	return cases
}

func mixedCase(s string) string {
	sb := strings.Builder{}
	for _, r := range s {
		// Flip a coin and make the character upper or lower case
		if rand.Intn(2) == 0 {
			sb.WriteRune(unicode.ToLower(r))
		} else {
			sb.WriteRune(unicode.ToUpper(r))
		}
	}
	return sb.String()
}

func addSpaces(s string) string {
	fn := strings.Repeat(" ", rand.Intn(8))
	bn := strings.Repeat(" ", rand.Intn(4))
	return fn + s + bn
}

func mangle(s string) string {
	sb := strings.Builder{}
	mangled := false
	for _, c := range s {
		// 40% chance to mangle the character
		if rand.Intn(100) < 40 {
			r := rune(rand.Intn(93) + 33) // 33-126
			sb.WriteRune(r)
			mangled = r != c && unicode.ToLower(r) != unicode.ToLower(c)
		} else {
			sb.WriteRune(c)
		}
	}

	if !mangled {
		for i := 0; i < rand.Intn(4)+1; i++ {
			sb.WriteRune(rune(rand.Intn(93) + 33)) // 33-126
		}
	}
	return sb.String()
}
