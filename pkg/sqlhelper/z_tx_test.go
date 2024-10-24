package sqlhelper

import (
	"context"
	sqldb "database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"sqlhelper/internal/utils"
	"sqlhelper/pkg/testhelper"
	"testing"
)

func TestSqlHelper_Transaction_Commit(t *testing.T) {
	testhelper.HelperSetT(t)

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 10))

	assert.Equal(t, 10, utils.Must(helper.Count(context.Background(), "")))

	err := helper.RunInTransaction(
		context.Background(),
		sqldb.TxOptions{
			Isolation: sqldb.LevelReadCommitted,
			ReadOnly:  false,
		},
		func(tx context.Context) error {
			_, err := helper.DeleteBySql(tx, "delete from testTable")
			assert.NoError(t, err)
			assert.Equal(t, 0, utils.Must(helper.Count(tx, "")))
			return nil
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, 0, utils.Must(helper.Count(context.Background(), "")))
}

func TestSqlHelper_Transaction_Rollback(t *testing.T) {
	testhelper.HelperSetT(t)

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 10))

	assert.Equal(t, 10, utils.Must(helper.Count(context.Background(), "")))

	err := helper.RunInTransaction(
		context.Background(),
		sqldb.TxOptions{
			Isolation: sqldb.LevelReadCommitted,
			ReadOnly:  false,
		},
		func(tx context.Context) error {
			_, err := helper.DeleteBySql(tx, "delete from testTable")
			assert.NoError(t, err)
			assert.Equal(t, 0, utils.Must(helper.Count(tx, "")))
			return errors.New("some error")
		},
	)
	assert.Error(t, err)
	assert.Equal(t, 10, utils.Must(helper.Count(context.Background(), "")))
}

func TestSqlHelper_Transaction_RollbackOnPanic(t *testing.T) {
	testhelper.HelperSetT(t)

	helper := testhelper.NoError(createAndOpenDb())
	testhelper.NoError0(createTestRecords(t, helper, 10))

	assert.Equal(t, 10, utils.Must(helper.Count(context.Background(), "")))

	err := helper.RunInTransaction(
		context.Background(),
		sqldb.TxOptions{
			Isolation: sqldb.LevelReadCommitted,
			ReadOnly:  false,
		},
		func(tx context.Context) error {
			_, err := helper.DeleteBySql(tx, "delete from testTable")
			assert.NoError(t, err)
			assert.Equal(t, 0, utils.Must(helper.Count(tx, "")))
			panic("panic message")
		},
	)
	assert.Error(t, err)
	var expectedError *PanicInTransactionError
	assert.ErrorAs(t, err, &expectedError)
	assert.Equal(t, 10, utils.Must(helper.Count(context.Background(), "")))
}
