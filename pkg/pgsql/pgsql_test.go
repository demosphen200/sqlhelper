package pgsql

import (
	"context"
	sqldb "database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"sqlhelper/internal/utils"
	"sqlhelper/pkg/sqlhelper"
	"sqlhelper/pkg/testhelper"
	"strconv"
	"testing"
	"time"
)

const DbUrl = "postgres://sqlhelper:sqlhelper@localhost:5432/sqlhelper"

var ctx = context.Background()

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

const initDbSql = `
drop table if exists testTable;
create table testTable(
	id serial primary key,
	id2 integer,
	name varchar(100),
	cnt int
);
delete from testTable;
`

func connectToDb() *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), DbUrl)
	if err != nil {
		panic(fmt.Sprintf("Unable to create connection pool: %v\n", err))
	}
	_, err = dbpool.Exec(ctx, initDbSql)
	if err != nil {
		panic(fmt.Sprintf("cannot init db structure: %v\n", err))
	}
	return dbpool
}

func createHelper() *sqlhelper.SqlHelper {
	pool := connectToDb()
	adapter := NewPgxPoolAdapter(pool)
	return sqlhelper.NewSqlHelper(adapter, "testTable")
}

func createTestModel(n int) testModel {
	return testModel{
		Id:   n,
		Name: "username" + strconv.Itoa(n),
		Cnt:  n,
	}
}

func createTestModels(cnt int) []testModel {
	models := make([]testModel, 0)
	for n := 1; n <= cnt; n++ {
		models = append(models, createTestModel(n))
	}
	return models
}

func createTestRecords(t *testing.T, helper *sqlhelper.SqlHelper, n int) {
	testhelper.HelperSetT(t)
	if n < 0 {
		panic("n must be > 0")
	}
	var model = testModel{}
	//helper.Db = db

	cnt := 1
	for cnt <= n {
		model = createTestModel(cnt)
		model.Id = 0
		//model.Name = "username" + strconv.Itoa(cnt)
		//model.Cnt = cnt
		_ = testhelper.NoError(helper.Db.Exec(
			ctx,
			"insert into testtable (id2,name,cnt) values ($1, $2, $3)",
			//@cnt,@name,@cnt2
			//sql.Named("cnt", cnt),
			//sql.Named("name", model.Name),
			//sql.Named("cnt2", cnt),
			cnt,
			model.Name,
			model.Cnt,
		))
		cnt++
	}
}

func TestPgx_SelectSingle(t *testing.T) {
	dbpool := connectToDb()
	defer dbpool.Close()

	h := createHelper()
	createTestRecords(t, h, 10)

	model := testModel{}
	assert.NoError(t, h.SelectById(ctx, &model, 3))
	assert.Equal(t, createTestModel(3), model)
}

func TestPgx_SelectTime(t *testing.T) {
	dbpool := connectToDb()
	defer dbpool.Close()

	h := createHelper()
	createTestRecords(t, h, 10)

	model := struct {
		Tm    time.Time  `db:"tm"`
		TmPtr *time.Time `db:"tm2"`
	}{}
	assert.Equal(t, true, model.Tm.Nanosecond() == 0)

	assert.NoError(t, h.SelectBySql(ctx, &model, "select now() as tm, now() as tm2"))
	assert.Equal(t, true, model.Tm.Nanosecond() > 0)
	assert.NotNil(t, model.TmPtr)
	assert.Equal(t, true, model.TmPtr.Nanosecond() > 0)
}

func TestPgx_SelectMany(t *testing.T) {
	dbpool := connectToDb()
	defer dbpool.Close()

	h := createHelper()
	createTestRecords(t, h, 10)

	var model []testModel
	assert.NoError(t, h.Select(ctx, &model, ""))
	assert.Equal(t, createTestModels(10), model)
}

func TestPgx_Update(t *testing.T) {
	dbpool := connectToDb()
	defer dbpool.Close()

	h := createHelper()
	createTestRecords(t, h, 10)

	var model = createTestModel(3)
	model.Name = "new name"
	model.Cnt = 1000000
	_, err := h.Update(ctx, &model)
	assert.NoError(t, err)

	expected := createTestModels(10)
	expected[2] = model

	var models []testModel
	assert.NoError(t, h.Select(ctx, &models, " order by id"))

	assert.Equal(t, expected, models)
}

