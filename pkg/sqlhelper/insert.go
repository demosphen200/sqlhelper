package sqlhelper

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sqlhelper/pkg/smapper"
	"strings"
)

func (helper *SqlHelper) Insert(
	ctx context.Context,
	ptrToModel interface{},
	withId bool,
) (sql.Result, error) {
	if !smapper.IsPtrToStruct(ptrToModel) {
		return nil, &MustBePtrToStructError{}
	}
	if helper.Db == nil {
		return nil, &DbNotSetError{}
	}
	if insertSql, params, err := helper.getInsertSql(ptrToModel, withId); err != nil {
		return nil, err
	} else {
		result, err := helper.getDb(ctx).Exec(ctx, insertSql, params...)
		return result, err
	}
}

func (helper *SqlHelper) getInsertSql(ptrToModel interface{}, withId bool) (string, []any, error) {
	if !smapper.IsPtrToStruct(ptrToModel) {
		return "", nil, &MustBePtrToStructError{}
	}
	if helper.TableName == "" {
		return "", nil, &TableNameNotSetError{}
	}

	cacheItem, err := helper.findOrCreateInsertCacheItem(reflect.TypeOf(ptrToModel).Elem(), withId)
	if err != nil {
		return "", nil, err
	}

	params := make([]any, 0)
	structValue := reflect.ValueOf(ptrToModel).Elem()
	for _, fieldIndexAndConverter := range cacheItem.structFieldIndexesAndConverters {
		index := fieldIndexAndConverter.Index
		field := structValue.Field(index)
		param := field.Interface()
		if fieldIndexAndConverter.Converter != nil {
			param, err = fieldIndexAndConverter.Converter.LocalToDb(param)
			if err != nil {
				return "", nil, fmt.Errorf("cannot convert param value for field index %d: %w", index, err)
			}
		}
		params = append(params, param)
	}

	return cacheItem.sql, params, nil
}

func (helper *SqlHelper) findOrCreateInsertCacheItem(structType reflect.Type, withId bool) (*sqlAndParamFieldIndexes, error) {
	if item, found := helper.insertCache.Get(structType, withId); found {
		return item, nil
	}

	cacheItem := &sqlAndParamFieldIndexes{
		structFieldIndexesAndConverters: make([]indexAndConverter, 0),
		//structFieldIndexesAndConverters: make([]int, 0),
	}

	fieldNames := make([]string, 0)
	questions := make([]string, 0)

	err := smapper.EnumTaggedStructFields(structType, helper.DbFieldTagName,
		func(index int, structField *reflect.StructField, tagValue string) error {
			if !withId && structField.Tag.Get(helper.IsIdTagName) != "" {
				return nil
			}
			fieldNames = append(fieldNames, tagValue)
			converterName := structField.Tag.Get(helper.ConverterTagName)
			converter, err := helper.findTypeConverter(converterName)
			if err != nil {
				return fmt.Errorf("cannot build insert cache item (field index %d): %w", index, err)
			}
			fieldIndexAndConverter := indexAndConverter{
				Index:     index,
				Converter: converter,
			}
			cacheItem.structFieldIndexesAndConverters = append(cacheItem.structFieldIndexesAndConverters, fieldIndexAndConverter)

			questions = append(questions, helper.Db.ParamPlaceholder(len(questions)))
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	modifier := helper.InsertModifier
	if modifier != "" {
		modifier = fmt.Sprintf("%s ", modifier)
	}
	cacheItem.sql = fmt.Sprintf(
		"insert %sinto"+" %s (%s) values (%s)",
		modifier,
		helper.TableName,
		strings.Join(fieldNames, ","),
		strings.Join(questions, ","),
	)
	helper.insertCache.Put(structType, withId, cacheItem)
	return cacheItem, nil
}
