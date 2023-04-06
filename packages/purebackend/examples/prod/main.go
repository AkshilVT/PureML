package main

import (
	"log"
	"os"
	"strconv"

	"github.com/PureMLHQ/PureML/packages/purebackend"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/config"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/settings"
)

func main() {
	port, err := strconv.ParseInt(os.Getenv("PURE_MAIL_PORT"), 10, 64)
	if err != nil {
		log.Fatal("PURE_MAIL_PORT is not a number")
	}
	app := purebackend.NewWithConfig(&purebackend.Config{
		DatabaseType: config.GetDatabaseType(),
		DatabaseUrl:  config.GetDatabaseURL(),
		Settings: &settings.Settings{
			R2: settings.R2Config{
				Enabled:   true,
				AccountId: os.Getenv("PURE_R2_ACCOUNT_ID"),
				Bucket:    os.Getenv("PURE_R2_BUCKET"),
				Endpoint:  os.Getenv("PURE_R2_ENDPOINT"),
				AccessKey: os.Getenv("PURE_R2_ACCESS_KEY"),
				Secret:    os.Getenv("PURE_R2_SECRET"),
			},
			// Search: settings.SearchConfig{
			// 	Enabled:     true,
			// 	Host:        os.Getenv("PURE_SEARCH_HOST"),
			// 	AdminAPIKey: os.Getenv("PURE_SEARCH_ADMIN_API_KEY"),
			// },
			AdminAuthToken: settings.TokenConfig{
				Secret: os.Getenv("PURE_ADMIN_AUTH_TOKEN_SECRET"),
			},
			MailService: settings.MailServiceConfig{
				Enabled:  true,
				Host:     os.Getenv("PURE_MAIL_HOST"),
				Port:     int(port),
				Username: os.Getenv("PURE_MAIL_USER"),
				Password: os.Getenv("PURE_MAIL_PASS"),
			},
			Site: settings.SiteConfig{
				BaseURL: os.Getenv("PURE_SITE_BASE_URL"),
			},
		},
	})
	if err := app.Bootstrap(); err != nil {
		log.Fatal(err)
	}
	// Manually seed db if admin does not exist
	// Set admin details in .env file
	//
	// Refer:= URL_TO_README_GOES_HERE
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
