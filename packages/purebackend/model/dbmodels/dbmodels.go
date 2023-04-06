package dbmodels

import (
	commondbmodels "github.com/PureMLHQ/PureML/packages/purebackend/core/common/dbmodels"
	userorgdbmodels "github.com/PureMLHQ/PureML/packages/purebackend/user_org/dbmodels"
	uuid "github.com/satori/go.uuid"
)

type Model struct {
	commondbmodels.BaseModel `gorm:"embedded"`
	Name                     string        `json:"name" gorm:"not null;index:idx_org_model_name,unique"`
	Wiki                     string        `json:"wiki"`
	OrganizationUUID         uuid.UUID     `json:"org_uuid" gorm:"type:uuid;not null;index:idx_org_model_name,unique"`
	CreatedBy                uuid.UUID     `json:"created_by" gorm:"type:uuid;not null"`
	UpdatedBy                uuid.NullUUID `json:"updated_by" gorm:"type:uuid;"`
	IsPublic                 bool          `json:"is_public" default:"false"`

	Org           userorgdbmodels.Organization `gorm:"foreignKey:OrganizationUUID"`
	CreatedByUser userorgdbmodels.User         `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser userorgdbmodels.User         `gorm:"foreignKey:UpdatedBy"`

	Readme   commondbmodels.Readme  `gorm:"foreignKey:ModelUUID"`
	Branches []ModelBranch          `gorm:"foreignKey:ModelUUID"`
	Users    []userorgdbmodels.User `gorm:"many2many:model_users;"`
}

type ModelUser struct {
	ModelUUID uuid.UUID `json:"model_uuid" gorm:"type:uuid;primaryKey"`
	UserUUID  uuid.UUID `json:"user_uuid" gorm:"type:uuid;primaryKey"`
	Role      string    `json:"role" gorm:"not null;default:member"`
}

type ModelBranch struct {
	commondbmodels.BaseModel `gorm:"embedded"`
	Name                     string    `json:"name" gorm:"not null;index:idx_model_branch,unique"`
	ModelUUID                uuid.UUID `json:"model_uuid" gorm:"type:uuid;not null;index:idx_model_branch,unique"`
	IsDefault                bool      `json:"is_default" default:"false"`
	// IsProtected bool      `json:"is_protected" default:"false"`

	Model Model `gorm:"foreignKey:ModelUUID"`

	Versions []ModelVersion `gorm:"foreignKey:BranchUUID"`
}

type ModelVersion struct {
	commondbmodels.BaseModel `gorm:"embedded"`
	Version                  string        `json:"version" gorm:"not null;index:idx_model_branch_version,unique"`
	Hash                     string        `json:"hash" gorm:"not null;index:idx_model_branch_hash,unique"`
	IsEmpty                  bool          `json:"is_empty"`
	BranchUUID               uuid.UUID     `json:"branch_uuid" gorm:"type:uuid;not null;index:idx_model_branch_version,unique;index:idx_model_branch_hash,unique"`
	PathUUID                 uuid.NullUUID `json:"path_uuid" gorm:"type:uuid;"`
	Path                     string        `json:"path"`
	SourceType               string        `json:"source_type"`
	CreatedBy                uuid.UUID     `json:"created_by" gorm:"type:uuid;"`

	Branch        ModelBranch          `gorm:"foreignKey:BranchUUID"`
	CreatedByUser userorgdbmodels.User `gorm:"foreignKey:CreatedBy"`
}

type ModelReview struct {
	commondbmodels.BaseModel `gorm:"embedded"`
	ModelUUID                uuid.UUID     `json:"model_uuid" gorm:"type:uuid;not null"`
	FromBranchUUID           uuid.UUID     `json:"from_branch_uuid" gorm:"type:uuid;not null"`
	FromBranchVersionUUID    uuid.UUID     `json:"from_branch_version_uuid" gorm:"type:uuid;not null"`
	ToBranchUUID             uuid.UUID     `json:"to_branch_uuid" gorm:"type:uuid;not null"`
	Title                    string        `json:"title" gorm:"not null"`
	Description              string        `json:"description"`
	CreatedBy                uuid.UUID     `json:"created_by" gorm:"type:uuid;not null"`
	AssignedTo               uuid.NullUUID `json:"assigned_to" gorm:"type:uuid;"`
	IsComplete               bool          `json:"is_complete" default:"false"`
	IsAccepted               bool          `json:"is_accepted" default:"false"`

	Model             Model                `gorm:"foreignKey:ModelUUID"`
	FromBranch        ModelBranch          `gorm:"foreignKey:FromBranchUUID"`
	FromBranchVersion ModelVersion         `gorm:"foreignKey:FromBranchVersionUUID"`
	ToBranch          ModelBranch          `gorm:"foreignKey:ToBranchUUID"`
	CreatedByUser     userorgdbmodels.User `gorm:"foreignKey:CreatedBy"`
	AssignedToUser    userorgdbmodels.User `gorm:"foreignKey:AssignedTo"`
}
