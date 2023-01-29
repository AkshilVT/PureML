package apis

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	_ "github.com/PureML-Inc/PureML/server/docs"
	"github.com/PureML-Inc/PureML/server/handler"
	"github.com/PureML-Inc/PureML/server/middlewares"
	"github.com/PureML-Inc/PureML/server/apis/service"
	"github.com/PureML-Inc/PureML/server/ui"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const trailedPath = "/_/"

// @title PureML API Documentation
// @version 1.0

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email contact@pureml.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

/////////////////////////// CHANGE BELOW HOST & SCHEME FOR SWAGGER DOCUMENTATION ///////////////////////////
// @host pureml-development.up.railway.app
// @BasePath /
// @schemes https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func InitApi() (*echo.Echo, error) {
	// Echo instance
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	bindStaticUI(e)

	api := e.Group("/api")

	//Health API
	api.GET("/health", handler.DefaultHandler(service.HealthCheck))

	//Swagger API
	api.GET("/swagger/*", SwaggerHandler)

	//Org APIs
	api.GET("/org/handle/:orgHandle", handler.DefaultHandler(service.GetOrgByHandle))

	orgGroup := api.Group("/org", middlewares.AuthenticateJWT)
	orgGroup.GET("/all", handler.DefaultHandler(service.GetAllAdminOrgs))
	orgGroup.GET("/id/:orgId", handler.DefaultHandler(service.GetOrgByID), middlewares.ValidateOrg)
	orgGroup.GET("/", handler.DefaultHandler(service.GetOrgsForUser))
	orgGroup.POST("/create", handler.DefaultHandler(service.CreateOrg))
	orgGroup.POST("/:orgId/update", handler.DefaultHandler(service.UpdateOrg), middlewares.ValidateOrg)
	orgGroup.POST("/:orgId/add", handler.DefaultHandler(service.AddUsersToOrg), middlewares.ValidateOrg)
	orgGroup.POST("/join", handler.DefaultHandler(service.JoinOrg))
	orgGroup.POST("/:orgId/remove", handler.DefaultHandler(service.RemoveOrg), middlewares.ValidateOrg)
	orgGroup.POST("/:orgId/leave", handler.DefaultHandler(service.LeaveOrg), middlewares.ValidateOrg)

	//User APIs
	userGroup := api.Group("user")
	userGroup.GET("/profile", handler.DefaultHandler(service.GetProfile), middlewares.AuthenticateJWT)
	userGroup.GET("/profile/:userHandle", handler.DefaultHandler(service.GetProfileByHandle))
	userGroup.POST("/profile", handler.DefaultHandler(service.UpdateProfile), middlewares.AuthenticateJWT)
	userGroup.POST("/signup", handler.DefaultHandler(service.UserSignUp))
	userGroup.POST("/login", handler.DefaultHandler(service.UserLogin))
	userGroup.POST("/forgot-password", handler.DefaultHandler(service.UserForgotPassword))
	userGroup.POST("/reset-password", handler.DefaultHandler(service.UserResetPassword)) //TODO To complete the logic here and update middlewares

	//Model APIs
	modelGroup := api.Group("/org/:orgId/model", middlewares.AuthenticateJWT, middlewares.ValidateOrg)
	modelGroup.GET("/all", handler.DefaultHandler(service.GetAllModels))
	modelGroup.GET("/:modelName", handler.DefaultHandler(service.GetModel), middlewares.ValidateModel)
	modelGroup.POST("/:modelName/create", handler.DefaultHandler(service.CreateModel))
	modelGroup.GET("/:modelName/branch", handler.DefaultHandler(service.GetModelAllBranches), middlewares.ValidateModel)
	modelGroup.GET("/:modelName/readme/version/:version", handler.DefaultHandler(service.GetModelReadmeVersion), middlewares.ValidateModel)
	modelGroup.GET("/:modelName/readme/version", handler.DefaultHandler(service.GetModelReadmeAllVersions), middlewares.ValidateModel)
	modelGroup.POST("/:modelName/readme", handler.DefaultHandler(service.UpdateModelReadme), middlewares.ValidateModel)
	modelGroup.POST("/:modelName/branch/create", handler.DefaultHandler(service.CreateModelBranch), middlewares.ValidateModel)
	modelGroup.GET("/:modelName/branch/:branchName", handler.DefaultHandler(service.GetModelBranch), middlewares.ValidateModel, middlewares.ValidateModelBranch)
	modelGroup.POST("/:modelName/branch/:branchName/update", handler.DefaultHandler(service.UpdateModelBranch), middlewares.ValidateModel, middlewares.ValidateModelBranch)
	modelGroup.POST("/:modelName/branch/:branchName/hash-status", handler.DefaultHandler(service.VerifyModelBranchHashStatus), middlewares.ValidateModel)
	modelGroup.POST("/:modelName/branch/:branchName/register", handler.DefaultHandler(service.RegisterModel), middlewares.ValidateModel, middlewares.ValidateModelBranch)
	modelGroup.DELETE("/:modelName/branch/:branchName/delete", handler.DefaultHandler(service.DeleteModelBranch), middlewares.ValidateModel, middlewares.ValidateModelBranch)
	modelGroup.GET("/:modelName/branch/:branchName/version", handler.DefaultHandler(service.GetModelBranchAllVersions), middlewares.ValidateModel, middlewares.ValidateModelBranch)
	modelGroup.GET("/:modelName/branch/:branchName/version/:version", handler.DefaultHandler(service.GetModelBranchVersion), middlewares.ValidateModel, middlewares.ValidateModelBranch)

	//Model Log APIs
	modelGroup.GET("/:modelName/branch/:branchName/version/:version/log", handler.DefaultHandler(service.GetLogModel), middlewares.ValidateModel, middlewares.ValidateModelBranch)
	modelGroup.POST("/:modelName/branch/:branchName/version/:version/log", handler.DefaultHandler(service.LogModel), middlewares.ValidateModel, middlewares.ValidateModelBranch)

	//Model Activity APIs
	modelGroup.GET("/:modelName/activity/:category", handler.DefaultHandler(service.GetModelActivity), middlewares.ValidateModel)
	modelGroup.POST("/:modelName/activity/:category", handler.DefaultHandler(service.CreateModelActivity), middlewares.ValidateModel)
	modelGroup.POST("/:modelName/activity/:category/:activityUUID", handler.DefaultHandler(service.UpdateModelActivity), middlewares.ValidateModel)
	modelGroup.DELETE("/:modelName/activity/:category/:activityUUID/delete", handler.DefaultHandler(service.DeleteModelActivity), middlewares.ValidateModel)

	//Dataset APIs
	datasetGroup := api.Group("/org/:orgId/dataset", middlewares.AuthenticateJWT, middlewares.ValidateOrg)
	datasetGroup.GET("/all", handler.DefaultHandler(service.GetAllDatasets))
	datasetGroup.GET("/:datasetName", handler.DefaultHandler(service.GetDataset), middlewares.ValidateDataset)
	datasetGroup.POST("/:datasetName/create", handler.DefaultHandler(service.CreateDataset))
	datasetGroup.GET("/:datasetName/branch", handler.DefaultHandler(service.GetDatasetAllBranches), middlewares.ValidateDataset)
	datasetGroup.GET("/:datasetName/readme/version/:version", handler.DefaultHandler(service.GetDatasetReadmeVersion), middlewares.ValidateDataset)
	datasetGroup.GET("/:datasetName/readme/version", handler.DefaultHandler(service.GetDatasetReadmeAllVersions), middlewares.ValidateDataset)
	datasetGroup.POST("/:datasetName/readme", handler.DefaultHandler(service.UpdateDatasetReadme), middlewares.ValidateDataset)
	datasetGroup.POST("/:datasetName/branch/create", handler.DefaultHandler(service.CreateDatasetBranch), middlewares.ValidateDataset)
	datasetGroup.GET("/:datasetName/branch/:branchName", handler.DefaultHandler(service.GetDatasetBranch), middlewares.ValidateDataset, middlewares.ValidateDatasetBranch)
	datasetGroup.POST("/:datasetName/branch/:branchName/update", handler.DefaultHandler(service.UpdateDatasetBranch), middlewares.ValidateDataset, middlewares.ValidateDatasetBranch)
	datasetGroup.DELETE("/:datasetName/branch/:branchName/delete", handler.DefaultHandler(service.DeleteDatasetBranch), middlewares.ValidateDataset, middlewares.ValidateDatasetBranch)
	datasetGroup.POST("/:datasetName/branch/:branchName/hash-status", handler.DefaultHandler(service.VerifyDatasetBranchHashStatus), middlewares.ValidateDataset)
	datasetGroup.POST("/:datasetName/branch/:branchName/register", handler.DefaultHandler(service.RegisterDataset), middlewares.ValidateDataset, middlewares.ValidateDatasetBranch)
	datasetGroup.GET("/:datasetName/branch/:branchName/version", handler.DefaultHandler(service.GetDatasetBranchAllVersions), middlewares.ValidateDataset, middlewares.ValidateDatasetBranch)
	datasetGroup.GET("/:datasetName/branch/:branchName/version/:version", handler.DefaultHandler(service.GetDatasetBranchVersion), middlewares.ValidateDataset, middlewares.ValidateDatasetBranch)

	//Dataset Log APIs
	datasetGroup.GET("/:datasetName/branch/:branchName/version/:version/log", handler.DefaultHandler(service.GetLogDataset), middlewares.ValidateDataset, middlewares.ValidateDatasetBranch)
	datasetGroup.POST("/:datasetName/branch/:branchName/version/:version/log", handler.DefaultHandler(service.LogDataset), middlewares.ValidateDataset, middlewares.ValidateDatasetBranch)

	//Dataset Activity APIs
	datasetGroup.GET("/:datasetName/activity/:category", handler.DefaultHandler(service.GetDatasetActivity), middlewares.ValidateDataset)
	datasetGroup.POST("/:datasetName/activity/:category", handler.DefaultHandler(service.CreateDatasetActivity), middlewares.ValidateDataset)
	datasetGroup.POST("/:datasetName/activity/:category/:activityUUID", handler.DefaultHandler(service.UpdateDatasetActivity), middlewares.ValidateDataset)
	datasetGroup.DELETE("/:datasetName/activity/:category/:activityUUID/delete", handler.DefaultHandler(service.DeleteDatasetActivity), middlewares.ValidateDataset)

	//Secret APIs
	secretGroup := api.Group("/org/:orgId/secret", middlewares.AuthenticateJWT, middlewares.ValidateOrg)
	// secretGroup.GET("/all", handler.DefaultHandler(service.GetAllSecrets))
	secretGroup.GET("/r2", handler.DefaultHandler(service.GetR2Secret))
	secretGroup.POST("/r2/connect", handler.DefaultHandler(service.ConnectR2Secret))
	secretGroup.DELETE("/r2/delete", handler.DefaultHandler(service.DeleteR2Secrets))
	secretGroup.GET("/s3", handler.DefaultHandler(service.GetS3Secret))
	secretGroup.POST("/s3/connect", handler.DefaultHandler(service.ConnectS3Secret))
	secretGroup.DELETE("/s3/delete", handler.DefaultHandler(service.DeleteS3Secrets))

	return e, nil
}

