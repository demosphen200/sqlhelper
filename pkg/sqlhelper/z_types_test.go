package sqlhelper

import (
	"context"
	"fmt"
	"sqlhelper/pkg/testhelper"
	"testing"
	"time"
)

type timeModel struct {
	Time time.Time `db:"time_value"`
}

func TestSqlHelper_types_time(t *testing.T) {
	testhelper.HelperSetT(t)
	var model = timeModel{
		Time: time.Now(),
	}
	var model2 = timeModel{}

	helper := testhelper.NoError(createAndOpenDb())
	helper.TableName = "testTimeTable"
	_ = testhelper.NoError(helper.Insert(context.Background(), &model, false))
	//testhelper.NoError0(helper.SelectSingleValue(context.Background(), &model2.Time, "select * from testTimeTable"))
	testhelper.NoError0(helper.SelectSingleValueC(context.Background(), &model2.Time, "select CURRENT_TIMESTAMP", "datetime"))
	fmt.Printf("%v\n", model2.Time)

	//var qq sql.DB{}
	//r, e := qq.QueryRow()
	//r.Scan()
	//assert.Equal(t, "insert or replace into testTable (name,cnt) values (?,?)", insertSql)
	//assert.Equal(t, []any{"username", 11}, params)
}
