package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	commonmodels "github.com/PureMLHQ/PureML/packages/purebackend/core/common/models"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/config"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/daos"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/settings"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/tools/filesystem"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/tools/mailer"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/tools/search"
	_ "github.com/joho/godotenv/autoload"
)

func Cleanup(optDataDir ...string) error {
	// remove data dir and all its contents
	var dataDir string
	if len(optDataDir) == 0 || optDataDir[0] == "" {
		// fallback to the default test data directory
		dataDir = config.GetDataDir()
	} else {
		dataDir = optDataDir[0]
	}
	if err := os.RemoveAll(dataDir); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

var _ App = (*BaseApp)(nil)

// BaseApp implements core.App and defines the base PureBackend app structure.
type BaseApp struct {
	// configurable parameters
	isDebug      bool
	dataDir      string
	databaseType string
	databaseUrl  string

	// internals
	settings *settings.Settings
	dao      *daos.Dao
}

// BaseAppConfig defines a BaseApp configuration option
type BaseAppConfig struct {
	DataDir      string
	IsDebug      bool
	DatabaseType string
	DatabaseUrl  string
	Settings     *settings.Settings
}

// NewBaseApp creates and returns a new BaseApp instance
// configured with the provided arguments.
//
// To initialize the app, you need to call `app.Bootstrap()`.
func NewBaseApp(appConfig *BaseAppConfig) *BaseApp {
	app := &BaseApp{
		dataDir:      appConfig.DataDir,
		isDebug:      appConfig.IsDebug,
		settings:     settings.New(),
		databaseType: appConfig.DatabaseType,
		databaseUrl:  appConfig.DatabaseUrl,
	}

	if appConfig.Settings != nil {
		if appConfig.Settings.S3.Enabled {
			app.settings.S3 = appConfig.Settings.S3
		}
		if appConfig.Settings.R2.Enabled {
			app.settings.R2 = appConfig.Settings.R2
		}
		if appConfig.Settings.Search.Enabled {
			app.settings.Search = appConfig.Settings.Search
		}
		if appConfig.Settings.AdminAuthToken.Secret != "" {
			app.settings.AdminAuthToken = appConfig.Settings.AdminAuthToken
		}
		if appConfig.Settings.MailVerifificationAuthToken.Secret != "" {
			app.settings.MailVerifificationAuthToken = appConfig.Settings.MailVerifificationAuthToken
		}
		if appConfig.Settings.PasswordResetAuthToken.Secret != "" {
			app.settings.PasswordResetAuthToken = appConfig.Settings.PasswordResetAuthToken
		}
		if appConfig.Settings.MailService.Enabled {
			app.settings.MailService = appConfig.Settings.MailService
		}
		if appConfig.Settings.Site.BaseURL != "" {
			app.settings.Site = appConfig.Settings.Site
		}
	}

	return app
}

// IsBootstrapped checks if the application was initialized
// (aka. whether Bootstrap() was called).
func (app *BaseApp) IsBootstrapped() bool {
	return app.dao != nil && app.settings != nil
}

// Bootstrap initializes the application
// (aka. create data dir, open db connections, load settings, etc.).
//
// It will call ResetBootstrapState() if the application was already bootstrapped.
func (app *BaseApp) Bootstrap() error {
	// clear resources of previous core state (if any)
	if err := app.ResetBootstrapState(); err != nil {
		return err
	}

	// ensure that data dir exist
	if err := os.MkdirAll(app.DataDir(), os.ModePerm); err != nil {
		return err
	}

	if err := app.initDataDB(); err != nil {
		return err
	}

	if err := app.Dao().Datastore().SeedAdminIfNotExists(); err != nil {
		return err
	}

	if app.settings.Search.Enabled {
		app.Dao().Datastore().SeedSearchClient()
	}

	if app.settings.S3.Enabled {
		if err := app.Dao().Datastore().SeedAdminS3SecretsIfNotExists(&app.settings.S3); err != nil {
			return err
		}
	}

	if app.settings.R2.Enabled {
		if err := app.Dao().Datastore().SeedAdminR2SecretsIfNotExists(&app.settings.R2); err != nil {
			return err
		}
	}

	// app.RefreshSettings()

	return nil
}

// ResetBootstrapState takes care for releasing initialized app resources
// (eg. closing db connections).
func (app *BaseApp) ResetBootstrapState() error {
	if app.Dao() != nil {
		if err := app.Dao().Close(); err != nil {
			return err
		}
		if err := app.Dao().Close(); err != nil {
			return err
		}
	}

	app.dao = nil
	// app.settings = nil

	return nil
}

// Dao returns the default app Dao instance.
func (app *BaseApp) Dao() *daos.Dao {
	return app.dao
}

// DataDir returns the app data directory path.
func (app *BaseApp) DataDir() string {
	return app.dataDir
}

// DatabaseType returns the database type (eg. "sqlite3").
func (app *BaseApp) DatabaseType() string {
	return app.databaseType
}

// DatabaseUrl returns the database connection url.
func (app *BaseApp) DatabaseUrl() string {
	return app.databaseUrl
}

// IsDebug returns whether the app is in debug mode
// (showing more detailed error logs, executed sql statements, etc.).
func (app *BaseApp) IsDebug() bool {
	return app.isDebug
}

// Settings returns the loaded app settings.
func (app *BaseApp) Settings() *settings.Settings {
	return app.settings
}

// NewFilesystem creates a new local or Object storage filesystem instance
// based on the current app settings.
//
// NB! Make sure to call `Close()` on the returned result
// after you are done working with it.
func (app *BaseApp) NewFilesystem(sourceSecrets *commonmodels.SourceSecrets) (*filesystem.System, error) {
	switch strings.ToUpper(sourceSecrets.SourceType) {
	case "S3":
		return filesystem.NewS3(
			app.settings.S3.Bucket,
			app.settings.S3.Region,
			app.settings.S3.Endpoint,
			app.settings.S3.AccessKey,
			app.settings.S3.Secret,
			app.settings.S3.ForcePathStyle,
		)
	case "R2":
		return filesystem.NewR2(
			app.settings.R2.AccountId,
			app.settings.R2.Bucket,
			app.settings.R2.Endpoint,
			app.settings.R2.AccessKey,
			app.settings.R2.Secret,
			app.settings.R2.ForcePathStyle,
		)
	default:
		// fallback to local filesystem
		return filesystem.NewLocal(filepath.Join(app.DataDir(), "storage"))
	}
}

// UploadFile uploads a file to the app storage.
func (app *BaseApp) UploadFile(file *filesystem.File, basePath string, sourceSecrets *commonmodels.SourceSecrets) (string, error) {
	fs, err := app.NewFilesystem(sourceSecrets)
	if err != nil {
		return "", err
	}
	defer fs.Close()

	path := basePath + "/" + file.Name
	if err := fs.UploadFile(file, path); err != nil {
		return "", err
	}
	return path, nil
}

// NewSearchClient creates a new search client instance
func (app *BaseApp) NewSearchClient() *search.SearchClient {
	return search.NewSearchClient(app.settings.Search.Host, app.settings.Search.AdminAPIKey)
}

// RefreshSettings reinitializes and reloads the stored application settings.
func (app *BaseApp) RefreshSettings() error {
	if app.settings == nil {
		app.settings = settings.New()
	}

	// Load S3 settings from db
	// if err := app.settings.LoadFromDB(app.dao, "S3"); err != nil {
	// 	return err
	// }

	return nil
}

// initDB initializes the app database connection.
func (app *BaseApp) initDataDB() error {
	var searchClient *search.SearchClient
	if app.settings.Search.Enabled {
		searchClient = app.NewSearchClient()
	}
	dao, err := daos.InitDB(app.DataDir(), app.DatabaseType(), app.DatabaseUrl(), app.settings.Search.Enabled, searchClient)
	if err != nil {
		return err
	}
	app.dao = dao
	return nil
}

func (app *BaseApp) SendMail(to string, subject string, body string) error {
	return mailer.SendMail(app.Settings(), to, subject, body)
}
