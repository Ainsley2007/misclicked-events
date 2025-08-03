package sqlite

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupCompetitionDB(t *testing.T) (*sql.DB, CompetitionDataSource) {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}
	create := `
CREATE TABLE competition (
    server_id TEXT PRIMARY KEY,
    current_boss TEXT,
    password TEXT
);`
	if _, err := db.Exec(create); err != nil {
		t.Fatalf("failed to create competition table: %v", err)
	}
	ds := NewCompetitionDataSource(db)
	return db, ds
}

func TestGetCompetition_Default(t *testing.T) {
	_, ds := setupCompetitionDB(t)
	comp, err := ds.GetCompetition("srv1")
	if err != nil {
		t.Fatalf("GetCompetition error: %v", err)
	}
	if comp != nil {
		t.Errorf("expected nil competition, got %+v", comp)
	}
}

func TestUpsertAndGetCompetition(t *testing.T) {
	_, ds := setupCompetitionDB(t)
	input := &Competition{ServerID: "srv1", CurrentBoss: "Zulrah", Password: "pwd123"}
	if err := ds.UpsertCompetition(input); err != nil {
		t.Fatalf("UpsertCompetition error: %v", err)
	}
	got, err := ds.GetCompetition("srv1")
	if err != nil {
		t.Fatalf("GetCompetition error after upsert: %v", err)
	}
	if !reflect.DeepEqual(got, input) {
		t.Errorf("expected %+v, got %+v", input, got)
	}
}

func TestDeleteCompetition(t *testing.T) {
	_, ds := setupCompetitionDB(t)
	_ = ds.UpsertCompetition(&Competition{ServerID: "srv1", CurrentBoss: "Barrows", Password: "pw"})
	if err := ds.DeleteCompetition("srv1"); err != nil {
		t.Fatalf("DeleteCompetition error: %v", err)
	}
	comp, err := ds.GetCompetition("srv1")
	if err != nil {
		t.Fatalf("GetCompetition error after delete: %v", err)
	}
	if comp != nil {
		t.Errorf("expected nil after delete, got %+v", comp)
	}
}
