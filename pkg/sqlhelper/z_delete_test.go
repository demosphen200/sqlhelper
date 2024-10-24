package sqlhelper

import (
	"github.com/stretchr/testify/assert"
	"sqlhelper/pkg/testhelper"
	"testing"
)

func TestSqlHelper_getDeleteSql(t *testing.T) {
	testhelper.HelperSetT(t)
	var model = testModel{
		Id:   1,
		Name: "username",
		Cnt:  11,
	}
	/*
		var builder = SqlHelper{
			TableName: "testTable",
		}*/
	helper := testhelper.NoError(createAndOpenDb())

	updateSql, params := testhelper.NoError2(helper.getDeleteSql(&model))
	assert.Equal(t, "delete from testTable where id=?", updateSql)
	assert.Equal(t, []any{1}, params)
}

func TestSqlHelper_getDeleteSqlWithComplexId(t *testing.T) {
	testhelper.HelperSetT(t)
	var model = testModel2{
		Id:   1,
		Id2:  2,
		Name: "username",
		Cnt:  11,
	}
	/*
		var builder = SqlHelper{
			TableName: "testTable",
		}*/
	helper := testhelper.NoError(createAndOpenDb())
	updateSql, params := testhelper.NoError2(helper.getDeleteSql(&model))
	assert.Equal(t, "delete from testTable where id=? and id2=?", updateSql)
	assert.Equal(t, []any{1, 2}, params)
}

func TestSqlHelper_Delete_ShouldDeleteByIdAndIgnoreOther(t *testing.T) {
	testhelper.HelperSetT(t)

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 10))

	var model = testModel{}

	model.Id = 3
	model.Name = "not used"
	model.Cnt = 333

	var res = testhelper.NoError(helper.Delete(ctx, &model))
	var rowsAffected = testhelper.NoError(res.RowsAffected())
	assert.Equal(t, int64(1), rowsAffected)

	var models []testModel

	testhelper.NoError0(helper.Select(ctx, &models, "order by id", 4))
	assert.Equal(t, 9, len(models))

	for index, model := range models {
		if index == 2 {
			continue
		}
		if index < 2 {
			validateModel(t, &model, index+1, index+1)
		} else {
			validateModel(t, &model, index+2, index+2)
		}
	}
}
