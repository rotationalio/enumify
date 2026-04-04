package example

// Status is an enum type that should be implemented by the enumify generator.
// It uses a 2D array of strings for the names, which should be discovered by the
// enumify generator due to the go:generate directive above the Status type declaration.
//
//go:generate go run ../cmd/enumify -names statusTable
type Status uint8

// Constants for the Status enum values.
// These values should be discovered by the enumify generator since they use the same
// type as the Status enum, which is the type being generated.
const (
	StatusUnknown Status = iota
	StatusPending
	StatusRunning
	StatusFailed
	StatusSuccess
	StatusCancelled
)

// 2D array of strings for the names of the Status enum values.
// This should be discovered by the enumify generator due to the go:generate directive
// that specifies this variable as the names for the Status enum.
//
//lint:ignore U1000 this is used by the enumify generator
var statusTable = [][]string{
	{"unknown", "pending", "running", "failed", "success", "cancelled"},
	{"text-secondary", "text-info", "text-primary", "text-danger", "text-success", "text-warning"},
}
