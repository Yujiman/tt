package postgresql

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm/logger"
	"time"

	"git-ffd.kz/pkg/golog"
	gormlog "git-ffd.kz/pkg/golog/contrib/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func NewDB(
	sources []string,
	replicas []string,
	appName string,
	isLocal bool,
	log golog.ContextLogger,
) (*sql.DB, *gorm.DB, error) {
	log.Infow(
		"Start connecting to DB",
		"sources", sources,
		"replicas", replicas,
	)

	if len(sources) == 0 {
		return nil, nil, fmt.Errorf("sources database connetions is empty")
	}

	var gormLogger logger.Interface = gormlog.NewGormLogAdapter(log)
	if isLocal {
		gormLogger = gormlog.NewGormLogAdapter(log).LogMode(logger.Info)
	}

	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  fmt.Sprintf("%s?application_name=%s", sources[0], appName),
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}),
		&gorm.Config{
			Logger:      gormLogger,
			PrepareStmt: false,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	if len(sources) > 1 || len(replicas) > 0 {
		log.Infow("Database: setup multi source/replica")

		sources = sources[1:]
		sourceConnections := make([]gorm.Dialector, len(sources))
		replicaConnections := make([]gorm.Dialector, len(replicas))

		for i, sourceDSN := range sources {
			sourceConnections[i] = postgres.New(postgres.Config{
				DSN:                  fmt.Sprintf("%s?application_name=%s", sourceDSN, appName),
				PreferSimpleProtocol: true,
			})
		}

		for i, replicaDSN := range replicas {
			replicaConnections[i] = postgres.New(postgres.Config{
				DSN:                  fmt.Sprintf("%s?application_name=%s", replicaDSN, appName),
				PreferSimpleProtocol: true,
			})
		}

		err = db.Use(dbresolver.Register(
			dbresolver.Config{
				Sources:  sourceConnections,
				Replicas: replicaConnections,
				Policy:   dbresolver.RandomPolicy{},
			},
		))
		if err != nil {
			return nil, nil, err
		}
	}

	connection, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	connection.SetConnMaxLifetime(time.Hour)
	connection.SetMaxIdleConns(15)
	connection.SetMaxOpenConns(15)

	log.Infow("DB connected")

	return connection, db, nil
}
