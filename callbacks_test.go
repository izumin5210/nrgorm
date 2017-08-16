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

	err = db.Create(&Post{Title: "awesome"}).Error
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	err = db.Model(&Post{}).Find(&Post{}, &Post{Title: "awesome"}).Error
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if got, want := len(calls), 2; got != want {
		t.Errorf("reporter#report() called %d times, want %d times", got, want)
	}
}
