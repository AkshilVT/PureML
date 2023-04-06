package main

import (
	"log"
	"os"
	"strconv"

	"github.com/PureMLHQ/PureML/packages/purebackend"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/settings"
)

func main() {
	databaseType := os.Getenv("DATABASE")
	if databaseType == "" {
		databaseType = "sqlite3"
	}
	portEnv := os.Getenv("PURE_MAIL_PORT")
	var port int64
	var err error
	if portEnv != "" {
		port, err = strconv.ParseInt(portEnv, 10, 64)
		if err != nil {
			log.Fatal("PURE_MAIL_PORT is not a number")
		}
	}

	appSettings := &settings.Settings{
		// set the default settings
		S3: settings.S3Config{
			Enabled:   os.Getenv("PURE_S3_ENABLE") == "true",
			Bucket:    os.Getenv("PURE_S3_BUCKET"),
			Region:    os.Getenv("PURE_S3_REGION"),
			Endpoint:  os.Getenv("PURE_S3_ENDPOINT"),
			AccessKey: os.Getenv("PURE_S3_ACCESS_KEY"),
			Secret:    os.Getenv("PURE_S3_SECRET"),
		},
		R2: settings.R2Config{
			Enabled:   os.Getenv("PURE_R2_ENABLE") == "true",
			AccountId: os.Getenv("PURE_R2_ACCOUNT_ID"),
			Bucket:    os.Getenv("PURE_R2_BUCKET"),
			Endpoint:  os.Getenv("PURE_R2_ENDPOINT"),
			AccessKey: os.Getenv("PURE_R2_ACCESS_KEY"),
			Secret:    os.Getenv("PURE_R2_SECRET"),
		},
		Search: settings.SearchConfig{
			Enabled: false,
		},
		Site: settings.SiteConfig{
			BaseURL: os.Getenv("PURE_SITE_BASE_URL"),
		},
		MailService: settings.MailServiceConfig{
			Enabled:  os.Getenv("PURE_MAIL_ENABLE") == "true",
			Host:     os.Getenv("PURE_MAIL_HOST"),
			Port:     int(port),
			Username: os.Getenv("PURE_MAIL_USER"),
			Password: os.Getenv("PURE_MAIL_PASS"),
		},
	}

	adminAuthTokenSecret := os.Getenv("PURE_ADMIN_AUTH_TOKEN_SECRET")
	if adminAuthTokenSecret != "" {
		adminAuthTokenSettings := settings.TokenConfig{
			Secret: os.Getenv("PURE_ADMIN_AUTH_TOKEN_SECRET"),
		}
		adminAuthTokenDuration := os.Getenv("PURE_ADMIN_AUTH_TOKEN_DURATION")
		if adminAuthTokenDuration != "" {
			duration, err := strconv.ParseInt(portEnv, 10, 64)
			if err != nil {
				log.Fatal("PURE_ADMIN_AUTH_TOKEN_DURATION is not a number")
			}
			adminAuthTokenSettings.Duration = duration
		}
		appSettings.AdminAuthToken = adminAuthTokenSettings
	}
	
	passwordResetAuthTokenSecret := os.Getenv("PURE_PASSWORD_RESET_TOKEN_SECRET")
	if passwordResetAuthTokenSecret != "" {
		passwordResetAuthTokenSettings := settings.TokenConfig{
			Secret: os.Getenv("PURE_PASSWORD_RESET_TOKEN_SECRET"),
		}
		passwordResetAuthTokenDuration := os.Getenv("PURE_PASSWORD_RESET_TOKEN_DURATION")
		if passwordResetAuthTokenDuration != "" {
			duration, err := strconv.ParseInt(portEnv, 10, 64)
			if err != nil {
				log.Fatal("PURE_PASSWORD_RESET_TOKEN_DURATION is not a number")
			}
			passwordResetAuthTokenSettings.Duration = duration
		}
		appSettings.PasswordResetAuthToken = passwordResetAuthTokenSettings
	}

	app := purebackend.NewWithConfig(&purebackend.Config{
		DefaultDebug: os.Getenv("DEBUG") == "true",
		DatabaseType: databaseType,
		DatabaseUrl:  os.Getenv("DATABASE_URL"),
		Settings:     appSettings,
	})
	if err := app.Bootstrap(); err != nil {
		log.Fatal(err)
	}
	// Set admin details in .env file
	//
	// Refer:= https://pureml.mintlify.app/self-hosting#backend-environment
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
