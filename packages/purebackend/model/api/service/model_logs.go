package service

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	authmiddlewares "github.com/PureMLHQ/PureML/packages/purebackend/auth/middlewares"
	"github.com/PureMLHQ/PureML/packages/purebackend/core"
	commonmodels "github.com/PureMLHQ/PureML/packages/purebackend/core/common/models"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/models"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/tools/filesystem"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/tools/inflector"
	"github.com/PureMLHQ/PureML/packages/purebackend/model/middlewares"
	orgmiddlewares "github.com/PureMLHQ/PureML/packages/purebackend/user_org/middlewares"
	"github.com/labstack/echo/v4"
)

// BindModelLogsApi registers the admin api endpoints and the corresponding handlers.
func BindModelLogsApi(app core.App, rg *echo.Group) {
	api := Api{app: app}

	modelGroup := rg.Group("/org/:orgId/model", authmiddlewares.RequireAuthContext, orgmiddlewares.ValidateOrg(api.app))
	modelGroup.GET("/:modelName/branch/:branchName/version/:version/log", api.DefaultHandler(GetAllLogsModel), middlewares.ValidateModel(api.app), middlewares.ValidateModelBranch(api.app), middlewares.ValidateModelBranchVersion(api.app))
	modelGroup.GET("/:modelName/branch/:branchName/version/:version/log/:key", api.DefaultHandler(GetKeyLogsModel), middlewares.ValidateModel(api.app), middlewares.ValidateModelBranch(api.app), middlewares.ValidateModelBranchVersion(api.app))
	modelGroup.POST("/:modelName/branch/:branchName/version/:version/log", api.DefaultHandler(LogModel), middlewares.ValidateModel(api.app), middlewares.ValidateModelBranch(api.app), middlewares.ValidateModelBranchVersion(api.app))
	modelGroup.POST("/:modelName/branch/:branchName/version/:version/logfile", api.DefaultHandler(LogFileModel), middlewares.ValidateModel(api.app), middlewares.ValidateModelBranch(api.app), middlewares.ValidateModelBranchVersion(api.app))
}

// GetAllLogsModel godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Get Log data for model
//	@Description	Get Log data for model
//	@Tags			Model
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/model/{modelName}/branch/{branchName}/version/{version}/log [get]
//	@Param			orgId		path	string	true	"Organization Id"
//	@Param			modelName	path	string	true	"Model Name"
//	@Param			branchName	path	string	true	"Branch Name"
//	@Param			version		path	string	true	"Version"
func (api *Api) GetAllLogsModel(request *models.Request) *models.Response {
	versionUUID := request.GetModelBranchVersionUUID()
	result, err := api.app.Dao().GetLogForModelVersion(versionUUID)
	if err != nil {
		return models.NewServerErrorResponse(err)
	}
	response := models.NewDataResponse(http.StatusOK, result, "Logs for model version")
	return response
}

// GetKeyLogsModel godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Get Log data for model with specific key
//	@Description	Get Log data for model with specific key
//	@Tags			Model
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/model/{modelName}/branch/{branchName}/version/{version}/log/{key} [get]
//	@Param			orgId		path	string	true	"Organization Id"
//	@Param			modelName	path	string	true	"Model Name"
//	@Param			branchName	path	string	true	"Branch Name"
//	@Param			version		path	string	true	"Version"
//	@Param			key			path	string	true	"Key"
func (api *Api) GetKeyLogsModel(request *models.Request) *models.Response {
	versionUUID := request.GetModelBranchVersionUUID()
	key := request.PathParams["key"]
	result, err := api.app.Dao().GetKeyLogForModelVersion(versionUUID, key)
	if err != nil {
		return models.NewServerErrorResponse(err)
	}
	response := models.NewDataResponse(http.StatusOK, result, "Specific Key Logs for model version")
	return response
}

