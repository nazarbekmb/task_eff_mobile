package handler

import (
	"net/http"
	"strconv"
	"task_eff_mobile/internal/entity"
	"task_eff_mobile/internal/repository"
	"task_eff_mobile/internal/usecase"
	"task_eff_mobile/pkg"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type PersonHandler struct {
	UC *usecase.PersonUseCase
}

func NewPersonHandler(uc *usecase.PersonUseCase) *PersonHandler {
	return &PersonHandler{UC: uc}
}

// Create godoc
// @Summary Добавить нового человека
// @Description Добавляет человека и обогащает его через внешние API
// @Tags people
// @Accept json
// @Produce json
// @Param person body entity.CreatePersonRequest true "Данные человека"
// @Success 201 {object} entity.Person
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people [post]
func (h *PersonHandler) Create(c *gin.Context) {
	var req entity.CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("Invalid request format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}

	if err := pkg.Validate.Struct(req); err != nil {
		errors := pkg.ParseValidationErrors(err)

		log.WithField("errors", errors).Warn("Validation error")
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	log.WithFields(log.Fields{
		"name":    req.Name,
		"surname": req.Surname,
	}).Info("User creation request received")

	person, err := h.UC.CreatePerson(c.Request.Context(), req)
	if err != nil {
		log.WithError(err).Error("User creation error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать человека"})
		return
	}

	log.WithField("id", person.ID).Info("User successfully created")
	c.JSON(http.StatusCreated, person)
}

// GetAll godoc
// @Summary Получить список людей
// @Description Получение списка людей с фильтрами и пагинацией
// @Tags people
// @Produce json
// @Param name query string false "Фильтр по имени"
// @Param surname query string false "Фильтр по фамилии"
// @Param gender query string false "Пол"
// @Param age_min query int false "Мин. возраст"
// @Param age_max query int false "Макс. возраст"
// @Param page query int false "Страница"
// @Param limit query int false "Лимит"
// @Success 200 {array} entity.Person
// @Router /people [get]
func (h *PersonHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	ageMin, _ := strconv.Atoi(c.DefaultQuery("age_min", "0"))
	ageMax, _ := strconv.Atoi(c.DefaultQuery("age_max", "0"))

	filter := repository.PeopleFilter{
		Name:    c.Query("name"),
		Surname: c.Query("surname"),
		Gender:  c.Query("gender"),
		AgeMin:  ageMin,
		AgeMax:  ageMax,
		Page:    page,
		Limit:   limit,
	}
	if err := pkg.Validate.Struct(filter); err != nil {
		errors := pkg.ParseValidationErrors(err)

		log.WithField("errors", errors).Warn("Filter validation error")
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	log.Info("Request to get the list of users")

	people, err := h.UC.GetPeople(c.Request.Context(), filter)
	if err != nil {
		log.WithError(err).Error("Error when retrieving the list of users")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении данных"})
		return
	}
	if len(people) == 0 {
		log.Warn("Users not found by filter")
		c.JSON(http.StatusOK, gin.H{
			"data":    people,
			"message": "Ничего не найдено",
		})

		return
	}
	c.JSON(http.StatusOK, people)
}

// Update godoc
// @Summary Обновить данные человека
// @Description Обновляет данные человека по ID
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Param person body entity.Person true "Обновлённые данные"
// @Success 200 {object} map[string]string
// @Failure 400,404,500 {object} map[string]string
// @Router /people/{id} [put]
func (h *PersonHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID"})
		return
	}
	log.WithField("id", id).Info("User data update initiated")

	var person entity.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}
	if err := pkg.Validate.Struct(&person); err != nil {
		errors := pkg.ParseValidationErrors(err)

		log.WithField("errors", errors).Warn("Validation error")
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}
	person.ID = id

	if err := h.UC.UpdatePerson(c.Request.Context(), &person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "обновлено"})
}

// Delete godoc
// @Summary Удалить человека
// @Description Удаляет человека по ID
// @Tags people
// @Param id path int true "ID человека"
// @Success 200 {object} map[string]string
// @Failure 400,404,500 {object} map[string]string
// @Router /people/{id} [delete]
func (h *PersonHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID"})
		return
	}

	if err := h.UC.DeletePerson(c.Request.Context(), id); err != nil {
		log.WithFields(log.Fields{
			"id":    id,
			"error": err,
		}).Error("User deletion error")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.WithField("id", id).Info("User deleted")
	c.JSON(http.StatusOK, gin.H{"message": "удалено"})
}
