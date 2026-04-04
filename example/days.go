package example

// Day is an enum type that should be implemented by the enumify generator.
// It uses a 1D array of strings for the names, which should be discovered by the
// enumify generator due to the go:generate directive above the Day type declaration.
//
//go:generate go run ../cmd/enumify
type Day uint8

// Constants for the Day enum values.
// These values should be discovered by the enumify generator since they use the same
// type as the Day enum, which is the type being generated.
const (
	Unknown Day = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

// 1D array of strings for the names of the Day enum values.
// This should be discovered by the enumify generator due to the go:generate directive
// and because it matches the dayNames pattern to connect it with the Day enum.
//
//lint:ignore U1000 this is used by the enumify generator
var dayNames = []string{
	"unknown",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
	"Sunday",
}
