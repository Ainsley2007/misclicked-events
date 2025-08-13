package sqlite

import (
	"database/sql"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupBotmDB(t *testing.T) (*sql.DB, BotmDataSource) {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("opening sqlite: %v", err)
	}
	create := `
CREATE TABLE botm (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    server_id TEXT NOT NULL,
    activity_id INTEGER NOT NULL,
    password TEXT NOT NULL,
    status TEXT NOT NULL
);`
	if _, err := db.Exec(create); err != nil {
		t.Fatalf("creating botm table: %v", err)
	}
	return db, NewBotmDataSource(db)
}

func TestGetCurrentBotm_NoTable(t *testing.T) {
	_, ds := func() (*sql.DB, BotmDataSource) {
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("open sqlite: %v", err)
		}
		return db, NewBotmDataSource(db)
	}()
	_, err := ds.GetCurrentBotm("srv")
	if err == nil {
		t.Fatal("expected error when table botm does not exist")
	}
	if !strings.Contains(err.Error(), "no such table") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetCurrentBotm_NoActive(t *testing.T) {
	_, ds := setupBotmDB(t)
	b, err := ds.GetCurrentBotm("srv")
	if err != nil {
		t.Fatalf("GetCurrentBotm error: %v", err)
	}
	if b != nil {
		t.Errorf("expected nil Botm when no active run, got %+v", b)
	}
}

func TestStart_CreatesActive(t *testing.T) {
	_, ds := setupBotmDB(t)
	if err := ds.Start("srv1", 1, "pwdA"); err != nil {
		t.Fatalf("Start error: %v", err)
	}
	b, err := ds.GetCurrentBotm("srv1")
	if err != nil {
		t.Fatalf("GetCurrentBotm error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Botm after Start")
	}
	if b.ID != 1 ||
		b.ServerID != "srv1" ||
		b.ActivityID != 1 ||
		b.Password != "pwdA" ||
		b.Status != "active" {
		t.Errorf("got %+v, want ID=1 srv1/1/pwdA/active", b)
	}
}

func TestStart_AllowsMultipleActive(t *testing.T) {
	_, ds := setupBotmDB(t)
	if err := ds.Start("srvX", 1, "pass1"); err != nil {
		t.Fatalf("first Start: %v", err)
	}
	if err := ds.Start("srvX", 2, "pass2"); err != nil {
		t.Fatalf("second Start: %v", err)
	}

	b, err := ds.GetCurrentBotm("srvX")
	if err != nil {
		t.Fatalf("GetCurrentBotm after two starts: %v", err)
	}
	// Should get the most recent one (highest ID)
	if b.ID != 2 || b.ActivityID != 2 {
		t.Errorf("expected newest active ID=2 ActivityID=2, got %+v", b)
	}
}

func TestStop_MarksActiveAsDone(t *testing.T) {
	_, ds := setupBotmDB(t)
	if err := ds.Start("srv1", 1, "pwdA"); err != nil {
		t.Fatalf("Start error: %v", err)
	}

	if err := ds.Stop("srv1"); err != nil {
		t.Fatalf("Stop error: %v", err)
	}

	b, err := ds.GetCurrentBotm("srv1")
	if err != nil {
		t.Fatalf("GetCurrentBotm error: %v", err)
	}
	if b != nil {
		t.Errorf("expected nil Botm after Stop, got %+v", b)
	}
}

func TestStop_NoActiveToStop(t *testing.T) {
	_, ds := setupBotmDB(t)
	if err := ds.Stop("srv1"); err != nil {
		t.Fatalf("Stop error when no active BOTM: %v", err)
	}
}