// LogModel godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Log data for model
//	@Description	Log data for model
//	@Tags			Model
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/model/{modelName}/branch/{branchName}/version/{version}/log [post]
//	@Param			orgId		path	string				true	"Organization Id"
//	@Param			modelName	path	string				true	"Model Name"
//	@Param			branchName	path	string				true	"Branch Name"
//	@Param			version		path	string				true	"Version"
//	@Param			data		body	models.LogRequest	true	"Data to log"
func (api *Api) LogModel(request *models.Request) *models.Response {
	request.ParseJsonBody()
	key := request.GetParsedBodyAttribute("key")
	var keyData string
	if key != nil {
		keyData = key.(string)
	} else {
		keyData = ""
	}
	data := request.GetParsedBodyAttribute("data")
	var dataData string
	if data != nil {
		dataData = data.(string)
	} else {
		dataData = ""
	}
	if keyData == "" {
		return models.NewErrorResponse(http.StatusBadRequest, "Key is required")
	}
	versionUUID := request.GetModelBranchVersionUUID()
	result, err := api.app.Dao().CreateLogForModelVersion(keyData, dataData, versionUUID)
	if err != nil {
		return models.NewServerErrorResponse(err)
	}
	response := models.NewDataResponse(http.StatusOK, result, "Log created")
	return response
}

// LogFileModel godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Log file data for model
//	@Description	Log file data for model
//	@Tags			Model
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/model/{modelName}/branch/{branchName}/version/{version}/logfile [post]
//	@Param			file		formData	[]file					true	"Model files (multiple supported)"
//	@Param			orgId		path		string					true	"Organization Id"
//	@Param			modelName	path		string					true	"Model Name"
//	@Param			branchName	path		string					true	"Branch Name"
//	@Param			version		path		string					true	"Version"
//	@Param			data		formData	models.LogFileRequest	true	"Data to log"
func (api *Api) LogFileModel(request *models.Request) *models.Response {
	orgId := request.GetOrgId()
	modelUUID := request.GetModelUUID()
	modelBranchUUID := request.GetModelBranchUUID()
	var modelSourceSecretName string
	if request.FormValues["storage"] != nil && len(request.FormValues["storage"]) > 0 {
		modelSourceSecretName = strings.ToUpper(request.FormValues["storage"][0])
	}
	fileHeaders := request.GetFormMultipleFiles("file")
	if fileHeaders == nil {
		return models.NewErrorResponse(http.StatusBadRequest, "File is required")
	}
	versionUUID := request.GetModelBranchVersionUUID()
	sourceSecrets, errresp := api.ValidateSourceTypeAndGetSourceSecrets(modelSourceSecretName, orgId)
	if errresp != nil {
		return errresp
	}
	modelSourceType := sourceSecrets.SourceType
	sourceValid := false
	for source := range commonmodels.SupportedSources {
		if commonmodels.SupportedSources[source] == modelSourceType {
			sourceValid = true
			break
		}
	}
	if !sourceValid {
		return models.NewErrorResponse(http.StatusBadRequest, "Unsupported model storage")
	}
	logs := make(map[string]string)
	for _, fileHeader := range fileHeaders {
		name := fileHeader.Filename
		originalExt := filepath.Ext(name)
		key := inflector.Snakecase(strings.TrimSuffix(name, originalExt))
		file, err := filesystem.NewFileFromMultipart(fileHeader)
		if err != nil {
			return models.NewServerErrorResponse(err)
		}
		filePath, err := api.app.UploadFile(file, fmt.Sprintf("model-registry/%s/models/%s/%s/logs", orgId, modelUUID, modelBranchUUID), sourceSecrets)
		if err != nil {
			return models.NewServerErrorResponse(err)
		}
		logs[key] = fmt.Sprintf("%s/%s", sourceSecrets.PublicURL, filePath)
	}
	var results []*models.LogResponse
	for key, data := range logs {
		result, err := api.app.Dao().CreateLogForModelVersion(key, data, versionUUID)
		if err != nil {
			return models.NewServerErrorResponse(err)
		}
		results = append(results, result)
	}
	response := models.NewDataResponse(http.StatusOK, results, "Logs created")
	return response
}

var GetAllLogsModel ServiceFunc = (*Api).GetAllLogsModel
var GetKeyLogsModel ServiceFunc = (*Api).GetKeyLogsModel
var LogModel ServiceFunc = (*Api).LogModel
var LogFileModel ServiceFunc = (*Api).LogFileModel
