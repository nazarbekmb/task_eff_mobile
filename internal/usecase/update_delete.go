package usecase

import (
	"context"
	"task_eff_mobile/internal/entity"
)

func (uc *PersonUseCase) UpdatePerson(ctx context.Context, p *entity.Person) error {
	return uc.Repo.Update(ctx, p)
}

func (uc *PersonUseCase) DeletePerson(ctx context.Context, id int) error {
	return uc.Repo.Delete(ctx, id)
}
