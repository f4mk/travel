package utils

import (
	"log"

	"github.com/f4mk/travel/backend/travel-api/config"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/database"
	"github.com/f4mk/travel/backend/travel-api/pkg/utils"
	"github.com/golang-migrate/migrate/v4"

	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	//for using file source
	_ "github.com/golang-migrate/migrate/v4/source/file"
	//for using postgres driver
	_ "github.com/lib/pq"
)

func RunMigration(cfg *config.Config, mp string) {
	db, err := database.Open(database.Config{
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		// TODO: shouldnt be hardcoded
		Host:       utils.GetHost("0.0.0.0", cfg.DB.Port),
		Name:       cfg.DB.DBName,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	driver, err := postgresMigrate.WithInstance(db.DB, &postgresMigrate.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		mp,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalln(err)
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalln(err)
	}
}
