package sqlite

import (
	"database/sql"
	"time"
)

type KotsDataSource interface {
	InsertNewKots(
		serverID, skill string,
		kingParticipant int64,
		streak int,
		start, end time.Time,
		status string,
	) (int64, error)

	GetCurrentKots(serverID string) (*KotsModel, error)
}

func NewKotsDataSource(db *sql.DB) KotsDataSource {
	return &kotsDS{db: db}
}

type kotsDS struct{ db *sql.DB }

func (ds *kotsDS) InsertNewKots(
	serverID, skill string,
	kingParticipant int64,
	streak int,
	start, end time.Time,
	status string,
) (int64, error) {
	res, err := ds.db.Exec(
		`INSERT INTO kots
		(server_id, current_skill, current_king_participant, streak, start_date, end_date, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		serverID, skill, kingParticipant, streak,
		start.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano),
		status,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (ds *kotsDS) GetCurrentKots(serverID string) (*KotsModel, error) {
	row := ds.db.QueryRow(
		`SELECT id, server_id, current_skill, current_king_participant, streak, start_date, end_date, status
		   FROM kots
		  WHERE server_id = ? AND status = 'active'
		  ORDER BY start_date DESC
		  LIMIT 1`,
		serverID,
	)
	var k KotsModel
	var startStr, endStr string
	if err := row.Scan(
		&k.ID, &k.ServerID, &k.CurrentSkill, &k.CurrentKingParticipant,
		&k.Streak, &startStr, &endStr, &k.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	k.StartDate, _ = time.Parse(time.RFC3339Nano, startStr)
	if endStr != "" {
		ee, _ := time.Parse(time.RFC3339Nano, endStr)
		k.EndDate = &ee
	}
	return &k, nil
}
