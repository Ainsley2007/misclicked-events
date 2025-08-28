package sqlite

import (
	"database/sql"
)

type BotmDataSource interface {
	Start(serverID string, activityID int64, password string) error
	Stop(serverID string) error
	GetCurrentBotm(serverID string) (*BotmModel, error)
	GetCurrentBotmWithActivity(serverID string) (*BotmWithActivityModel, error)
}

func NewBotmDataSource(db *sql.DB) BotmDataSource {
	return &botmDS{db: db}
}

type botmDS struct{ db *sql.DB }

func (ds *botmDS) Start(serverID string, activityID int64, password string) error {
	_, err := ds.db.Exec(
		`INSERT INTO botm
            (server_id, activity_id, password, status)
          VALUES (?, ?, ?, ?)`,
		serverID,
		activityID,
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

func (ds *botmDS) GetCurrentBotm(serverID string) (*BotmModel, error) {
	row := ds.db.QueryRow(`
        SELECT id, server_id, activity_id, password, status
        FROM botm
        WHERE server_id = ? AND status = 'active'
      	ORDER BY id DESC
        LIMIT 1`,
		serverID,
	)

	var b BotmModel
	if err := row.Scan(
		&b.ID, &b.ServerID, &b.ActivityID, &b.Password, &b.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &b, nil
}

func (ds *botmDS) GetCurrentBotmWithActivity(serverID string) (*BotmWithActivityModel, error) {
	row := ds.db.QueryRow(`
        SELECT b.id, b.server_id, b.activity_id, b.password, b.status,
               a.id, a.name, a.type, a.hiscore_names, a.threshold
        FROM botm b
        JOIN activity a ON b.activity_id = a.id
        WHERE b.server_id = ? AND b.status = 'active'
      	ORDER BY b.id DESC
        LIMIT 1`,
		serverID,
	)

	var b BotmWithActivityModel
	var a ActivityModel
	if err := row.Scan(
		&b.ID, &b.ServerID, &b.ActivityID, &b.Password, &b.Status,
		&a.ID, &a.Name, &a.Type, &a.HiscoreNames, &a.Threshold,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	b.Activity = &a
	return &b, nil
}
