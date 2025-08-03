package repository

import (
	"misclicked-events/internal/data/datasource/sqlite"
)

type ServerRepository struct {
	ds sqlite.ServerDataSource
}

func NewServerRepository(ds sqlite.ServerDataSource) *ServerRepository {
	return &ServerRepository{ds}
}

func (r *ServerRepository) RegisterServer(id, name string) error {
	return r.ds.CreateServer(&sqlite.ServerModel{ID: id, Name: name})
}

func (r *ServerRepository) GetServer(id string) (*sqlite.ServerModel, error) {
	return r.ds.GetServer(id)
}

func (r *ServerRepository) RemoveServer(id string) error {
	return r.ds.DeleteServer(id)
}

func (r *ServerRepository) ListServers() ([]*sqlite.ServerModel, error) {
	return r.ds.ListServers()
}
