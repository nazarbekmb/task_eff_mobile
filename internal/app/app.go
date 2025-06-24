package app

import (
	"context"
	"os"
	"task_eff_mobile/internal/config"
	"task_eff_mobile/internal/handler"
	"task_eff_mobile/internal/repository"
	"task_eff_mobile/internal/service"
	"task_eff_mobile/internal/usecase"
	"task_eff_mobile/pkg"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start() error {
	if err := godotenv.Load(); err != nil {
		log.WithError(err).Fatal("Failed to load .env file")
	}
	log.Info("Configuration successfully loaded")

	ctx := context.Background()
	db, err := config.ConnectDB(ctx)
	if err != nil {
		log.WithError(err).Fatal("Database connection error")
	}
	defer db.Close(ctx)

	repo := repository.NewPersonRepository(db)
	enricher := service.NewEnricherService()
	useCase := usecase.NewPersonUseCase(repo, enricher)
	personHandler := handler.NewPersonHandler(useCase)

	router := gin.Default()
	router.Use(pkg.RequestLogger())
	router.POST("/people", personHandler.Create)
	router.GET("/people", personHandler.GetAll)
	router.PUT("/people/:id", personHandler.Update)
	router.DELETE("/people/:id", personHandler.Delete)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("SERVER_PORT")
	host := os.Getenv("SERVER_HOST")
	if port == "" || host == "" {
		log.Infof("port=%s host=%s", port, host)
	}
	if err := router.Run(host + ":" + port); err != nil {
		log.WithError(err).Fatal("Failed to launch")
	}

	return nil
}
