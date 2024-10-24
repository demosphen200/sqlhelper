package sqlhelper

import (
	"slices"
)

var DefaultDbTypeConverters = NewDbTypeConverters()

type DbTypeConverter interface {
	Name() string
	LocalToDb(src any) (dst any, err error)
	DbToLocal(src any, ptrToDst any) error
	IsSupportedLocal(src any) bool
}

type DbTypeConverters struct {
	items []DbTypeConverter
}

func NewDbTypeConverters() *DbTypeConverters {
	return &DbTypeConverters{
		items: make([]DbTypeConverter, 0),
	}
}

func (cs *DbTypeConverters) Find(name string) (DbTypeConverter, error) {
	for _, item := range cs.items {
		if item.Name() == name {
			return item, nil
		}
	}
	return nil, &TypeConverterNotFoundError{Name: name}
}

func (cs *DbTypeConverters) TryFindByLocalType(local any) DbTypeConverter {
	for _, item := range cs.items {
		if item.IsSupportedLocal(local) {
			return item
		}
	}
	return nil
}

func (cs *DbTypeConverters) Register(
	converter DbTypeConverter,
) {
	cs.Unregister(converter.Name())
	cs.items = append(cs.items, converter)
}

func (cs *DbTypeConverters) Unregister(name string) {
	cs.items = slices.DeleteFunc(cs.items, func(converter DbTypeConverter) bool {
		return converter.Name() == name
	})
}
