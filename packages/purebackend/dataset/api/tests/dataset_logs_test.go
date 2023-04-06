package tests

import (
	"net/http"
	"strings"
	"testing"

	"github.com/PureMLHQ/PureML/packages/purebackend/test"
	"github.com/labstack/echo/v4"
)

func TestGetAllLogsDataset(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "get all logs of dataset + unauthorized",
			Method:         http.MethodGet,
			Url:            "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "get all logs of dataset + invalid token",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "get all logs of dataset + valid token + user not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidTokenNoUser,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`User not found`,
			},
		},
		{
			Name:   "get all logs of dataset + valid token + invalid org uuid",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.InvalidOrgUuidString + "/dataset/Demo%20Dataset/branch/dev/version/v1/log",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:   "get all logs of dataset + valid token + org not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidNoOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:   "get all logs of dataset + valid token + dataset not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/NoDataset/branch/dev/version/v1/log",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset not found`,
			},
		},
		{
			Name:   "get all logs of dataset + valid token + dataset branch not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/nobranch/version/v1/log",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset Branch not found`,
			},
		},
		{
			Name:   "get all logs of dataset + valid token + dataset branch version not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v2/log",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset Branch Version not found`,
			},
		},
		{
			Name:   "get all logs of dataset + valid token + dataset branch version found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Make notadmin a "member" of the admin user org
				_, err := app.Dao().CreateLogForDatasetVersion("accuracy", "accuracyData", test.ValidAdminUserOrgUuid)
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"key":"accuracy"`,
				`"data":"accuracyData"`,
				`"dataset_version":{`,
				`"uuid":"` + test.ValidAdminUserUuid.String() + `"`,
				`"version":"v1"`,
				`"message":"Logs for dataset version"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestGetKeyLogsDataset(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "get log by key of dataset + unauthorized",
			Method:         http.MethodGet,
			Url:            "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log/accuracy",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "get log by key of dataset + invalid token",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log/accuracy",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "get log by key of dataset + valid token + user not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log/accuracy",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidTokenNoUser,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`User not found`,
			},
		},
		{
			Name:   "get log by key of dataset + valid token + invalid org uuid",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.InvalidOrgUuidString + "/dataset/Demo%20Dataset/branch/dev/version/v1/log/accuracy",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:   "get log by key of dataset + valid token + org not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidNoOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log/accuracy",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:   "get log by key of dataset + valid token + dataset not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/NoDataset/branch/dev/version/v1/log/accuracy",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset not found`,
			},
		},
		{
			Name:   "get log by key of dataset + valid token + dataset branch not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/nobranch/version/v1/log/accuracy",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset Branch not found`,
			},
		},
		{
			Name:   "get log by key of dataset + valid token + dataset branch version not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v2/log/accuracy",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset Branch Version not found`,
			},
		},
		{
			Name:   "get log by key of dataset + valid token + dataset branch version found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log/accuracy",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Make notadmin a "member" of the admin user org
				_, err := app.Dao().CreateLogForDatasetVersion("accuracy", "accuracyData", test.ValidAdminUserOrgUuid)
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"key":"accuracy"`,
				`"data":"accuracyData"`,
				`"dataset_version":{`,
				`"uuid":"` + test.ValidAdminUserUuid.String() + `"`,
				`"version":"v1"`,
				`"message":"Specific Key Logs for dataset version"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestLogDataset(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "create dataset branch version log + unauthorized",
			Method:         http.MethodPost,
			Url:            "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "create dataset branch version log + invalid token",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "create dataset branch version log + valid token + invalid org uuid",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.InvalidOrgUuidString + "/dataset/Demo%20Dataset/branch/dev/version/v1/log",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:   "create dataset branch version log + valid token + org not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidNoOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidUserToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:   "create dataset branch version log + valid token + dataset not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/NoDataset/branch/dev/version/v1/log",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset not found`,
			},
		},
		{
			Name:   "create dataset branch version log + valid token + dataset branch not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/nobranch/version/v1/log",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset Branch not found`,
			},
		},
		{
			Name:   "create dataset branch version log + valid token + dataset branch version not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v2/log",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset Branch Version not found`,
			},
		},
		{
			Name:   "create dataset branch version log + valid token + key empty",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"key":"",
				"data":""
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Key is required"`,
			},
		},
		{
			Name:   "create dataset branch version log + valid token + log created successfully",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1/log",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"key":"test",
				"data":"test"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"key":"test"`,
				`"data":"test"`,
				`"message":"Log created"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

// TODO
func TestLogFileDataset(t *testing.T) {

}
