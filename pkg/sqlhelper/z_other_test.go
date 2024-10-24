package sqlhelper

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"sqlhelper/pkg/testhelper"
	"testing"
)

func TestSqlHelper_scan(t *testing.T) {
	testhelper.HelperSetT(t)

	var queryUser = testModel{}

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 1))
	db := helper.Db

	selectSql := testhelper.NoError(helper.getSelectSql(&queryUser, "where id=1"))
	rows := testhelper.NoError(db.Query(ctx, selectSql.sql))
	assert.Equal(t, true, rows.Next())
	testhelper.NoError0(helper.scan(rows, nil, nil, make([]DbTypeConverter, selectSql.selectedFieldsCount()), &queryUser, "id"))

	assert.Equal(t, 1, queryUser.Id)
	assert.Equal(t, "username1", queryUser.Name)
	assert.Equal(t, 1, queryUser.Cnt)
}
