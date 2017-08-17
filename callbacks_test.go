package nrgorm

import (
	"database/sql/driver"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/erikstmartin/go-testdb"
	newrelic "github.com/newrelic/go-agent"
)

type mockReporter struct {
	cb func(tableName string, sql string, op operation)
}

func (r *mockReporter) Report(startTime *newrelic.SegmentStartTime, tableName string, sql string, op operation) error {
	r.cb(tableName, sql, op)
	return nil
}

func Test_Callbacks_registerCallbacks(t *testing.T) {
	db, err := gorm.Open("testdb", "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	type call struct {
		tableName string
		sql       string
		op        operation
	}

	calls := []*call{}

	cb := newCallbacks(&mockReporter{cb: func(tableName, sql string, op operation) {
		calls = append(calls, &call{tableName: tableName, sql: sql, op: op})
	}})

	cb.registerCallbacks(db, "testdb")

	type Post struct{ Title string }

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (driver.Result, error) {
		return testdb.NewResult(1, nil, 1, nil), nil
	})
	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		return testdb.RowsFromCSVString([]string{"title"}, "awesome"), nil
	})

	cases := []struct {
		run func() error
		op  operation
	}{
		{
			run: func() error { return db.Create(&Post{Title: "awesome"}).Error },
			op:  operationCreate,
		},
		{
			run: func() error { return db.Model(&Post{}).Find(&Post{}, &Post{Title: "awesome"}).Error },
			op:  operationQuery,
		},
		{
			run: func() error { return db.Update(&Post{Title: "awesomeawesome"}).Error },
			op:  operationUpdate,
		},
		{
			run: func() error { return db.Delete(&Post{Title: "awesomeawesome"}).Error },
			op:  operationDelete,
		},
	}

	for i, c := range cases {
		if err := c.run(); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if got, want := calls[i].op, c.op; got != want {
			t.Errorf("Report() received op %v, want %v", got, want)
		}
	}

	if got, want := len(calls), len(cases); got != want {
		t.Errorf("Report() calls %d times, want %d times", got, want)
	}
}
