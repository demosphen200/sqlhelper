package sqlhelper

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"sqlhelper/pkg/ex"
	"sqlhelper/pkg/smapper"
	"strings"
)

func (helper *SqlHelper) Select(
	ctx context.Context,
	ptrToResult interface{},
	sqlEnding string,
	params ...any,
) error {
	if helper.AutoConvertParams {
		if err := helper.convertParams(&params); err != nil {
			return err
		}
	}

	if smapper.IsPtrToStruct(ptrToResult) {
		return helper.selectSingleRow(ctx, ptrToResult, sqlEnding, params...)
	} else if smapper.IsPtrToSliceOfStruct(ptrToResult) {
		return helper.selectManyRows(ctx, ptrToResult, sqlEnding, params...)
	} else {
		return &MustBePtrToStructOrSliceError{}
	}
}

func (helper *SqlHelper) selectSingleRow(
	ctx context.Context,
	ptrToModel interface{},
	sqlEnding string,
	params ...interface{},
) error {
	if !smapper.IsPtrToStruct(ptrToModel) {
		return &MustBePtrToStructError{}
	}
	modelType := reflect.TypeOf(ptrToModel).Elem()
	if selectSql, err := helper.buildSelectSql(modelType, sqlEnding); err != nil {
		return err
	} else {
		return helper.selectSingleRowInternal(
			ctx,
			ptrToModel,
			selectSql.sql,
			selectSql.mapperId,
			nil,
			selectSql.resultTypes,
			nil,
			selectSql.resultConverters,
			nil,
			params...,
		)
	}
}

func (helper *SqlHelper) selectSingleRowInternal(
	ctx context.Context,
	ptrToModel interface{},
	sql string,
	mapperId string,
	paramConverters []DbTypeConverter,
	resultTypes []reflect.Type,
	resultColumnIndexes []int,
	resultConverters []DbTypeConverter,
	onLateResultParams func(rows DbRows) ([]reflect.Type, []int, []DbTypeConverter, error),
	params ...interface{},
) error {
	if helper.Db == nil {
		return &DbNotSetError{}
	}
	if rows, err := helper.getDb(ctx).Query(ctx, sql, params...); err != nil {
		return err
	} else {
		defer func() { _ = rows.Close() }()

		if onLateResultParams != nil {
			if types, indexes, converters, err := onLateResultParams(rows); err != nil {
				return fmt.Errorf("lateinit error %w", err)
			} else {
				resultTypes = types
				resultColumnIndexes = indexes
				resultConverters = converters
			}
		}

		if !rows.Next() {
			return &NoRowsReturnedError{}
		}
		if err = helper.scan(rows, resultTypes, resultColumnIndexes, resultConverters, ptrToModel, mapperId); err != nil {
			return err
		}
		if rows.Next() {
			return &MoreThanOneRowReturnedError{}
		}
		return nil
	}

}

func (helper *SqlHelper) selectManyRows(
	ctx context.Context,
	ptrToSlice interface{},
	sqlEnding string,
	params ...interface{},
) error {
	if !smapper.IsPtrToSliceOfStruct(ptrToSlice) {
		return &MustBePtrToSliceError{}
	}
	sliceType := reflect.TypeOf(ptrToSlice).Elem()
	modelType := sliceType.Elem()
	if selectSql, err := helper.buildSelectSql(modelType, sqlEnding); err != nil {
		return err
	} else {
		return helper.selectManyRowsInternal(
			ctx,
			ptrToSlice,
			selectSql.sql,
			selectSql.mapperId,
			nil,
			selectSql.resultTypes,
			nil,
			selectSql.resultConverters,
			nil,
			params...,
		)
	}
}

func (helper *SqlHelper) selectManyRowsInternal(
	ctx context.Context,
	ptrToSliceOfModel interface{},
	sql string,
	mapperId string,
	paramConverters []DbTypeConverter,
	resultTypes []reflect.Type,
	resultColumnIndexes []int,
	resultConverters []DbTypeConverter,
	onLateResultParams func(rows DbRows) ([]reflect.Type, []int, []DbTypeConverter, error),
	params ...any,
) error {
	if helper.Db == nil {
		return &DbNotSetError{}
	}
	sliceType := reflect.TypeOf(ptrToSliceOfModel).Elem()
	modelType := sliceType.Elem()
	if rows, err := helper.getDb(ctx).Query(ctx, sql, params...); err != nil {
		return err
	} else {
		defer func() { _ = rows.Close() }()

		if onLateResultParams != nil {
			if types, indexes, converters, err := onLateResultParams(rows); err != nil {
				return fmt.Errorf("lateinit error %w", err)
			} else {
				resultTypes = types
				resultColumnIndexes = indexes
				resultConverters = converters
			}
		}

		newSlice := reflect.MakeSlice(sliceType, 0, 0)
		for rows.Next() {
			dataPtr := reflect.New(modelType)
			if err := helper.scan(rows, resultTypes, resultColumnIndexes, resultConverters, dataPtr.Interface(), mapperId); err != nil {
				return err
			}
			newSlice = reflect.Append(newSlice, dataPtr.Elem())
		}
		reflect.ValueOf(ptrToSliceOfModel).Elem().Set(newSlice)
		return nil
	}
}