// StaticDirectoryHandler is similar to `echo.StaticDirectoryHandler`
// but without the directory redirect which conflicts with RemoveTrailingSlash middleware.
//
// If a file resource is missing and indexFallback is set, the request
// will be forwarded to the base index.html (useful also for SPA).
//
// @see https://github.com/labstack/echo/issues/2211
func StaticDirectoryHandler(fileSystem fs.FS, number int, indexFallback bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		p := c.PathParam("*")

		// escape url path
		tmpPath, err := url.PathUnescape(p)
		if err != nil {
			return fmt.Errorf("failed to unescape path variable: %w", err)
		}
		p = tmpPath

		// fs.FS.Open() already assumes that file names are relative to FS root path and considers name with prefix `/` as invalid
		name := filepath.ToSlash(filepath.Clean(strings.TrimPrefix(p, "/")))

		fileErr := c.FileFS(name, fileSystem)

		if fileErr != nil && indexFallback && errors.Is(fileErr, echo.ErrNotFound) {
			return c.FileFS("index.html", fileSystem)
		}

		return fileErr
	}
}

// bindStaticUI registers the endpoints that serves the static  UI.
func bindStaticUI(e *echo.Echo) error {
	// redirect to trailing slash to ensure that relative urls will still work properly
	e.GET(
		strings.TrimRight(trailedPath, "/"),
		func(c echo.Context) error {
			return c.Redirect(http.StatusTemporaryRedirect, trailedPath)
		},
	)

	// e.FileFS(trailedPath+"*", "index.html", ui.BuildIndexFS)
	// serves static files from the /ui/public/build directory
	// (similar to echo.StaticFS but with gzip middleware enabled)
	e.GET(
		trailedPath+"build/*",
		StaticDirectoryHandler(ui.BuildDirFS, 0, false),
		middleware.Gzip(),
	)
	e.GET(
		trailedPath+"*",
		StaticDirectoryHandler(ui.BuildIndexFS, 1, false),
		middleware.Gzip(),
	)

	return nil
}
