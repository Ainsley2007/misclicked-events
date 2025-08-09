package sqlite

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupServerDB(t *testing.T) (*sql.DB, ServerDataSource) {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}
	create := `
CREATE TABLE server (
    id TEXT PRIMARY KEY,
    name TEXT
);`
	if _, err := db.Exec(create); err != nil {
		t.Fatalf("failed to create server table: %v", err)
	}
	ds := NewServerDataSource(db)
	return db, ds
}

func TestCreateAndGetServer(t *testing.T) {
	_, ds := setupServerDB(t)
	err := ds.CreateServer(&ServerModel{ID: "srv1", Name: "Guild One"})
	if err != nil {
		t.Fatalf("CreateServer error: %v", err)
	}
	srv, err := ds.GetServer("srv1")
	if err != nil {
		t.Fatalf("GetServer error: %v", err)
	}
	expected := &ServerModel{ID: "srv1", Name: "Guild One"}
	if !reflect.DeepEqual(srv, expected) {
		t.Errorf("expected %+v, got %+v", expected, srv)
	}
}

func TestListAndDeleteServer(t *testing.T) {
	_, ds := setupServerDB(t)
	_ = ds.CreateServer(&ServerModel{ID: "srv1", Name: "Guild One"})
	_ = ds.CreateServer(&ServerModel{ID: "srv2", Name: "Guild Two"})
	list, err := ds.ListServers()
	if err != nil {
		t.Fatalf("ListServers error: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 servers, got %d", len(list))
	}
	if err := ds.DeleteServer("srv1"); err != nil {
		t.Fatalf("DeleteServer error: %v", err)
	}
	list, err = ds.ListServers()
	if err != nil {
		t.Fatalf("ListServers error after delete: %v", err)
	}
	if len(list) != 1 || list[0].ID != "srv2" {
		t.Errorf("expected only srv2 after delete, got %v", list)
	}
}
