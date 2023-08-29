package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/f4mk/travel/backend/travel-api/cmd/admin/utils"
	"github.com/f4mk/travel/backend/travel-api/config"
)

var configPath = "config/.env"
var migrationsPath = "file://cmd/admin/migration/sql"

func main() {

	cfg, err := config.New(configPath)
	if err != nil {
		log.Printf("could not read config: %v", err.Error())
		os.Exit(1)
	}

	keyFlag := flag.Bool("keygen", false, "To generate secure key for JWT")
	tokenFlag := flag.Bool("tokengen", false, "To generate JWT token")
	kidFlag := flag.String("kid", "", "KID for generate JWT token")
	roleFlag := flag.String("role", "", "Role for generate JWT token")
	tokenAllFlag := flag.Bool("tokengen-all", false, "To generate JWT token")
	migrateFlag := flag.Bool("migrate", false, "To perform migration")

	flag.Parse()

	if *keyFlag {
		utils.GenerateKey(cfg)
	}

	if *tokenFlag {

		if *kidFlag == "" {
			log.Println("kid is required")
			os.Exit(1)
		}

		if *roleFlag == "" {
			log.Println("role is required")
			os.Exit(1)
		}

		token, _ := utils.GenerateToken(cfg, *kidFlag, []string{*roleFlag})
		fmt.Println("===== Token =====")
		fmt.Println(token[*roleFlag])
	}

	if *tokenAllFlag {

		if *roleFlag == "" {
			log.Println("role is required")
			os.Exit(1)
		}

		tokens, _ := utils.GenerateAllTokens(cfg, []string{*roleFlag})
		for _, tokenMap := range tokens {
			for kid, token := range tokenMap {
				fmt.Println("====== KID ======")
				fmt.Println(kid)
				fmt.Println("===== Token =====")
				fmt.Println(token)
			}
		}
	}

	if *migrateFlag {
		utils.RunMigration(cfg, migrationsPath)
	}

}
