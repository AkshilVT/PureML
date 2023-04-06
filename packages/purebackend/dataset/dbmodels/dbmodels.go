package dbmodels

import (
	commondbmodels "github.com/PureMLHQ/PureML/packages/purebackend/core/common/dbmodels"
	userorgdbmodels "github.com/PureMLHQ/PureML/packages/purebackend/user_org/dbmodels"
	uuid "github.com/satori/go.uuid"
)

type Dataset struct {
	commondbmodels.BaseModel `gorm:"embedded"`
	Name                     string        `json:"name" gorm:"not null;index:idx_org_dataset_name,unique"`
	Wiki                     string        `json:"wiki"`
	OrganizationUUID         uuid.UUID     `json:"org_uuid" gorm:"type:uuid;not null;index:idx_org_dataset_name,unique"`
	CreatedBy                uuid.UUID     `json:"created_by" gorm:"type:uuid;not null"`
	UpdatedBy                uuid.NullUUID `json:"updated_by" gorm:"type:uuid;"`
	IsPublic                 bool          `json:"is_public" default:"false"`

	Org           userorgdbmodels.Organization `gorm:"foreignKey:OrganizationUUID"`
	CreatedByUser userorgdbmodels.User         `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser userorgdbmodels.User         `gorm:"foreignKey:UpdatedBy"`

	Readme   commondbmodels.Readme  `gorm:"foreignKey:DatasetUUID"`
	Branches []DatasetBranch        `gorm:"foreignKey:DatasetUUID"`
	Users    []userorgdbmodels.User `gorm:"many2many:dataset_users;"`
}

type DatasetUser struct {
	DatasetUUID uuid.UUID `json:"dataset_uuid" gorm:"type:uuid;primaryKey"`
	UserUUID    uuid.UUID `json:"user_uuid" gorm:"type:uuid;primaryKey"`
	Role        string    `json:"role" gorm:"not null;default:member"`
}

type DatasetBranch struct {
	commondbmodels.BaseModel `gorm:"embedded"`
	Name                     string    `json:"name" gorm:"not null;index:idx_dataset_branch,unique"`
	DatasetUUID              uuid.UUID `json:"dataset_uuid" gorm:"type:uuid;not null;index:idx_dataset_branch,unique"`
	IsDefault                bool      `json:"is_default" default:"false"`

	Dataset Dataset `gorm:"foreignKey:DatasetUUID"`

	Versions []DatasetVersion `gorm:"foreignKey:BranchUUID"`
}

type DatasetVersion struct {
	commondbmodels.BaseModel `gorm:"embedded"`
	Version                  string        `json:"version" gorm:"not null;index:idx_dataset_branch_version,unique"`
	Hash                     string        `json:"hash" gorm:"not null;index:idx_dataset_branch_hash,unique"`
	IsEmpty                  bool          `json:"is_empty"`
	BranchUUID               uuid.UUID     `json:"branch_uuid" gorm:"type:uuid;not null;index:idx_dataset_branch_version,unique;index:idx_dataset_branch_hash,unique"`
	LineageUUID              uuid.NullUUID `json:"lineage_uuid" gorm:"type:uuid;"`
	PathUUID                 uuid.NullUUID `json:"path_uuid" gorm:"type:uuid;"`
	Path                     string        `json:"path"`
	SourceType               string        `json:"source_type"`
	CreatedBy                uuid.UUID     `json:"created_by" gorm:"type:uuid;"`

	Branch        DatasetBranch        `gorm:"foreignKey:BranchUUID"`
	Lineage       Lineage              `gorm:"foreignKey:LineageUUID"`
	CreatedByUser userorgdbmodels.User `gorm:"foreignKey:CreatedBy"`
}

type Lineage struct {
	commondbmodels.BaseModel `gorm:"embedded"`
	Lineage                  string `json:"lineage"`
}

type DatasetReview struct {
	commondbmodels.BaseModel `gorm:"embedded"`
	DatasetUUID              uuid.UUID     `json:"dataset_uuid" gorm:"type:uuid;not null"`
	FromBranchUUID           uuid.UUID     `json:"from_branch_uuid" gorm:"type:uuid;not null"`
	FromBranchVersionUUID    uuid.UUID     `json:"from_branch_version_uuid" gorm:"type:uuid;not null"`
	ToBranchUUID             uuid.UUID     `json:"to_branch_uuid" gorm:"type:uuid;not null"`
	Title                    string        `json:"title" gorm:"not null"`
	Description              string        `json:"description"`
	CreatedBy                uuid.UUID     `json:"created_by" gorm:"type:uuid;not null"`
	AssignedTo               uuid.NullUUID `json:"assigned_to" gorm:"type:uuid;"`
	IsComplete               bool          `json:"is_complete" default:"false"`
	IsAccepted               bool          `json:"is_accepted" default:"false"`

	Dataset           Dataset              `gorm:"foreignKey:DatasetUUID"`
	FromBranch        DatasetBranch        `gorm:"foreignKey:FromBranchUUID"`
	FromBranchVersion DatasetVersion       `gorm:"foreignKey:FromBranchVersionUUID"`
	ToBranch          DatasetBranch        `gorm:"foreignKey:ToBranchUUID"`
	CreatedByUser     userorgdbmodels.User `gorm:"foreignKey:CreatedBy"`
	AssignedToUser    userorgdbmodels.User `gorm:"foreignKey:AssignedTo"`
}
