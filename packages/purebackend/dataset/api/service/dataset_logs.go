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
	"github.com/PureMLHQ/PureML/packages/purebackend/dataset/middlewares"
	orgmiddlewares "github.com/PureMLHQ/PureML/packages/purebackend/user_org/middlewares"
	"github.com/labstack/echo/v4"
)

// BindDatasetLogsApi registers the admin api endpoints and the corresponding handlers.
func BindDatasetLogsApi(app core.App, rg *echo.Group) {
	api := Api{app: app}

	datasetGroup := rg.Group("/org/:orgId/dataset", authmiddlewares.RequireAuthContext, orgmiddlewares.ValidateOrg(api.app))
	datasetGroup.GET("/:datasetName/branch/:branchName/version/:version/log", api.DefaultHandler(GetAllLogsDataset), middlewares.ValidateDataset(api.app), middlewares.ValidateDatasetBranch(api.app), middlewares.ValidateDatasetBranchVersion(api.app))
	datasetGroup.GET("/:datasetName/branch/:branchName/version/:version/log/:key", api.DefaultHandler(GetKeyLogsDataset), middlewares.ValidateDataset(api.app), middlewares.ValidateDatasetBranch(api.app), middlewares.ValidateDatasetBranchVersion(api.app))
	datasetGroup.POST("/:datasetName/branch/:branchName/version/:version/log", api.DefaultHandler(LogDataset), middlewares.ValidateDataset(api.app), middlewares.ValidateDatasetBranch(api.app), middlewares.ValidateDatasetBranchVersion(api.app))
	datasetGroup.POST("/:datasetName/branch/:branchName/version/:version/logfile", api.DefaultHandler(LogFileDataset), middlewares.ValidateDataset(api.app), middlewares.ValidateDatasetBranch(api.app), middlewares.ValidateDatasetBranchVersion(api.app))
}

// GetAllLogsDataset godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Get Log data for dataset
//	@Description	Get Log data for dataset
//	@Tags			Dataset
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/dataset/{datasetName}/branch/{branchName}/version/{version}/log [get]
//	@Param			orgId		path	string	true	"Organization Id"
//	@Param			datasetName	path	string	true	"Dataset Name"
//	@Param			branchName	path	string	true	"Branch Name"
//	@Param			version		path	string	true	"Version"
func (api *Api) GetAllLogsDataset(request *models.Request) *models.Response {
	versionUUID := request.GetDatasetBranchVersionUUID()
	result, err := api.app.Dao().GetLogForDatasetVersion(versionUUID)
	if err != nil {
		return models.NewServerErrorResponse(err)
	}
	response := models.NewDataResponse(http.StatusOK, result, "Logs for dataset version")
	return response
}

// GetKeyLogsDataset godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Get Log data for dataset with specific key
//	@Description	Get Log data for dataset with specific key
//	@Tags			Dataset
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/dataset/{datasetName}/branch/{branchName}/version/{version}/log/{key} [get]
//	@Param			orgId		path	string	true	"Organization Id"
//	@Param			datasetName	path	string	true	"Dataset Name"
//	@Param			branchName	path	string	true	"Branch Name"
//	@Param			version		path	string	true	"Version"
//	@Param			key			path	string	true	"Key"
func (api *Api) GetKeyLogsDataset(request *models.Request) *models.Response {
	versionUUID := request.GetDatasetBranchVersionUUID()
	key := request.PathParams["key"]
	result, err := api.app.Dao().GetKeyLogForDatasetVersion(versionUUID, key)
	if err != nil {
		return models.NewServerErrorResponse(err)
	}
	response := models.NewDataResponse(http.StatusOK, result, "Specific Key Logs for dataset version")
	return response
}

// LogDataset godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Log data for dataset
//	@Description	Log data for dataset
//	@Tags			Dataset
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/dataset/{datasetName}/branch/{branchName}/version/{version}/log [post]
//	@Param			orgId		path	string				true	"Organization Id"
//	@Param			datasetName	path	string				true	"Dataset Name"
//	@Param			branchName	path	string				true	"Branch Name"
//	@Param			version		path	string				true	"Version"
//	@Param			data		body	models.LogRequest	true	"Data to log"
func (api *Api) LogDataset(request *models.Request) *models.Response {
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
	versionUUID := request.GetDatasetBranchVersionUUID()
	result, err := api.app.Dao().CreateLogForDatasetVersion(keyData, dataData, versionUUID)
	if err != nil {
		return models.NewServerErrorResponse(err)
	}
	response := models.NewDataResponse(http.StatusOK, result, "Log created")
	return response
}

// LogFileDataset godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Log file data for dataset
//	@Description	Log file data for dataset
//	@Tags			Dataset
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/dataset/{datasetName}/branch/{branchName}/version/{version}/logfile [post]
//	@Param			file		formData	[]file					true	"Dataset files (multiple supported)"
//	@Param			orgId		path		string					true	"Organization Id"
//	@Param			datasetName	path		string					true	"Dataset Name"
//	@Param			branchName	path		string					true	"Branch Name"
//	@Param			version		path		string					true	"Version"
//	@Param			data		formData	models.LogFileRequest	true	"Data to log"
func (api *Api) LogFileDataset(request *models.Request) *models.Response {
	orgId := request.GetOrgId()
	datasetUUID := request.GetDatasetUUID()
	datasetBranchUUID := request.GetDatasetBranchUUID()
	var datasetSourceSecretName string
	if request.FormValues["storage"] != nil && len(request.FormValues["storage"]) > 0 {
		datasetSourceSecretName = strings.ToUpper(request.FormValues["storage"][0])
	}
	fileHeaders := request.GetFormMultipleFiles("file")
	if fileHeaders == nil {
		return models.NewErrorResponse(http.StatusBadRequest, "File is required")
	}
	versionUUID := request.GetDatasetBranchVersionUUID()
	sourceSecrets, errresp := api.ValidateSourceTypeAndGetSourceSecrets(datasetSourceSecretName, orgId)
	if errresp != nil {
		return errresp
	}
	datasetSourceType := sourceSecrets.SourceType
	sourceValid := false
	for source := range commonmodels.SupportedSources {
		if commonmodels.SupportedSources[source] == datasetSourceType {
			sourceValid = true
			break
		}
	}
	if !sourceValid {
		return models.NewErrorResponse(http.StatusBadRequest, "Unsupported dataset storage")
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
		filePath, err := api.app.UploadFile(file, fmt.Sprintf("dataset-registry/%s/datasets/%s/%s/logs", orgId, datasetUUID, datasetBranchUUID), sourceSecrets)
		if err != nil {
			return models.NewServerErrorResponse(err)
		}
		logs[key] = fmt.Sprintf("%s/%s", sourceSecrets.PublicURL, filePath)
	}
	var results []*models.LogResponse
	for key, data := range logs {
		result, err := api.app.Dao().CreateLogForDatasetVersion(key, data, versionUUID)
		if err != nil {
			return models.NewServerErrorResponse(err)
		}
		results = append(results, result)
	}
	response := models.NewDataResponse(http.StatusOK, results, "Logs created")
	return response
}

var GetAllLogsDataset ServiceFunc = (*Api).GetAllLogsDataset
var GetKeyLogsDataset ServiceFunc = (*Api).GetKeyLogsDataset
var LogDataset ServiceFunc = (*Api).LogDataset
var LogFileDataset ServiceFunc = (*Api).LogFileDataset
