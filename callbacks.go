package nrgorm

import (
	"github.com/jinzhu/gorm"
)

type callbacks struct {
	reporter reporter
}

// RegisterCallbacks register before and after callbacks of each operations to gorm.DB object
func RegisterCallbacks(db *gorm.DB, dbName string) {
	reporter := newReporter(db, dbName)
	c := newCallbacks(reporter)
	c.registerCallbacks(db, dbName)
}

func newCallbacks(reporter reporter) *callbacks {
	return &callbacks{reporter: reporter}
}

func (c *callbacks) registerCallbacks(db *gorm.DB, dbName string) {
	for _, op := range operations() {
		op.registerBeforeCallback(db, dbName, c.beforeFunc())
		op.registerAfterCallback(db, dbName, c.afterFunc(op))
	}
}

func (c *callbacks) beforeFunc() func(scope *gorm.Scope) {
	return func(scope *gorm.Scope) {
		setStartTimeToScope(scope)
	}
}

func (c *callbacks) afterFunc(op operation) func(scope *gorm.Scope) {
	return func(scope *gorm.Scope) {
		if scope.HasError() {
			return
		}
		startTime := getStartTimeFromScope(scope)
		c.reporter.Report(startTime, scope.TableName(), scope.SQL, op)
	}
}
