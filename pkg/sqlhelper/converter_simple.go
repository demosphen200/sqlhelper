package sqlhelper

import (
	"fmt"
)

type SimpleTypeConverter struct {
	name             string
	localToDb        func(src any) (dst any, err error)
	dbToLocal        func(src any, ptrToDst any) error
	isSupportedLocal func(src any) bool
}

func (c *SimpleTypeConverter) Name() string {
	return c.name
}

func (c *SimpleTypeConverter) DbToLocal(src any, ptrToDst any) error {
	return c.dbToLocal(src, ptrToDst)
}

func (c *SimpleTypeConverter) LocalToDb(src any) (dst any, err error) {
	return c.localToDb(src)
}

func (c *SimpleTypeConverter) IsSupportedLocal(src any) bool {
	return c.isSupportedLocal(src)
}

func NewSimpleDbTypeConverter[LOCAL any, DB any](
	name string,
	localToDb func(local LOCAL) (db DB, err error),
	dbToLocal func(db DB) (LOCAL, error),
) *SimpleTypeConverter {
	return &SimpleTypeConverter{
		name: name,
		localToDb: func(src any) (dst any, err error) {
			switch s := src.(type) {
			case LOCAL:
				return localToDb(s)
			case *LOCAL:
				return localToDb(*s)
			default:
				return nil, fmt.Errorf(`converter "%s" (local to ob) cannot convert source type: %T`, name, src)
			}
		},
		dbToLocal: func(src any, ptrToDst any) error {
			var result LOCAL
			var err error
			switch s := src.(type) {
			case DB:
				result, err = dbToLocal(s)
				if err != nil {
					return err
				}
			case *DB:
				result, err = dbToLocal(*s)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf(`converter "%s" (db to local) cannot convert source type: %T`, name, src)
			}

			switch d := ptrToDst.(type) {
			case *LOCAL:
				*d = result
			case **LOCAL:
				*d = &result
			case *any:
				switch d2 := (*d).(type) {
				case LOCAL:
					*d = result
				case *LOCAL:
					*d2 = result
				case **LOCAL:
					*d2 = &result
				default:
					*d = result
				}
			default:
				return fmt.Errorf(`converter "%s" (db to local) cannot convert to destination type: %T`, name, ptrToDst)
			}
			return nil
		},
		isSupportedLocal: func(local any) bool {
			switch local.(type) {
			case LOCAL, *LOCAL:
				return true
			default:
				return false
			}
		},
	}
}
