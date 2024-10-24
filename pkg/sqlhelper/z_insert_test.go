package sqlhelper

import (
	"github.com/stretchr/testify/assert"
	"sqlhelper/pkg/testhelper"
	"testing"
)

func TestSqlHelper_getInsertSql(t *testing.T) {
	testhelper.HelperSetT(t)
	var model = testModel{
		Id:   1,
		Name: "username",
		Cnt:  11,
	}
	helper := testhelper.NoError(createAndOpenDb())
	insertSql, params := testhelper.NoError2(helper.getInsertSql(&model, false))
	assert.Equal(t, "insert into testTable (name,cnt) values (?,?)", insertSql)
	assert.Equal(t, []any{"username", 11}, params)
}

func TestSqlHelper_getInsertSqlWithModifier(t *testing.T) {
	testhelper.HelperSetT(t)
	var model = testModel{
		Id:   1,
		Name: "username",
		Cnt:  11,
	}
	helper := testhelper.NoError(createAndOpenDb())
	helper.InsertModifier = "or replace"
	insertSql, params := testhelper.NoError2(helper.getInsertSql(&model, false))
	assert.Equal(t, "insert or replace into testTable (name,cnt) values (?,?)", insertSql)
	assert.Equal(t, []any{"username", 11}, params)
}

func TestSqlHelper_Insert_ShouldInsert(t *testing.T) {
	testhelper.HelperSetT(t)
	var user = testModel{
		Id:   1,
		Name: "username",
		Cnt:  11,
	}
	helper := testhelper.NoError(createAndOpenDb())
	_ = testhelper.NoError(helper.Insert(ctx, &user, false))
	var users []testModel
	testhelper.NoError0(helper.SelectAll(ctx, &users))
	assert.Equal(t, 1, len(users))
	assert.Equal(t, user, users[0])
}

func TestSqlHelper_Insert_CanReturnConverterNotFoundError(t *testing.T) {
	testhelper.HelperSetT(t)
	var model = converterNotFoundModel{}
	helper := testhelper.NoError(createAndOpenDb())
	_, err := helper.Insert(ctx, &model, false)
	assert.ErrorContains(t, err, "not found")
	assert.ErrorContains(t, err, "not_existing")
}
