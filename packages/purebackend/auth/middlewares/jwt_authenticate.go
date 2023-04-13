package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/PureMLHQ/PureML/packages/purebackend/core"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/middlewares"
	"github.com/PureMLHQ/PureML/packages/purebackend/user_org/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

const (
	AuthHeaderName   = "Authorization"
	APIIdHeaderName  = "X-API-ID"
	APIKeyHeaderName = "X-API-KEY"
	ContextAuthKey   = "User"
)

func AuthenticateJWT(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			authHeaderValue := middlewares.ExtractRequestHeader(AuthHeaderName, context)
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
				_, err = context.Response().Writer.Write([]byte("Could not parse authentication token"))
				if err != nil {
					return err
				}
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
			}
			return next(context)
		}
	}
}

func AuthenticateAPI(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			apiIdHeaderValue := middlewares.ExtractRequestHeader(APIIdHeaderName, context)
			if apiIdHeaderValue == "" {
				return next(context)
			}
			apiKeyHeaderValue := middlewares.ExtractRequestHeader(APIKeyHeaderName, context)
			if apiKeyHeaderValue == "" {
				return next(context)
			}
			tokenUUID, err := uuid.FromString(apiIdHeaderValue)
			if err != nil {
				context.Response().WriteHeader(http.StatusForbidden)
				_, err := context.Response().Writer.Write([]byte("API Id is invalid"))
				if err != nil {
					return err
				}
				return nil
			}
			tokenSecretHMAC := hmac.New(sha256.New, []byte(app.Settings().APITokenSecretKey))
			tokenSecretHMAC.Write([]byte(tokenUUID.String()))
			tokenSecret := hex.EncodeToString(tokenSecretHMAC.Sum(nil))
			if tokenSecret != apiKeyHeaderValue {
				context.Response().WriteHeader(http.StatusForbidden)
				_, err := context.Response().Writer.Write([]byte("API key is invalid"))
				if err != nil {
					return err
				}
				return nil
			}
			user, err := app.Dao().GetUserByAPIKeyId(tokenUUID)
			if err != nil || user != nil {
				context.Set(ContextAuthKey, &models.UserClaims{
					UUID:   user.UUID,
					Email:  user.Email,
					Handle: user.Handle,
				})
			} else {
				context.Set(ContextAuthKey, &models.UserClaims{
					UUID:   uuid.Nil,
					Email:  "",
					Handle: "",
				})
			}
			return next(context)
		}
	}
}
