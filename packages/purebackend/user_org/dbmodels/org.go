package dbmodels

import (
	commondbmodels "github.com/PureMLHQ/PureML/packages/purebackend/core/common/dbmodels"
)

type Organization struct {
	commondbmodels.BaseModel `gorm:"embedded"`
	Name                     string `json:"name" gorm:"not null"`
	Handle                   string `json:"handle" gorm:"unique"`
	Avatar                   string `json:"avatar"`
	Description              string `json:"description"`
	APITokenHash             string `json:"api_token_hash"`
	JoinCode                 string `json:"join_code" gorm:"not null"`

	Users   []User   `gorm:"many2many:user_organizations;"` // many to many
	Secrets []Secret `gorm:"foreignKey:OrgUUID"`
}

type Secret struct {
	commondbmodels.BaseModel `gorm:"embedded"`
	OrgUUID                  string `json:"org_uuid" gorm:"not null;index:idx_org_name_key_secret,unique"`
	Name                     string `json:"name" gorm:"not null;index:idx_org_name_key_secret,unique"`
	Key                      string `json:"key" gorm:"not null;index:idx_org_name_key_secret,unique"`
	Value                    string `json:"value" gorm:"not null"`

	Org Organization `gorm:"foreignKey:OrgUUID"`
}
