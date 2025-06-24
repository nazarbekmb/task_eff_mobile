package usecase

import (
	"context"
	"task_eff_mobile/internal/entity"
	"task_eff_mobile/internal/repository"
	"task_eff_mobile/internal/service"

	log "github.com/sirupsen/logrus"
)

type PersonUseCase struct {
	Repo     *repository.PersonRepository
	Enricher *service.EnricherService
}

func NewPersonUseCase(repo *repository.PersonRepository, enricher *service.EnricherService) *PersonUseCase {
	return &PersonUseCase{
		Repo:     repo,
		Enricher: enricher,
	}
}

func (uc *PersonUseCase) CreatePerson(ctx context.Context, req entity.CreatePersonRequest) (*entity.Person, error) {
	log.WithField("name", req.Name).Info("Let's start enriching the data")

	age, gender, nations, err := uc.Enricher.Enrich(req.Name)
	if err != nil {
		log.WithError(err).Error("Enrichment error")
		return nil, err
	}

	person := &entity.Person{
		Name:          req.Name,
		Surname:       req.Surname,
		Patronymic:    req.Patronymic,
		Age:           age,
		Gender:        gender,
		Nationalities: nations,
	}

	if err := uc.Repo.Create(ctx, person); err != nil {
		return nil, err
	}

	return person, nil
}
