package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/f4mk/api/config"
	"github.com/f4mk/api/internal/pkg/database"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

var configPath = "config/.env"

const batchSize = 100

var db *sqlx.DB

func main() {

	cfg, err := config.New(configPath)
	if err != nil {
		log.Printf("could not read config: %v", err.Error())
		os.Exit(1)
	}

	db, err = database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       getHost(cfg.DB.HostName, cfg.DB.Port),
		Name:       cfg.DB.DBName,
		DisableTLS: cfg.DB.DisableTLS,
	})

	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	c := cron.New()

	_, err = c.AddFunc("0 0 * * *", removeExpiredRecords)
	if err != nil {
		fmt.Println("error scheduling task:", err)
		return
	}

	c.Start()
	fmt.Println("revoke cron has starter")
	select {}
}

func removeExpiredRecords() {
	for {
		q := `
		DELETE FROM revoked_tokens
		WHERE token_id IN (
				SELECT token_id FROM revoked_tokens
				WHERE expires_at < $1
				LIMIT $2
		);
		`
		result, err := db.Exec(q, time.Now(), batchSize)
		if err != nil {
			fmt.Println("error removing records:", err)
			return
		}
		removed, err := result.RowsAffected()
		if err != nil {
			fmt.Println("error getting rows affected:", err)
			return
		}

		fmt.Println("cron removed entries:", removed)

		if removed < batchSize {
			break
		}

		time.Sleep(2 * time.Second)
	}
}
func getHost(hostName string, port string) string {
	return hostName + ":" + port
}
