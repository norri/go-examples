package database

import (
	"reflect"
	"testing"
)

func TestNewDatabase_MemoryDB(t *testing.T) {
	db := NewDatabase()
	if reflect.TypeOf(db).String() != "*database.DB" {
		t.Fatalf("Expected *database.DB, but got %v", reflect.TypeOf(db).String())
	}
	if reflect.TypeOf(db.impl).String() != "*database.memoryDB" {
		t.Fatalf("Expected *database.memoryDB, but got %v", reflect.TypeOf(db).String())
	}
}
