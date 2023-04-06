package seeds

import (
	"fmt"

	authdbmodels "github.com/PureMLHQ/PureML/packages/purebackend/auth/dbmodels"
	commondbmodels "github.com/PureMLHQ/PureML/packages/purebackend/core/common/dbmodels"
	datasetdbmodels "github.com/PureMLHQ/PureML/packages/purebackend/dataset/dbmodels"
	modeldbmodels "github.com/PureMLHQ/PureML/packages/purebackend/model/dbmodels"
	userorgdbmodels "github.com/PureMLHQ/PureML/packages/purebackend/user_org/dbmodels"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var defaultUUID = uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))
var defaultUUID2 = uuid.Must(uuid.FromString("22222222-2222-2222-2222-222222222222"))

func All() []Seed {
	return []Seed{
		{
			Name: "CreateDemoAdminUserAndOrg",
			Run: func(d *gorm.DB) error {
				return CreateUser(d, defaultUUID, "Demo User", "demo@aztlan.in", "demo", "demo", "Demo User Bio", "")
			},
		},
		{
			Name: "CreateNonAdminUserAndOrg",
			Run: func(d *gorm.DB) error {
				return CreateUser(d, defaultUUID2, "Normal User", "notadmin@aztlan.in", "notadmin", "notadmin", "User Bio", "")
			},
		},
		{
			Name: "CreateDemoModel",
			Run: func(d *gorm.DB) error {
				return CreateModel(d, defaultUUID, "Demo Model", "Demo Model Wiki", true)
			},
		},
		{
			Name: "CreateDemoPrivateModel",
			Run: func(d *gorm.DB) error {
				return CreateModel(d, defaultUUID2, "Demo Private Model", "Demo Private Model Wiki", false)
			},
		},
		{
			Name: "CreateDemoDataset",
			Run: func(d *gorm.DB) error {
				return CreateDataset(d, defaultUUID, "Demo Dataset", "Demo Dataset Wiki", true)
			},
		},
		{
			Name: "CreateDemoPrivateDataset",
			Run: func(d *gorm.DB) error {
				return CreateDataset(d, defaultUUID2, "Demo Private Dataset", "Demo Private Dataset Wiki", false)
			},
		},
		{
			Name: "CreateDemoSession",
			Run: func(d *gorm.DB) error {
				return CreateSession(d, defaultUUID, defaultUUID, "device", "123456", "India")
			},
		},
	}
}

func CreateUser(db *gorm.DB, uuid uuid.UUID, name string, email string, handle string, password string, bio string, avatar string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	err = db.Create(&userorgdbmodels.User{
		BaseModel: commondbmodels.BaseModel{
			UUID: uuid,
		},
		Name:       name,
		Email:      email,
		Handle:     handle,
		Password:   string(hashedPassword),
		Bio:        bio,
		Avatar:     avatar,
		IsVerified: true,
		Orgs: []userorgdbmodels.Organization{
			{
				BaseModel: commondbmodels.BaseModel{
					UUID: uuid,
				},
				Name:        "Demo Org",
				Handle:      handle,
				Avatar:      "",
				Description: "Demo Org Description",
				JoinCode:    fmt.Sprintf("iwanttojoin%s", handle),
			},
		},
	}).Error
	if err != nil {
		return err
	}
	return db.Table("user_organizations").Where("user_uuid = ?", uuid).Where("organization_uuid = ?", uuid).Update("role", "owner").Error
}

