package sqlhelper

import (
	"github.com/stretchr/testify/assert"
	"sqlhelper/pkg/testhelper"
	"testing"
)

func TestSqlHelper_getSelectSql(t *testing.T) {
	testhelper.HelperSetT(t)
	var model testModel
	//var helper = &SqlHelper{
	//	TableName: "testTable",
	//}
	//helper.Init()
	helper := testhelper.NoError(createAndOpenDb())
	selectSql := testhelper.NoError(helper.getSelectSql(&model, ""))
	assert.Equal(t, "select id,name,cnt from testTable", selectSql.sql)
	assert.Equal(t, 1, selectSql.idFieldsCount())
	assert.Equal(t, 3, selectSql.selectedFieldsCount())
}

func TestSqlHelper_selectManyRows_ShouldReturnEmptySliceOnNoRows(t *testing.T) {
	testhelper.HelperSetT(t)
	var users []testModel

	helper := testhelper.NoError(createAndOpenDb())

	testhelper.NoError0(helper.selectManyRows(ctx, &users, ""))
	assert.NotNil(t, users)
	assert.Equal(t, 0, len(users))
}

func TestSqlHelper_selectManyRows_ShouldGetSlice(t *testing.T) {
	testhelper.HelperSetT(t)
	var users []testModel

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 3))

	testhelper.NoError0(helper.selectManyRows(ctx, &users, ""))
	assert.Equal(t, 3, len(users))
}

func TestSqlHelper_selectManyRows_CanReturnMustBePtrToSliceError(t *testing.T) {
	helper := SqlHelper{}
	var i int
	err := helper.selectManyRows(ctx, &i, "")
	assert.Equal(t, true, assert.Error(t, err))
	assert.ErrorIs(t, &MustBePtrToSliceError{}, err)
}

func TestSqlHelper_selectSingleRow_ShouldGetModel(t *testing.T) {
	testhelper.HelperSetT(t)
	var user = testModel{}
	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 3))
	testhelper.NoError0(helper.selectSingleRow(ctx, &user, "where id=1", 1))
}

func TestSqlHelper_Select_ShouldSelectRowBySql(t *testing.T) {
	testhelper.HelperSetT(t)

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 10))

	for n := 1; n < 10; n++ {
		var model testModel
		testhelper.NoError0(helper.SelectBySql(ctx, &model, "select id,name,cnt from testTable where id=?", n))
		validateModel(t, &model, n, n)
	}
}

func TestSqlHelper_Select_ShouldSelectDifferentModelsBySql(t *testing.T) {
	testhelper.HelperSetT(t)

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 10))

	for n := 1; n < 10; n++ {
		var model testModel
		testhelper.NoError0(helper.SelectBySql(ctx, &model, "select id,name,cnt from testTable where id=?", n))
		validateModel(t, &model, n, n)

		var model2 testModel2
		testhelper.NoError0(helper.SelectBySql(ctx, &model2, "select id,id2,name,cnt from testTable where id=?", n))
		validateModel2(t, &model2, n, n)
	}
}

func TestSqlHelper_Select_ShouldSelectManyRowsBySql(t *testing.T) {
	testhelper.HelperSetT(t)

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 10))

	var models []testModel
	testhelper.NoError0(helper.SelectBySql(ctx, &models, "select id,name,cnt from testTable where id < ? order by id", 4))

	assert.Equal(t, 3, len(models), 3)
	for index, model := range models {
		validateModel(t, &model, index+1, index+1)
	}
}

func TestSqlHelper_Select_ShouldSelectBySqlWithoutTableName(t *testing.T) {
	testhelper.HelperSetT(t)

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 10))

	var helper2 = NewSqlHelper(helper.Db, "")

	var models []testModel
	testhelper.NoError0(helper2.SelectBySql(ctx, &models, "select id,name,cnt from testTable where id < ? order by id", 4))

	assert.Equal(t, 3, len(models), 3)
	for index, model := range models {
		validateModel(t, &model, index+1, index+1)
	}
}

func TestSqlHelper_Select_CanReturnConverterNotFoundError(t *testing.T) {
	testhelper.HelperSetT(t)
	var model = converterNotFoundModel{}
	helper := testhelper.NoError(createAndOpenDb())
	_, err := helper.Update(ctx, &model)
	assert.ErrorContains(t, err, "not found")
	assert.ErrorContains(t, err, "not_existing")
}

func TestSqlHelper_selectSingleRow_CanReturnMustBePtrToStructError(t *testing.T) {
	helper := SqlHelper{}
	var i int
	err := helper.selectSingleRow(ctx, &i, "")
	assert.Equal(t, true, assert.Error(t, err))
	assert.ErrorIs(t, &MustBePtrToStructError{}, err)
}

func TestSqlHelper_selectSingleValue(t *testing.T) {
	testhelper.HelperSetT(t)

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 10))

	var cnt int
	testhelper.NoError0(helper.SelectSingleValue(ctx, &cnt, "select count(*) from testTable"))
	assert.Equal(t, 10, cnt)
}

func TestSqlHelper_SelectById_ShouldSelectById(t *testing.T) {
	testhelper.HelperSetT(t)

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 10))

	var model = testModel{}

	testhelper.NoError0(helper.SelectById(ctx, &model, 3))
	assert.Equal(t, 3, model.Id, 3)
	assert.Equal(t, "username3", model.Name)
	assert.Equal(t, 3, model.Cnt)
}

func TestSqlHelper_Select_ShouldSelectRow(t *testing.T) {
	testhelper.HelperSetT(t)

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 10))

	for n := 1; n < 10; n++ {
		selectAndValidateModel(t, helper, n, n)
	}
}

func TestSqlHelper_Select_ShouldSelectManyRow(t *testing.T) {
	testhelper.HelperSetT(t)

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 10))

	var models []testModel

	testhelper.NoError0(helper.Select(ctx, &models, "where id < ? order by id", 4))
	assert.Equal(t, 3, len(models), 3)

	for index, model := range models {
		validateModel(t, &model, index+1, index+1)
	}
}

func TestSqlHelper_Select_CanReturnMustBePtrToStructOrSliceError(t *testing.T) {
	helper := SqlHelper{}
	var i int
	err := helper.Select(ctx, &i, "")
	assert.Equal(t, true, assert.Error(t, err))
	assert.ErrorIs(t, &MustBePtrToStructOrSliceError{}, err)
}
