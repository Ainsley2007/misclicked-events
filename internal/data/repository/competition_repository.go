package repository

import (
	"misclicked-events/internal/data/datasource/sqlite"
)

type CompetitionRepository struct {
	ds sqlite.CompetitionDataSource
}

func NewCompetitionRepository(ds sqlite.CompetitionDataSource) *CompetitionRepository {
	return &CompetitionRepository{ds}
}

func (r *CompetitionRepository) FetchCompetition(serverID string) (*sqlite.Competition, error) {
	return r.ds.GetCompetition(serverID)
}

func (r *CompetitionRepository) StartCompetition(c *sqlite.Competition) error {
	return r.ds.UpsertCompetition(c)
}

func (r *CompetitionRepository) EndCompetition(serverID string) error {
	return r.ds.DeleteCompetition(serverID)
}
