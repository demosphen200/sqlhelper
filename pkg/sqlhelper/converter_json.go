package sqlhelper

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type JsonConverter struct {
	DbTypeConverter
}

func NewJsonConverter() *JsonConverter {
	return &JsonConverter{}
}

func (c *JsonConverter) Name() string {
	return "json"
}

func (c *JsonConverter) LocalToDb(src any) (dst any, err error) {
	str, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	return str, nil
}

func (c *JsonConverter) IsSupportedLocal(local any) bool {
	localValue := reflect.ValueOf(local)
	kind := localValue.Kind()
	switch kind {
	case reflect.Pointer:
		return isJsonableKind(localValue.Elem().Kind())
	default:
		return isJsonableKind(kind)
	}
}

func isJsonableKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice:
		return true
	default:
		return false
	}
}

func (c *JsonConverter) DbToLocal(src any, dst any) error {
	var sourceBytes []byte
	switch s := src.(type) {
	case string:
		sourceBytes = []byte(s)
	case *string:
		sourceBytes = []byte(*s)
	case []byte:
		sourceBytes = s
	case *[]byte:
		sourceBytes = *s
	default:
		return fmt.Errorf("abnormal: invalid source type (%T) in DbToLocal", src)
	}

	rfDst := reflect.ValueOf(dst)

	if rfDst.Kind() == reflect.Ptr && isJsonableKind(rfDst.Elem().Kind()) {
		err := json.Unmarshal(sourceBytes, dst)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("not supported dst type: %s", rfDst.Type().String())
}

func RegisterJsonConverter(converters *DbTypeConverters) {
	converters.Register(NewJsonConverter())
}
