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
	"github.com/PureMLHQ/PureML/packages/purebackend/model/middlewares"
	orgmiddlewares "github.com/PureMLHQ/PureML/packages/purebackend/user_org/middlewares"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

// BindModelBranchVersionApi registers the admin api endpoints and the corresponding handlers.
func BindModelBranchVersionApi(app core.App, rg *echo.Group) {
	api := Api{app: app}

	modelGroup := rg.Group("/org/:orgId/model", authmiddlewares.RequireAuthContext, orgmiddlewares.ValidateOrg(api.app))
	modelGroup.POST("/:modelName/branch/:branchName/hash-status", api.DefaultHandler(VerifyModelBranchHashStatus), middlewares.ValidateModel(api.app))
	modelGroup.POST("/:modelName/branch/:branchName/register", api.DefaultHandler(RegisterModel), middlewares.ValidateModel(api.app), middlewares.ValidateModelBranch(api.app))
	modelGroup.GET("/:modelName/branch/:branchName/version", api.DefaultHandler(GetModelBranchAllVersions), middlewares.ValidateModel(api.app), middlewares.ValidateModelBranch(api.app))
	modelGroup.GET("/:modelName/branch/:branchName/version/:version", api.DefaultHandler(GetModelBranchVersion), middlewares.ValidateModel(api.app), middlewares.ValidateModelBranch(api.app), middlewares.ValidateModelBranchVersion(api.app))
}

// GetModelBranchAllVersions godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Get all branch versions of a model
//	@Description	Get all branch versions of a model
//	@Tags			Model
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/model/{modelName}/branch/{branchName}/version [get]
//	@Param			orgId		path	string	true	"Organization Id"
//	@Param			modelName	path	string	true	"Model Name"
//	@Param			branchName	path	string	true	"Branch Name"
func (api *Api) GetModelBranchAllVersions(request *models.Request) *models.Response {
	var response *models.Response
	branchUUID := request.GetModelBranchUUID()
	withLogs := strings.ToLower(request.GetQueryParam("withLogs")) == "true"
	allVersions, err := api.app.Dao().GetModelBranchAllVersions(branchUUID, withLogs)
	if err != nil {
		return models.NewServerErrorResponse(err)
	} else {
		response = models.NewDataResponse(http.StatusOK, allVersions, "All model branch versions")
	}
	return response
}

// GetModelBranchVersion godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Get specific branch version of a model
//	@Description	Get specific branch version of a model
//	@Tags			Model
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/model/{modelName}/branch/{branchName}/version/{version} [get]
//	@Param			orgId		path	string	true	"Organization Id"
//	@Param			modelName	path	string	true	"Model Name"
//	@Param			branchName	path	string	true	"Branch Name"
//	@Param			version		path	string	true	"Version"
func (api *Api) GetModelBranchVersion(request *models.Request) *models.Response {
	branchUUID := request.GetModelBranchUUID()
	versionName := request.PathParams["version"]
	version, err := api.app.Dao().GetModelBranchVersion(branchUUID, versionName)
	if err != nil {
		return models.NewErrorResponse(http.StatusInternalServerError, err.Error())
	}
	if version == nil {
		return models.NewErrorResponse(http.StatusNotFound, "Version not found")
	}
	return models.NewDataResponse(http.StatusOK, version, "Model branch version details")
}

