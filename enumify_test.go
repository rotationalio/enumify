package enumify_test

import (
	"database/sql/driver"
	"encoding/json"

	"go.rtnl.ai/enumify"
)

//============================================================================
// Enum type for testing the library
//============================================================================

// Status is an enum type that is used for library testing. This code is not generated
// but must match what the generated code would produce.
type Status uint8

const (
	StatusUnknown Status = iota
	StatusDraft
	StatusReview
	StatusPublished
	StatusArchived
)

var StatusNames = []string{
	"unknown",
	"draft",
	"review",
	"published",
	"archived",
}

var StatusNames2D = [][]string{
	{"unknown", "draft", "review", "published", "archived"},
	{"Unknown", "Draft", "Needs Review", "Published", "Archived"},
	{"Unbekannt", "Entwurf", "Überprüfung", "Veröffentlicht", "Archiviert"},
}

func (s Status) String() string {
	if s >= Status(len(StatusNames)) {
		return StatusNames[StatusUnknown]
	}
	return StatusNames[s]
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Status) UnmarshalJSON(data []byte) (err error) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s, err = enumify.ParseFactory[Status](StatusNames)(v)
	return err
}

func (s Status) MarshalYAML() (any, error) {
	return s.String(), nil
}

func (s *Status) UnmarshalYAML(unmarshal func(any) error) (err error) {
	var v string
	if err = unmarshal(&v); err != nil {
		return err
	}
	*s, err = enumify.ParseFactory[Status](StatusNames)(v)
	return err
}

func (s *Status) Scan(src any) (err error) {
	switch v := src.(type) {
	case nil:
		*s = StatusUnknown
		return nil
	case []byte:
		*s, err = enumify.ParseFactory[Status](StatusNames)(string(v))
		return err
	default:
		*s, err = enumify.ParseFactory[Status](StatusNames)(v)
		return err
	}
}

func (s Status) Value() (driver.Value, error) {
	return s.String(), nil
}
