package enumify

import (
	"fmt"
	"strings"
)

// Generic parser function that can be used to parse a specific enum type.
type Parser[T ~uint8] func(any) (T, error)

// ParseFactory creates a parser function for a specific enum type. The parser can parse
// the enum from a string or a numeric type that can be converted into a uint8. In order
// to support string parsing, a names array is required. It can either be a single
// array where the index of the name is the value of the enum, or a 2D array where the
// first column contains the names indexed according to the enum value.
//
// NOTE: string parsing is case-insensitive and leading and trailing whitespace is
// ignored. Names should be slug values to ensure consistent parsing.
func ParseFactory[T ~uint8, Names []string | [][]string](names Names) Parser[T] {
	var normalizedNames []string
	switch col := any(names).(type) {
	case []string:
		if len(col) < 1 {
			panic(fmt.Errorf("names array must contain at least one name"))
		}

		normalizedNames = make([]string, len(col))
		for i, name := range col {
			normalizedNames[i] = normalize(name)
		}
	case [][]string:
		if len(col) < 1 {
			panic(fmt.Errorf("names array must contain at least one column"))
		}

		if len(col[0]) < 1 {
			panic(fmt.Errorf("names array must contain at least one name"))
		}

		normalizedNames = make([]string, len(col[0]))
		for i, name := range col[0] {
			normalizedNames[i] = normalize(name)
		}
	}

	// The "unknown" value is the zero-valued T.
	unknown := T(0)

	return func(val any) (T, error) {
		switch v := val.(type) {
		case string:
			v = normalize(v)

			// For an empty string, return the "unknown" or zero-valued T.
			if v == "" {
				return unknown, nil
			}

			// Iterate over the normalized names and return the value of the first match.
			for i, name := range normalizedNames {
				if name == v {
					return T(i), nil
				}
			}

			// If no match is found, return an error.
			return unknown, fmt.Errorf("invalid %T value: %q", unknown, v)
		case T:
			if v >= T(len(normalizedNames)) {
				return unknown, fmt.Errorf("invalid %T value: %d", unknown, v)
			}
			return v, nil
		case uint:
			if v >= uint(len(normalizedNames)) {
				return unknown, fmt.Errorf("invalid %T value: %d", unknown, v)
			}
			return T(v), nil
		case uint8:
			if v >= uint8(len(normalizedNames)) {
				return unknown, fmt.Errorf("invalid %T value: %d", unknown, v)
			}
			return T(v), nil
		case uint16:
			if v >= uint16(len(normalizedNames)) {
				return unknown, fmt.Errorf("invalid %T value: %d", unknown, v)
			}
			return T(v), nil
		case uint32:
			if v >= uint32(len(normalizedNames)) {
				return unknown, fmt.Errorf("invalid %T value: %d", unknown, v)
			}
			return T(v), nil
		case uint64:
			if v >= uint64(len(normalizedNames)) {
				return unknown, fmt.Errorf("invalid %T value: %d", unknown, v)
			}
			return T(v), nil
		case int:
			if v < 0 || v >= len(normalizedNames) {
				return unknown, fmt.Errorf("invalid %T value: %d", unknown, v)
			}
			return T(v), nil
		case int8:
			if v < 0 || v >= int8(len(normalizedNames)) {
				return unknown, fmt.Errorf("invalid %T value: %d", unknown, v)
			}
			return T(v), nil
		case int16:
			if v < 0 || v >= int16(len(normalizedNames)) {
				return unknown, fmt.Errorf("invalid %T value: %d", unknown, v)
			}
			return T(v), nil
		case int32:
			if v < 0 || v >= int32(len(normalizedNames)) {
				return unknown, fmt.Errorf("invalid %T value: %d", unknown, v)
			}
			return T(v), nil
		case int64:
			if v < 0 || v >= int64(len(normalizedNames)) {
				return unknown, fmt.Errorf("invalid %T value: %d", unknown, v)
			}
			return T(v), nil
		case float32:
			if v < 0 || v >= float32(len(normalizedNames)) {
				return unknown, fmt.Errorf("invalid %T value: %f", unknown, v)
			}
			return T(v), nil
		case float64:
			if v < 0 || v >= float64(len(normalizedNames)) {
				return unknown, fmt.Errorf("invalid %T value: %f", unknown, v)
			}
			return T(v), nil
		default:
			return unknown, fmt.Errorf("cannot parse %T into %T", v, unknown)
		}
	}
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
