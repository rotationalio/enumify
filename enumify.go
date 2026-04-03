package enumify

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Enum interface {
	fmt.Stringer
	json.Marshaler
	json.Unmarshaler
	sql.Scanner
	driver.Valuer
}
