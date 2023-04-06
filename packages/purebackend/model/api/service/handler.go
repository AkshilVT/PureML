package service

import (
	"fmt"
	"net/http"
	"strings"

	coreservice "github.com/PureMLHQ/PureML/packages/purebackend/core/apis/service"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/models"
	commonmodels "github.com/PureMLHQ/PureML/packages/purebackend/core/common/models"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type ServiceFunc func(*Api, *models.Request) *models.Response

func (api *Api) DefaultHandler(f ServiceFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		request := coreservice.ExtractRequest(context)
		response := f(api, request)
		responseWriter := context.Response().Writer
		if response.Error != nil {
			populateErrorResponse(context, response, responseWriter)
		} else {
			populateSuccessResponse(context, response, responseWriter)
		}
		return nil
	}
}

func populateSuccessResponse(context echo.Context, response *models.Response, responseWriter http.ResponseWriter) {
	context.Response().WriteHeader(response.StatusCode)
	_, err := responseWriter.Write(coreservice.ConvertToBytes(response.Body))
	if err != nil {
		panic(fmt.Sprintf("Error writing response: %v \n", err.Error()))
	}
}

func populateErrorResponse(context echo.Context, response *models.Response, responseWriter http.ResponseWriter) {
	context.Response().WriteHeader(http.StatusInternalServerError)
	_, err := responseWriter.Write(coreservice.ConvertToBytes(map[string]interface{}{
		"error": "Internal server error - " + response.Error.Error(),
	}))
	if err != nil {
		panic(fmt.Sprintf("Error writing response: %v \n", err.Error()))
	}
}

func (api *Api) ValidateSourceTypeAndGetSourceSecrets(modelSourceSecretName string, orgId uuid.UUID) (*commonmodels.SourceSecrets, *models.Response) {
	var sourceSecrets *commonmodels.SourceSecrets
	var err error
	if strings.ToUpper(modelSourceSecretName) == "PUREML-STORAGE" {
		sourceSecrets, err = api.app.Dao().GetSecretByName(uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111")), "admin")
		sourceSecrets.SourceType = "PUREML-STORAGE"
	} else {
		sourceSecrets, err = api.app.Dao().GetSecretByName(orgId, modelSourceSecretName)
	}
	if sourceSecrets == nil || err != nil {
		return nil, models.NewErrorResponse(http.StatusBadRequest, fmt.Sprintf("Source %s not connected properly to organization", modelSourceSecretName))
	}
	return sourceSecrets, nil
}
