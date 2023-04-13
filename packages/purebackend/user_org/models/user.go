package models

import (
	uuid "github.com/satori/go.uuid"
)

// Request models

type UserSignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Handle   string `json:"handle"`
	Avatar   string `json:"avatar"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Handle   string `json:"handle"`
	Password string `json:"password"`
}

type UserEmailRequest struct {
	Email string `json:"email"`
}

type UserUpdateRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Bio    string `json:"bio"`
}

type UserOrgJoin struct {
	JoinCode string `json:"join_code"`
}

type UserVerifyTokenRequest struct {
	Token string `json:"token"`
}

type UserResetPasswordRequest struct {
	Token       string `json:"token"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UserRoleRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Response models

type UserClaims struct {
	UUID   uuid.UUID `json:"uuid"`
	Email  string    `json:"email"`
	Handle string    `json:"handle"`
}

type UserHandleResponse struct {
	UUID   uuid.UUID `json:"uuid"`
	Handle string    `json:"handle"`
	Name   string    `json:"name"`
	Avatar string    `json:"avatar"`
	Email  string    `json:"email"`
}

type UserResponse struct {
	UUID       uuid.UUID `json:"uuid"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Handle     string    `json:"handle"`
	Bio        string    `json:"bio"`
	Avatar     string    `json:"avatar"`
	Password   string    `json:"-"`
	IsVerified bool      `json:"is_verified"`
}

type UserProfileResponse struct {
	UUID             uuid.UUID `json:"uuid"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Handle           string    `json:"handle"`
	Bio              string    `json:"bio"`
	Avatar           string    `json:"avatar"`
	NumberOfModels   int64     `json:"number_of_models"`
	NumberOfDatasets int64     `json:"number_of_datasets"`

	Orgs []UserOrganizationsResponse `json:"orgs"`
}

type UserOrganizationsResponse struct {
	Org  OrganizationHandleResponse `json:"org"`
	Role string                     `json:"role"`
}

type UserOrganizationsRoleResponse struct {
	UserUUID uuid.UUID `json:"user_uuid"`
	OrgUUID  uuid.UUID `json:"org_uuid"`
	Role     string    `json:"role"`
}
