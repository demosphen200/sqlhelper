package sqlhelper

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sqlhelper/pkg/smapper"
)

func (helper *SqlHelper) Delete(
	ctx context.Context,
	ptrToModel interface{},
) (sql.Result, error) {
	if !smapper.IsPtrToStruct(ptrToModel) {
		return nil, &MustBePtrToStructError{}
	}
	if helper.Db == nil {
		return nil, &DbNotSetError{}
	}
	if deleteSql, params, err := helper.getDeleteSql(ptrToModel); err != nil {
		return nil, err
	} else {
		result, err := helper.getDb(ctx).Exec(ctx, deleteSql, params...)
		//result, err := helper.Db.Exec(ctx, deleteSql, params...)
		return result, err
	}
}

func (helper *SqlHelper) DeleteBySql(
	ctx context.Context,
	deleteSql string,
	params ...interface{},
) (sql.Result, error) {
	if helper.AutoConvertParams {
		if err := helper.convertParams(&params); err != nil {
			return nil, err
		}
	}

	if helper.Db == nil {
		return nil, &DbNotSetError{}
	}
	return helper.getDb(ctx).Exec(ctx, deleteSql, params...)
}

func (helper *SqlHelper) getDeleteSql(ptrToModel interface{}) (string, []any, error) {
	if !smapper.IsPtrToStruct(ptrToModel) {
		return "", nil, &MustBePtrToStructError{}
	}
	if helper.TableName == "" {
		return "", nil, &TableNameNotSetError{}
	}

	cacheItem, err := helper.findOrCreateDeleteCacheItem(reflect.TypeOf(ptrToModel).Elem())
	if err != nil {
		return "", nil, err
	}

	params := make([]any, 0)
	structValue := reflect.ValueOf(ptrToModel).Elem()
	for _, fieldIndexAndConverter := range cacheItem.structFieldIndexesAndConverters {
		fieldValue := structValue.Field(fieldIndexAndConverter.Index).Interface()
		params = append(params, fieldValue)
	}
	return cacheItem.sql, params, nil
}

func (helper *SqlHelper) findOrCreateDeleteCacheItem(structType reflect.Type) (*sqlAndParamFieldIndexes, error) {
	if item, found := helper.deleteCache.Get(structType); found {
		return item, nil
	}

	whereFieldNames := make([]string, 0)
	whereParamIndexes := make([]indexAndConverter, 0)

	_ = smapper.EnumTaggedStructFields(structType, helper.DbFieldTagName,
		func(index int, structField *reflect.StructField, tagValue string) error {
			isIdField := structField.Tag.Get(helper.IsIdTagName) != ""
			if !isIdField {
				return nil
			}
			whereFieldNames = append(whereFieldNames, tagValue)
			whereParamIndexes = append(whereParamIndexes, indexAndConverter{
				Index:     index,
				Converter: nil,
			})
			return nil
		},
	)
	cacheItem := &sqlAndParamFieldIndexes{
		sql: fmt.Sprintf(
			"delete"+" from %s where %s",
			helper.TableName,
			helper.buildIdFilterString(whereFieldNames),
			//strings.Join(whereFieldNames, "=? and ")+"=?",
		),
		structFieldIndexesAndConverters: whereParamIndexes,
	}
	helper.deleteCache.Put(structType, cacheItem)
	return cacheItem, nil
}
