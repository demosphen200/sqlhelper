package sqlhelper

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sqlhelper/pkg/KvCache"
	"sqlhelper/pkg/smapper"
	"strconv"
	"sync/atomic"
)

type SqlHelper struct {
	DbFieldTagName         string
	IsIdTagName            string
	ConverterTagName       string
	TableName              string
	Db                     DbAdapter
	InsertModifier         string
	selectCache            KvCache.Cache[reflect.Type, *selectByIdSqlAndIdFieldIndexes]
	selectBySqlCache       KvCache.Cache2[reflect.Type, string, *selectBySqlCacheItem]
	insertCache            KvCache.Cache2[reflect.Type, bool, *sqlAndParamFieldIndexes]
	updateCache            KvCache.Cache[reflect.Type, *sqlAndParamFieldIndexes]
	deleteCache            KvCache.Cache[reflect.Type, *sqlAndParamFieldIndexes]
	mapper                 smapper.Mapper
	Converters             *DbTypeConverters
	selectBySqlCacheNextId atomic.Int64
	AutoConvertParams      bool
}

func NewSqlHelper(
	db DbAdapter,
	tableName string,
) *SqlHelper {
	helper := &SqlHelper{
		TableName:         tableName,
		DbFieldTagName:    "db",
		IsIdTagName:       "id",
		ConverterTagName:  "converter",
		Db:                db,
		Converters:        DefaultDbTypeConverters,
		AutoConvertParams: true,

		selectCache:      KvCache.MakeCache[reflect.Type, *selectByIdSqlAndIdFieldIndexes](),
		selectBySqlCache: KvCache.MakeCache2[reflect.Type, string, *selectBySqlCacheItem](),
		insertCache:      KvCache.MakeCache2[reflect.Type, bool, *sqlAndParamFieldIndexes](),
		updateCache:      KvCache.MakeCache[reflect.Type, *sqlAndParamFieldIndexes](),
		deleteCache:      KvCache.MakeCache[reflect.Type, *sqlAndParamFieldIndexes](),

		mapper: smapper.MakeMapper(),

		selectBySqlCacheNextId: atomic.Int64{},
	}
	helper.mapper.TagKey = helper.DbFieldTagName
	return helper
}

type indexAndConverter struct {
	Index     int
	Converter DbTypeConverter
}

type sqlAndParamFieldIndexes struct {
	sql                             string
	structFieldIndexesAndConverters []indexAndConverter
}

type selectBySqlCacheItem struct {
	resultTypes          []reflect.Type
	resultColumnsIndexes []int
	resultConverters     []DbTypeConverter
	mapperId             string
}

type selectByIdSqlAndIdFieldIndexes struct {
	sql              string
	resultTypes      []reflect.Type
	resultConverters []DbTypeConverter
	idConverters     []DbTypeConverter
	whereIdEnding    string
	mapperId         string
}

func (s *selectByIdSqlAndIdFieldIndexes) selectedFieldsCount() int {
	return len(s.resultConverters)
}

func (s *selectByIdSqlAndIdFieldIndexes) idFieldsCount() int {
	return len(s.idConverters)
}

func (helper *SqlHelper) buildIdFilterString(fields []string) string {
	str := ""
	for index, field := range fields {
		if str != "" {
			str += " and "
		}
		str += field + "=" + helper.Db.ParamPlaceholder(index)
	}
	return str
}

func (helper *SqlHelper) findTypeConverter(name string) (DbTypeConverter, error) {
	if name == "" {
		return nil, nil
	}
	return helper.Converters.Find(name)
}

func (helper *SqlHelper) Count(
	ctx context.Context,
	sqlEnding string,
	params ...interface{},
) (int, error) {
	var count int
	if err := helper.SelectSingleValue(
		ctx,
		&count,
		fmt.Sprintf("select count(*) from"+" %s %s", helper.TableName, sqlEnding),
		params...,
	); err != nil {
		return -1, err
	} else {
		return count, nil
	}
}

