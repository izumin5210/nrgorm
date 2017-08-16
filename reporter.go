package nrgorm

import (
	"github.com/jinzhu/gorm"
	"github.com/newrelic/go-agent"
)

type reporter interface {
	Report(startTime *newrelic.SegmentStartTime, tableName string, sql string, op operation) error
}

type repoImpl struct {
	product newrelic.DatastoreProduct
	dbName  string
}

func newReporter(db *gorm.DB, dbName string) reporter {
	var product newrelic.DatastoreProduct
	switch db.Dialect().GetName() {
	case "postgres":
		product = newrelic.DatastorePostgres
	case "mysql":
		product = newrelic.DatastoreMySQL
	case "sqlite3":
		product = newrelic.DatastoreSQLite
	case "mssql":
		product = newrelic.DatastoreMSSQL
	default:
		// TODO: Should return an error
	}
	return &repoImpl{
		product: product,
		dbName:  dbName,
	}
}

func (r *repoImpl) Report(startTime *newrelic.SegmentStartTime, tableName string, sql string, op operation) error {
	seg := newrelic.DatastoreSegment{
		StartTime:          *startTime,
		Product:            r.product,
		Collection:         tableName,
		Operation:          op.Name(sql),
		ParameterizedQuery: sql,
		DatabaseName:       r.dbName,
	}
	return seg.End()
}