func (helper *SqlHelper) SelectBySql(
	ctx context.Context,
	result interface{},
	sql string,
	params ...interface{},
) error {
	if helper.AutoConvertParams {
		if err := helper.convertParams(&params); err != nil {
			return err
		}
	}

	var structType reflect.Type

	if smapper.IsPtrToStruct(result) {
		structType = reflect.TypeOf(result).Elem()
	} else if smapper.IsPtrToSliceOfStruct(result) {
		structType = reflect.TypeOf(result).Elem().Elem()
	} else {
		return &MustBePtrToStructOrSliceError{}
	}

	cacheItem, err := helper.findOrCreateSelectBySqlCacheItem(structType, sql)
	if err != nil {
		return err
	}

	var lateInit func(rows DbRows) ([]reflect.Type, []int, []DbTypeConverter, error)
	if cacheItem.resultColumnsIndexes == nil || cacheItem.resultConverters == nil {
		lateInit = func(rows DbRows) ([]reflect.Type, []int, []DbTypeConverter, error) {
			columns, err := rows.Columns()
			if err != nil {
				return nil, nil, nil, err
			}
			err = helper.fillSelectBySqlCacheItem(cacheItem, structType, columns)
			if err != nil {
				return nil, nil, nil, err
			}
			return cacheItem.resultTypes, cacheItem.resultColumnsIndexes, cacheItem.resultConverters, nil
		}
	}

	if smapper.IsPtrToStruct(result) {
		return helper.selectSingleRowInternal(
			ctx,
			result,
			sql,
			cacheItem.mapperId,
			nil,
			cacheItem.resultTypes,
			cacheItem.resultColumnsIndexes,
			cacheItem.resultConverters,
			lateInit,
			params...,
		)
	} else { //if smapper.IsPtrToSliceOfStruct(result) {
		return helper.selectManyRowsInternal(
			ctx,
			result,
			sql,
			cacheItem.mapperId,
			nil,
			cacheItem.resultTypes,
			cacheItem.resultColumnsIndexes,
			cacheItem.resultConverters,
			lateInit,
			params...,
		)
	}
}

func (helper *SqlHelper) fillSelectBySqlCacheItem(
	item *selectBySqlCacheItem,
	structType reflect.Type,
	columns []string,
) error {
	resultTypes := make([]reflect.Type, 0)
	resultColumnsIndexes := make([]int, 0)
	converters := make([]DbTypeConverter, 0)

	err := smapper.EnumTaggedStructFields(structType, helper.DbFieldTagName,
		func(index int, structField *reflect.StructField, tagValue string) error {
			columnIndex := slices.Index(columns, tagValue)
			if columnIndex == -1 {
				//resultColumnsIndexes = append(resultColumnsIndexes, -1)
				//converters = append(converters, nil)
				return nil
			}

			converterName := structField.Tag.Get(helper.ConverterTagName)
			converter, err := helper.findTypeConverter(converterName)
			if err != nil {
				return err
			}

			resultTypes = append(resultTypes, structField.Type)
			resultColumnsIndexes = append(resultColumnsIndexes, columnIndex)
			converters = append(converters, converter)
			return nil
		},
	)
	if err != nil {
		return err
	}
	item.resultTypes = resultTypes
	item.resultColumnsIndexes = resultColumnsIndexes
	item.resultConverters = converters
	return nil
}

