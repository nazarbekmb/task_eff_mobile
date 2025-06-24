package entity

type Person struct {
	ID            int      `json:"id"`
	Name          string   `json:"name" validate:"required,alpha"`
	Surname       string   `json:"surname" validate:"required,alpha"`
	Patronymic    *string  `json:"patronymic,omitempty"`
	Age           int      `json:"age" validate:"required,gte=0,lte=120"`
	Gender        string   `json:"gender" validate:"required,alpha"`
	Nationalities []string `json:"nationalities" validate:"required,dive,alpha,len=2"`
}

type CreatePersonRequest struct {
	Name       string  `json:"name" validate:"required,alpha"`
	Surname    string  `json:"surname" validate:"required,alpha"`
	Patronymic *string `json:"patronymic,omitempty"`
}
