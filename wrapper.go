package nrgorm

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/newrelic/go-agent"
)

const (
	txnKeyForGorm = "github.com/izumin5210/nrgorm.NewrelicTransaction"
	startTimeKey  = "github.com/izumin5210/nrgorm.TransactionStartTime"
)

// Wrap sets newrelic transaction instance to gorm.DB object
func Wrap(txn newrelic.Transaction, db *gorm.DB) *gorm.DB {
	if txn == nil {
		return db
	}
	return db.Set(txnKeyForGorm, txn)
}

// Wrapped returns true if given gorm.DB object has been set newrelic transaction
func Wrapped(db *gorm.DB) bool {
	_, ok := db.Get(txnKeyForGorm)
	return ok
}

func getTxn(scope *gorm.Scope) newrelic.Transaction {
	if v, ok := scope.Get(txnKeyForGorm); !ok {
		return nil
	} else if txn, ok := v.(newrelic.Transaction); ok {
		return txn
	}
	return nil
}

func setStartTimeToScope(scope *gorm.Scope) {
	startTime := newrelic.StartSegmentNow(getTxn(scope))
	scope.Set(getStartTimeKey(scope), &startTime)
}

func getStartTimeFromScope(scope *gorm.Scope) *newrelic.SegmentStartTime {
	if v, ok := scope.Get(getStartTimeKey(scope)); !ok {
		return nil
	} else if startTime, ok := v.(*newrelic.SegmentStartTime); ok {
		return startTime
	}
	return nil
}

func getStartTimeKey(scope *gorm.Scope) string {
	return fmt.Sprintf("%s#%s", startTimeKey, scope.InstanceID())
}