func (helper *SqlHelper) findOrCreateSelectBySqlCacheItem(structType reflect.Type, sql string) (*selectBySqlCacheItem, error) {
	if cacheItem, found := helper.selectBySqlCache.Get(structType, sql); found {
		return cacheItem, nil
	}
	idFieldNames := make([]string, 0)
	dbFieldNames := make([]string, 0)

	converters := make([]DbTypeConverter, 0)
	idConverters := make([]DbTypeConverter, 0)

	err := smapper.EnumTaggedStructFields(structType, helper.DbFieldTagName,
		func(index int, structField *reflect.StructField, tagValue string) error {

			converterName := structField.Tag.Get(helper.ConverterTagName)
			converter, err := helper.findTypeConverter(converterName)
			if err != nil {
				return err
			}
			converters = append(converters, converter)

			isIdField := structField.Tag.Get(helper.IsIdTagName) != ""
			if isIdField {
				idFieldNames = append(idFieldNames, tagValue)
				idConverters = append(idConverters, converter)
			}
			dbFieldNames = append(dbFieldNames, tagValue)

			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot build select cache item: %w", err)
	}
	cacheItem := &selectBySqlCacheItem{
		mapperId:             helper.getNextMapperId(),
		resultTypes:          nil,
		resultColumnsIndexes: nil,
		resultConverters:     nil,
	}
	helper.selectBySqlCache.Put(structType, sql, cacheItem)
	return cacheItem, nil
}

func (helper *SqlHelper) SelectAll(ctx context.Context, ptrToSlice interface{}) error {
	return helper.Select(ctx, ptrToSlice, "")
}

func (helper *SqlHelper) getSelectSql(
	ptrToModel interface{},
	sqlEnding string,
) (*selectByIdSqlAndIdFieldIndexes, error) {
	//) (string, error) {
	if !smapper.IsPtrToStruct(ptrToModel) {
		return nil, &MustBePtrToStructError{}
	}
	if helper.TableName == "" {
		return nil, &DbNotSetError{}
	}
	return helper.buildSelectSql(reflect.TypeOf(ptrToModel).Elem(), sqlEnding)
}

func (helper *SqlHelper) buildSelectSql(
	structType reflect.Type,
	sqlEnding string,
) (
	*selectByIdSqlAndIdFieldIndexes,
	error,
) {
	cacheItem, err := helper.findOrCreateSelectCacheItem(structType)
	if err != nil {
		return nil, err
	}
	if sqlEnding != "" {
		tempCacheItem := *cacheItem
		tempCacheItem.sql = fmt.Sprintf("%s %s", cacheItem.sql, sqlEnding)
		return &tempCacheItem, nil
	} else {
		return cacheItem, nil
	}
}

func (helper *SqlHelper) findOrCreateSelectCacheItem(structType reflect.Type) (*selectByIdSqlAndIdFieldIndexes, error) {
	if cacheItem, found := helper.selectCache.Get(structType); found {
		return cacheItem, nil
	}
	idFieldNames := make([]string, 0)
	dbFieldNames := make([]string, 0)

	resultTypes := make([]reflect.Type, 0)
	converters := make([]DbTypeConverter, 0)
	idConverters := make([]DbTypeConverter, 0)

	err := smapper.EnumTaggedStructFields(structType, helper.DbFieldTagName,
		func(index int, structField *reflect.StructField, tagValue string) error {

			converterName := structField.Tag.Get(helper.ConverterTagName)
			converter, err := helper.findTypeConverter(converterName)
			if err != nil {
				return err
			}
			converters = append(converters, converter)

			resultTypes = append(resultTypes, structField.Type)

			isIdField := structField.Tag.Get(helper.IsIdTagName) != ""
			if isIdField {
				idFieldNames = append(idFieldNames, tagValue)
				idConverters = append(idConverters, converter)
			}
			dbFieldNames = append(dbFieldNames, tagValue)

			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot build select cache item: %w", err)
	}
	cacheItem := &selectByIdSqlAndIdFieldIndexes{
		sql: fmt.Sprintf(
			"select"+" %s from %s",
			strings.Join(dbFieldNames, ","),
			helper.TableName,
		),
		resultTypes:      resultTypes,
		resultConverters: converters,
		idConverters:     idConverters,
		whereIdEnding: fmt.Sprintf(
			"where %s",
			helper.buildIdFilterString(idFieldNames),
		),
		mapperId: helper.getNextMapperId(),
	}
	helper.selectCache.Put(structType, cacheItem)
	return cacheItem, nil
}

func (helper *SqlHelper) SelectSingleValue(
	ctx context.Context,
	ptrToResult interface{},
	selectSql string,
	params ...interface{},
) error {
	return helper.SelectSingleValueC(ctx, ptrToResult, selectSql, "", params...)
}

func (helper *SqlHelper) SelectSingleValueC(
	ctx context.Context,
	ptrToResult interface{},
	selectSql string,
	converterName string,
	params ...interface{},
) error {
	if helper.AutoConvertParams {
		if err := helper.convertParams(&params); err != nil {
			return err
		}
	}

	var converter DbTypeConverter
	var err error
	if converterName != "" {
		converter, err = helper.findTypeConverter(converterName)
		if err != nil {
			return err
		}
	}

	if helper.Db == nil {
		return &DbNotSetError{}
	}
	if rows, err := helper.getDb(ctx).Query(ctx, selectSql, params...); err != nil {
		return err
	} else {
		defer ex.CloseSilent(rows)
		if !rows.Next() {
			return &NoRowsReturnedError{}
		}

		if converter != nil {
			var placeholder any
			if err = rows.Scan(&placeholder); err != nil {
				return err
			}
			if err = converter.DbToLocal(placeholder, ptrToResult); err != nil {
				return err
			}
		} else {
			if err = rows.Scan(ptrToResult); err != nil {
				return err
			}
		}
		if rows.Next() {
			return &MoreThanOneRowReturnedError{}
		}
		return nil
	}
}

func (helper *SqlHelper) SelectById(
	ctx context.Context,
	ptrToModel interface{},
	id ...interface{},
) error {
	if !smapper.IsPtrToStruct(ptrToModel) {
		return &MustBePtrToStructError{}
	}
	if cacheItem, err := helper.findOrCreateSelectCacheItem(reflect.TypeOf(ptrToModel).Elem()); err != nil {
		return err
	} else {
		if len(id) != cacheItem.idFieldsCount() {
			return &IdFieldAndArgCountNotMatchError{}
		}
		return helper.selectSingleRow(ctx, ptrToModel, cacheItem.whereIdEnding, id...)
	}
}
