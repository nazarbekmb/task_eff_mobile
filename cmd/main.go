// @title Task Eff Mobile API
// @version 1.0
// @description API для обработки и обогащения персональных данных
// @host localhost:8080
// @BasePath /
package main

import (
	_ "task_eff_mobile/docs"
	"task_eff_mobile/internal/app"

	log "github.com/sirupsen/logrus"
)

func main() {
	if err := app.Start(); err != nil {
		log.WithError(err).Fatal("The server terminated with an error")
		return
	}
}
