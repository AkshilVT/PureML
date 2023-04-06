package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/PureMLHQ/PureML/packages/purebackend/test"
	"github.com/labstack/echo/v4"
)

func TestGetOrgByHandle(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "get org by handle + unauthorized + org not found",
			Method:         http.MethodGet,
			Url:            "/api/org/handle/notfound",
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`"status":404`,
				`"data":null`,
				`"message":"Organization not found"`,
			},
		},
		{
			Name:           "get org by handle + unauthorized + org found",
			Method:         http.MethodGet,
			Url:            "/api/org/handle/demo",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserOrgUuid.String() + `"`,
				`"name":"Demo Org"`,
				`"handle":"demo"`,
				`"avatar":""`,
				`"description":"Demo Org Description"`,
				`"message":"Organization details"`,
			},
		},
		{
			Name:   "get org by handle + invalid token",
			Method: http.MethodGet,
			Url:    "/api/org/handle/demo",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "get org by handle + valid token + org not found",
			Method: http.MethodGet,
			Url:    "/api/org/handle/notfound",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidUserToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`"status":404`,
				`"data":null`,
				`"message":"Organization not found"`,
			},
		},
		{
			Name:   "get org by handle + valid token + org found",
			Method: http.MethodGet,
			Url:    "/api/org/handle/demo",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidUserToken,
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserOrgUuid.String() + `"`,
				`"name":"Demo Org"`,
				`"handle":"demo"`,
				`"avatar":""`,
				`"description":"Demo Org Description"`,
				`"message":"Organization details"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestGetOrgByID(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "get org by id + unauthorized",
			Method:         http.MethodGet,
			Url:            "/api/org/id/" + test.ValidAdminUserOrgUuid.String(),
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "get org by id + invalid token",
			Method: http.MethodGet,
			Url:    "/api/org/id/" + test.ValidAdminUserOrgUuid.String(),
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "get org by id + valid token + invalid org uuid",
			Method: http.MethodGet,
			Url:    "/api/org/id/" + test.InvalidOrgUuidString,
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:   "get org by id + valid token + org not found",
			Method: http.MethodGet,
			Url:    "/api/org/id/" + test.ValidNoOrgUuid.String(),
			RequestHeaders: map[string]string{
				"Authorization": test.ValidUserToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:   "get org by id + valid token + org found",
			Method: http.MethodGet,
			Url:    "/api/org/id/" + test.ValidAdminUserOrgUuid.String(),
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserOrgUuid.String() + `"`,
				`"name":"Demo Org"`,
				`"handle":"demo"`,
				`"avatar":""`,
				`"description":"Demo Org Description"`,
				`"join_code":"`,
				`"message":"Organization details"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestGetOrgAllPublicModels(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "get org all public models + unauthorized + invalid org uuid",
			Method:         http.MethodGet,
			Url:            "/api/org/" + test.InvalidOrgUuidString + "/public/model",
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:           "get org all public models + unauthorized + org not found",
			Method:         http.MethodGet,
			Url:            "/api/org/" + test.ValidNoOrgUuid.String() + "/public/model",
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:           "get org all public models + unauthorized",
			Method:         http.MethodGet,
			Url:            "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/public/model",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserOrgUuid.String() + `"`,
				`"name":"Demo Model"`,
				`"wiki":"Demo Model Wiki"`,
				`"is_public":true`,
				`"created_by":{`,
				`"updated_by":{`,
				`"message":"Public models of Organization"`,
			},
		},
		{
			Name:   "get org all public models + invalid token",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/public/model",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "get org all public models + valid token",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/public/model",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserOrgUuid.String() + `"`,
				`"name":"Demo Model"`,
				`"wiki":"Demo Model Wiki"`,
				`"is_public":true`,
				`"created_by":{`,
				`"updated_by":{`,
				`"message":"Public models of Organization"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestGetOrgAllPublicDatasets(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "get org all public datasets + unauthorized + invalid org uuid",
			Method:         http.MethodGet,
			Url:            "/api/org/" + test.InvalidOrgUuidString + "/public/dataset",
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:           "get org all public datasets + unauthorized + org not found",
			Method:         http.MethodGet,
			Url:            "/api/org/" + test.ValidNoOrgUuid.String() + "/public/dataset",
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:           "get org all public datasets + unauthorized",
			Method:         http.MethodGet,
			Url:            "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/public/dataset",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserOrgUuid.String() + `"`,
				`"name":"Demo Dataset"`,
				`"wiki":"Demo Dataset Wiki"`,
				`"is_public":true`,
				`"created_by":{`,
				`"updated_by":{`,
				`"message":"Public datasets of Organization"`,
			},
		},
		{
			Name:   "get org all public datasets + invalid token",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/public/dataset",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "get org all public datasets + valid token",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/public/dataset",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserOrgUuid.String() + `"`,
				`"name":"Demo Dataset"`,
				`"wiki":"Demo Dataset Wiki"`,
				`"is_public":true`,
				`"created_by":{`,
				`"updated_by":{`,
				`"message":"Public datasets of Organization"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestCreateOrg(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "create org + unauthorized",
			Method:         http.MethodPost,
			Url:            "/api/org/create",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "create org + invalid token",
			Method: http.MethodPost,
			Url:    "/api/org/create",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "create org + valid token + no handle",
			Method: http.MethodPost,
			Url:    "/api/org/create",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"handle":"",
				"description":"test"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Organization handle is required"`,
			},
		},
		// TODO: Handle this case without relying on unique constraint
		{
			Name:   "create org + valid token + organization with handle exists",
			Method: http.MethodPost,
			Url:    "/api/org/create",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"handle":"test",
				"description":"test"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Create user
				_, err := app.Dao().CreateOrgFromEmail("demo@aztlan.in", "test", "test", "test")
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 500,
			ExpectedContent: []string{
				`"error":"Internal server error`,
			},
		},
		{
			Name:   "create org + valid token + valid request",
			Method: http.MethodPost,
			Url:    "/api/org/create",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"handle":"test",
				"description":"test"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Delete user if exists
				err := app.Dao().ExecuteSQL(fmt.Sprintf("DELETE FROM organizations WHERE handle = '%s'", "test"))
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"`,
				`"name":"test"`,
				`"handle":"test"`,
				`"avatar":""`,
				`"description":"test"`,
				`"join_code":"`,
				`"message":"Organization created"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUpdateOrg(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "update org + unauthorized",
			Method:         http.MethodPost,
			Url:            "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/update",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "update org + invalid token",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/update",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "update org + valid token + invalid org uuid",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.InvalidOrgUuidString + "/update",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:   "update org + valid token + org not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidNoOrgUuid.String() + "/update",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:   "update org + valid token + org found + not member",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/update",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidUserToken,
			},
			Body: strings.NewReader(`{
				"name":"",
				"description":"",
				"avatar":""
			}`),
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`"status":403`,
				`"data":null`,
				`"message":"You are not authorized to update this organization"`,
			},
		},
		{
			Name:   "update org + valid token + org found + not owner",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/update",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"name":"",
				"description":"",
				"avatar":""
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Change role to member
				err := app.Dao().ExecuteSQL(fmt.Sprintf("UPDATE user_organizations SET role = '%s' WHERE user_uuid = '%s' AND organization_uuid = '%s'", "member", test.ValidAdminUserOrgUuid.String(), test.ValidAdminUserOrgUuid.String()))
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`"status":403`,
				`"data":null`,
				`"message":"You are not authorized to update this organization"`,
			},
		},
		{
			Name:   "update org + valid token + org found + empty name",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/update",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"name":"",
				"description":"",
				"avatar":""
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Name cannot be empty"`,
			},
		},
		{
			Name:   "update org + valid token + user found + update name",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/update",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"name":"Demo Org New"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserOrgUuid.String() + `"`,
				`"name":"Demo Org New"`,
				`"handle":"demo"`,
				`"avatar":""`,
				`"description":"Demo Org Description"`,
				`"join_code":"`,
				`"message":"Organization updated"`,
			},
		},
		{
			Name:   "update org + valid token + user found + update avatar",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/update",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"avatar":"test"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserOrgUuid.String() + `"`,
				`"name":"Demo Org"`,
				`"handle":"demo"`,
				`"avatar":"test"`,
				`"description":"Demo Org Description"`,
				`"join_code":"`,
				`"message":"Organization updated"`,
			},
		},
		{
			Name:   "update org + valid token + user found + update bio",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/update",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"description":"Demo Org Description New"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserOrgUuid.String() + `"`,
				`"name":"Demo Org"`,
				`"handle":"demo"`,
				`"avatar":""`,
				`"description":"Demo Org Description New"`,
				`"join_code":"`,
				`"message":"Organization updated"`,
			},
		},
		{
			Name:   "update org + valid token + user found + update all",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/update",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"name":"Demo Org New",
				"avatar":"test",
				"description":"Demo Org Description New"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserOrgUuid.String() + `"`,
				`"name":"Demo Org New"`,
				`"handle":"demo"`,
				`"avatar":"test"`,
				`"description":"Demo Org Description New"`,
				`"join_code":"`,
				`"message":"Organization updated"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
