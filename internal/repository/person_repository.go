package repository

import (
	"context"
	"fmt"

	"task_eff_mobile/internal/entity"

	"github.com/jackc/pgx/v5"
)

type PersonRepository struct {
	DB *pgx.Conn
}

type PeopleFilter struct {
	Name    string `validate:"omitempty,alpha"`
	Surname string `validate:"omitempty,alpha"`
	Gender  string `validate:"omitempty,alpha"`
	AgeMin  int    `validate:"omitempty,gte=0,lte=120"`
	AgeMax  int    `validate:"omitempty,gte=0,lte=120"`
	Page    int    `validate:"gte=1"`
	Limit   int    `validate:"gte=1,lte=100"`
}

func NewPersonRepository(conn *pgx.Conn) *PersonRepository {
	return &PersonRepository{DB: conn}
}

func (r *PersonRepository) Create(ctx context.Context, p *entity.Person) error {
	query := `
		INSERT INTO people (name, surname, patronymic, age, gender, nationalities)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	var id int
	err := r.DB.QueryRow(ctx, query,
		p.Name,
		p.Surname,
		p.Patronymic,
		p.Age,
		p.Gender,
		p.Nationalities,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("ошибка при вставке записи: %w", err)
	}

	p.ID = id
	return nil
}

func (r *PersonRepository) FindAll(ctx context.Context, filter PeopleFilter) ([]entity.Person, error) {
	query := `
	SELECT id, name, surname, patronymic, age, gender, nationalities
	FROM people
	WHERE 1=1
	`
	args := []any{}
	argID := 1

	if filter.Name != "" {
		query += fmt.Sprintf(" AND name LIKE $%d", argID)
		args = append(args, "%"+filter.Name+"%")
		argID++
	}
	if filter.Surname != "" {
		query += fmt.Sprintf(" AND surname LIKE $%d", argID)
		args = append(args, "%"+filter.Surname+"%")
		argID++
	}
	if filter.Gender != "" {
		query += fmt.Sprintf(" AND gender = $%d", argID)
		args = append(args, filter.Gender)
		argID++
	}
	if filter.AgeMin > 0 {
		query += fmt.Sprintf(" AND age >= $%d", argID)
		args = append(args, filter.AgeMin)
		argID++
	}
	if filter.AgeMax > filter.AgeMin {
		query += fmt.Sprintf(" AND age <= $%d", argID)
		args = append(args, filter.AgeMax)
		argID++
	}

	offset := (filter.Page - 1) * filter.Limit
	query += fmt.Sprintf(" ORDER BY id DESC LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, filter.Limit, offset)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []entity.Person
	for rows.Next() {
		var p entity.Person
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Surname, &p.Patronymic,
			&p.Age, &p.Gender, &p.Nationalities,
		); err != nil {
			return nil, err
		}
		people = append(people, p)
	}

	return people, nil
}

func (r *PersonRepository) Update(ctx context.Context, p *entity.Person) error {
	query := `
		UPDATE people
		SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationalities = $6
		WHERE id = $7
	`

	cmdTag, err := r.DB.Exec(ctx, query,
		p.Name,
		p.Surname,
		p.Patronymic,
		p.Age,
		p.Gender,
		p.Nationalities,
		p.ID,
	)

	if err != nil {
		return fmt.Errorf("ошибка обновления: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("человек с id %d не найден", p.ID)
	}
	return nil
}

func (r *PersonRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM people WHERE id = $1`

	cmdTag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("человек с id %d не найден", id)
	}
	return nil
}
