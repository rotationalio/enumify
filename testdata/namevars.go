package testdata

var namesArray = [...]string{
	"foo",
	"bar",
	"baz",
}

var namesSlice = []string{
	"foo",
	"bar",
	"baz",
}

var namesTable = [][]string{
	{"foo", "bar", "baz"},
	{"foo", "bar", "baz"},
}

var names2DArray = [2][3]string{
	{"foo", "bar", "baz"},
	{"foo", "bar", "baz"},
}

var notNamesArray = [...]int{
	1,
	2,
	3,
}

var notNamesSlice = []int{
	1,
	2,
	3,
}

var notNamesTable = [][]int{
	{1, 2, 3},
	{1, 2, 3},
}

var notNames2DArray = [2][3]int{
	{1, 2, 3},
	{1, 2, 3},
}

const (
	Foo int32 = iota
	Bar
	Baz
)

var (
	debug    bool   = false
	progName string = "testing"
)