// VerifyModelBranchHashStatus godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Verify model hash status
//	@Description	Verify model hash status to determine if model is already uploaded
//	@Tags			Model
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/model/{modelName}/branch/{branchName}/hash-status [post]
//	@Param			orgId		path	string				true	"Organization Id"
//	@Param			modelName	path	string				true	"Model Name"
//	@Param			branchName	path	string				true	"Branch Name"
//	@Param			hash		body	models.HashRequest	true	"Hash value"
func (api *Api) VerifyModelBranchHashStatus(request *models.Request) *models.Response {
	modelName := request.GetModelName()
	modelBranchName := request.GetPathParam("branchName")
	orgId := uuid.Must(uuid.FromString(request.GetPathParam("orgId")))
	message := "Hash validity (False - does not exist in db)"
	model, err := api.app.Dao().GetModelBranchByName(orgId, modelName, modelBranchName)
	if err != nil {
		return models.NewServerErrorResponse(err)
	}
	if model == nil {
		return models.NewDataResponse(http.StatusOK, false, message)
	}
	modelBranchUUID := model.UUID
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
	versions, err := api.app.Dao().GetModelBranchAllVersions(modelBranchUUID, false)
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

// RegisterModel godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Register model
//	@Description	Register model file. Create model and default branches if not exists
//	@Tags			Model
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/org/{orgId}/model/{modelName}/branch/{branchName}/register [post]
//	@Param			file		formData	file						true	"Model file"
//	@Param			orgId		path		string						true	"Organization UUID"
//	@Param			modelName	path		string						true	"Model name"
//	@Param			branchName	path		string						true	"Branch name"
//	@Param			data		formData	models.RegisterModelRequest	true	"Model details"
func (api *Api) RegisterModel(request *models.Request) *models.Response {
	orgId := request.GetOrgId()
	userUUID := request.GetUserUUID()
	modelUUID := request.GetModelUUID()
	modelBranchUUID := request.GetModelBranchUUID()
	var modelHash string
	if request.FormValues["hash"] != nil && len(request.FormValues["hash"]) > 0 && request.FormValues["hash"][0] != "" {
		modelHash = request.FormValues["hash"][0]
	} else {
		return models.NewErrorResponse(http.StatusBadRequest, "Hash is required")
	}
	var modelSourceSecretName string
	if request.FormValues["storage"] != nil && len(request.FormValues["storage"]) > 0 {
		modelSourceSecretName = request.FormValues["storage"][0]
	}
	var modelIsEmpty bool
	if request.FormValues["is_empty"] != nil && len(request.FormValues["is_empty"]) > 0 {
		modelIsEmpty = request.FormValues["is_empty"][0] == "true"
	}
	fileHeader := request.GetFormFile("file")
	if fileHeader == nil {
		return models.NewErrorResponse(http.StatusBadRequest, "File is required")
	}
	modelBranchName := request.GetPathParam("branchName")
	if modelBranchName == "main" {
		return models.NewErrorResponse(http.StatusBadRequest, "Cannot register model directly to main branch")
	}
	versions, err := api.app.Dao().GetModelBranchAllVersions(modelBranchUUID, true)
	if err != nil {
		return models.NewServerErrorResponse(err)
	}
	response := false
	for _, version := range versions {
		if version.Hash == modelHash {
			response = true
			break
		}
	}
	if response {
		return models.NewErrorResponse(http.StatusBadRequest, "Model with this hash already exists")
	}
	var modelSourceType string
	var modelSourceSecrets *commonmodels.SourceSecrets
	var errresp *models.Response
	if strings.ToUpper(modelSourceSecretName) != "LOCAL" {
		modelSourceSecrets, errresp = api.ValidateSourceTypeAndGetSourceSecrets(modelSourceSecretName, orgId)
		if errresp != nil {
			return errresp
		}
		modelSourceType = modelSourceSecrets.SourceType
	} else {
		modelSourceType = "LOCAL"
		modelSourceSecrets.SourceType = "LOCAL"
	}
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
	var filePath string
	if !modelIsEmpty {
		file, err := filesystem.NewFileFromMultipart(fileHeader)
		if err != nil {
			return models.NewServerErrorResponse(err)
		}
		filePath, err = api.app.UploadFile(file, fmt.Sprintf("model-registry/%s/models/%s/%s", orgId, modelUUID, modelBranchUUID), modelSourceSecrets)
		if err != nil {
			return models.NewServerErrorResponse(err)
		}
	}
	modelVersion, err := api.app.Dao().RegisterModelFile(modelBranchUUID, modelSourceType, modelSourceSecrets.PublicURL, filePath, modelIsEmpty, modelHash, userUUID)
	if err != nil {
		return models.NewServerErrorResponse(err)
	}
	return models.NewDataResponse(http.StatusOK, modelVersion, "Model successfully registered")
}

var GetModelBranchAllVersions ServiceFunc = (*Api).GetModelBranchAllVersions
var GetModelBranchVersion ServiceFunc = (*Api).GetModelBranchVersion
var VerifyModelBranchHashStatus ServiceFunc = (*Api).VerifyModelBranchHashStatus
var RegisterModel ServiceFunc = (*Api).RegisterModel
