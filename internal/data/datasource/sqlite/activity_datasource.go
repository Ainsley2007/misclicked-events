package sqlite

import (
	"database/sql"
)

type ActivityDataSource interface {
	AddActivity(a *ActivityModel) error
	RemoveActivity(name string) error
	GetActivityListByType(activityType string) ([]*ActivityModel, error)
	GetAllActivities() ([]*ActivityModel, error)
	GetActivityByName(name string) (*ActivityModel, error)
}

func NewActivityDataSource(db *sql.DB) ActivityDataSource {
	return &activityDS{db}
}

type activityDS struct{ db *sql.DB }

func (ds *activityDS) AddActivity(a *ActivityModel) error {
	sqlStmt := `
		INSERT INTO activity(name, type, hiscore_names, threshold)
		VALUES(?, ?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET type = excluded.type, hiscore_names = excluded.hiscore_names, threshold = excluded.threshold`
	_, err := ds.db.Exec(sqlStmt, a.Name, a.Type, a.HiscoreNames, a.Threshold)
	return err
}

func (ds *activityDS) RemoveActivity(name string) error {
	_, err := ds.db.Exec(`DELETE FROM activity WHERE name = ?`, name)
	return err
}

func (ds *activityDS) GetActivityListByType(activityType string) ([]*ActivityModel, error) {
	rows, err := ds.db.Query(`SELECT id, name, type, hiscore_names, threshold FROM activity WHERE type = ?`, activityType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*ActivityModel
	for rows.Next() {
		a := &ActivityModel{}
		if err := rows.Scan(&a.ID, &a.Name, &a.Type, &a.HiscoreNames, &a.Threshold); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (ds *activityDS) GetAllActivities() ([]*ActivityModel, error) {
	rows, err := ds.db.Query(`SELECT id, name, type, hiscore_names, threshold FROM activity`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*ActivityModel
	for rows.Next() {
		a := &ActivityModel{}
		if err := rows.Scan(&a.ID, &a.Name, &a.Type, &a.HiscoreNames, &a.Threshold); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (ds *activityDS) GetActivityByName(name string) (*ActivityModel, error) {
	row := ds.db.QueryRow(`SELECT id, name, type, hiscore_names, threshold FROM activity WHERE name = ?`, name)
	var a ActivityModel
	if err := row.Scan(&a.ID, &a.Name, &a.Type, &a.HiscoreNames, &a.Threshold); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}
