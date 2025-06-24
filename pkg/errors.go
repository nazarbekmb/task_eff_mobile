package pkg

import (
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		errors["global"] = "ошибка валидации"
		return errors
	}

	for _, ve := range validationErrors {
		field := ve.Field()
		tag := ve.Tag()

		switch tag {
		case "required":
			errors[field] = "обязательное поле"
		case "alpha":
			errors[field] = "только латинские буквы"
		case "gte":
			errors[field] = "должно быть больше или равно минимальному"
		case "lte":
			errors[field] = "должно быть меньше или равно максимальному"
		case "len":
			errors[field] = "должно содержать ровно 2 символа (например, RU)"
		default:
			errors[field] = "некорректное значение"
		}
	}

	return errors
}
