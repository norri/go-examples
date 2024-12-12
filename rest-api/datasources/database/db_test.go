package database

import (
	"reflect"
	"testing"
)

func TestNewDatabase_MemoryDB(t *testing.T) {
	db := NewDatabase()
	if reflect.TypeOf(db).String() != "*database.memoryDB" {
		t.Fatalf("Expected *database.memoryDB, but got %v", reflect.TypeOf(db).String())
	}
}
