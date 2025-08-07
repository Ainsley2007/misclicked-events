package sqlite

import (
	"database/sql"
)

type BotmDataSource interface {
	StartNewBotm(serverID, boss, password string) error
	Start(serverID, boss, password string) error
	Stop(serverID string) error
	GetCurrentBotm(serverID string) (*Botm, error)
}

func NewBotmDataSource(db *sql.DB) BotmDataSource {
	return &botmDS{db: db}
}

type botmDS struct{ db *sql.DB }

func (ds *botmDS) GetCurrentBotm(serverID string) (*Botm, error) {
	row := ds.db.QueryRow(`
        SELECT id, server_id, current_boss, password, status
        FROM botm
        WHERE server_id = ? AND status = 'active'
      	ORDER BY id DESC
        LIMIT 1`,
		serverID,
	)
	var b Botm
	if err := row.Scan(
		&b.ID, &b.ServerID, &b.CurrentBoss, &b.Password, &b.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

func (ds *botmDS) Start(serverID, boss, password string) error {
	_, err := ds.db.Exec(
		`INSERT INTO botm
            (server_id, current_boss, password, status)
          VALUES (?, ?, ?, ?)`,
		serverID,
		boss,
		password,
		"active",
	)
	return err
}

func (ds *botmDS) Stop(serverID string) error {
	_, err := ds.db.Exec(
		`UPDATE botm
            SET status = ?
          WHERE server_id = ? 
            AND status = ?`,
		"done",
		serverID,
		"active",
	)
	return err
}

func (ds *botmDS) StartNewBotm(serverID, boss, password string) error {
	tx, err := ds.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		`UPDATE botm
            SET status = ?
          WHERE server_id = ? 
            AND status = ?`,
		"done",
		serverID,
		"active",
	); err != nil {
		return err
	}

	if _, err := tx.Exec(
		`INSERT INTO botm
            (server_id, current_boss, password, status)
          VALUES (?, ?, ?, ?)`,
		serverID,
		boss,
		password,
		"active",
	); err != nil {
		return err
	}

	return tx.Commit()
}
