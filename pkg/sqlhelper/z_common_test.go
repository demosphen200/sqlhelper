package sqlhelper

import (
	"context"
	sqldb "database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	"sqlhelper/pkg/testhelper"
	"strconv"
	"testing"
)

type testModel struct {
	Id   int    `db:"id" id:"true"`
	Name string `db:"name"`
	Cnt  int    `db:"cnt"`
}

type testModel2 struct {
	Id   int    `db:"id" id:"true"`
	Id2  int    `db:"id2" id:"true"`
	Name string `db:"name"`
	Cnt  int    `db:"cnt"`
}

type converterNotFoundModel struct {
	Id   int    `db:"id" id:"true"`
	Name string `db:"name" converter:"not_existing"`
}

type testModelConv struct {
	Id   int    `db:"id" id:"true"`
	Name string `db:"name" converter:"brackets"`
	Cnt  int    `db:"cnt"`
}

const initDbSql = `
create table testTable(
	id integer not null primary key autoincrement,
	id2 integer,
	name string,
	cnt int
);
create table testTimeTable(
    time_value time
)
`

var ctx = context.Background()

func createAndOpenDb() (*SqlHelper, error) {
	dbFile := "./test.db"
	_ = os.Remove(dbFile)
	if db, err := sqldb.Open("sqlite3", dbFile); err != nil {
		return nil, fmt.Errorf("cannot open db: %w", err)
	} else {
		if _, err = db.Exec(initDbSql); err != nil {
			return nil, fmt.Errorf("cannot create db structure: %w", err)
		}
		adapter := NewSqlAdapter(db)
		helper := NewSqlHelper(adapter, "testTable")
		helper.Converters.Register(NewAddBracketsConverter())
		RegisterTimeConverters(helper.Converters)
		RegisterJsonConverter(helper.Converters)
		return helper, nil
	}
}

func createTestRecords(t *testing.T, helper *SqlHelper, n int) error {
	testhelper.HelperSetT(t)
	if n < 0 {
		return errors.New("n must be > 0")
	}
	var model = testModel{}
	//helper.Db = db

	cnt := 1
	for cnt <= n {
		model.Name = "username" + strconv.Itoa(cnt)
		model.Cnt = cnt
		_ = testhelper.NoError(helper.Db.Exec(
			ctx,
			"insert into testTable (id2,name, cnt) values (?,?,?)",
			cnt,
			&model.Name,
			&model.Cnt,
		))
		cnt++
	}
	return nil
}

func validateModel(t *testing.T, model *testModel, id int, n int) {
	assert.Equal(t, id, model.Id)
	assert.Equal(t, "username"+strconv.Itoa(n), model.Name)
	assert.Equal(t, n, model.Cnt)
}

func validateModel2(t *testing.T, model2 *testModel2, id int, n int) {
	assert.Equal(t, id, model2.Id)
	assert.Equal(t, id, model2.Id2)
	assert.Equal(t, "username"+strconv.Itoa(n), model2.Name)
	assert.Equal(t, n, model2.Cnt)
}

func selectAndValidateModel(t *testing.T, helper *SqlHelper, id int, n int) {
	var model = testModel{}
	testhelper.NoError0(helper.Select(ctx, &model, "where id=?", id))
	validateModel(t, &model, id, n)
}
