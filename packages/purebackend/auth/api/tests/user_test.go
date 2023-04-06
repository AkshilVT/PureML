package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/PureMLHQ/PureML/packages/purebackend/test"
	"github.com/labstack/echo/v4"
)

func TestUserSignUp(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:   "signup + no email",
			Method: http.MethodPost,
			Url:    "/api/user/signup",
			Body: strings.NewReader(`{
				"handle":"test",
				"name":"test",
				"password":"test"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Email is required"`,
			},
		},
		{
			Name:   "signup + invalid email",
			Method: http.MethodPost,
			Url:    "/api/user/signup",
			Body: strings.NewReader(`{
				"email":"invalidEmail",
				"handle":"test",
				"name":"test",
				"password":"test"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Email is invalid"`,
			},
		},
		{
			Name:   "signup + no handle",
			Method: http.MethodPost,
			Url:    "/api/user/signup",
			Body: strings.NewReader(`{
				"email":"test@test.com",
				"name":"test",
				"password":"test"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Handle is required"`,
			},
		},
		{
			Name:   "signup + no password",
			Method: http.MethodPost,
			Url:    "/api/user/signup",
			Body: strings.NewReader(`{
				"email":"test@test.com",
				"handle":"test",
				"name":"test"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Password is required"`,
			},
		},
		{
			Name:   "signup + user exists with same email or handle",
			Method: http.MethodPost,
			Url:    "/api/user/signup",
			Body: strings.NewReader(`{
				"email":"test@test.com",
				"handle":"test",
				"name":"test",
				"password":"test"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Create user
				_, err := app.Dao().CreateUser("test", "test@test.com", "test", "", "", "$2a$10$N..OOp8lPw0fRGCXT.HxH.LO8BUKwlncI/ufXK/bLTEvyeFmdCun.", true)
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 409,
			ExpectedContent: []string{
				`"status":409`,
				`"data":null`,
				`"message":"User with email already exists"`,
			},
		},
		{
			Name:   "signup successful",
			Method: http.MethodPost,
			Url:    "/api/user/signup",
			Body: strings.NewReader(`{
				"email":"test@test.com",
				"handle":"test",
				"name":"test",
				"password":"test"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Delete user if exists
				err := app.Dao().ExecuteSQL(fmt.Sprintf("DELETE FROM users WHERE email = '%s'", "test@test.com"))
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"name":"test",`,
				`"email":"test@test.com",`,
				`"handle":"test",`,
				`"message":"User created"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUserLogin(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:   "login + no email & no handle",
			Method: http.MethodPost,
			Url:    "/api/user/login",
			Body: strings.NewReader(`{
				"password":"test"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Email or handle is required"`,
			},
		},
		{
			Name:   "login + no password",
			Method: http.MethodPost,
			Url:    "/api/user/login",
			Body: strings.NewReader(`{
				"email":"test@test.com",
				"handle":"test"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Password is required"`,
			},
		},
		{
			Name:   "login + invalid email",
			Method: http.MethodPost,
			Url:    "/api/user/login",
			Body: strings.NewReader(`{
				"email":"invalid",
				"password":"test"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Email is invalid"`,
			},
		},
		{
			Name:   "login + invalid credentials with email",
			Method: http.MethodPost,
			Url:    "/api/user/login",
			Body: strings.NewReader(`{
				"email":"demo@aztlan.in",
				"password":"wrongpassword"
			}`),
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`"status":401`,
				`"data":null`,
				`"message":"Invalid credentials"`,
			},
		},
		{
			Name:   "login + invalid credentials with handle",
			Method: http.MethodPost,
			Url:    "/api/user/login",
			Body: strings.NewReader(`{
				"handle":"demo",
				"password":"wrongpassword"
			}`),
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`"status":401`,
				`"data":null`,
				`"message":"Invalid credentials"`,
			},
		},
		{
			Name:   "login + email not verified",
			Method: http.MethodPost,
			Url:    "/api/user/login",
			Body: strings.NewReader(`{
				"handle":"test",
				"password":"test"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Create user
				_, err := app.Dao().CreateUser("test", "test@test.com", "test", "", "", "$2a$10$N..OOp8lPw0fRGCXT.HxH.LO8BUKwlncI/ufXK/bLTEvyeFmdCun.", false)
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`"status":401`,
				`"data":null`,
				`"message":"User email is not verified"`,
			},
		},
		{
			Name:   "login successful with email",
			Method: http.MethodPost,
			Url:    "/api/user/login",
			Body: strings.NewReader(`{
				"email":"demo@aztlan.in",
				"password":"demo"
			}`),
			ExpectedStatus: 202,
			ExpectedContent: []string{
				`"status":202`,
				`"data":[{`,
				`"email":"demo@aztlan.in"`,
				`"accessToken":`,
				`"message":"User logged in"`,
			},
		},
		{
			Name:   "login successful with handle",
			Method: http.MethodPost,
			Url:    "/api/user/login",
			Body: strings.NewReader(`{
				"handle":"demo",
				"password":"demo"
			}`),
			ExpectedStatus: 202,
			ExpectedContent: []string{
				`"status":202`,
				`"data":[{`,
				`"email":"demo@aztlan.in"`,
				`"accessToken":`,
				`"message":"User logged in"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

// TODO
func TestUserResetPassword(t *testing.T) {
}

// TODO
func TestUserForgotPassword(t *testing.T) {
}

func TestGetProfile(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "get profile + unauthorized",
			Method:         http.MethodGet,
			Url:            "/api/user/profile",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "get profile + invalid token",
			Method: http.MethodGet,
			Url:    "/api/user/profile",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "get profile + valid token + user not found",
			Method: http.MethodGet,
			Url:    "/api/user/profile",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidTokenNoUser,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`User not found`,
			},
		},
		{
			Name:   "get profile + valid token + user found",
			Method: http.MethodGet,
			Url:    "/api/user/profile",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserUuid.String() + `"`,
				`"email":"demo@aztlan.in"`,
				`"handle":"demo"`,
				`"name":"Demo User"`,
				`"message":"User profile"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestGetProfileByHandle(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "get profile by handle + unauthorized + user not found",
			Method:         http.MethodGet,
			Url:            "/api/user/profile/noone",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":null`,
				`"message":"Public User profile"`,
			},
		},
		{
			Name:           "get profile by handle + unauthorized + user found",
			Method:         http.MethodGet,
			Url:            "/api/user/profile/demo",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"email":"demo@aztlan.in"`,
				`"handle":"demo"`,
				`"name":"Demo User"`,
				`"message":"Public User profile"`,
			},
		},
		{
			Name:   "get profile by handle + invalid token",
			Method: http.MethodGet,
			Url:    "/api/user/profile/demo",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "get profile by handle + valid token + user not found",
			Method: http.MethodGet,
			Url:    "/api/user/profile/noone",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidTokenNoUser,
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":null`,
				`"message":"Public User profile"`,
			},
		},
		{
			Name:   "get profile + valid token + user found",
			Method: http.MethodGet,
			Url:    "/api/user/profile/demo",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"email":"demo@aztlan.in"`,
				`"handle":"demo"`,
				`"name":"Demo User"`,
				`"message":"Public User profile"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestUpdateProfile(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "update profile + unauthorized",
			Method:         http.MethodPost,
			Url:            "/api/user/profile",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "update profile + invalid token",
			Method: http.MethodPost,
			Url:    "/api/user/profile",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "update profile + valid token + user not found",
			Method: http.MethodPost,
			Url:    "/api/user/profile",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidTokenNoUser,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`User not found`,
			},
		},
		{
			Name:   "update profile + valid token + user found + empty name",
			Method: http.MethodPost,
			Url:    "/api/user/profile",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"name":"",
				"avatar":"",
				"bio":""
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Name cannot be empty"`,
			},
		},
		{
			Name:   "update profile + valid token + user found + update name",
			Method: http.MethodPost,
			Url:    "/api/user/profile",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"name":"Demo User New"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"email":"demo@aztlan.in"`,
				`"handle":"demo"`,
				`"name":"Demo User New"`,
				`"avatar":""`,
				`"bio":"Demo User Bio"`,
				`"message":"User profile updated"`,
			},
		},
		{
			Name:   "update profile + valid token + user found + update avatar",
			Method: http.MethodPost,
			Url:    "/api/user/profile",
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
				`"email":"demo@aztlan.in"`,
				`"handle":"demo"`,
				`"name":"Demo User"`,
				`"avatar":"test"`,
				`"bio":"Demo User Bio"`,
				`"message":"User profile updated"`,
			},
		},
		{
			Name:   "update profile + valid token + user found + update bio",
			Method: http.MethodPost,
			Url:    "/api/user/profile",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"bio":"Demo User Bio New"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"email":"demo@aztlan.in"`,
				`"handle":"demo"`,
				`"name":"Demo User"`,
				`"avatar":""`,
				`"bio":"Demo User Bio New"`,
				`"message":"User profile updated"`,
			},
		},
		{
			Name:   "update profile + valid token + user found + update all",
			Method: http.MethodPost,
			Url:    "/api/user/profile",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"name":"Demo User New",
				"avatar":"avatar",
				"bio":"bio"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"email":"demo@aztlan.in"`,
				`"handle":"demo"`,
				`"name":"Demo User New"`,
				`"avatar":"avatar"`,
				`"bio":"bio"`,
				`"message":"User profile updated"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

// TODO
func TestDeleteProfile(t *testing.T) {
}

func TestCreateSession(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:   "create session + no device id",
			Method: http.MethodPost,
			Url:    "/api/user/create-session",
			Body: strings.NewReader(`{
				"device":"",
				"device_id":"",
				"device_location":""
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Device Id is required"`,
			},
		},
		{
			Name:   "create session successful",
			Method: http.MethodPost,
			Url:    "/api/user/create-session",
			Body: strings.NewReader(`{
				"device":"mydevice",
				"device_id":"123",
				"device_location":"India"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"session_id":"`,
				`"created_at":"`,
				`"message":"Session created"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestGetSessionToken(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:   "get session token + no session id",
			Method: http.MethodPost,
			Url:    "/api/user/session-token",
			Body: strings.NewReader(`{
				"session_id":"",
				"device":"",
				"device_id":"",
				"device_location":""
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Session Id is required"`,
			},
		},
		{
			Name:   "get session token + no device id",
			Method: http.MethodPost,
			Url:    "/api/user/session-token",
			Body: strings.NewReader(`{
				"session_id":"` + test.ValidSessionUuid.String() + `",
				"device":"",
				"device_id":"",
				"device_location":""
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Device Id is required"`,
			},
		},
		{
			Name:   "get session token + not found",
			Method: http.MethodPost,
			Url:    "/api/user/session-token",
			Body: strings.NewReader(`{
				"session_id":"` + test.ValidNoSessionUuid.String() + `",
				"device":"device",
				"device_id":"123456",
				"device_location":"India"
			}`),
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`"status":404`,
				`"data":null`,
				`"message":"Session not found"`,
			},
		},
		{
			Name:   "get session token + not approved",
			Method: http.MethodPost,
			Url:    "/api/user/session-token",
			Body: strings.NewReader(`{
				"session_id":"` + test.ValidSessionUuid.String() + `",
				"device":"device",
				"device_id":"123456",
				"device_location":"India"
			}`),
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`"status":401`,
				`"data":null`,
				`"message":"Session not approved"`,
			},
		},
		{
			Name:   "get session token + invalid session",
			Method: http.MethodPost,
			Url:    "/api/user/session-token",
			Body: strings.NewReader(`{
				"session_id":"` + test.ValidSessionUuid.String() + `",
				"device":"device",
				"device_id":"123456",
				"device_location":"India"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Make changes to session in db to simulate conditions
				_, err := app.Dao().UpdateSession(test.ValidSessionUuid, test.ValidAdminUserUuid, map[string]interface{}{
					"approved": true,
					"invalid":  true,
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`"status":403`,
				`"data":null`,
				`"message":"Session invalid"`,
			},
		},
		{
			Name:   "get session token + invalid session device",
			Method: http.MethodPost,
			Url:    "/api/user/session-token",
			Body: strings.NewReader(`{
				"session_id":"` + test.ValidSessionUuid.String() + `",
				"device":"device",
				"device_id":"not123456",
				"device_location":"India"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Make changes to session in db to simulate conditions
				_, err := app.Dao().UpdateSession(test.ValidSessionUuid, test.ValidAdminUserUuid, map[string]interface{}{
					"approved": true,
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`"status":403`,
				`"data":null`,
				`"message":"Session invalid device"`,
			},
		},
		{
			Name:   "get session token successful",
			Method: http.MethodPost,
			Url:    "/api/user/session-token",
			Body: strings.NewReader(`{
				"session_id":"` + test.ValidSessionUuid.String() + `",
				"device":"device",
				"device_id":"123456",
				"device_location":"India"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Make changes to session in db to simulate conditions
				_, err := app.Dao().UpdateSession(test.ValidSessionUuid, test.ValidAdminUserUuid, map[string]interface{}{
					"approved": true,
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"email":"`,
				`"accessToken":"`,
				`"message":"Session Token created"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestVerifySession(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "verify session + unauthorized",
			Method:         http.MethodPost,
			Url:            "/api/user/verify-session",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "verify session + invalid token",
			Method: http.MethodPost,
			Url:    "/api/user/verify-session",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "verify session + valid token + user not found",
			Method: http.MethodPost,
			Url:    "/api/user/verify-session",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidTokenNoUser,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`User not found`,
			},
		},
		{
			Name:   "verify session + no session id",
			Method: http.MethodPost,
			Url:    "/api/user/verify-session",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"session_id":""
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Session Id is required"`,
			},
		},
		{
			Name:   "session verified and approved",
			Method: http.MethodPost,
			Url:    "/api/user/verify-session",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"session_id":"` + test.ValidSessionUuid.String() + `"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Make changes to session in db to simulate conditions
				_, err := app.Dao().UpdateSession(test.ValidSessionUuid, test.ValidAdminUserUuid, map[string]interface{}{
					"approved": true,
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Session already approved"`,
			},
		},
		{
			Name:   "session verified and invalid",
			Method: http.MethodPost,
			Url:    "/api/user/verify-session",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"session_id":"` + test.ValidSessionUuid.String() + `"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Make changes to session in db to simulate conditions
				_, err := app.Dao().UpdateSession(test.ValidSessionUuid, test.ValidAdminUserUuid, map[string]interface{}{
					"invalid": true,
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`"status":403`,
				`"data":null`,
				`"message":"Session invalid"`,
			},
		},
		{
			Name:   "session verified and expired",
			Method: http.MethodPost,
			Url:    "/api/user/verify-session",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"session_id":"` + test.ValidSessionUuid.String() + `"
			}`),
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`"status":403`,
				`"data":null`,
				`"message":"Session expired"`,
			},
		},
		{
			Name:   "session verified and approved successfully",
			Method: http.MethodPost,
			Url:    "/api/user/verify-session",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"session_id":"` + test.ValidSessionUuid.String() + `"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Make changes to session in db to simulate conditions
				_, err := app.Dao().UpdateSession(test.ValidSessionUuid, test.ValidAdminUserUuid, map[string]interface{}{
					"created_at": time.Now().Add(-time.Minute * 5),
				})
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":null`,
				`"message":"Session approved"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
