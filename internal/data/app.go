package data

import (
	"database/sql"
	"fmt"

	"misclicked-events/internal/data/datasource/api"
	"misclicked-events/internal/data/datasource/sqlite"
	"misclicked-events/internal/data/repository"
)

var (
	DB              *sql.DB
	ServerRepo      *repository.ServerRepository
	ConfigRepo      *repository.ConfigRepository
	CompetitionRepo *repository.CompetitionRepository
	HiscoreRepo     *repository.HiscoreRepository
	ParticipantRepo *repository.ParticipantRepository
)

func Init(dbPath string) error {
	var err error
	DB, err = sqlite.Init(dbPath)
	if err != nil {
		return fmt.Errorf("data.Init: %w", err)
	}

	sDS := sqlite.NewServerDataSource(DB)
	ServerRepo = repository.NewServerRepository(sDS)

	cDS := sqlite.NewConfigDataSource(DB)
	ConfigRepo = repository.NewConfigRepository(cDS)

	botmDS := sqlite.NewBotmDataSource(DB)
	kotsDS := sqlite.NewKotsDataSource(DB)
	CompetitionRepo = repository.NewCompetitionRepository(botmDS, kotsDS)

	hiscoreDS := api.NewHiscoreDataSource()
	HiscoreRepo = repository.NewHiscoreRepository(hiscoreDS)

	participantDS := sqlite.NewParticipantDataSource(DB)
	ParticipantRepo = repository.NewParticipantRepository(participantDS)

	return nil
}
