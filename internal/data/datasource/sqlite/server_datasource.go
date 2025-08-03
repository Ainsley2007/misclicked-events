package sqlite

import (
	"database/sql"
)

type ServerDataSource interface {
	CreateServer(s *ServerModel) error
	GetServer(id string) (*ServerModel, error)
	DeleteServer(id string) error
	ListServers() ([]*ServerModel, error)
}

func NewServerDataSource(db *sql.DB) ServerDataSource {
	return &serverDS{db}
}

type serverDS struct{ db *sql.DB }

func (ds *serverDS) CreateServer(s *ServerModel) error {
	sqlStmt := `
		INSERT INTO server(id, name)
		VALUES(?, ?)
		ON CONFLICT(id) DO UPDATE SET name = excluded.name`
	_, err := ds.db.Exec(sqlStmt, s.ID, s.Name)
	return err
}

func (ds *serverDS) GetServer(id string) (*ServerModel, error) {
	query := `SELECT name FROM server WHERE id = ?`
	row := ds.db.QueryRow(query, id)
	s := &ServerModel{ID: id}
	err := row.Scan(&s.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func (ds *serverDS) DeleteServer(id string) error {
	_, err := ds.db.Exec(`DELETE FROM server WHERE id = ?`, id)
	return err
}

func (ds *serverDS) ListServers() ([]*ServerModel, error) {
	rows, err := ds.db.Query(`SELECT id, name FROM server`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*ServerModel
	for rows.Next() {
		s := &ServerModel{}
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}
