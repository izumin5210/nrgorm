package nrgorm

import (
	"testing"
)

func containOp(ops []operation, op operation) bool {
	for _, o := range ops {
		if o == op {
			return true
		}
	}
	return false
}

func Test_operations(t *testing.T) {
	ops := operations()

	wants := []operation{
		operationUnknown, // row_query
		operationQuery,   // query
		operationCreate,  // insert
		operationUpdate,  // update
		operationDelete,  // delete
	}

	if got, want := len(ops), len(wants); got != want {
		t.Errorf("operations() returned %d items, want %d items", got, want)
	}

	for _, want := range wants {
		if !containOp(ops, want) {
			t.Errorf("operations() returned slice and it should contain %v", want)
		}
	}
}

func Test_operation_String(t *testing.T) {
	cases := []struct {
		in  operation
		out string
	}{
		{in: operationUnknown, out: ""},
		{in: operationQuery, out: "SELECT"},
		{in: operationCreate, out: "INSERT"},
		{in: operationUpdate, out: "UPDATE"},
		{in: operationDelete, out: "DELETE"},
	}

	for _, c := range cases {
		if got, want := c.out, c.in.String(); got != want {
			t.Errorf("%v.String() returned %v, want %v", c.in, got, want)
		}
	}
}

func Test_operation_Kind(t *testing.T) {
	cases := []struct {
		in  operation
		out string
	}{
		{in: operationUnknown, out: "row_query"},
		{in: operationQuery, out: "query"},
		{in: operationCreate, out: "create"},
		{in: operationUpdate, out: "update"},
		{in: operationDelete, out: "delete"},
	}

	for _, c := range cases {
		if got, want := c.out, c.in.Kind(); got != want {
			t.Errorf("%v.Kind() returned %v, want %v", c.in, got, want)
		}
	}
}
