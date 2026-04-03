package enumify_test

//============================================================================
// Test Enum Type
//============================================================================

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
