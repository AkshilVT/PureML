package service

import (
	"fmt"
	"net/http"
	"strings"

	authmiddlewares "github.com/PureMLHQ/PureML/packages/purebackend/auth/middlewares"
	"github.com/PureMLHQ/PureML/packages/purebackend/core"
	commonmodels "github.com/PureMLHQ/PureML/packages/purebackend/core/common/models"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/models"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/tools/filesystem"
	"github.com/PureMLHQ/PureML/packages/purebackend/dataset/middlewares"
	orgmiddlewares "github.com/PureMLHQ/PureML/packages/purebackend/user_org/middlewares"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

// BindDatasetBranchVersionApi registers the admin api endpoints and the corresponding handlers.
func BindDatasetBranchVersionApi(app core.App, rg *echo.Group) {
	api := Api{app: app}

	datasetGroup := rg.Group("/org/:orgId/dataset", authmiddlewares.RequireAuthContext, orgmiddlewares.ValidateOrg(api.app))
	datasetGroup.POST("/:datasetName/branch/:branchName/hash-status", api.DefaultHandler(VerifyDatasetBranchHashStatus), middlewares.ValidateDataset(api.app))
	datasetGroup.POST("/:datasetName/branch/:branchName/register", api.DefaultHandler(RegisterDataset), middlewares.ValidateDataset(api.app), middlewares.ValidateDatasetBranch(api.app))
	datasetGroup.GET("/:datasetName/branch/:branchName/version", api.DefaultHandler(GetDatasetBranchAllVersions), middlewares.ValidateDataset(api.app), middlewares.ValidateDatasetBranch(api.app))
	datasetGroup.GET("/:datasetName/branch/:branchName/version/:version", api.DefaultHandler(GetDatasetBranchVersion), middlewares.ValidateDataset(api.app), middlewares.ValidateDatasetBranch(api.app), middlewares.ValidateDatasetBranchVersion(api.app))
}

// GetDatasetBranchAllVersions godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Get all branch versions of a dataset
//	@Description	Get all branch versions of a dataset
//	@Tags			Dataset
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/dataset/{datasetName}/branch/{branchName}/version [get]
//	@Param			orgId		path	string	true	"Organization Id"
//	@Param			datasetName	path	string	true	"Dataset Name"
//	@Param			branchName	path	string	true	"Branch Name"
func (api *Api) GetDatasetBranchAllVersions(request *models.Request) *models.Response {
	var response *models.Response
	branchUUID := request.GetDatasetBranchUUID()
	allVersions, err := api.app.Dao().GetDatasetBranchAllVersions(branchUUID)
	if err != nil {
		return models.NewServerErrorResponse(err)
	} else {
		response = models.NewDataResponse(http.StatusOK, allVersions, "All dataset branch versions")
	}
	return response
}

// GetDatasetBranchVersion godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Get specific branch version of a dataset
//	@Description	Get specific branch version of a dataset
//	@Tags			Dataset
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/dataset/{datasetName}/branch/{branchName}/version/{version} [get]
//	@Param			orgId		path	string	true	"Organization Id"
//	@Param			datasetName	path	string	true	"Dataset Name"
//	@Param			branchName	path	string	true	"Branch Name"
//	@Param			version		path	string	true	"Version"
func (api *Api) GetDatasetBranchVersion(request *models.Request) *models.Response {
	branchUUID := request.GetDatasetBranchUUID()
	versionName := request.PathParams["version"]
	version, err := api.app.Dao().GetDatasetBranchVersion(branchUUID, versionName)
	if err != nil {
		return models.NewErrorResponse(http.StatusInternalServerError, err.Error())
	}
	if version == nil {
		return models.NewErrorResponse(http.StatusNotFound, "Version not found")
	}
	return models.NewDataResponse(http.StatusOK, version, "Dataset branch version details")
}

// VerifyDatasetBranchHashStatus godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Verify dataset hash status
//	@Description	Verify dataset hash status to determine if dataset is already uploaded
//	@Tags			Dataset
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/dataset/{datasetName}/branch/{branchName}/hash-status [post]
//	@Param			orgId		path	string				true	"Organization Id"
//	@Param			datasetName	path	string				true	"Dataset Name"
//	@Param			branchName	path	string				true	"Branch Name"
//	@Param			hash		body	models.HashRequest	true	"Hash value"
func (api *Api) VerifyDatasetBranchHashStatus(request *models.Request) *models.Response {
	datasetName := request.GetDatasetName()
	datasetBranchName := request.GetPathParam("branchName")
	orgId := uuid.Must(uuid.FromString(request.GetPathParam("orgId")))
	message := "Hash validity (False - does not exist in db)"
	dataset, err := api.app.Dao().GetDatasetBranchByName(orgId, datasetName, datasetBranchName)
	if err != nil {
		return models.NewServerErrorResponse(err)
	}
	if dataset == nil {
		return models.NewDataResponse(http.StatusOK, false, message)
	}
	datasetBranchUUID := dataset.UUID
	request.ParseJsonBody()
	hashValue := request.GetParsedBodyAttribute("hash")
	var hashValueData string
	if hashValue == nil {
		hashValueData = ""
	} else {
		hashValueData = hashValue.(string)
	}
	if hashValueData == "" {
		return models.NewErrorResponse(http.StatusBadRequest, "Hash value is empty")
	}
	versions, err := api.app.Dao().GetDatasetBranchAllVersions(datasetBranchUUID)
	if err != nil {
		return models.NewServerErrorResponse(err)
	}
	response := false
	for _, version := range versions {
		if version.Hash == hashValue {
			response = true
			message = "Hash validity (True - exists in db)"
			break
		}
	}
	return models.NewDataResponse(http.StatusOK, response, message)
}

// RegisterDataset godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Register dataset
//	@Description	Register dataset file. Create dataset and default branches if not exists
//	@Tags			Dataset
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/dataset/{datasetName}/branch/{branchName}/register [post]
//	@Param			file		formData	file							true	"Dataset file"
//	@Param			orgId		path		string							true	"Organization UUID"
//	@Param			datasetName	path		string							true	"Dataset name"
//	@Param			branchName	path		string							true	"Branch name"
//	@Param			data		formData	models.RegisterDatasetRequest	true	"Dataset details"
func (api *Api) RegisterDataset(request *models.Request) *models.Response {
	orgId := request.GetOrgId()
	userUUID := request.GetUserUUID()
	datasetUUID := request.GetDatasetUUID()
	datasetBranchUUID := request.GetDatasetBranchUUID()
	var datasetHash string
	if request.FormValues["hash"] != nil && len(request.FormValues["hash"]) > 0 && request.FormValues["hash"][0] != "" {
		datasetHash = request.FormValues["hash"][0]
	} else {
		return models.NewErrorResponse(http.StatusBadRequest, "Hash is required")
	}
	var datasetSourceSecretName string
	if request.FormValues["storage"] != nil && len(request.FormValues["storage"]) > 0 {
		datasetSourceSecretName = request.FormValues["storage"][0]
	}
	var datasetIsEmpty bool
	if request.FormValues["is_empty"] != nil && len(request.FormValues["is_empty"]) > 0 {
		datasetIsEmpty = request.FormValues["is_empty"][0] == "true"
	}
	var datasetLineage string
	if request.FormValues["lineage"] != nil && len(request.FormValues["lineage"]) > 0 {
		datasetLineage = request.FormValues["lineage"][0]
	}
	fileHeader := request.GetFormFile("file")
	if fileHeader == nil {
		return models.NewErrorResponse(http.StatusBadRequest, "File is required")
	}
	datasetBranchName := request.GetPathParam("branchName")
	if datasetBranchName == "main" {
		return models.NewErrorResponse(http.StatusBadRequest, "Cannot register dataset directly to main branch")
	}
	versions, err := api.app.Dao().GetDatasetBranchAllVersions(datasetBranchUUID)
	if err != nil {
		return models.NewServerErrorResponse(err)
	}
	response := false
	for _, version := range versions {
		if version.Hash == datasetHash {
			response = true
			break
		}
	}
	if response {
		return models.NewErrorResponse(http.StatusBadRequest, "Dataset with this hash already exists")
	}
	var datasetSourceType string
	var datasetSourceSecrets *commonmodels.SourceSecrets
	var errresp *models.Response
	if strings.ToUpper(datasetSourceSecretName) != "LOCAL" {
		datasetSourceSecrets, errresp = api.ValidateSourceTypeAndGetSourceSecrets(datasetSourceSecretName, orgId)
		if errresp != nil {
			return errresp
		}
		datasetSourceType = datasetSourceSecrets.SourceType
	} else {
		datasetSourceType = "LOCAL"
		datasetSourceSecrets.SourceType = "LOCAL"
	}
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
	var filePath string
	if !datasetIsEmpty {
		file, err := filesystem.NewFileFromMultipart(fileHeader)
		if err != nil {
			return models.NewServerErrorResponse(err)
		}
		filePath, err = api.app.UploadFile(file, fmt.Sprintf("dataset-registry/%s/datasets/%s/%s", orgId, datasetUUID, datasetBranchUUID), datasetSourceSecrets)
		if err != nil {
			return models.NewServerErrorResponse(err)
		}
	}
	datasetVersion, err := api.app.Dao().RegisterDatasetFile(datasetBranchUUID, datasetSourceType, datasetSourceSecrets.PublicURL, filePath, datasetIsEmpty, datasetHash, datasetLineage, userUUID)
	if err != nil {
		return models.NewServerErrorResponse(err)
	}
	return models.NewDataResponse(http.StatusOK, datasetVersion, "Dataset successfully registered")
}

var GetDatasetBranchAllVersions ServiceFunc = (*Api).GetDatasetBranchAllVersions
var GetDatasetBranchVersion ServiceFunc = (*Api).GetDatasetBranchVersion
var VerifyDatasetBranchHashStatus ServiceFunc = (*Api).VerifyDatasetBranchHashStatus
var RegisterDataset ServiceFunc = (*Api).RegisterDataset
