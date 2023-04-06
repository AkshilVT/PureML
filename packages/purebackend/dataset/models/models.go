package models

import (
	"time"

	commonmodels "github.com/PureMLHQ/PureML/packages/purebackend/core/common/models"
	userorgmodels "github.com/PureMLHQ/PureML/packages/purebackend/user_org/models"
	uuid "github.com/satori/go.uuid"
)

// Request models

type CreateDatasetRequest struct {
	Wiki        string                     `json:"wiki"`
	IsPublic    bool                       `json:"is_public"`
	Readme      commonmodels.ReadmeRequest `json:"readme"`
	BranchNames []string                   `json:"branch_names"`
}

type CreateDatasetBranchRequest struct {
	BranchName string `json:"branch_name"`
}

type RegisterDatasetRequest struct {
	Hash       string `json:"hash"`
	Lineage    string `json:"lineage"`
	Storage    string `json:"storage"`
	IsEmpty    bool   `json:"is_empty"`
}

type LogFileRequest struct {
	Storage    string `json:"storage"`
}

type DatasetReviewRequest struct {
	FromBranch        string `json:"from_branch"`
	FromBranchVersion string `json:"from_branch_version"`
	ToBranch          string `json:"to_branch"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	IsComplete        bool   `json:"is_complete"`
	IsAccepted        bool   `json:"is_accepted"`
}

type DatasetReviewUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsComplete  bool   `json:"is_complete"`
	IsAccepted  bool   `json:"is_accepted"`
}

// Response models

type DatasetNameResponse struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

type DatasetResponse struct {
	UUID      uuid.UUID                                `json:"uuid"`
	Name      string                                   `json:"name"`
	Wiki      string                                   `json:"wiki"`
	Org       userorgmodels.OrganizationHandleResponse `json:"org"`
	CreatedBy userorgmodels.UserHandleResponse         `json:"created_by"`
	UpdatedBy userorgmodels.UserHandleResponse         `json:"updated_by"`
	Readme    commonmodels.ReadmeResponse              `json:"readme"`
	IsPublic  bool                                     `json:"is_public"`
}

type DatasetUserResponse struct {
	User userorgmodels.UserHandleResponse `json:"user"`
	Role string                           `json:"role"`
}

type DatasetBranchNameResponse struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

type DatasetBranchResponse struct {
	UUID      uuid.UUID           `json:"uuid"`
	Name      string              `json:"name"`
	Dataset   DatasetNameResponse `json:"dataset"`
	IsDefault bool                `json:"is_default"`
}

type DatasetBranchVersionNameResponse struct {
	UUID    uuid.UUID `json:"uuid"`
	Version string    `json:"version"`
}

type DatasetBranchVersionResponse struct {
	UUID       uuid.UUID                        `json:"uuid"`
	Version    string                           `json:"version"`
	Branch     DatasetBranchNameResponse        `json:"branch"`
	Lineage    LineageResponse                  `json:"lineage"`
	Hash       string                           `json:"hash"`
	Path       string                           `json:"path"`
	SourceType string                           `json:"source_type"`
	IsEmpty    bool                             `json:"is_empty"`
	CreatedBy  userorgmodels.UserHandleResponse `json:"created_by"`
	CreatedAt  time.Time                        `json:"created_at"`
}

type LineageResponse struct {
	UUID    uuid.UUID `json:"uuid"`
	Lineage string    `json:"lineage"`
}

type DatasetReviewResponse struct {
	UUID              uuid.UUID                        `json:"uuid"`
	Dataset           DatasetNameResponse              `json:"dataset"`
	FromBranch        DatasetBranchNameResponse        `json:"from_branch"`
	FromBranchVersion DatasetBranchVersionNameResponse `json:"from_branch_version"`
	ToBranch          DatasetBranchNameResponse        `json:"to_branch"`
	Title             string                           `json:"title"`
	Description       string                           `json:"description"`
	CreatedBy         userorgmodels.UserHandleResponse `json:"created_by"`
	AssignedTo        userorgmodels.UserHandleResponse `json:"assigned_to"`
	IsComplete        bool                             `json:"is_complete"`
	IsAccepted        bool                             `json:"is_accepted"`
}
