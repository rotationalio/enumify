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

// This enum should be discovered by the enumify generator due to the go:generate
// directive on line 7 of this file and because it matches the Enum spec pattern.
type Month uint8

const (
	Monthless Month = iota
	January
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

// This 2D array of strings should match the names pattern for the color enum without
// having to specify it using the -names flag.
//
//lint:ignore U1000 this is used by the enumify generator
var monthNames = [][]string{
	{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
	{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
}

// This additional enum method should be ignored by the enumify generator.
func (m Month) Abbreviation() string {
	if m >= Month(len(monthNames[1])) {
		return monthNames[1][Monthless]
	}
	return monthNames[1][m]
}
