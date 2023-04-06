package models

import (
	"time"

	userorgmodels "github.com/PureMLHQ/PureML/packages/purebackend/user_org/models"
	uuid "github.com/satori/go.uuid"
)

// Request models

type CreateSessionRequest struct {
	Device         string `json:"device"`
	DeviceId       string `json:"device_id"`
	DeviceLocation string `json:"device_location"`
}

type VerifySessionRequest struct {
	SessionUUID uuid.UUID `json:"session_id"`
}

type SessionTokenRequest struct {
	SessionUUID uuid.UUID `json:"session_id"`
	DeviceId    string    `json:"device_id"`
}

// Response models

type SessionResponse struct {
	SessionUUID    uuid.UUID `json:"session_id"`
	Device         string    `json:"device"`
	DeviceId       string    `json:"device_id"`
	DeviceLocation string    `json:"device_location"`
	UserUUID       uuid.UUID `json:"user_id"`
	Approved       bool      `json:"approved"`
	Invalid        bool      `json:"invalid"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateSessionResponse struct {
	SessionUUID uuid.UUID `json:"session_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type VerifySessionResponse struct {
	SessionUUID uuid.UUID                        `json:"session_uuid"`
	User        userorgmodels.UserHandleResponse `json:"user"`
	Approved    bool                             `json:"approved"`
	Invalid     bool                             `json:"invalid"`
	CreatedAt   time.Time                        `json:"created_at"`
}
