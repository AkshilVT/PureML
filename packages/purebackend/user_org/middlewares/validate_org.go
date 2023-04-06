package middlewares

import (
	"net/http"

	"github.com/PureMLHQ/PureML/packages/purebackend/core"
	userorgmodels "github.com/PureMLHQ/PureML/packages/purebackend/user_org/models"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

const (
	ContextOrgKey = "Org"
)

func ValidateOrg(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			orgId := context.Param("orgId")
			if orgId == "" {
				context.Response().WriteHeader(http.StatusBadRequest)
				_, err := context.Response().Writer.Write([]byte("Organization Id required"))
				if err != nil {
					return err
				}
				return nil
			}
			orgUUID, err := uuid.FromString(orgId)
			if err != nil {
				context.Response().WriteHeader(http.StatusBadRequest)
				_, err = context.Response().Writer.Write([]byte("Invalid UUID format"))
				if err != nil {
					return err
				}
				return nil
			}
			org, err := app.Dao().GetOrgById(orgUUID)
			if err != nil {
				context.Response().WriteHeader(http.StatusInternalServerError)
				_, err = context.Response().Writer.Write([]byte(err.Error()))
				if err != nil {
					return err
				}
				return nil
			}
			if org == nil {
				context.Response().WriteHeader(http.StatusNotFound)
				_, err = context.Response().Writer.Write([]byte("Organization not found"))
				if err != nil {
					return err
				}
				return nil
			}
			userUUID := context.Get("User").(*userorgmodels.UserClaims).UUID
			if userUUID != uuid.Nil {
				res, err := app.Dao().GetUserOrganizationByOrgIdAndUserUUID(orgUUID, userUUID)
				if err != nil {
					context.Response().WriteHeader(http.StatusInternalServerError)
					_, err = context.Response().Writer.Write([]byte(err.Error()))
					if err != nil {
						return err
					}
					return nil
				}
				if res == nil {
					context.Response().WriteHeader(http.StatusForbidden)
					_, err = context.Response().Writer.Write([]byte("You are not a member of this organization"))
					if err != nil {
						return err
					}
					return nil
				}
			}
			context.Set(ContextOrgKey, &userorgmodels.OrganizationHandleResponse{
				Name:        org.Name,
				UUID:        org.UUID,
				Handle:      org.Handle,
				Avatar:      org.Avatar,
				Description: org.Description,
			})
			return next(context)
		}
	}
}
