package sqlite

import (
	"database/sql"
	"testing"
	"time"
)

func setupKotsDB(t *testing.T) (*sql.DB, KotsDataSource) {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("opening sqlite: %v", err)
	}
	// create kots table
	create := `
CREATE TABLE kots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    server_id TEXT NOT NULL,
    current_skill TEXT NOT NULL,
    current_king_participant INTEGER NOT NULL,
    streak INTEGER NOT NULL,
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    status TEXT NOT NULL
);`
	if _, err := db.Exec(create); err != nil {
		t.Fatalf("creating kots table: %v", err)
	}
	ds := NewKotsDataSource(db)
	return db, ds
}

func TestGetCurrentKots_NoActive(t *testing.T) {
	_, ds := setupKotsDB(t)
	k, err := ds.GetCurrentKots("srv")
	if err != nil {
		t.Fatalf("GetCurrentKots error: %v", err)
	}
	if k != nil {
		t.Errorf("expected nil when no active runs, got %+v", k)
	}
}

func TestInsertAndGetCurrentKots(t *testing.T) {
	_, ds := setupKotsDB(t)
	start := time.Now().UTC().Truncate(time.Second)
	end := start.Add(2 * time.Hour)
	// pending insert
	if _, err := ds.InsertNewKots("srv", "Skill1", 123, 5, start, end, "pending"); err != nil {
		t.Fatalf("InsertNewKots pending error: %v", err)
	}
	// active insert
	id, err := ds.InsertNewKots("srv", "Skill2", 456, 10, start.Add(time.Minute), end.Add(time.Minute), "active")
	if err != nil {
		t.Fatalf("InsertNewKots active error: %v", err)
	}
	got, err := ds.GetCurrentKots("srv")
	if err != nil {
		t.Fatalf("GetCurrentKots error: %v", err)
	}
	if got.ID != id || got.ServerID != "srv" || got.CurrentSkill != "Skill2" ||
		got.CurrentKingParticipant != 456 || got.Streak != 10 || got.Status != "active" {
		t.Errorf("GetCurrentKots = %+v", got)
	}
}

func TestMultipleActiveKots_SelectLatest(t *testing.T) {
	_, ds := setupKotsDB(t)
	base := time.Now().UTC().Truncate(time.Second)
	_, err := ds.InsertNewKots("srv", "S1", 1, 1, base, base.Add(time.Hour), "active")
	if err != nil {
		t.Fatalf("InsertNewKots #1: %v", err)
	}
	id2, err := ds.InsertNewKots("srv", "S2", 2, 2, base.Add(10*time.Minute), base.Add(2*time.Hour), "active")
	if err != nil {
		t.Fatalf("InsertNewKots #2: %v", err)
	}
	got, err := ds.GetCurrentKots("srv")
	if err != nil {
		t.Fatalf("GetCurrentKots: %v", err)
	}
	if got.ID != id2 {
		t.Errorf("expected ID %d, got %d", id2, got.ID)
	}
}
