package main

import (
	"embed"
	"time"

	"github.com/getsentry/sentry-go"

	_ "git-ffd.kz/{{cookiecutter.namespace}}/{{cookiecutter.project_name | slugify | lower}}/docs"
	"git-ffd.kz/{{cookiecutter.namespace}}/{{cookiecutter.project_name | slugify | lower}}/internal/app"
	"git-ffd.kz/{{cookiecutter.namespace}}/{{cookiecutter.project_name | slugify | lower}}/internal/config"
	"git-ffd.kz/{{cookiecutter.namespace}}/{{cookiecutter.project_name | slugify | lower}}/pkg/db/migrate"
)

//go:generate go run github.com/swaggo/swag/cmd/swag init

//go:embed migrations/*.sql
var embedMigrations embed.FS

// @title {{cookiecutter.project_name | slugify | lower}}
// @версия 1.0.0
// @description {{cookiecutter.project_description}}
//
// @host api-dev.fmobile.kz
// @BasePath /{{cookiecutter.project_name | slugify | lower}}
// @schemes https http
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	defer func() {
		err := recover()
		hub := sentry.CurrentHub()

		if err != nil && hub != nil {
			hub.Recover(err)
			sentry.Flush(time.Second * 5)
		}

		if err != nil {
			panic(err)
		}
	}()

	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	migrate.SetBaseFS(embedMigrations)

	app.Run(cfg)
}
