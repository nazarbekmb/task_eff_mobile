package usecase

import (
	"context"
	"task_eff_mobile/internal/entity"
	"task_eff_mobile/internal/repository"
)

func (uc *PersonUseCase) GetPeople(ctx context.Context, filter repository.PeopleFilter) ([]entity.Person, error) {
	return uc.Repo.FindAll(ctx, filter)
}
