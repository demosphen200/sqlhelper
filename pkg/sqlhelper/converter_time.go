package sqlhelper

import (
	"fmt"
	"time"
)

func NewTimeConverter(
	name string,
	timeLayout string,
) *SimpleTypeConverter {
	if timeLayout == "" {
		timeLayout = time.DateTime
	}
	return NewSimpleDbTypeConverter[time.Time, string](
		name,
		func(local time.Time) (db string, err error) {
			return local.Format(timeLayout), nil
		},
		func(db string) (local time.Time, err error) {
			var tm time.Time
			tm, err = time.Parse(timeLayout, db)
			if err != nil {
				return tm, fmt.Errorf(
					`cannot parse time: string="%s" layout="%s": %w`,
					local,
					timeLayout,
					err,
				)
			}
			return tm, nil
		},
	)
}

func RegisterTimeConverters(converters *DbTypeConverters) {
	converters.Register(NewTimeConverter("datetime", time.DateTime))
	converters.Register(NewTimeConverter("date", time.DateOnly))
	converters.Register(NewTimeConverter("time", time.TimeOnly))
}
