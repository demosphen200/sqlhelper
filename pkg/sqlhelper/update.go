package sqlhelper

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sqlhelper/pkg/smapper"
)

func (helper *SqlHelper) Update(
	ctx context.Context,
	ptrToModel interface{},
) (sql.Result, error) {
	if !smapper.IsPtrToStruct(ptrToModel) {
		return nil, &MustBePtrToStructError{}
	}
	if helper.Db == nil {
		return nil, &DbNotSetError{}
	}
	if updateSql, params, err := helper.getUpdateSql(ptrToModel); err != nil {
		return nil, err
	} else {
		result, err := helper.getDb(ctx).Exec(ctx, updateSql, params...)
		return result, err
	}
}

func (helper *SqlHelper) UpdateNoResult(
	ctx context.Context,
	ptrToModel interface{},
) error {
	_, err := helper.Update(ctx, ptrToModel)
	return err
}

func (helper *SqlHelper) UpdateBySql(
	ctx context.Context,
	updateSql string,
	params ...interface{},
) (sql.Result, error) {
	if helper.AutoConvertParams {
		if err := helper.convertParams(&params); err != nil {
			return nil, err
		}
	}

	return helper.getDb(ctx).Exec(ctx, updateSql, params...)
}

func (helper *SqlHelper) getUpdateSql(ptrToModel interface{}) (string, []any, error) {
	if !smapper.IsPtrToStruct(ptrToModel) {
		return "", nil, &MustBePtrToStructError{}
	}
	if helper.TableName == "" {
		return "", nil, &TableNameNotSetError{}
	}

	cacheItem, err := helper.findOrCreateUpdateCacheItem(reflect.TypeOf(ptrToModel).Elem())
	if err != nil {
		return "", nil, err
	}

	params := make([]any, 0)
	structValue := reflect.ValueOf(ptrToModel).Elem()
	for _, fieldIndexAndConverter := range cacheItem.structFieldIndexesAndConverters {
		fieldIndex := fieldIndexAndConverter.Index
		field := structValue.Field(fieldIndex)
		param := field.Interface()
		if fieldIndexAndConverter.Converter != nil {
			param, err = fieldIndexAndConverter.Converter.LocalToDb(param)
			if err != nil {
				return "", nil, fmt.Errorf("cannot convert param value for field index %d: %w", fieldIndex, err)
			}
		}
		params = append(params, param)
	}
	return cacheItem.sql, params, nil
}

func (helper *SqlHelper) findOrCreateUpdateCacheItem(structType reflect.Type) (*sqlAndParamFieldIndexes, error) {
	if item, found := helper.updateCache.Get(structType); found {
		return item, nil
	}

	updateFieldNames := make([]string, 0)
	updateFieldIndexesAndConverters := make([]indexAndConverter, 0)

	whereFieldNames := make([]string, 0)
	whereFieldIndexesAndConverters := make([]indexAndConverter, 0)

	err := smapper.EnumTaggedStructFields(structType, helper.DbFieldTagName,
		func(index int, structField *reflect.StructField, tagValue string) error {
			isIdField := structField.Tag.Get(helper.IsIdTagName) != ""

			converter, err := helper.findTypeConverter(structField.Tag.Get(helper.ConverterTagName))
			if err != nil {
				return err
			}

			fieldIndexAndConverter := indexAndConverter{
				Index:     index,
				Converter: converter,
			}
			if isIdField {
				whereFieldNames = append(whereFieldNames, tagValue)
				whereFieldIndexesAndConverters = append(whereFieldIndexesAndConverters, fieldIndexAndConverter)
			} else {
				updateFieldNames = append(updateFieldNames, tagValue)
				updateFieldIndexesAndConverters = append(updateFieldIndexesAndConverters, fieldIndexAndConverter)
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	updateStr := ""
	whereStr := ""
	paramIndex := 0

	for _, fieldName := range updateFieldNames {
		updateStr = joinNotEmpty(updateStr, fieldName+"="+helper.Db.ParamPlaceholder(paramIndex), ",")
		paramIndex++
	}
	for _, fieldName := range whereFieldNames {
		whereStr = joinNotEmpty(whereStr, fieldName+"="+helper.Db.ParamPlaceholder(paramIndex), " and ")
		paramIndex++
	}

	cacheItem := &sqlAndParamFieldIndexes{
		sql: fmt.Sprintf(
			"update"+" %s set %s where %s",
			helper.TableName,
			updateStr,
			whereStr,
		),
		structFieldIndexesAndConverters: append(updateFieldIndexesAndConverters, whereFieldIndexesAndConverters...),
	}

	helper.updateCache.Put(structType, cacheItem)

	return cacheItem, nil
}
