package nrgorm

import "strings"
import "github.com/jinzhu/gorm"
import "fmt"

type operation int

const (
	operationUnknown operation = iota
	operationQuery
	operationCreate
	operationUpdate
	operationDelete
)

const (
	namespace = "nrgorm"
)

func operations() []operation {
	return []operation{
		operationQuery,
		operationCreate,
		operationUpdate,
		operationDelete,
		operationUnknown,
	}
}

func (op operation) String() string {
	switch op {
	case operationQuery:
		return "SELECT"
	case operationCreate:
		return "INSERT"
	case operationUpdate:
		return "UPDATE"
	case operationDelete:
		return "DELETE"
	default:
		return ""
	}
}

func (op operation) Kind() string {
	switch op {
	case operationQuery:
		return "query"
	case operationCreate:
		return "create"
	case operationUpdate:
		return "update"
	case operationDelete:
		return "delete"
	default:
		return "row_query"
	}
}

func (op operation) Name(sql string) string {
	if op == operationUnknown {
		return strings.Split(sql, " ")[0]
	}
	return op.String()
}

func (op operation) callbackProcessor(db *gorm.DB) *gorm.CallbackProcessor {
	switch op {
	case operationQuery:
		return db.Callback().Query()
	case operationCreate:
		return db.Callback().Create()
	case operationUpdate:
		return db.Callback().Update()
	case operationDelete:
		return db.Callback().Delete()
	default:
		return db.Callback().RowQuery()
	}
}

func (op operation) registerBeforeCallback(db *gorm.DB, dbName string, callback func(*gorm.Scope)) {
	op.callbackProcessor(db).Before(op.callbackName()).Register(op.beforeCallbackName(dbName), callback)
}

func (op operation) registerAfterCallback(db *gorm.DB, dbName string, callback func(*gorm.Scope)) {
	op.callbackProcessor(db).After(op.callbackName()).Register(op.afterCallbackName(dbName), callback)
}

func (op operation) callbackName() string {
	return fmt.Sprintf("gorm:%s", op.Kind())
}

func (op operation) beforeCallbackName(dbName string) string {
	return fmt.Sprintf("%s:%s:%s_%s", namespace, dbName, op.callbackName(), "before")
}

func (op operation) afterCallbackName(dbName string) string {
	return fmt.Sprintf("%s:%s:%s_%s", namespace, dbName, op.callbackName(), "after")
}
