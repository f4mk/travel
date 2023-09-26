package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/f4mk/travel/backend/pkg/utils"
	"github.com/f4mk/travel/backend/travel-api/config"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/database"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/images"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"

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
		Host:       utils.GetHost(cfg.DB.HostName, cfg.DB.Port),
		Name:       cfg.DB.DBName,
		DisableTLS: cfg.DB.DisableTLS,
	})

	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	c := cron.New()

	_, err = c.AddFunc("0 0 * * *", removeExpiredRevokedTokens)
	if err != nil {
		fmt.Println("error scheduling removeExpiredRevokedTokens task:", err)
		return
	}

	_, err = c.AddFunc("0 1 * * *", removeExpiredResetTokens)
	if err != nil {
		fmt.Println("error scheduling removeExpiredResetTokens task:", err)
		return
	}

	_, err = c.AddFunc("0 2 * * *", removeExpiredVerificationTokens)
	if err != nil {
		fmt.Println("error scheduling removeExpiredVerificationTokens task:", err)
		return
	}

	_, err = c.AddFunc("0 3 * * *", removeUncommitedImages)
	if err != nil {
		fmt.Println("error scheduling removeUncommitedImages task:", err)
		return
	}

	c.Start()
	fmt.Println("cron has starter")
	select {}
}

func removeExpiredRevokedTokens() {
	for {
		q := `
		DELETE FROM revoked_tokens
		WHERE token_id IN 
				(SELECT token_id FROM revoked_tokens
				WHERE expires_at < $1
				LIMIT $2);`
		result, err := db.Exec(q, time.Now(), batchSize)
		if err != nil {
			fmt.Println("error removing revoked_tokens records:", err)
			return
		}
		removed, err := result.RowsAffected()
		if err != nil {
			fmt.Println("error getting revoked_tokens rows affected:", err)
			return
		}

		fmt.Println("cron removed entries from revoked_tokens:", removed)

		if removed < batchSize {
			break
		}

		time.Sleep(2 * time.Second)
	}
}

func removeExpiredResetTokens() {
	for {
		q := `
		DELETE FROM reset_tokens
		WHERE token_id IN 
				(SELECT token_id FROM reset_tokens
				WHERE expires_at < $1
				LIMIT $2);`
		result, err := db.Exec(q, time.Now(), batchSize)
		if err != nil {
			fmt.Println("error removing reset_tokens records:", err)
			return
		}
		removed, err := result.RowsAffected()
		if err != nil {
			fmt.Println("error getting reset_tokens rows affected:", err)
			return
		}

		fmt.Println("cron removed entries from reset_tokens:", removed)

		if removed < batchSize {
			break
		}

		time.Sleep(2 * time.Second)
	}
}

func removeExpiredVerificationTokens() {
	for {
		q := `
		DELETE FROM verify_tokens
		WHERE token_id IN 
				(SELECT token_id FROM verify_tokens
				WHERE expires_at < $1
				LIMIT $2);`
		result, err := db.Exec(q, time.Now(), batchSize)
		if err != nil {
			fmt.Println("error removing verify_tokens records:", err)
			return
		}
		removed, err := result.RowsAffected()
		if err != nil {
			fmt.Println("error getting verify_tokens rows affected:", err)
			return
		}

		fmt.Println("cron removed entries from verify_tokens:", removed)

		if removed < batchSize {
			break
		}

		time.Sleep(2 * time.Second)
	}
}

func removeUncommitedImages() {
	for {
		dayBefore := time.Now().Add(-24 * time.Hour)
		q := `
		DELETE FROM images
		WHERE image_id IN 
				(SELECT token_id FROM images
				WHERE status = $1 AND date_created < $2
				LIMIT $3);`
		result, err := db.Exec(q, images.Pending, dayBefore, batchSize)
		if err != nil {
			fmt.Println("error removing images records:", err)
			return
		}
		removed, err := result.RowsAffected()
		if err != nil {
			fmt.Println("error getting images rows affected:", err)
			return
		}

		fmt.Println("cron removed entries from images:", removed)

		if removed < batchSize {
			break
		}

		time.Sleep(2 * time.Second)
	}
}