func CreateModel(db *gorm.DB, uuid uuid.UUID, name string, wiki string, isPublic bool) error {
	var branches []modeldbmodels.ModelBranch
	var readme commondbmodels.Readme
	if isPublic {
		branches = []modeldbmodels.ModelBranch{
			{
				BaseModel: commondbmodels.BaseModel{
					UUID: defaultUUID,
				},
				Name:      "main",
				IsDefault: true,
			},
			{
				BaseModel: commondbmodels.BaseModel{
					UUID: defaultUUID2,
				},
				Name:      "dev",
				IsDefault: false,
				Versions: []modeldbmodels.ModelVersion{
					{
						BaseModel: commondbmodels.BaseModel{
							UUID: defaultUUID,
						},
						Version: "v1",
						Hash:    "1234567890",
						IsEmpty: true,
						CreatedByUser: userorgdbmodels.User{
							BaseModel: commondbmodels.BaseModel{
								UUID: defaultUUID,
							},
						},
					},
				},
			},
		}
		readme = commondbmodels.Readme{
			BaseModel: commondbmodels.BaseModel{
				UUID: defaultUUID,
			},
			ReadmeVersions: []commondbmodels.ReadmeVersion{
				{
					BaseModel: commondbmodels.BaseModel{
						UUID: defaultUUID,
					},
					Content:  "Demo Readme",
					FileType: "md",
					Version:  "v1",
				},
			},
		}
	}
	err := db.Create(&modeldbmodels.Model{
		BaseModel: commondbmodels.BaseModel{
			UUID: uuid,
		},
		Name: name,
		Wiki: wiki,
		Org: userorgdbmodels.Organization{
			BaseModel: commondbmodels.BaseModel{
				UUID: defaultUUID,
			},
		},
		CreatedByUser: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: defaultUUID,
			},
		},
		UpdatedByUser: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: defaultUUID,
			},
		},
		IsPublic: isPublic,
		Branches: branches,
		Readme:   readme,
		Users: []userorgdbmodels.User{
			{
				BaseModel: commondbmodels.BaseModel{
					UUID: defaultUUID,
				},
			},
		},
	}).Error
	if err != nil {
		return err
	}
	err = db.Table("model_users").Where("user_uuid = ?", defaultUUID).Where("model_uuid =  ?", uuid).Update("role", "owner").Error
	if err != nil {
		return err
	}
	return nil
}

func CreateDataset(db *gorm.DB, uuid uuid.UUID, name string, wiki string, isPublic bool) error {
	var branches []datasetdbmodels.DatasetBranch
	var readme commondbmodels.Readme
	if isPublic {
		branches = []datasetdbmodels.DatasetBranch{
			{
				BaseModel: commondbmodels.BaseModel{
					UUID: defaultUUID,
				},
				Name:      "main",
				IsDefault: true,
			},
			{
				BaseModel: commondbmodels.BaseModel{
					UUID: defaultUUID2,
				},
				Name:      "dev",
				IsDefault: false,
				Versions: []datasetdbmodels.DatasetVersion{
					{
						BaseModel: commondbmodels.BaseModel{
							UUID: defaultUUID,
						},
						Version: "v1",
						Hash:    "1234567890",
						IsEmpty: true,
						Lineage: datasetdbmodels.Lineage{
							BaseModel: commondbmodels.BaseModel{
								UUID: defaultUUID,
							},
							Lineage: "{}",
						},
						CreatedByUser: userorgdbmodels.User{
							BaseModel: commondbmodels.BaseModel{
								UUID: defaultUUID,
							},
						},
					},
				},
			},
		}
		readme = commondbmodels.Readme{
			BaseModel: commondbmodels.BaseModel{
				UUID: defaultUUID,
			},
			ReadmeVersions: []commondbmodels.ReadmeVersion{
				{
					BaseModel: commondbmodels.BaseModel{
						UUID: defaultUUID,
					},
					Content:  "Demo Readme",
					FileType: "md",
					Version:  "v1",
				},
			},
		}
	}
	err := db.Create(&datasetdbmodels.Dataset{
		BaseModel: commondbmodels.BaseModel{
			UUID: uuid,
		},
		Name: name,
		Wiki: wiki,
		Org: userorgdbmodels.Organization{
			BaseModel: commondbmodels.BaseModel{
				UUID: defaultUUID,
			},
		},
		CreatedByUser: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: defaultUUID,
			},
		},
		UpdatedByUser: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: defaultUUID,
			},
		},
		IsPublic: isPublic,
		Branches: branches,
		Readme:   readme,
		Users: []userorgdbmodels.User{
			{
				BaseModel: commondbmodels.BaseModel{
					UUID: defaultUUID,
				},
			},
		},
	}).Error
	if err != nil {
		return err
	}
	err = db.Table("dataset_users").Where("user_uuid = ?", defaultUUID).Where("dataset_uuid =  ?", uuid).Update("role", "owner").Error
	if err != nil {
		return err
	}
	return nil
}

func CreateSession(db *gorm.DB, uuid uuid.UUID, userUUID uuid.UUID, device string, deviceId string, deviceLocation string) error {
	err := db.Create(&authdbmodels.Session{
		BaseModel: commondbmodels.BaseModel{
			UUID: uuid,
		},
		User: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: userUUID,
			},
		},
		Device:         device,
		DeviceId:       deviceId,
		DeviceLocation: deviceLocation,
	}).Error
	if err != nil {
		return err
	}
	return nil
}
