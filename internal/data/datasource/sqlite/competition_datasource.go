package sqlite

import (
	"database/sql"
)

type CompetitionDataSource interface {
	GetCompetition(serverID string) (*Competition, error)
	UpsertCompetition(c *Competition) error
	DeleteCompetition(serverID string) error
}

// NewCompetitionDataSource creates a SQLite CompetitionDataSource.
func NewCompetitionDataSource(db *sql.DB) CompetitionDataSource {
	return &competitionDS{db}
}

type competitionDS struct{ db *sql.DB }

// GetCompetition retrieves active competition or nil if none.
func (ds *competitionDS) GetCompetition(serverID string) (*Competition, error) {
	row := ds.db.QueryRow(
		`SELECT current_boss, password FROM competition WHERE server_id = ?`,
		serverID,
	)
	c := &Competition{ServerID: serverID}
	err := row.Scan(&c.CurrentBoss, &c.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

// UpsertCompetition creates or updates the competition row.
func (ds *competitionDS) UpsertCompetition(c *Competition) error {
	sqlStmt := `
		INSERT INTO competition(server_id, current_boss, password)
		VALUES(?, ?, ?)
		ON CONFLICT(server_id) DO UPDATE SET
		  current_boss = excluded.current_boss,
		  password     = excluded.password`
	_, err := ds.db.Exec(sqlStmt, c.ServerID, c.CurrentBoss, c.Password)
	return err
}

// DeleteCompetition removes the competition for a server.
func (ds *competitionDS) DeleteCompetition(serverID string) error {
	_, err := ds.db.Exec(
		`DELETE FROM competition WHERE server_id = ?`,
		serverID,
	)
	return err
}
