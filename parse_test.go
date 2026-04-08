package enumify_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/enumify"
)

func TestParseFactory(t *testing.T) {
	validTestCases := []struct {
		value    any
		expected Status
	}{
		{value: "", expected: StatusUnknown},
		{value: "unknown", expected: StatusUnknown},
		{value: "draft", expected: StatusDraft},
		{value: "review", expected: StatusReview},
		{value: "published", expected: StatusPublished},
		{value: "archived", expected: StatusArchived},
		{value: "    ", expected: StatusUnknown},
		{value: "\n", expected: StatusUnknown},
		{value: "\t", expected: StatusUnknown},
		{value: "  unknown ", expected: StatusUnknown},
		{value: "draft  ", expected: StatusDraft},
		{value: "  review", expected: StatusReview},
		{value: "\tpublished\t", expected: StatusPublished},
		{value: "\narchived\t", expected: StatusArchived},
		{value: "UNKNOWN", expected: StatusUnknown},
		{value: "DRAFT", expected: StatusDraft},
		{value: "REVIEW", expected: StatusReview},
		{value: "PUBLISHED", expected: StatusPublished},
		{value: "ARCHIVED", expected: StatusArchived},
		{value: byte(0), expected: StatusUnknown},
		{value: byte(1), expected: StatusDraft},
		{value: byte(2), expected: StatusReview},
		{value: byte(3), expected: StatusPublished},
		{value: byte(4), expected: StatusArchived},
		{value: []byte{0}, expected: StatusUnknown},
		{value: []byte{1}, expected: StatusDraft},
		{value: []byte{2}, expected: StatusReview},
		{value: []byte{3}, expected: StatusPublished},
		{value: []byte{4}, expected: StatusArchived},
		{value: []byte(""), expected: StatusUnknown},
		{value: []byte("unknown"), expected: StatusUnknown},
		{value: []byte("draft"), expected: StatusDraft},
		{value: []byte("review"), expected: StatusReview},
		{value: []byte("published"), expected: StatusPublished},
		{value: []byte("archived"), expected: StatusArchived},
		{value: uint(0), expected: StatusUnknown},
		{value: uint(1), expected: StatusDraft},
		{value: uint(2), expected: StatusReview},
		{value: uint(3), expected: StatusPublished},
		{value: uint(4), expected: StatusArchived},
		{value: uint8(0), expected: StatusUnknown},
		{value: uint8(1), expected: StatusDraft},
		{value: uint8(2), expected: StatusReview},
		{value: uint8(3), expected: StatusPublished},
		{value: uint8(4), expected: StatusArchived},
		{value: uint16(0), expected: StatusUnknown},
		{value: uint16(1), expected: StatusDraft},
		{value: uint16(2), expected: StatusReview},
		{value: uint16(3), expected: StatusPublished},
		{value: uint16(4), expected: StatusArchived},
		{value: uint32(0), expected: StatusUnknown},
		{value: uint32(1), expected: StatusDraft},
		{value: uint32(2), expected: StatusReview},
		{value: uint32(3), expected: StatusPublished},
		{value: uint32(4), expected: StatusArchived},
		{value: uint64(0), expected: StatusUnknown},
		{value: uint64(1), expected: StatusDraft},
		{value: uint64(2), expected: StatusReview},
		{value: uint64(3), expected: StatusPublished},
		{value: uint64(4), expected: StatusArchived},
		{value: int(0), expected: StatusUnknown},
		{value: int(1), expected: StatusDraft},
		{value: int(2), expected: StatusReview},
		{value: int(3), expected: StatusPublished},
		{value: int(4), expected: StatusArchived},
		{value: int8(0), expected: StatusUnknown},
		{value: int8(1), expected: StatusDraft},
		{value: int8(2), expected: StatusReview},
		{value: int8(3), expected: StatusPublished},
		{value: int8(4), expected: StatusArchived},
		{value: int16(0), expected: StatusUnknown},
		{value: int16(1), expected: StatusDraft},
		{value: int16(2), expected: StatusReview},
		{value: int16(3), expected: StatusPublished},
		{value: int16(4), expected: StatusArchived},
		{value: int32(0), expected: StatusUnknown},
		{value: int32(1), expected: StatusDraft},
		{value: int32(2), expected: StatusReview},
		{value: int32(3), expected: StatusPublished},
		{value: int32(4), expected: StatusArchived},
		{value: int64(0), expected: StatusUnknown},
		{value: int64(1), expected: StatusDraft},
		{value: int64(2), expected: StatusReview},
		{value: int64(3), expected: StatusPublished},
		{value: int64(4), expected: StatusArchived},
		{value: float32(0), expected: StatusUnknown},
		{value: float32(1), expected: StatusDraft},
		{value: float32(2), expected: StatusReview},
		{value: float32(3), expected: StatusPublished},
		{value: float32(4), expected: StatusArchived},
		{value: float64(0), expected: StatusUnknown},
		{value: float64(1), expected: StatusDraft},
		{value: float64(2), expected: StatusReview},
		{value: float64(3), expected: StatusPublished},
		{value: float64(4), expected: StatusArchived},
		{value: StatusUnknown, expected: StatusUnknown},
		{value: StatusDraft, expected: StatusDraft},
		{value: StatusReview, expected: StatusReview},
		{value: StatusPublished, expected: StatusPublished},
		{value: StatusArchived, expected: StatusArchived},
		{value: Status(0), expected: StatusUnknown},
		{value: Status(1), expected: StatusDraft},
		{value: Status(2), expected: StatusReview},
		{value: Status(3), expected: StatusPublished},
		{value: Status(4), expected: StatusArchived},
	}

	invalidTestCases := []struct {
		value    any
		expected string
	}{
		{value: "foo", expected: `invalid enumify_test.Status value: "foo"`},
		{value: "Unbekannt", expected: `invalid enumify_test.Status value: "unbekannt"`},
		{value: "Entwurf", expected: `invalid enumify_test.Status value: "entwurf"`},
		{value: "Überprüfung", expected: `invalid enumify_test.Status value: "überprüfung"`},
		{value: "Veröffentlicht", expected: `invalid enumify_test.Status value: "veröffentlicht"`},
		{value: "Archiviert", expected: `invalid enumify_test.Status value: "archiviert"`},
		{value: "Needs Review", expected: `invalid enumify_test.Status value: "needs review"`},
		{value: nil, expected: `cannot parse <nil> into enumify_test.Status`},
		{value: true, expected: `cannot parse bool into enumify_test.Status`},
		{value: Status(42), expected: `invalid enumify_test.Status value: 42`},
		{value: uint(42), expected: `invalid enumify_test.Status value: 42`},
		{value: uint8(42), expected: `invalid enumify_test.Status value: 42`},
		{value: uint16(42), expected: `invalid enumify_test.Status value: 42`},
		{value: uint32(42), expected: `invalid enumify_test.Status value: 42`},
		{value: uint64(42), expected: `invalid enumify_test.Status value: 42`},
		{value: int(-42), expected: `invalid enumify_test.Status value: -42`},
		{value: int(42), expected: `invalid enumify_test.Status value: 42`},
		{value: int8(-42), expected: `invalid enumify_test.Status value: -42`},
		{value: int8(42), expected: `invalid enumify_test.Status value: 42`},
		{value: int16(-42), expected: `invalid enumify_test.Status value: -42`},
		{value: int16(42), expected: `invalid enumify_test.Status value: 42`},
		{value: int32(-42), expected: `invalid enumify_test.Status value: -42`},
		{value: int32(42), expected: `invalid enumify_test.Status value: 42`},
		{value: int64(-42), expected: `invalid enumify_test.Status value: -42`},
		{value: int64(42), expected: `invalid enumify_test.Status value: 42`},
		{value: float32(-42), expected: `invalid enumify_test.Status value: -42.000000`},
		{value: float32(42), expected: `invalid enumify_test.Status value: 42.000000`},
		{value: float64(-42), expected: `invalid enumify_test.Status value: -42.000000`},
		{value: float64(42), expected: `invalid enumify_test.Status value: 42.000000`},
	}

	t.Run("Names", func(t *testing.T) {
		parse := enumify.ParseFactory[Status](StatusNames)
		for i, tc := range validTestCases {
			t.Run(fmt.Sprintf("Valid/%02d", i+1), func(t *testing.T) {
				actual, err := parse(tc.value)
				require.NoError(t, err)
				require.Equal(t, tc.expected, actual)
			})
		}

		for i, tc := range invalidTestCases {
			t.Run(fmt.Sprintf("Invalid/%02d", i+1), func(t *testing.T) {
				actual, err := parse(tc.value)
				require.EqualError(t, err, tc.expected)
				require.Equal(t, StatusUnknown, actual)
			})
		}
	})

	t.Run("Names2D", func(t *testing.T) {
		parse := enumify.ParseFactory[Status](StatusNames2D)
		for i, tc := range validTestCases {
			t.Run(fmt.Sprintf("Valid/%02d", i+1), func(t *testing.T) {
				actual, err := parse(tc.value)
				require.NoError(t, err)
				require.Equal(t, tc.expected, actual)
			})
		}

		for i, tc := range invalidTestCases {
			t.Run(fmt.Sprintf("Invalid/%02d", i+1), func(t *testing.T) {
				actual, err := parse(tc.value)
				require.EqualError(t, err, tc.expected)
				require.Equal(t, StatusUnknown, actual)
			})
		}
	})

	t.Run("InvalidNames", func(t *testing.T) {
		require.Panics(t, func() {
			enumify.ParseFactory[Status]([]string{})
		})
		require.Panics(t, func() {
			enumify.ParseFactory[Status]([][]string{})
		})
		require.Panics(t, func() {
			enumify.ParseFactory[Status]([][]string{{}, {}})
		})
	})
}
