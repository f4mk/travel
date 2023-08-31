package utils

import (
	"log"

	"github.com/f4mk/travel/backend/pkg/utils"
	"github.com/f4mk/travel/backend/travel-api/config"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/database"
	"github.com/golang-migrate/migrate/v4"

	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	//for using file source
	_ "github.com/golang-migrate/migrate/v4/source/file"
	//for using postgres driver
	_ "github.com/jackc/pgx/v5/stdlib"
)

func RunMigration(cfg *config.Config, mp string) {
	db, err := database.Open(database.Config{
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		// TODO: need to provied hostname when run as a standalone script
		Host:       utils.GetHost(cfg.DB.HostName, cfg.DB.Port),
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

	log.Println("applying migrations")
	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalln(err)
	}
}
