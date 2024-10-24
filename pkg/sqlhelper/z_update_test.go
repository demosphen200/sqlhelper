package sqlhelper

import (
	"github.com/stretchr/testify/assert"
	"sqlhelper/pkg/testhelper"
	"testing"
)

func TestSqlHelper_getUpdateSql(t *testing.T) {
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
	updateSql, params := testhelper.NoError2(helper.getUpdateSql(&model))
	assert.Equal(t, "update testTable set name=?,cnt=? where id=?", updateSql)
	assert.Equal(t, []any{"username", 11, 1}, params)
}

func TestSqlHelper_getUpdateSqlWithComplexId(t *testing.T) {
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

	updateSql, params := testhelper.NoError2(helper.getUpdateSql(&model))
	assert.Equal(t, "update testTable set name=?,cnt=? where id=? and id2=?", updateSql)
	assert.Equal(t, []any{"username", 11, 1, 2}, params)
}

func TestSqlHelper_Update_ShouldUpdateByIdAndIgnoreOther(t *testing.T) {
	testhelper.HelperSetT(t)

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 10))

	var model = testModel{}

	model.Id = 3
	model.Name = "username333"
	model.Cnt = 333
	var res = testhelper.NoError(helper.Update(ctx, &model))
	var rowsAffected = testhelper.NoError(res.RowsAffected())
	assert.Equal(t, int64(1), rowsAffected)

	for n := 1; n < 10; n++ {
		if n != 3 {
			selectAndValidateModel(t, helper, n, n)
		} else {
			selectAndValidateModel(t, helper, 3, 333)
		}
	}
}
