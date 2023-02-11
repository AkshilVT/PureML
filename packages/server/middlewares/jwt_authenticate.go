package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PureML-Inc/PureML/server/core"
	"github.com/PureML-Inc/PureML/server/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

const (
	AuthHeaderName                 = "Authorization"
	ContextAuthKey                 = "User"
	ContextOrgKey                  = "Org"
	ContextModelKey                = "Model"
	ContextModelBranchKey          = "ModelBranch"
	ContextModelBranchVersionKey   = "ModelBranchVersion"
	ContextDatasetKey              = "Dataset"
	ContextDatasetBranchKey        = "DatasetBranch"
	ContextDatasetBranchVersionKey = "DatasetBranchVersion"
)

func AuthenticateJWT(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			authHeaderValue := extractRequestHeader(AuthHeaderName, context)
			if authHeaderValue == "" {
				return next(context)
			}
			// the schema is not required and it is only for
			// compatibility with the defaults of some HTTP clients
			authHeaderValue = strings.TrimPrefix(authHeaderValue, "Bearer ")

			token, err := jwt.Parse(authHeaderValue, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("invalid token signing algorithm %v", t.Method.Alg())
				}
				return []byte(app.Settings().AdminAuthToken.Secret), nil
			})
			if err != nil {
				// fmt.Println(err)
				context.Response().WriteHeader(http.StatusForbidden)
				context.Response().Writer.Write([]byte("Could not parse authentication token"))
				return nil
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userUUID := uuid.Must(uuid.FromString(claims["uuid"].(string)))
				user, err := app.Dao().GetUserByUUID(userUUID)
				if err != nil || user != nil {
					context.Set(ContextAuthKey, &models.UserClaims{
						UUID:   user.UUID,
						Email:  claims["email"].(string),
						Handle: claims["handle"].(string),
					})
				} else {
					context.Set(ContextAuthKey, &models.UserClaims{
						UUID:   uuid.Nil,
						Email:  "",
						Handle: "",
					})
				}
			} else {
				// context.Response().WriteHeader(http.StatusForbidden)
				// context.Response().Writer.Write([]byte("Invalid authentication token"))
			}
			return next(context)
		}
	}
}
