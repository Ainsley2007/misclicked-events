package repository

import (
	"misclicked-events/internal/data/datasource/sqlite"
	"misclicked-events/internal/data/mappers"
	"misclicked-events/internal/domain"
)

type ServerRepository struct {
	ds     sqlite.ServerDataSource
	mapper *mappers.ServerMapper
}

func NewServerRepository(ds sqlite.ServerDataSource) *ServerRepository {
	return &ServerRepository{
		ds:     ds,
		mapper: mappers.NewServerMapper(),
	}
}

func (r *ServerRepository) RegisterServer(id, name string) error {
	server := &domain.Server{ID: id, Name: name}
	serverModel := r.mapper.ToModel(server)
	return r.ds.CreateServer(serverModel)
}

func (r *ServerRepository) GetServer(id string) (*domain.Server, error) {
	serverModel, err := r.ds.GetServer(id)
	if err != nil {
		return nil, err
	}
	return r.mapper.ToDomain(serverModel), nil
}

func (r *ServerRepository) RemoveServer(id string) error {
	return r.ds.DeleteServer(id)
}

func (r *ServerRepository) ListServers() ([]*domain.Server, error) {
	serverModels, err := r.ds.ListServers()
	if err != nil {
		return nil, err
	}

	servers := make([]*domain.Server, len(serverModels))
	for i, serverModel := range serverModels {
		servers[i] = r.mapper.ToDomain(serverModel)
	}
	return servers, nil
}
