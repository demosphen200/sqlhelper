package sqlhelper

import (
	"github.com/stretchr/testify/assert"
	"sqlhelper/pkg/testhelper"
	"testing"
)

func TestSqlHelper_TypeConverters_Insert_LocalToDb(t *testing.T) {
	testhelper.HelperSetT(t)
	var modelConv = testModelConv{
		Name: "username",
		Cnt:  11,
	}
	var expected = testModel{
		Id:   1,
		Name: "[username]",
		Cnt:  11,
	}
	var model testModel

	helper := testhelper.NoError(createAndOpenDb())
	_ = testhelper.NoError(helper.Insert(ctx, &modelConv, false))
	testhelper.NoError0(helper.SelectById(ctx, &model, 1))
	assert.Equal(t, expected, model)
}

func TestSqlHelper_TypeConverters_Update_LocalToDb(t *testing.T) {
	testhelper.HelperSetT(t)
	var initial = testModel{
		Name: "username",
		Cnt:  11,
	}
	var updated = testModelConv{
		Id:   1,
		Name: "new username",
		Cnt:  11,
	}
	var expected = testModel{
		Id:   1,
		Name: "[new username]",
		Cnt:  11,
	}
	var model testModel

	helper := testhelper.NoError(createAndOpenDb())
	_ = testhelper.NoError(helper.Insert(ctx, &initial, false))
	_ = testhelper.NoError(helper.Update(ctx, &updated))

	testhelper.NoError0(helper.SelectById(ctx, &model, 1))
	assert.Equal(t, expected, model)
}

func TestSqlHelper_TypeConverters_SelectById_DbToLocal(t *testing.T) {
	testhelper.HelperSetT(t)
	var model = testModel{
		Name: "username",
		Cnt:  11,
	}
	var expected = testModelConv{
		Id:   1,
		Name: "(username)",
		Cnt:  11,
	}
	var modelConv testModelConv

	helper := testhelper.NoError(createAndOpenDb())
	_ = testhelper.NoError(helper.Insert(ctx, &model, false))
	testhelper.NoError0(helper.SelectById(ctx, &modelConv, 1))
	assert.Equal(t, expected, modelConv)
}

func TestSqlHelper_TypeConverters_SelectSingle_DbToLocal(t *testing.T) {
	testhelper.HelperSetT(t)
	var model = testModel{
		Name: "username",
		Cnt:  11,
	}
	var expected = testModelConv{
		Id:   1,
		Name: "(username)",
		Cnt:  11,
	}

	var result testModelConv

	helper := testhelper.NoError(createAndOpenDb())
	_ = testhelper.NoError(helper.Insert(ctx, &model, false))
	testhelper.NoError0(helper.Select(ctx, &result, ""))
	assert.Equal(t, expected, result)
}

func TestSqlHelper_TypeConverters_SelectMany_DbToLocal(t *testing.T) {
	testhelper.HelperSetT(t)
	var model = testModel{
		Name: "username",
		Cnt:  11,
	}
	var expected = []testModelConv{
		{
			Id:   1,
			Name: "(username)",
			Cnt:  11,
		},
		{
			Id:   2,
			Name: "(username)",
			Cnt:  11,
		},
	}
	var result []testModelConv

	helper := testhelper.NoError(createAndOpenDb())
	_ = testhelper.NoError(helper.Insert(ctx, &model, false))
	_ = testhelper.NoError(helper.Insert(ctx, &model, false))
	testhelper.NoError0(helper.Select(ctx, &result, ""))
	assert.Equal(t, expected, result)
}

func TestSqlHelper_TypeConverters_SelectBySqlSingle_DbToLocal(t *testing.T) {
	testhelper.HelperSetT(t)
	var model = testModel{
		Name: "username",
		Cnt:  11,
	}
	var expected = testModelConv{
		Id:   1,
		Name: "(username)",
		Cnt:  11,
	}

	var result testModelConv

	helper := testhelper.NoError(createAndOpenDb())
	_ = testhelper.NoError(helper.Insert(ctx, &model, false))
	testhelper.NoError0(helper.SelectBySql(ctx, &result, "select * from testTable"))
	assert.Equal(t, expected, result)
}

func TestSqlHelper_TypeConverters_SelectBySqlMany_DbToLocal(t *testing.T) {
	testhelper.HelperSetT(t)
	var model = testModel{
		Name: "username",
		Cnt:  11,
	}
	var expected = []testModelConv{
		{
			Id:   1,
			Name: "(username)",
			Cnt:  11,
		},
	}
	var result []testModelConv

	helper := testhelper.NoError(createAndOpenDb())
	_ = testhelper.NoError(helper.Insert(ctx, &model, false))
	testhelper.NoError0(helper.SelectBySql(ctx, &result, "select * from testTable"))
	assert.Equal(t, expected, result)
}

type testName struct {
	NameA string `json:"nameA"`
	NameB string `json:"nameB"`
	Age   int    `json:"age"`
}

type testModelJson struct {
	Id   int      `db:"id" id:"true"`
	Name testName `db:"name" converter:"json"`
	Cnt  int      `db:"cnt"`
}

func TestSqlHelper_TypeConverters_Json_LocalToDb(t *testing.T) {
	testhelper.HelperSetT(t)
	var modelJson = testModelJson{
		Name: testName{
			NameA: "aaaaa",
			NameB: "bbbbb",
			Age:   123,
		},
		Cnt: 11,
	}
	model := testModel{}
	var expected = testModel{
		Id:   1,
		Name: `{"nameA":"aaaaa","nameB":"bbbbb","age":123}`,
		Cnt:  11,
	}

	helper := testhelper.NoError(createAndOpenDb())
	_ = testhelper.NoError(helper.Insert(ctx, &modelJson, false))
	testhelper.NoError0(helper.SelectBySql(ctx, &model, "select * from testTable"))
	assert.Equal(t, expected, model)
}

func TestSqlHelper_TypeConverters_Json_DbToLocal(t *testing.T) {
	testhelper.HelperSetT(t)
	var model = testModel{
		Name: `{"nameA":"aaaaa","nameB":"bbbbb","age":123}`,
		Cnt:  11,
	}
	modelJson := testModelJson{}
	var expected = testModelJson{
		Id: 1,
		Name: testName{
			NameA: "aaaaa",
			NameB: "bbbbb",
			Age:   123,
		},
		Cnt: 11,
	}

	helper := testhelper.NoError(createAndOpenDb())
	_ = testhelper.NoError(helper.Insert(ctx, &model, false))
	testhelper.NoError0(helper.SelectBySql(ctx, &modelJson, "select * from testTable"))
	assert.Equal(t, expected, modelJson)
}

func TestSqlHelper_TypeConverters_UpdateBySql_ParamsConvert(t *testing.T) {
	testhelper.HelperSetT(t)
	var model = testModel{
		Name: ``,
		Cnt:  11,
	}

	obj := struct {
		A int `json:"a"`
		B int `json:"b"`
	}{
		A: 1,
		B: 2,
	}

	var expected = testModel{
		Id:   1,
		Name: `{"a":1,"b":2}`,
		Cnt:  11,
	}

	helper := testhelper.NoError(createAndOpenDb())
	_ = testhelper.NoError(helper.Insert(ctx, &model, false))
	_ = testhelper.NoError(helper.UpdateBySql(ctx, "update testTable set name=?", obj))
	testhelper.NoError0(helper.SelectById(ctx, &model, 1))

	assert.Equal(t, expected, model)
}
