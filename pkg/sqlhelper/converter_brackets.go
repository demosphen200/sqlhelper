package sqlhelper

import (
	"fmt"
)

func NewAddBracketsConverter() *SimpleTypeConverter {
	return NewSimpleDbTypeConverter[string, string](
		"brackets",
		func(local string) (db string, err error) {
			return fmt.Sprintf("[%s]", local), nil
		},
		func(db string) (local string, err error) {
			return fmt.Sprintf("(%s)", db), nil
		},
	)
}