func TestPgx_Delete(t *testing.T) {
	dbpool := connectToDb()
	defer dbpool.Close()

	h := createHelper()
	createTestRecords(t, h, 10)

	var model = testModel{
		Id: 3,
	}
	_, err := h.Delete(ctx, &model)
	assert.NoError(t, err)

	allModels := createTestModels(10)
	expected := append(make([]testModel, 0), allModels[:2]...)
	expected = append(expected, allModels[3:]...)

	var models []testModel
	assert.NoError(t, h.Select(ctx, &models, ""))

	assert.Equal(t, expected, models)
}

func TestPgx_Transaction_Commit(t *testing.T) {
	dbpool := connectToDb()
	defer dbpool.Close()

	h := createHelper()
	createTestRecords(t, h, 10)

	assert.Equal(t, 10, utils.Must(h.Count(context.Background(), "")))

	err := h.RunInTransaction(
		context.Background(),
		sqldb.TxOptions{
			Isolation: sqldb.LevelReadCommitted,
			ReadOnly:  false,
		},
		func(tx context.Context) error {
			_, err := h.DeleteBySql(tx, "delete from testTable")
			assert.NoError(t, err)
			assert.Equal(t, 0, utils.Must(h.Count(tx, "")))
			return nil
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, 0, utils.Must(h.Count(context.Background(), "")))
}

func TestPgx_Transaction_Rollback(t *testing.T) {
	dbpool := connectToDb()
	defer dbpool.Close()

	h := createHelper()
	createTestRecords(t, h, 10)

	assert.Equal(t, 10, utils.Must(h.Count(context.Background(), "")))

	err := h.RunInTransaction(
		context.Background(),
		sqldb.TxOptions{
			Isolation: sqldb.LevelReadCommitted,
			ReadOnly:  false,
		},
		func(tx context.Context) error {
			_, err := h.DeleteBySql(tx, "delete from testTable")
			assert.NoError(t, err)
			assert.Equal(t, 0, utils.Must(h.Count(tx, "")))
			return errors.New("some error")
		},
	)
	assert.Error(t, err)
	assert.Equal(t, 10, utils.Must(h.Count(context.Background(), "")))
}

func TestPgx_Transaction_RollbackOnPanic(t *testing.T) {
	dbpool := connectToDb()
	defer dbpool.Close()

	h := createHelper()
	createTestRecords(t, h, 10)

	assert.Equal(t, 10, utils.Must(h.Count(context.Background(), "")))

	err := h.RunInTransaction(
		context.Background(),
		sqldb.TxOptions{
			Isolation: sqldb.LevelReadCommitted,
			ReadOnly:  false,
		},
		func(tx context.Context) error {
			_, err := h.DeleteBySql(tx, "delete from testTable")
			assert.NoError(t, err)
			assert.Equal(t, 0, utils.Must(h.Count(tx, "")))
			panic("panic message")
		},
	)
	assert.Error(t, err)
	var expectedError *sqlhelper.PanicInTransactionError
	assert.ErrorAs(t, err, &expectedError)
	assert.Equal(t, 10, utils.Must(h.Count(context.Background(), "")))
}

func TestPgx_Transaction_CanReturnUnsupportedIsolationLevelError(t *testing.T) {
	dbpool := connectToDb()
	defer dbpool.Close()
	h := createHelper()
	err := h.RunInTransaction(
		context.Background(),
		sqldb.TxOptions{
			Isolation: sqldb.LevelWriteCommitted,
			ReadOnly:  false,
		},
		func(tx context.Context) error {
			return nil
		},
	)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sqlhelper.UnsupportedIsolationLevel)
}

func TestPgx_Transaction_RowsAffected(t *testing.T) {
	dbpool := connectToDb()
	defer dbpool.Close()

	h := createHelper()
	createTestRecords(t, h, 10)

	assert.Equal(t, 10, utils.Must(h.Count(context.Background(), "")))

	result, err := h.DeleteBySql(context.Background(), "delete from testTable")
	assert.NoError(t, err)

	rowsAffected, err := result.RowsAffected()
	assert.Equal(t, int64(10), rowsAffected)
}

func TestPgx_Transaction_LastInsertId(t *testing.T) {
	dbpool := connectToDb()
	defer dbpool.Close()

	h := createHelper()
	createTestRecords(t, h, 10)

	model := createTestModel(11)

	err := h.RunInTransaction(
		context.Background(),
		sqldb.TxOptions{
			Isolation: sqldb.LevelReadCommitted,
			ReadOnly:  false,
		},
		func(tx context.Context) error {
			result, err := h.Insert(tx, &model, false)
			assert.NoError(t, err)
			lastInsertId, err := result.LastInsertId()
			assert.NoError(t, err)
			assert.Equal(t, int64(11), lastInsertId)
			return nil
		},
	)
	assert.NoError(t, err)
}
