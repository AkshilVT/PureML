package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	uuid "github.com/satori/go.uuid"
)

func init() {
	// fmt.Println("Loading environment variables from .env file")
	_ = godotenv.Load("../.env")
}

var adminAccess = map[string]bool{
	"priyavkkaneria@gmail.com": true,
	"kaneriakesha7@gmail.com":  true,
	"akshilvthumar@gmail.com":  true,
	"test.pureml@gmail.com":    true,
	"demo@aztlan.in":           true,
}

// Development: localhost:8080
// Production: host_url
func GetHost() string {
	hostURL := os.Getenv("HOST_URL")
	port := GetPort()
	if hostURL == "" {
		hostURL = fmt.Sprintf("localhost:%s", port)
	}
	return hostURL
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// Development: http
// Production: https
func GetScheme() string {
	scheme := os.Getenv("SCHEME")
	host := GetHost()
	if scheme == "" && (strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1")) {
		scheme = "http"
	} else if scheme == "" {
		scheme = "https"
	}
	return scheme
}

func GetPureMLR2Secrets() map[string]string {
	return map[string]string{
		"R2_ACCOUNT_ID":        os.Getenv("R2_ACCOUNT_ID"),
		"R2_ACCESS_KEY_ID":     os.Getenv("R2_ACCESS_KEY_ID"),
		"R2_ACCESS_KEY_SECRET": os.Getenv("R2_ACCESS_KEY_SECRET"),
		"R2_BUCKET_NAME":       os.Getenv("R2_BUCKET_NAME"),
		"R2_PUBLIC_URL":        os.Getenv("R2_PUBLIC_URL"),
	}
}

func IsCGOEnabled() bool {
	return os.Getenv("CGO_ENABLED") == "1"
}

func GetDatabaseType() string {
	databaseType := os.Getenv("DATABASE")
	if databaseType == "" {
		databaseType = "sqlite3"
	}
	return databaseType
}

func GetDatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}

func GetDataDir() string {
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./pureml_data"
	}
	return dataDir
}

func GetHttpAddr() string {
	httpAddr := os.Getenv("HTTP_ADDRESS")
	if httpAddr == "" {
		httpAddr = "0.0.0.0:8080"
	}
	return httpAddr
}

func GetHttpsAddr() string {
	return os.Getenv("HTTPS_ADDRESS")
}

func GetAdminDetails() map[string]interface{} {
	adminUUID := os.Getenv("ADMIN_UUID")
	if adminUUID == "" {
		return nil
	}
	adminDetails := map[string]interface{}{
		"uuid":       uuid.Must(uuid.FromString(adminUUID)),
		"email":      os.Getenv("ADMIN_EMAIL"),
		"password":   os.Getenv("ADMIN_PASSWORD"),
		"handle":     os.Getenv("ADMIN_HANDLE"),
		"org_name":   os.Getenv("ADMIN_ORG_NAME"),
		"org_handle": os.Getenv("ADMIN_ORG_HANDLE"),
	}
	return adminDetails
}

func HasAdminAccess(email string) bool {
	_, ok := adminAccess[email]
	return ok
}