func joinNotEmpty(value1, value2, separator string) string {
	if value1 != "" {
		return value1 + separator + value2
	}
	return value2
}

func (helper *SqlHelper) scan(
	rows DbRows,
	resultTypes []reflect.Type,
	resultColumnIndexes []int,
	converters []DbTypeConverter,
	structPtr any,
	mapperId string,
) error {
	if mapperId == "" {
		return fmt.Errorf("abnormal: cannot scan: mapperId is empty")
	}

	columnNames, err := rows.Columns()

	if len(columnNames) != len(converters) {
		if resultColumnIndexes == nil {
			return fmt.Errorf(
				"scan: columns count (%d) and converter count (%d) mismatch",
				len(columnNames),
				len(converters),
			)
		}
	}

	if err != nil {
		return err
	}
	values := make([]any, len(columnNames))
	pointers := make([]any, len(columnNames))
	pointersForConverters := make([]any, 0)

	for t := 0; t < len(columnNames); t++ {
		pointers[t] = &values[t]
	}

	if resultTypes != nil {
		if resultColumnIndexes != nil {
			for index, columnIndex := range resultColumnIndexes {
				value := reflect.New(resultTypes[index])
				if converters[index] == nil {
					pointers[columnIndex] = value.Interface()
				} else {
					pointersForConverters = append(pointersForConverters, value.Interface())
				}
			}
		} else {
			for index, resultType := range resultTypes {
				value := reflect.New(resultType)
				if converters[index] == nil {
					pointers[index] = value.Interface()
				} else {
					pointersForConverters = append(pointersForConverters, value.Interface())
				}
			}

		}
	}

	if err = rows.Scan(pointers...); err != nil {
		return err
	}

	for index, pointer := range pointers {
		values[index] = reflect.ValueOf(pointer).Elem().Interface()
	}

	if resultColumnIndexes != nil {
		newValues := make([]any, len(resultColumnIndexes))
		for index, columnIndex := range resultColumnIndexes {
			newValues[index] = values[columnIndex]
		}
		values = newValues
	}

	for index, converter := range converters {
		if converter != nil {
			pointer := pointersForConverters[0]
			pointersForConverters = pointersForConverters[1:]

			err = converter.DbToLocal(values[index], pointer)
			if err != nil {
				return err
			}
			values[index] = pointer
		}
	}

	return helper.mapper.MapFromSlice(mapperId, values, structPtr)
}

func (helper *SqlHelper) getDb(ctx context.Context) DbAdapter {
	adapter := ctx.Value(contextTxAdapterKey)
	if adapter != nil {
		return adapter.(DbAdapter)
	}
	return helper.Db
}

func (helper *SqlHelper) RunInTransaction(
	ctx context.Context,
	options sql.TxOptions,
	block func(tx context.Context) error,
) error {
	return helper.getDb(ctx).RunInTransaction(ctx, options, func(adapter DbAdapter) error {
		//txHelper := NewSqlHelper(adapter, helper.TableName)
		cancelableCtx, cancel := context.WithCancel(ctx)
		txCtx := context.WithValue(cancelableCtx, contextTxAdapterKey, adapter)
		err := block(txCtx)
		cancel()
		return err
	})
}

func (helper *SqlHelper) RequireTransaction(
	ctx context.Context,
	options sql.TxOptions,
	block func(tx context.Context) error,
) error {
	if ctx.Value(contextTxAdapterKey) != nil {
		return block(ctx)
	} else {
		return helper.RunInTransaction(ctx, options, block)
	}
}

func (helper *SqlHelper) getNextMapperId() string {
	return strconv.FormatInt(helper.selectBySqlCacheNextId.Add(1), 16)
}

func (helper *SqlHelper) convertParams(params *[]any) error {
	for index, param := range *params {
		converter := helper.Converters.TryFindByLocalType(param)
		if converter != nil {
			newValue, err := converter.LocalToDb(param)
			if err != nil {
				return err
			}
			(*params)[index] = newValue
		}
	}
	return nil
}
