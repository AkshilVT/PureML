package tests

import (
	"net/http"
	"strings"
	"testing"

	"github.com/PureMLHQ/PureML/packages/purebackend/test"
)

func TestGetDatasetBranchAllVersions(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "get specific version of dataset branch + unauthorized",
			Method:         http.MethodGet,
			Url:            "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "get all versions of dataset branch + invalid token",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "get all versions of dataset branch + valid token + user not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidTokenNoUser,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`User not found`,
			},
		},
		{
			Name:   "get all versions of dataset branch + valid token + invalid org uuid",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.InvalidOrgUuidString + "/dataset/Demo%20Dataset/branch/dev/version",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:   "get all versions of dataset branch + valid token + org not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidNoOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:   "get all versions of dataset branch + valid token + dataset not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/NoDataset/branch/dev/version",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset not found`,
			},
		},
		{
			Name:   "get all versions of dataset branch + valid token + branch not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/nobranch/version",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset Branch not found`,
			},
		},
		{
			Name:   "get all versions of dataset branch + valid token + dataset branch found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserUuid.String() + `"`,
				`"hash":"1234567890"`,
				`"version":"v1"`,
				`"branch":{`,
				`"name":"dev"`,
				`"path":{`,
				`"source_path":`,
				`"source_type":{`,
				`"lineage":{`,
				`"lineage":"{}"`,
				`"created_by":{`,
				`"name":"Demo User"`,
				`"avatar":`,
				`"email":"demo@aztlan.in"`,
				`"handle":"demo"`,
				`"created_at":`,
				`"is_empty":true`,
				`"message":"All dataset branch versions"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestGetDatasetBranchVersion(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "get specific version of dataset branch + unauthorized",
			Method:         http.MethodGet,
			Url:            "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "get specific version of dataset branch + invalid token",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "get specific version of dataset branch + valid token + user not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidTokenNoUser,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`User not found`,
			},
		},
		{
			Name:   "get specific version of dataset branch + valid token + invalid org uuid",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.InvalidOrgUuidString + "/dataset/Demo%20Dataset/branch/dev/version/v1",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:   "get specific version of dataset branch + valid token + org not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidNoOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:   "get specific version of dataset branch + valid token + dataset not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/NoDataset/branch/dev/version/v1",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset not found`,
			},
		},
		{
			Name:   "get specific version of dataset branch + valid token + branch not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/nobranch/version",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset Branch not found`,
			},
		},
		{
			Name:   "get specific version of dataset branch + valid token + version not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v2",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset Branch Version not found`,
			},
		},
		{
			Name:   "get specific version of dataset branch + valid token + dataset branch version found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/version/v1",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserUuid.String() + `"`,
				`"hash":"1234567890"`,
				`"version":"v1"`,
				`"branch":{`,
				`"name":"dev"`,
				`"path":{`,
				`"source_path":`,
				`"source_type":{`,
				`"lineage":{`,
				`"lineage":"{}"`,
				`"created_by":{`,
				`"name":"Demo User"`,
				`"avatar":`,
				`"email":"demo@aztlan.in"`,
				`"handle":"demo"`,
				`"created_at":`,
				`"is_empty":true`,
				`"message":"Dataset branch version details"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestVerifyDatasetBranchHashStatus(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "verify dataset hash status + unauthorized",
			Method:         http.MethodPost,
			Url:            "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/hash-status",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "verify dataset hash status + invalid token",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/hash-status",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "verify dataset hash status + valid token + invalid org uuid",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.InvalidOrgUuidString + "/dataset/Demo%20Dataset/branch/dev/hash-status",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:   "verify dataset hash status + valid token + org not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidNoOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/hash-status",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidUserToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:   "verify dataset hash status + valid token + dataset not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/noDataset/branch/dev/hash-status",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidUserToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset not found`,
			},
		},
		{
			Name:   "verify dataset hash status + valid token + hash empty",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/hash-status",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"hash":""
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Hash value is empty"`,
			},
		},
		{
			Name:   "verify dataset hash status + valid token + branch not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/nobranch/hash-status",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"hash":"test"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[false]`,
				`"message":"Hash validity (False - does not exist in db)"`,
			},
		},
		{
			Name:   "verify dataset hash status + valid token + hash does not exist",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/hash-status",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"hash":"test"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[false]`,
				`"message":"Hash validity (False - does not exist in db)"`,
			},
		},
		{
			Name:   "verify dataset hash status + valid token + hash exists",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/hash-status",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"hash":"1234567890"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[true]`,
				`"message":"Hash validity (True - exists in db)"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

// TODO
func TestRegisterDataset(t *testing.T) {
	emptyHashMultipartBody, emptyHashMultipartContentType, err := test.MockMultipartData(map[string]string{
		"hash":     "",
		"storage":  "local",
		"is_empty": "true",
		"lineage":  "{}",
	}, "file")
	if err != nil {
		t.Fatal(err)
	}
	emptyFileMultipartBody, emptyFileMultipartContentType, err := test.MockMultipartData(map[string]string{
		"hash":     "uniquehash",
		"storage":  "local",
		"is_empty": "true",
		"lineage":  "{}",
	})
	if err != nil {
		t.Fatal(err)
	}
	invalidStorageMultipartBody, invalidStorageMultipartContentType, err := test.MockMultipartData(map[string]string{
		"hash":     "uniquehash",
		"storage":  "nope",
		"is_empty": "true",
		"lineage":  "{}",
	}, "file")
	if err != nil {
		t.Fatal(err)
	}
	invalidHashMultipartBody, invalidHashMultipartContentType, err := test.MockMultipartData(map[string]string{
		"hash":     "1234567890",
		"storage":  "local",
		"is_empty": "true",
		"lineage":  "{}",
	}, "file")
	if err != nil {
		t.Fatal(err)
	}
	validMultipartBody, validMultipartContentType, err := test.MockMultipartData(map[string]string{
		"hash":     "uniquehash",
		"storage":  "local",
		"is_empty": "true",
		"lineage":  "{}",
	}, "file")
	if err != nil {
		t.Fatal(err)
	}
	validProtectedMultipartBody, validProtectedMultipartContentType, err := test.MockMultipartData(map[string]string{
		"hash":     "uniquehash",
		"storage":  "local",
		"is_empty": "true",
		"lineage":  "{}",
	}, "file")
	if err != nil {
		t.Fatal(err)
	}
	validS3MultipartBody, validS3MultipartContentType, err := test.MockMultipartData(map[string]string{
		"hash":     "uniquehash",
		"storage":  "s3",
		"is_empty": "true",
		"lineage":  "{}",
	}, "file")
	if err != nil {
		t.Fatal(err)
	}
	validR2MultipartBody, validR2MultipartContentType, err := test.MockMultipartData(map[string]string{
		"hash":     "uniquehash",
		"storage":  "r2",
		"is_empty": "true",
		"lineage":  "{}",
	}, "file")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []test.ApiScenario{
		{
			Name:           "register dataset + unauthorized",
			Method:         http.MethodPost,
			Url:            "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/register",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "register dataset + invalid token",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/register",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "register dataset + valid token + invalid org uuid",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.InvalidOrgUuidString + "/dataset/Demo%20Dataset/branch/dev/register",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:   "register dataset + valid token + org not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidNoOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/register",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidUserToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:   "register dataset + valid token + dataset not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/noDataset/branch/dev/register",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidUserToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset not found`,
			},
		},
		{
			Name:   "register dataset + valid token + branch not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/nobranch/register",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidUserToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Dataset Branch not found`,
			},
		},
		{
			Name:   "register dataset + valid token + hash empty",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/register",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
				"Content-Type":  emptyHashMultipartContentType.FormDataContentType(),
			},
			Body:           emptyHashMultipartBody,
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Hash is required"`,
			},
		},
		{
			Name:   "register dataset + valid token + file empty",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/register",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
				"Content-Type":  emptyFileMultipartContentType.FormDataContentType(),
			},
			Body:           emptyFileMultipartBody,
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"File is required"`,
			},
		},
		{
			Name:   "register dataset + valid token + protected main branch",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/main/register",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
				"Content-Type":  validProtectedMultipartContentType.FormDataContentType(),
			},
			Body:           validProtectedMultipartBody,
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Cannot register dataset directly to main branch"`,
			},
		},
		{
			Name:   "register dataset + valid token + unsupported storage",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/register",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
				"Content-Type":  invalidStorageMultipartContentType.FormDataContentType(),
			},
			Body:           invalidStorageMultipartBody,
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Unsupported dataset storage"`,
			},
		},
		{
			Name:   "register dataset + valid token + hash already exists",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/register",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
				"Content-Type":  invalidHashMultipartContentType.FormDataContentType(),
			},
			Body:           invalidHashMultipartBody,
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Dataset with this hash already exists"`,
			},
		},
		{
			Name:   "register dataset + valid token + s3 not enabled",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/register",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
				"Content-Type":  validS3MultipartContentType.FormDataContentType(),
			},
			Body:           validS3MultipartBody,
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"S3 source not enabled"`,
			},
		},
		{
			Name:   "register dataset + valid token + r2 not enabled",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/register",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
				"Content-Type":  validR2MultipartContentType.FormDataContentType(),
			},
			Body:           validR2MultipartBody,
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"R2 source not enabled"`,
			},
		},
		{
			Name:   "register dataset + valid token + registered successfully",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/dataset/Demo%20Dataset/branch/dev/register",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
				"Content-Type":  validMultipartContentType.FormDataContentType(),
			},
			Body:           validMultipartBody,
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"`,
				`"hash":"uniquehash"`,
				`"version":"v2"`,
				`"branch":{`,
				`"name":"dev"`,
				`"path":{`,
				`"source_path":`,
				`"source_type":{`,
				`"lineage":{`,
				`"lineage":"{}"`,
				`"created_by":{`,
				`"name":"Demo User"`,
				`"avatar":`,
				`"email":"demo@aztlan.in"`,
				`"handle":"demo"`,
				`"created_at":`,
				`"is_empty":true`,
				`"message":"Dataset successfully registered"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
