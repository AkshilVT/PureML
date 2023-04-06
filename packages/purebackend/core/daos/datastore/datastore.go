package impl

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	authdbmodels "github.com/PureMLHQ/PureML/packages/purebackend/auth/dbmodels"
	authmodels "github.com/PureMLHQ/PureML/packages/purebackend/auth/models"
	commondbmodels "github.com/PureMLHQ/PureML/packages/purebackend/core/common/dbmodels"
	commonmodels "github.com/PureMLHQ/PureML/packages/purebackend/core/common/models"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/config"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/dbmodels"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/models"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/settings"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/tools/search"
	datasetdbmodels "github.com/PureMLHQ/PureML/packages/purebackend/dataset/dbmodels"
	datasetmodels "github.com/PureMLHQ/PureML/packages/purebackend/dataset/models"
	modeldbmodels "github.com/PureMLHQ/PureML/packages/purebackend/model/dbmodels"
	modelmodels "github.com/PureMLHQ/PureML/packages/purebackend/model/models"
	userorgdbmodels "github.com/PureMLHQ/PureML/packages/purebackend/user_org/dbmodels"
	userorgmodels "github.com/PureMLHQ/PureML/packages/purebackend/user_org/models"
	puregosqlite "github.com/glebarez/sqlite"
	uuid "github.com/satori/go.uuid"
	"github.com/teris-io/shortid"
	"gorm.io/driver/postgres"
	cgosqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var gormConfig = &gorm.Config{
	Logger: logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second * 2, // Slow SQL threshold
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	),
}

func NewSQLiteDatastore(optDataDir ...string) *Datastore {
	var dialector gorm.Dialector
	var err error
	var dataDir string
	if len(optDataDir) == 0 || optDataDir[0] == "" {
		// fallback to the default test data directory
		dataDir = config.GetDataDir()
	} else {
		dataDir = optDataDir[0]
	}
	databasePath := fmt.Sprintf("%s/pureml.db", dataDir)
	if config.IsCGOEnabled() {
		dialector = cgosqlite.Open(databasePath)
	} else {
		dialector = puregosqlite.Open(databasePath)
	}
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		fmt.Println(err)
		panic("Error connecting to database")
	}
	err = db.AutoMigrate(
		dbmodels.Activity{},
		modeldbmodels.Model{},
		datasetdbmodels.Dataset{},
		datasetdbmodels.DatasetBranch{},
		datasetdbmodels.DatasetReview{},
		datasetdbmodels.DatasetUser{},
		datasetdbmodels.DatasetVersion{},
		datasetdbmodels.Lineage{},
		dbmodels.Log{},
		// dbmodels.Tag{},
		datasetdbmodels.Dataset{},
		modeldbmodels.ModelBranch{},
		modeldbmodels.ModelReview{},
		modeldbmodels.ModelUser{},
		modeldbmodels.ModelVersion{},
		userorgdbmodels.Organization{},
		userorgdbmodels.User{},
		userorgdbmodels.UserOrganizations{},
		userorgdbmodels.Secret{},
		commondbmodels.Readme{},
		commondbmodels.ReadmeVersion{},
		authdbmodels.Session{},
	)
	if err != nil {
		return &Datastore{}
	}
	return &Datastore{
		DB: db,
	}
}

func NewPostgresDatastore(databaseUrl string) *Datastore {
	dsn := databaseUrl
	// fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		fmt.Println(err)
		panic("Error connecting to database")
	}
	err = db.AutoMigrate(
		dbmodels.Activity{},
		modeldbmodels.Model{},
		datasetdbmodels.Dataset{},
		datasetdbmodels.DatasetBranch{},
		datasetdbmodels.DatasetReview{},
		datasetdbmodels.DatasetUser{},
		datasetdbmodels.DatasetVersion{},
		datasetdbmodels.Lineage{},
		dbmodels.Log{},
		// dbmodels.Tag{},
		datasetdbmodels.Dataset{},
		modeldbmodels.ModelBranch{},
		modeldbmodels.ModelReview{},
		modeldbmodels.ModelUser{},
		modeldbmodels.ModelVersion{},
		userorgdbmodels.Organization{},
		userorgdbmodels.User{},
		userorgdbmodels.UserOrganizations{},
		userorgdbmodels.Secret{},
		commondbmodels.Readme{},
		commondbmodels.ReadmeVersion{},
		authdbmodels.Session{},
	)
	if err != nil {
		return &Datastore{}
	}
	return &Datastore{
		DB: db,
	}
}

func (ds *Datastore) SeedSearchClient() {
	// Query database
	db, err := ds.DB.DB()
	if err != nil {
		fmt.Println(err)
	}
	// Users
	users, err := db.Query("SELECT uuid, name, email, handle FROM users")
	if err != nil {
		fmt.Println(err)
	}
	var userDocs []interface{}
	for users.Next() {
		var myuuid uuid.UUID
		var name string
		var email string
		var handle string

		err = users.Scan(&myuuid, &name, &email, &handle)
		if err != nil {
			fmt.Println(err)
		}
		userDocs = append(userDocs, map[string]interface{}{
			"uuid":   myuuid,
			"name":   name,
			"email":  email,
			"handle": handle,
		})
	}

	// Organizations
	orgs, err := db.Query("SELECT uuid, name, handle, description FROM organizations")
	if err != nil {
		fmt.Println(err)
	}
	var orgDocs []interface{}
	for orgs.Next() {
		var myuuid uuid.UUID
		var name string
		var handle string
		var description string

		err = orgs.Scan(&myuuid, &name, &handle, &description)
		if err != nil {
			fmt.Println(err)
		}
		orgDocs = append(orgDocs, map[string]interface{}{
			"uuid":        myuuid,
			"name":        name,
			"handle":      handle,
			"description": description,
		})
	}

	// Models
	models, err := db.Query("SELECT uuid, name, wiki, organization_uuid, is_public FROM models")
	if err != nil {
		fmt.Println(err)
	}
	var modelDocs []interface{}
	for models.Next() {
		var myuuid uuid.UUID
		var name string
		var wiki string
		var organizationUuid uuid.UUID
		var is_public bool

		err = models.Scan(&myuuid, &name, &wiki, &organizationUuid, &is_public)
		if err != nil {
			fmt.Println(err)
		}
		modelDocs = append(modelDocs, map[string]interface{}{
			"uuid":              myuuid,
			"name":              name,
			"wiki":              wiki,
			"organization_uuid": organizationUuid,
			"is_public":         is_public,
		})
	}

	// Datasets
	datasets, err := db.Query("SELECT uuid, name, wiki, organization_uuid, is_public FROM datasets")
	if err != nil {
		fmt.Println(err)
	}
	var datasetDocs []interface{}
	for datasets.Next() {
		var myuuid uuid.UUID
		var name string
		var wiki string
		var organizationUuid uuid.UUID
		var is_public bool

		err = datasets.Scan(&myuuid, &name, &wiki, &organizationUuid, &is_public)
		if err != nil {
			fmt.Println(err)
		}
		datasetDocs = append(datasetDocs, map[string]interface{}{
			"uuid":              myuuid,
			"name":              name,
			"wiki":              wiki,
			"organization_uuid": organizationUuid,
			"is_public":         is_public,
		})
	}

	// Add documents to index
	_, err = ds.SearchClient.Client.Index("users").AddDocuments(userDocs, "uuid")
	if err != nil {
		panic(err)
	}
	_, err = ds.SearchClient.Client.Index("organizations").AddDocuments(orgDocs, "uuid")
	if err != nil {
		panic(err)
	}
	_, err = ds.SearchClient.Client.Index("models").AddDocuments(modelDocs, "uuid")
	if err != nil {
		panic(err)
	}
	_, err = ds.SearchClient.Client.Index("datasets").AddDocuments(datasetDocs, "uuid")
	if err != nil {
		panic(err)
	}
}

func (ds *Datastore) SeedAdminIfNotExists() error {
	var user userorgdbmodels.User
	adminDetails := config.GetAdminDetails()
	if adminDetails == nil {
		// return errors.New("admin details not found in env vars")
		// fmt.Println("admin details not found in env vars. using default values")
		adminDetails = map[string]interface{}{
			"uuid":       uuid.Must(uuid.FromString("99999999-9999-9999-9999-999999999999")),
			"email":      "admin@admin.com",
			"password":   "admin",
			"handle":     "admin",
			"org_name":   "admin-org",
			"org_handle": "adminorg",
		}
	}
	err := ds.DB.Where("uuid = ?", adminDetails["uuid"]).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		// admin user does not exist, create it
		adminUser := userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: adminDetails["uuid"].(uuid.UUID),
			},
			Email:    adminDetails["email"].(string),
			Password: adminDetails["password"].(string),
			Name:     adminDetails["handle"].(string),
			Handle:   adminDetails["handle"].(string),
			Orgs: []userorgdbmodels.Organization{
				{
					BaseModel: commondbmodels.BaseModel{
						UUID: adminDetails["uuid"].(uuid.UUID),
					},
					Name:     adminDetails["org_name"].(string),
					Handle:   adminDetails["org_handle"].(string),
					JoinCode: "",
				},
			},
		}
		ds.DB.Create(&adminUser)
		var userOrganization userorgdbmodels.UserOrganizations
		ds.DB.Where("user_uuid = ? AND organization_uuid = ?", adminUser.UUID, adminUser.UUID).First(&userOrganization)
		userOrganization.Role = "owner"
		ds.DB.Save(&userOrganization)
	} else if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (ds *Datastore) SeedAdminS3SecretsIfNotExists(s3config *settings.S3Config) error {
	adminDetails := config.GetAdminDetails()
	if adminDetails == nil {
		return errors.New("admin details not found in env vars")
	}
	if _, err := ds.CreateS3Secrets(adminDetails["uuid"].(uuid.UUID), "admin", s3config.AccessKey, s3config.Secret, s3config.Bucket, s3config.Region); err != nil {
		// secrets already exist
		return nil
	}
	return nil
}

func (ds *Datastore) SeedAdminR2SecretsIfNotExists(r2config *settings.R2Config) error {
	adminDetails := config.GetAdminDetails()
	if adminDetails == nil {
		return errors.New("admin details not found in env vars")
	}
	if _, err := ds.CreateR2Secrets(adminDetails["uuid"].(uuid.UUID), "admin", r2config.AccountId, r2config.AccessKey, r2config.Secret, r2config.Bucket, r2config.Endpoint); err != nil {
		// secrets already exist
		return nil
	}
	return nil
}

type Datastore struct {
	DB           *gorm.DB
	SearchClient *search.SearchClient
}

func (ds *Datastore) ExecuteSQL(sql string) error {
	return ds.DB.Exec(sql).Error
}

func (ds *Datastore) Close() error {
	sqlDB, err := ds.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

//////////////////////////// ORGANIZATION METHODS ////////////////////////////

func (ds *Datastore) GetAllAdminOrgs() ([]userorgmodels.OrganizationResponse, error) {
	var organizations []userorgdbmodels.Organization
	ds.DB.Find(&organizations)
	var responseOrganizations []userorgmodels.OrganizationResponse
	for _, org := range organizations {
		responseOrganizations = append(responseOrganizations, userorgmodels.OrganizationResponse{
			UUID:        org.UUID,
			Name:        org.Name,
			Handle:      org.Handle,
			Avatar:      org.Avatar,
			Description: org.Description,
			JoinCode:    org.JoinCode,
		})
	}
	return responseOrganizations, nil
}

func (ds *Datastore) GetOrgByID(orgId uuid.UUID) (*userorgmodels.OrganizationResponseWithMembers, error) {
	org := userorgdbmodels.Organization{
		BaseModel: commondbmodels.BaseModel{
			UUID: orgId,
		},
	}
	result := ds.DB.Limit(1).Preload("Users").Find(&org)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	var dbuserRoles []userorgdbmodels.UserOrganizations
	err := ds.DB.Where("organization_uuid = ?", orgId).Find(&dbuserRoles).Error
	if err != nil {
		return nil, err
	}
	userRoles := make(map[uuid.UUID]string)
	for userRole := range dbuserRoles {
		userRoles[dbuserRoles[userRole].UserUUID] = dbuserRoles[userRole].Role
	}
	var members []userorgmodels.UserHandleRoleResponse
	for _, user := range org.Users {
		members = append(members, userorgmodels.UserHandleRoleResponse{
			UUID:   user.UUID,
			Handle: user.Handle,
			Name:   user.Name,
			Avatar: user.Avatar,
			Email:  user.Email,
			Role:   userRoles[user.UUID],
		})
	}
	return &userorgmodels.OrganizationResponseWithMembers{
		UUID:        org.UUID,
		Name:        org.Name,
		Handle:      org.Handle,
		Avatar:      org.Avatar,
		Description: org.Description,

		Members: members,
	}, nil
}

func (ds *Datastore) GetOrgByJoinCode(joinCode string) (*userorgmodels.OrganizationResponse, error) {
	var org userorgdbmodels.Organization
	result := ds.DB.Where("join_code = ?", joinCode).Limit(1).Find(&org)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &userorgmodels.OrganizationResponse{
		UUID:        org.UUID,
		Name:        org.Name,
		Handle:      org.Handle,
		Avatar:      org.Avatar,
		Description: org.Description,
		JoinCode:    org.JoinCode,
	}, nil
}

func (ds *Datastore) CreateOrgFromEmail(email string, orgName string, orgDesc string, orgHandle string) (*userorgmodels.OrganizationResponse, error) {
	org := userorgdbmodels.Organization{
		Name:         orgName,
		Handle:       orgHandle,
		Avatar:       "",
		Description:  orgDesc,
		JoinCode:     shortid.MustGenerate(),
		APITokenHash: "",
	}
	user := userorgdbmodels.User{
		Email: email,
	}
	err := ds.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&org)
		if result.Error != nil {
			return result.Error
		}
		result = tx.Where("email = ?", email).First(&user)
		if result.Error != nil {
			return result.Error
		}
		userOrg := userorgdbmodels.UserOrganizations{
			UserUUID:         user.UUID,
			OrganizationUUID: org.UUID,
			Role:             "owner",
		}
		result = tx.Create(&userOrg)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if ds.SearchClient != nil {
		err := ds.SearchClient.AddDocument("organizations", map[string]interface{}{
			"uuid":        org.UUID,
			"name":        org.Name,
			"handle":      org.Handle,
			"description": org.Description,
		})
		if err != nil {
			return nil, err
		}
	}
	return &userorgmodels.OrganizationResponse{
		UUID:        org.UUID,
		Name:        org.Name,
		Handle:      org.Handle,
		Avatar:      org.Avatar,
		Description: org.Description,
		JoinCode:    org.JoinCode,
	}, nil
}

func (ds *Datastore) GetOrgByHandle(handle string) (*userorgmodels.OrganizationResponse, error) {
	var org userorgdbmodels.Organization
	result := ds.DB.Where("handle = ?", handle).Limit(1).Find(&org)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &userorgmodels.OrganizationResponse{
		UUID:        org.UUID,
		Name:        org.Name,
		Handle:      org.Handle,
		Avatar:      org.Avatar,
		Description: org.Description,
		// JoinCode:    org.JoinCode,
	}, nil
}

func (ds *Datastore) GetUserOrganizationsByEmail(email string) ([]userorgmodels.UserOrganizationsResponse, error) {
	var orgs []userorgmodels.UserOrganizationsResponse
	var tableOrgs []struct {
		UUID        uuid.UUID
		Handle      string
		Name        string
		Avatar      string
		Description string
		Role        string
	}
	result := ds.DB.Table("organizations").Select("organizations.uuid, organizations.handle, organizations.name, organizations.avatar, organizations.description, user_organizations.role").Joins("JOIN user_organizations ON user_organizations.organization_uuid = organizations.uuid").Joins("JOIN users ON users.uuid = user_organizations.user_uuid").Where("users.email = ?", email).Scan(&tableOrgs)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, org := range tableOrgs {
		orgs = append(orgs, userorgmodels.UserOrganizationsResponse{
			Org: userorgmodels.OrganizationHandleResponse{
				UUID:        org.UUID,
				Handle:      org.Handle,
				Name:        org.Name,
				Avatar:      org.Avatar,
				Description: org.Description,
			},
			Role: org.Role,
		})
	}
	return orgs, nil
}

func (ds *Datastore) GetUserOrganizationByOrgIdAndUserUUID(orgId uuid.UUID, userUUID uuid.UUID) (*userorgmodels.UserOrganizationsRoleResponse, error) {
	userOrganizations := userorgdbmodels.UserOrganizations{
		UserUUID:         userUUID,
		OrganizationUUID: orgId,
	}
	result := ds.DB.Limit(1).Find(&userOrganizations)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userOrgResponse := userorgmodels.UserOrganizationsRoleResponse{
		UserUUID: userOrganizations.UserUUID,
		OrgUUID:  userOrganizations.OrganizationUUID,
		Role:     userOrganizations.Role,
	}
	return &userOrgResponse, nil
}

func (ds *Datastore) CreateUserOrganizationFromEmailAndOrgId(email string, orgId uuid.UUID) (*userorgmodels.UserOrganizationsResponse, error) {
	var org userorgdbmodels.Organization
	result := ds.DB.First(&org, orgId)
	if result.Error != nil {
		return nil, result.Error
	}
	var user userorgdbmodels.User
	result = ds.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	userOrganization := userorgdbmodels.UserOrganizations{
		OrganizationUUID: org.UUID,
		UserUUID:         user.UUID,
		Role:             "member",
	}
	result = ds.DB.Create(&userOrganization)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userorgmodels.UserOrganizationsResponse{
		Org: userorgmodels.OrganizationHandleResponse{
			UUID:        org.UUID,
			Name:        org.Name,
			Handle:      org.Handle,
			Avatar:      org.Avatar,
			Description: org.Description,
		},
		Role: userOrganization.Role,
	}, nil
}

func (ds *Datastore) DeleteUserOrganizationFromEmailAndOrgId(email string, orgId uuid.UUID) error {
	var org userorgdbmodels.Organization
	result := ds.DB.First(&org, orgId)
	if result.Error != nil {
		return result.Error
	}
	var user userorgdbmodels.User
	result = ds.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return result.Error
	}
	result = ds.DB.Where("organization_uuid = ?", org.UUID).Where("user_uuid = ?", user.UUID).Delete(&userorgdbmodels.UserOrganizations{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ds *Datastore) CreateUserOrganizationFromEmailAndJoinCode(email string, joinCode string) (*userorgmodels.UserOrganizationsResponse, error) {
	var org userorgdbmodels.Organization
	result := ds.DB.Where("join_code = ?", joinCode).First(&org)
	if result.Error != nil {
		return nil, result.Error
	}
	var user userorgdbmodels.User
	result = ds.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	userOrganization := userorgdbmodels.UserOrganizations{
		OrganizationUUID: org.UUID,
		UserUUID:         user.UUID,
		Role:             "member",
	}
	result = ds.DB.Create(&userOrganization)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userorgmodels.UserOrganizationsResponse{
		Org: userorgmodels.OrganizationHandleResponse{
			UUID:        org.UUID,
			Name:        org.Name,
			Handle:      org.Handle,
			Avatar:      org.Avatar,
			Description: org.Description,
		},
		Role: userOrganization.Role,
	}, nil
}

func (ds *Datastore) UpdateUserRoleByOrgIdAndUserUUID(orgId uuid.UUID, userUUID uuid.UUID, role string) error {
	userOrganizations := userorgdbmodels.UserOrganizations{
		UserUUID:         userUUID,
		OrganizationUUID: orgId,
	}
	result := ds.DB.Limit(1).Find(&userOrganizations)
	if result.RowsAffected == 0 {
		return nil
	}
	if result.Error != nil {
		return result.Error
	}
	userOrganizations.Role = role
	result = ds.DB.Save(&userOrganizations)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ds *Datastore) UpdateOrg(orgId uuid.UUID, updatedAttributes map[string]interface{}) (*userorgmodels.OrganizationResponse, error) {
	var org userorgdbmodels.Organization
	result := ds.DB.First(&org, orgId)
	if result.Error != nil {
		return nil, result.Error
	}
	if updatedAttributes["name"] != nil {
		org.Name = updatedAttributes["name"].(string)
	}
	if updatedAttributes["desc"] != nil {
		org.Description = updatedAttributes["desc"].(string)
	}
	if updatedAttributes["avatar"] != nil {
		org.Avatar = updatedAttributes["avatar"].(string)
	}
	result = ds.DB.Save(&org)
	if result.Error != nil {
		return nil, result.Error
	}
	if ds.SearchClient != nil {
		err := ds.SearchClient.UpdateDocument("organizations", orgId.String(), map[string]interface{}{
			"uuid":        org.UUID,
			"name":        org.Name,
			"handle":      org.Handle,
			"description": org.Description,
		})
		if err != nil {
			return nil, err
		}
	}
	return &userorgmodels.OrganizationResponse{
		UUID:        org.UUID,
		Name:        org.Name,
		Handle:      org.Handle,
		Avatar:      org.Avatar,
		Description: org.Description,
		JoinCode:    org.JoinCode,
	}, nil
}

func (ds *Datastore) GetOrgAllPublicModels(orgId uuid.UUID) ([]modelmodels.ModelResponse, error) {
	var modelsdb []*modeldbmodels.Model
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Where("is_public = ?", true).Where("organization_uuid = ?", orgId).Find(&modelsdb)
	if result.Error != nil {
		return nil, result.Error
	}
	var returnModels []modelmodels.ModelResponse
	for _, model := range modelsdb {
		returnModels = append(returnModels, modelmodels.ModelResponse{
			UUID:     model.UUID,
			Name:     model.Name,
			Wiki:     model.Wiki,
			IsPublic: model.IsPublic,
			CreatedBy: userorgmodels.UserHandleResponse{
				UUID:   model.CreatedByUser.UUID,
				Handle: model.CreatedByUser.Handle,
				Avatar: model.CreatedByUser.Avatar,
				Name:   model.CreatedByUser.Name,
				Email:  model.CreatedByUser.Email,
			},
			UpdatedBy: userorgmodels.UserHandleResponse{
				UUID:   model.UpdatedByUser.UUID,
				Handle: model.UpdatedByUser.Handle,
				Avatar: model.UpdatedByUser.Avatar,
				Name:   model.UpdatedByUser.Name,
				Email:  model.UpdatedByUser.Email,
			},
		})
	}
	return returnModels, nil
}

func (ds *Datastore) GetOrgAllPublicDatasets(orgId uuid.UUID) ([]datasetmodels.DatasetResponse, error) {
	var datasetsdb []*datasetdbmodels.Dataset
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Where("is_public = ?", true).Where("organization_uuid = ?", orgId).Find(&datasetsdb)
	if result.Error != nil {
		return nil, result.Error
	}
	var returnDatasets []datasetmodels.DatasetResponse
	for _, dataset := range datasetsdb {
		returnDatasets = append(returnDatasets, datasetmodels.DatasetResponse{
			UUID:     dataset.UUID,
			Name:     dataset.Name,
			Wiki:     dataset.Wiki,
			IsPublic: dataset.IsPublic,
			CreatedBy: userorgmodels.UserHandleResponse{
				UUID:   dataset.CreatedByUser.UUID,
				Handle: dataset.CreatedByUser.Handle,
				Avatar: dataset.CreatedByUser.Avatar,
				Name:   dataset.CreatedByUser.Name,
				Email:  dataset.CreatedByUser.Email,
			},
			UpdatedBy: userorgmodels.UserHandleResponse{
				UUID:   dataset.UpdatedByUser.UUID,
				Handle: dataset.UpdatedByUser.Handle,
				Avatar: dataset.UpdatedByUser.Avatar,
				Name:   dataset.UpdatedByUser.Name,
				Email:  dataset.UpdatedByUser.Email,
			},
		})
	}
	return returnDatasets, nil
}

/////////////////////////////// USER METHODS /////////////////////////////////

func (ds *Datastore) GetUserByEmail(email string) (*userorgmodels.UserResponse, error) {
	var user userorgdbmodels.User
	result := ds.DB.Where("email = ?", email).Limit(1).Find(&user)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &userorgmodels.UserResponse{
		UUID:       user.UUID,
		Name:       user.Name,
		Email:      user.Email,
		Handle:     user.Handle,
		Bio:        user.Bio,
		Avatar:     user.Avatar,
		IsVerified: user.IsVerified,
	}, nil
}

func (ds *Datastore) GetUserByHandle(handle string) (*userorgmodels.UserProfileResponse, error) {
	var user userorgdbmodels.User
	result := ds.DB.Where("handle = ?", handle).Limit(1).Find(&user)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	numberOfDatasets := int64(0)
	ds.DB.Model(&datasetdbmodels.DatasetUser{}).Where("user_uuid = ?", user.BaseModel.UUID).Count(&numberOfDatasets)
	numberOfModel := int64(0)
	ds.DB.Model(&modeldbmodels.ModelUser{}).Where("user_uuid = ?", user.BaseModel.UUID).Count(&numberOfModel)
	return &userorgmodels.UserProfileResponse{
		Name:             user.Name,
		Email:            user.Email,
		Handle:           user.Handle,
		Bio:              user.Bio,
		Avatar:           user.Avatar,
		NumberOfModels:   numberOfModel,
		NumberOfDatasets: numberOfDatasets,
	}, nil
}

func (ds *Datastore) GetSecureUserByEmail(email string) (*userorgmodels.UserResponse, error) {
	var user userorgdbmodels.User
	result := ds.DB.Where("email = ?", email).Limit(1).Find(&user)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &userorgmodels.UserResponse{
		UUID:       user.UUID,
		Name:       user.Name,
		Email:      user.Email,
		Handle:     user.Handle,
		Bio:        user.Bio,
		Avatar:     user.Avatar,
		Password:   user.Password,
		IsVerified: user.IsVerified,
	}, nil
}

func (ds *Datastore) GetSecureUserByHandle(handle string) (*userorgmodels.UserResponse, error) {
	var user userorgdbmodels.User
	result := ds.DB.Where("handle = ?", handle).Limit(1).Find(&user)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &userorgmodels.UserResponse{
		UUID:       user.UUID,
		Name:       user.Name,
		Email:      user.Email,
		Handle:     user.Handle,
		Bio:        user.Bio,
		Avatar:     user.Avatar,
		Password:   user.Password,
		IsVerified: user.IsVerified,
	}, nil
}

func (ds *Datastore) GetSecureUserByUUID(userUUID uuid.UUID) (*userorgmodels.UserResponse, error) {
	var user userorgdbmodels.User
	result := ds.DB.Where("uuid = ?", userUUID).Limit(1).Find(&user)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &userorgmodels.UserResponse{
		UUID:       user.UUID,
		Name:       user.Name,
		Email:      user.Email,
		Handle:     user.Handle,
		Bio:        user.Bio,
		Avatar:     user.Avatar,
		Password:   user.Password,
		IsVerified: user.IsVerified,
	}, nil
}

func (ds *Datastore) GetUserByUUID(userUUID uuid.UUID) (*userorgmodels.UserResponse, error) {
	var user userorgdbmodels.User
	result := ds.DB.Limit(1).Find(&user, userUUID)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &userorgmodels.UserResponse{
		UUID:       user.UUID,
		Name:       user.Name,
		Email:      user.Email,
		Handle:     user.Handle,
		Bio:        user.Bio,
		Avatar:     user.Avatar,
		IsVerified: user.IsVerified,
	}, nil
}

func (ds *Datastore) GetUserProfileByUUID(userUUID uuid.UUID) (*userorgmodels.UserProfileResponse, error) {
	var user userorgdbmodels.User
	result := ds.DB.Limit(1).Find(&user, userUUID)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	numberOfDatasets := int64(0)
	ds.DB.Model(&datasetdbmodels.DatasetUser{}).Where("user_uuid = ?", userUUID).Count(&numberOfDatasets)
	numberOfModel := int64(0)
	ds.DB.Model(&modeldbmodels.ModelUser{}).Where("user_uuid = ?", userUUID).Count(&numberOfModel)
	return &userorgmodels.UserProfileResponse{
		UUID:             user.UUID,
		Name:             user.Name,
		Email:            user.Email,
		Handle:           user.Handle,
		Bio:              user.Bio,
		Avatar:           user.Avatar,
		NumberOfModels:   numberOfModel,
		NumberOfDatasets: numberOfDatasets,
	}, nil
}

func (ds *Datastore) CreateUser(name string, email string, handle string, bio string, avatar string, hashedPassword string, isVerified bool) (*userorgmodels.UserResponse, error) {
	user := userorgdbmodels.User{
		Name:       name,
		Email:      email,
		Password:   hashedPassword,
		Handle:     handle,
		Bio:        bio,
		Avatar:     avatar,
		IsVerified: isVerified,

		Orgs: []userorgdbmodels.Organization{
			{
				Name:        "Private",
				Handle:      handle,
				Avatar:      avatar,
				JoinCode:    shortid.MustGenerate(),
				Description: fmt.Sprintf("Private Organization for %s", handle),
			},
		},
	}
	err := ds.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&user)
		if result.Error != nil {
			return result.Error
		}
		result = tx.Table("user_organizations").Where("user_uuid = ?", user.UUID).Where("organization_uuid = ?", user.Orgs[0].UUID).Update("role", "owner")
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if ds.SearchClient != nil {
		err := ds.SearchClient.AddDocument("users", map[string]interface{}{
			"uuid":   user.UUID,
			"name":   user.Name,
			"email":  user.Email,
			"handle": user.Handle,
		})
		if err != nil {
			return nil, err
		}
	}
	return &userorgmodels.UserResponse{
		UUID:   user.UUID,
		Name:   user.Name,
		Email:  user.Email,
		Handle: user.Handle,
		Bio:    user.Bio,
		Avatar: user.Avatar,
	}, nil
}

func (ds *Datastore) VerifyUserEmail(userUUID uuid.UUID) error {
	result := ds.DB.Model(&userorgdbmodels.User{}).Where("uuid = ?", userUUID).Update("is_verified", true)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ds *Datastore) UpdateUser(email string, updatedAttributes map[string]interface{}) (*userorgmodels.UserResponse, error) {
	var user userorgdbmodels.User
	result := ds.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if updatedAttributes["name"] != nil {
		user.Name = updatedAttributes["name"].(string)
	}
	if updatedAttributes["bio"] != nil {
		user.Bio = updatedAttributes["bio"].(string)
	}
	if updatedAttributes["avatar"] != nil {
		user.Avatar = updatedAttributes["avatar"].(string)
	}
	result = ds.DB.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if ds.SearchClient != nil {
		err := ds.SearchClient.UpdateDocument("users", user.UUID.String(), map[string]interface{}{
			"uuid":   user.UUID,
			"name":   user.Name,
			"email":  user.Email,
			"handle": user.Handle,
		})
		if err != nil {
			return nil, err
		}
	}
	return &userorgmodels.UserResponse{
		UUID:   user.UUID,
		Name:   user.Name,
		Email:  user.Email,
		Handle: user.Handle,
		Bio:    user.Bio,
		Avatar: user.Avatar,
	}, nil
}

func (ds *Datastore) UpdateUserPassword(userUUID uuid.UUID, hashedPassword string) error {
	result := ds.DB.Model(&userorgdbmodels.User{}).Where("uuid = ?", userUUID).Update("password", hashedPassword)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/////////////////////////////// CLI SESSION METHODS /////////////////////////////////

func (ds *Datastore) GetSession(sessionUUID uuid.UUID) (*authmodels.SessionResponse, error) {
	session := &authdbmodels.Session{
		BaseModel: commondbmodels.BaseModel{
			UUID: sessionUUID,
		},
	}
	err := ds.DB.Model(&session).Preload("User").First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &authmodels.SessionResponse{
		SessionUUID:    session.UUID,
		Device:         session.Device,
		DeviceId:       session.DeviceId,
		DeviceLocation: session.DeviceLocation,
		UserUUID:       session.User.UUID,
		Approved:       session.Approved,
		Invalid:        session.Invalid,
		CreatedAt:      session.CreatedAt,
	}, nil
}

func (ds *Datastore) CreateSession(deviceData string, deviceId string, deviceLocation string) (*authmodels.CreateSessionResponse, error) {
	session := authdbmodels.Session{
		Device:         deviceData,
		DeviceId:       deviceId,
		DeviceLocation: deviceLocation,
	}
	err := ds.DB.Create(&session).Error
	if err != nil {
		return nil, err
	}
	return &authmodels.CreateSessionResponse{
		SessionUUID: session.UUID,
		CreatedAt:   session.CreatedAt,
	}, nil
}

func (ds *Datastore) UpdateSession(sessionUUID uuid.UUID, userUUID uuid.UUID, updatedAttributes map[string]interface{}) (*authmodels.SessionResponse, error) {
	var session authdbmodels.Session
	result := ds.DB.Where("uuid = ?", sessionUUID).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	if updatedAttributes["approved"] != nil {
		session.Approved = updatedAttributes["approved"].(bool)
	}
	if updatedAttributes["invalid"] != nil {
		session.Invalid = updatedAttributes["invalid"].(bool)
	}
	if updatedAttributes["created_at"] != nil {
		session.CreatedAt = updatedAttributes["created_at"].(time.Time)
	}
	session.User = userorgdbmodels.User{
		BaseModel: commondbmodels.BaseModel{
			UUID: userUUID,
		},
	}
	result = ds.DB.Save(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &authmodels.SessionResponse{
		SessionUUID:    session.UUID,
		Device:         session.Device,
		DeviceId:       session.DeviceId,
		DeviceLocation: session.DeviceLocation,
		UserUUID:       session.User.UUID,
		Approved:       session.Approved,
		Invalid:        session.Invalid,
		CreatedAt:      session.CreatedAt,
	}, nil
}

// Helper
func IncrementVersion(latestVersion string) string {
	version := strings.Split(latestVersion, "v")
	versionNumber, _ := strconv.Atoi(version[1])
	newVersionNumber := versionNumber + 1
	newVersion := fmt.Sprintf("v%d", newVersionNumber)
	return newVersion
}

/////////////////////////////// MODEL METHODS /////////////////////////////////

func (ds *Datastore) GetModelByName(orgId uuid.UUID, modelName string) (*modelmodels.ModelResponse, error) {
	var model modeldbmodels.Model
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Preload("Readme.ReadmeVersions", func(db *gorm.DB) *gorm.DB {
		return db.Order("LENGTH(readme_versions.version) DESC").Order("readme_versions.version DESC").Limit(1)
	}).Where("name = ?", modelName).Where("organization_uuid = ?", orgId).Limit(1).Find(&model)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &modelmodels.ModelResponse{
		UUID: model.UUID,
		Name: model.Name,
		Wiki: model.Wiki,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   model.CreatedByUser.UUID,
			Handle: model.CreatedByUser.Handle,
			Avatar: model.CreatedByUser.Avatar,
			Name:   model.CreatedByUser.Name,
			Email:  model.CreatedByUser.Email,
		},
		UpdatedBy: userorgmodels.UserHandleResponse{
			UUID:   model.UpdatedByUser.UUID,
			Handle: model.UpdatedByUser.Handle,
			Avatar: model.UpdatedByUser.Avatar,
			Name:   model.UpdatedByUser.Name,
			Email:  model.UpdatedByUser.Email,
		},
		IsPublic: model.IsPublic,
		Readme: commonmodels.ReadmeResponse{
			UUID: model.Readme.UUID,
			LatestVersion: commonmodels.ReadmeVersionResponse{
				UUID:     model.Readme.ReadmeVersions[0].UUID,
				Version:  model.Readme.ReadmeVersions[0].Version,
				FileType: model.Readme.ReadmeVersions[0].FileType,
				Content:  model.Readme.ReadmeVersions[0].Content,
			},
		},
	}, nil
}

func (ds *Datastore) GetModelByUUID(modelUUID uuid.UUID) (*modelmodels.ModelResponse, error) {
	var model modeldbmodels.Model
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Preload("Readme.ReadmeVersions", func(db *gorm.DB) *gorm.DB {
		return db.Order("LENGTH(readme_versions.version) DESC").Order("readme_versions.version DESC").Limit(1)
	}).Where("uuid = ?", modelUUID).Limit(1).Find(&model)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &modelmodels.ModelResponse{
		UUID: model.UUID,
		Name: model.Name,
		Wiki: model.Wiki,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   model.CreatedByUser.UUID,
			Handle: model.CreatedByUser.Handle,
			Avatar: model.CreatedByUser.Avatar,
			Name:   model.CreatedByUser.Name,
			Email:  model.CreatedByUser.Email,
		},
		UpdatedBy: userorgmodels.UserHandleResponse{
			UUID:   model.UpdatedByUser.UUID,
			Handle: model.UpdatedByUser.Handle,
			Avatar: model.UpdatedByUser.Avatar,
			Name:   model.UpdatedByUser.Name,
			Email:  model.UpdatedByUser.Email,
		},
		IsPublic: model.IsPublic,
		Readme: commonmodels.ReadmeResponse{
			UUID: model.Readme.UUID,
			LatestVersion: commonmodels.ReadmeVersionResponse{
				UUID:     model.Readme.ReadmeVersions[0].UUID,
				Version:  model.Readme.ReadmeVersions[0].Version,
				FileType: model.Readme.ReadmeVersions[0].FileType,
				Content:  model.Readme.ReadmeVersions[0].Content,
			},
		},
	}, nil
}

func (ds *Datastore) GetModelReadmeVersion(modelUUID uuid.UUID, version string) (*commonmodels.ReadmeVersionResponse, error) {
	var model modeldbmodels.Model
	result := ds.DB.Preload("Readme.ReadmeVersions", func(db *gorm.DB) *gorm.DB {
		return db.Where("version = ?", version).Limit(1)
	}).Where("uuid = ?", modelUUID).Limit(1).Find(&model)
	if result.RowsAffected == 0 || len(model.Readme.ReadmeVersions) == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &commonmodels.ReadmeVersionResponse{
		UUID:     model.Readme.ReadmeVersions[0].UUID,
		Version:  model.Readme.ReadmeVersions[0].Version,
		FileType: model.Readme.ReadmeVersions[0].FileType,
		Content:  model.Readme.ReadmeVersions[0].Content,
	}, nil
}

func (ds *Datastore) GetModelReadmeAllVersions(modelUUID uuid.UUID) ([]commonmodels.ReadmeVersionResponse, error) {
	var model modeldbmodels.Model
	result := ds.DB.Preload("Readme.ReadmeVersions").Where("uuid = ?", modelUUID).Limit(1).Find(&model)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	var versions []commonmodels.ReadmeVersionResponse
	for _, version := range model.Readme.ReadmeVersions {
		versions = append(versions, commonmodels.ReadmeVersionResponse{
			UUID:     version.UUID,
			Version:  version.Version,
			FileType: version.FileType,
			Content:  version.Content,
		})
	}
	return versions, nil
}

func (ds *Datastore) UpdateModelReadme(modelUUID uuid.UUID, fileType string, content string) (*commonmodels.ReadmeVersionResponse, error) {
	var model modeldbmodels.Model
	result := ds.DB.Preload("Readme.ReadmeVersions", func(db *gorm.DB) *gorm.DB {
		return db.Order("LENGTH(version) DESC").Order("version DESC").Limit(1)
	}).Where("uuid = ?", modelUUID).Limit(1).Find(&model)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	var version string
	if len(model.Readme.ReadmeVersions) == 0 {
		version = "v1"
	} else {
		version = IncrementVersion(model.Readme.ReadmeVersions[0].Version)
	}
	readmeVersion := commondbmodels.ReadmeVersion{
		Version:  version,
		FileType: fileType,
		Content:  content,
		Readme: commondbmodels.Readme{
			BaseModel: commondbmodels.BaseModel{
				UUID: model.Readme.UUID,
			},
		},
	}
	result = ds.DB.Create(&readmeVersion)
	if result.Error != nil {
		return nil, result.Error
	}
	return &commonmodels.ReadmeVersionResponse{
		UUID:     readmeVersion.UUID,
		Version:  readmeVersion.Version,
		FileType: readmeVersion.FileType,
		Content:  readmeVersion.Content,
	}, nil
}

func (ds *Datastore) CreateModel(orgId uuid.UUID, name string, wiki string, isPublic bool, readmeData *commonmodels.ReadmeRequest, createdByUser uuid.UUID) (*modelmodels.ModelResponse, error) {
	model := modeldbmodels.Model{
		Name: name,
		Wiki: wiki,
		Org: userorgdbmodels.Organization{
			BaseModel: commondbmodels.BaseModel{
				UUID: orgId,
			},
		},
		CreatedByUser: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: createdByUser,
			},
		},
		UpdatedByUser: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: createdByUser,
			},
		},
		IsPublic: isPublic,
		Readme: commondbmodels.Readme{
			ReadmeVersions: []commondbmodels.ReadmeVersion{
				{
					Version:  "v1",
					FileType: readmeData.FileType,
					Content:  readmeData.Content,
				},
			},
		},
	}
	var user userorgdbmodels.User
	err := ds.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&model)
		if result.Error != nil {
			return result.Error
		}
		result = tx.Where("uuid = ?", createdByUser).First(&user)
		if result.Error != nil {
			return result.Error
		}
		modelUser := modeldbmodels.ModelUser{
			UserUUID:  user.UUID,
			ModelUUID: model.UUID,
			Role:      "owner",
		}
		result = tx.Create(&modelUser)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if ds.SearchClient != nil {
		err := ds.SearchClient.AddDocument("models", map[string]interface{}{
			"uuid":              model.UUID,
			"name":              model.Name,
			"wiki":              model.Wiki,
			"organization_uuid": model.Org.UUID,
			"is_public":         model.IsPublic,
		})
		if err != nil {
			return nil, err
		}
	}
	return &modelmodels.ModelResponse{
		UUID: model.UUID,
		Name: model.Name,
		Wiki: model.Wiki,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   model.CreatedByUser.UUID,
			Handle: model.CreatedByUser.Handle,
			Avatar: model.CreatedByUser.Avatar,
			Name:   model.CreatedByUser.Name,
			Email:  model.CreatedByUser.Email,
		},
		UpdatedBy: userorgmodels.UserHandleResponse{
			UUID:   model.UpdatedByUser.UUID,
			Handle: model.UpdatedByUser.Handle,
			Avatar: model.UpdatedByUser.Avatar,
			Name:   model.UpdatedByUser.Name,
			Email:  model.UpdatedByUser.Email,
		},
		IsPublic: model.IsPublic,
		Readme: commonmodels.ReadmeResponse{
			UUID: model.Readme.UUID,
			LatestVersion: commonmodels.ReadmeVersionResponse{
				UUID:     model.Readme.ReadmeVersions[0].UUID,
				Version:  model.Readme.ReadmeVersions[0].Version,
				FileType: model.Readme.ReadmeVersions[0].FileType,
				Content:  model.Readme.ReadmeVersions[0].Content,
			},
		},
	}, nil
}

func (ds *Datastore) GetAllPublicModels() ([]modelmodels.ModelResponse, error) {
	var mymodels []modeldbmodels.Model
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Preload("Org").Where("is_public = ?", true).Find(&mymodels)
	if result.Error != nil {
		return nil, result.Error
	}
	modelResponses := make([]modelmodels.ModelResponse, len(mymodels))
	for i, model := range mymodels {
		modelResponses[i] = modelmodels.ModelResponse{
			UUID: model.UUID,
			Name: model.Name,
			Wiki: model.Wiki,
			CreatedBy: userorgmodels.UserHandleResponse{
				UUID:   model.CreatedByUser.UUID,
				Handle: model.CreatedByUser.Handle,
				Avatar: model.CreatedByUser.Avatar,
				Name:   model.CreatedByUser.Name,
				Email:  model.CreatedByUser.Email,
			},
			UpdatedBy: userorgmodels.UserHandleResponse{
				UUID:   model.UpdatedByUser.UUID,
				Handle: model.UpdatedByUser.Handle,
				Avatar: model.UpdatedByUser.Avatar,
				Name:   model.UpdatedByUser.Name,
				Email:  model.UpdatedByUser.Email,
			},
			Org: userorgmodels.OrganizationHandleResponse{
				UUID:        model.Org.UUID,
				Name:        model.Org.Name,
				Handle:      model.Org.Handle,
				Avatar:      model.Org.Avatar,
				Description: model.Org.Description,
			},
			IsPublic: model.IsPublic,
		}
	}
	return modelResponses, nil
}

func (ds *Datastore) GetAllModels(orgId uuid.UUID) ([]modelmodels.ModelResponse, error) {
	var mymodels []modeldbmodels.Model
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Where("organization_uuid = ?", orgId).Find(&mymodels)
	if result.Error != nil {
		return nil, result.Error
	}
	modelResponses := make([]modelmodels.ModelResponse, len(mymodels))
	for i, model := range mymodels {
		modelResponses[i] = modelmodels.ModelResponse{
			UUID: model.UUID,
			Name: model.Name,
			Wiki: model.Wiki,
			CreatedBy: userorgmodels.UserHandleResponse{
				UUID:   model.CreatedByUser.UUID,
				Handle: model.CreatedByUser.Handle,
				Avatar: model.CreatedByUser.Avatar,
				Name:   model.CreatedByUser.Name,
				Email:  model.CreatedByUser.Email,
			},
			UpdatedBy: userorgmodels.UserHandleResponse{
				UUID:   model.UpdatedByUser.UUID,
				Handle: model.UpdatedByUser.Handle,
				Avatar: model.UpdatedByUser.Avatar,
				Name:   model.UpdatedByUser.Name,
				Email:  model.UpdatedByUser.Email,
			},
			IsPublic: model.IsPublic,
		}
	}
	return modelResponses, nil
}

func (ds *Datastore) GetModelAllBranches(modelUUID uuid.UUID) ([]modelmodels.ModelBranchResponse, error) {
	var modelBranches []modeldbmodels.ModelBranch
	result := ds.DB.Preload("Model").Where("model_uuid = ?", modelUUID).Find(&modelBranches)
	if result.Error != nil {
		return nil, result.Error
	}
	branches := make([]modelmodels.ModelBranchResponse, len(modelBranches))
	for i, branch := range modelBranches {
		branches[i] = modelmodels.ModelBranchResponse{
			UUID: branch.UUID,
			Name: branch.Name,
			Model: modelmodels.ModelNameResponse{
				UUID: branch.Model.UUID,
				Name: branch.Model.Name,
			},
			IsDefault: branch.IsDefault,
		}
	}
	return branches, nil
}

func (ds *Datastore) CreateModelBranch(modelUUID uuid.UUID, modelBranchName string) (*modelmodels.ModelBranchResponse, error) {
	modelBranch := modeldbmodels.ModelBranch{
		Name: modelBranchName,
		Model: modeldbmodels.Model{
			BaseModel: commondbmodels.BaseModel{
				UUID: modelUUID,
			},
		},
	}
	err := ds.DB.Create(&modelBranch).Preload("Model").Error
	if err != nil {
		return nil, err
	}
	return &modelmodels.ModelBranchResponse{
		UUID: modelBranch.UUID,
		Name: modelBranch.Name,
		Model: modelmodels.ModelNameResponse{
			UUID: modelBranch.Model.UUID,
			Name: modelBranch.Model.Name,
		},
		IsDefault: modelBranch.IsDefault,
	}, nil
}

func (ds *Datastore) RegisterModelFile(modelBranchUUID uuid.UUID, sourceType string, sourcePublicURL string, filePath string, isEmpty bool, hash string, userUUID uuid.UUID) (*modelmodels.ModelBranchVersionResponse, error) {
	latestModelVersion := modeldbmodels.ModelVersion{
		BranchUUID: modelBranchUUID,
	}
	res := ds.DB.Where(&latestModelVersion).Order("created_at desc").Limit(1).Find(&latestModelVersion)

	var newVersion string
	if res.RowsAffected == 0 {
		newVersion = "v1"
	} else {
		latestVersion := latestModelVersion.Version
		newVersion = IncrementVersion(latestVersion)
	}

	modelVersion := modeldbmodels.ModelVersion{
		Hash:    hash,
		Version: newVersion,
		Branch: modeldbmodels.ModelBranch{
			BaseModel: commondbmodels.BaseModel{
				UUID: modelBranchUUID,
			},
		},
		CreatedByUser: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: userUUID,
			},
		},
		Path:       fmt.Sprintf("%s/%s", sourcePublicURL, filePath),
		SourceType: sourceType,
		IsEmpty:    isEmpty,
	}

	err := ds.DB.Create(&modelVersion).Error
	if err != nil {
		return nil, err
	}
	err = ds.DB.Preload("Branch").Preload("CreatedByUser").Find(&modelVersion).Error
	if err != nil {
		return nil, err
	}

	return &modelmodels.ModelBranchVersionResponse{
		UUID:    modelVersion.UUID,
		Hash:    modelVersion.Hash,
		Version: modelVersion.Version,
		Branch: modelmodels.ModelBranchNameResponse{
			UUID: modelVersion.Branch.UUID,
			Name: modelVersion.Branch.Name,
		},
		Path:       modelVersion.Path,
		SourceType: modelVersion.SourceType,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   modelVersion.CreatedByUser.UUID,
			Name:   modelVersion.CreatedByUser.Name,
			Avatar: modelVersion.CreatedByUser.Avatar,
			Email:  modelVersion.CreatedByUser.Email,
			Handle: modelVersion.CreatedByUser.Handle,
		},
		CreatedAt: modelVersion.CreatedAt,
		IsEmpty:   modelVersion.IsEmpty,
	}, nil
}

func (ds *Datastore) MigrateModelVersionBranch(modelVersion uuid.UUID, toBranch uuid.UUID) (*modelmodels.ModelBranchVersionResponse, error) {
	var modelVersionDB modeldbmodels.ModelVersion
	err := ds.DB.Preload("Branch").Preload("CreatedByUser").Preload("Path").Where("uuid = ?", modelVersion).First(&modelVersionDB).Error
	if err != nil {
		return nil, err
	}
	//Update the branch of the model version
	modelVersionDB.Branch.UUID = toBranch
	modelVersionDB.BaseModel.UUID = uuid.Nil
	modelVersionDB.CreatedAt = time.Now()
	modelVersionDB.UpdatedAt = time.Now()
	err = ds.DB.Save(&modelVersionDB).Error
	if err != nil {
		return nil, err
	}

	// Migrate logs
	var logs []dbmodels.Log
	err = ds.DB.Where("model_version_uuid = ?", modelVersion).Preload("ModelVersion").Find(&logs).Error
	if err != nil {
		return nil, err
	}
	err = ds.DB.Transaction(func(tx *gorm.DB) error {
		for _, log := range logs {
			log.BaseModel.UUID = uuid.Nil
			log.ModelVersion.BaseModel.UUID = modelVersionDB.UUID
			log.CreatedAt = time.Now()
			log.UpdatedAt = time.Now()
			err = tx.Save(&log).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &modelmodels.ModelBranchVersionResponse{
		UUID:    modelVersionDB.UUID,
		Hash:    modelVersionDB.Hash,
		Version: modelVersionDB.Version,
		Branch: modelmodels.ModelBranchNameResponse{
			UUID: modelVersionDB.BranchUUID,
			Name: modelVersionDB.Branch.Name,
		},
		Path:       modelVersionDB.Path,
		SourceType: modelVersionDB.SourceType,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   modelVersionDB.CreatedByUser.UUID,
			Handle: modelVersionDB.CreatedByUser.Handle,
			Name:   modelVersionDB.CreatedByUser.Name,
			Avatar: modelVersionDB.CreatedByUser.Avatar,
			Email:  modelVersionDB.CreatedByUser.Email,
		},
		CreatedAt: modelVersionDB.CreatedAt,
		IsEmpty:   modelVersionDB.IsEmpty,
	}, nil
}

func (ds *Datastore) GetModelAllVersions(modelUUID uuid.UUID) ([]modelmodels.ModelBranchVersionResponse, error) {
	var modelVersions []modeldbmodels.ModelVersion
	err := ds.DB.Select("model_versions.*").Joins("JOIN model_branches ON model_branches.uuid = model_versions.branch_uuid").Where("model_branches.model_uuid = ?", modelUUID).Preload("Branch").Preload("CreatedByUser").Preload("Path").Find(&modelVersions).Error
	if err != nil {
		return nil, err
	}
	var modelVersionsResponse []modelmodels.ModelBranchVersionResponse
	for _, modelVersion := range modelVersions {
		modelVersionsResponse = append(modelVersionsResponse, modelmodels.ModelBranchVersionResponse{
			UUID:    modelVersion.UUID,
			Hash:    modelVersion.Hash,
			Version: modelVersion.Version,
			Branch: modelmodels.ModelBranchNameResponse{
				UUID: modelVersion.Branch.UUID,
				Name: modelVersion.Branch.Name,
			},
			Path:       modelVersion.Path,
			SourceType: modelVersion.SourceType,
			CreatedBy: userorgmodels.UserHandleResponse{
				UUID:   modelVersion.CreatedByUser.UUID,
				Handle: modelVersion.CreatedByUser.Handle,
				Name:   modelVersion.CreatedByUser.Name,
				Avatar: modelVersion.CreatedByUser.Avatar,
				Email:  modelVersion.CreatedByUser.Email,
			},
			CreatedAt: modelVersion.CreatedAt,
			IsEmpty:   modelVersion.IsEmpty,
		})
	}
	return modelVersionsResponse, nil
}

func (ds *Datastore) GetModelBranchByName(orgId uuid.UUID, modelName string, modelBranchName string) (*modelmodels.ModelBranchResponse, error) {
	var modelBranch modeldbmodels.ModelBranch
	model, err := ds.GetModelByName(orgId, modelName)
	if err != nil {
		return nil, err
	}
	res := ds.DB.Where("name = ?", modelBranchName).Where("model_uuid = ?", model.UUID).Preload("Model").Limit(1).Find(&modelBranch)
	if res.RowsAffected == 0 {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &modelmodels.ModelBranchResponse{
		UUID: modelBranch.UUID,
		Name: modelBranch.Name,
		Model: modelmodels.ModelNameResponse{
			UUID: modelBranch.Model.UUID,
			Name: modelBranch.Model.Name,
		},
		IsDefault: modelBranch.IsDefault,
	}, nil
}

func (ds *Datastore) GetModelBranchByUUID(modelBranchUUID uuid.UUID) (*modelmodels.ModelBranchResponse, error) {
	var modelBranch modeldbmodels.ModelBranch
	res := ds.DB.Where("uuid = ?", modelBranchUUID).Preload("Model").Limit(1).Find(&modelBranch)
	if res.RowsAffected == 0 {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &modelmodels.ModelBranchResponse{
		UUID: modelBranch.UUID,
		Name: modelBranch.Name,
		Model: modelmodels.ModelNameResponse{
			UUID: modelBranch.Model.UUID,
			Name: modelBranch.Model.Name,
		},
		IsDefault: modelBranch.IsDefault,
	}, nil
}

func (ds *Datastore) GetModelBranchAllVersions(modelBranchUUID uuid.UUID, withLogs bool) ([]modelmodels.ModelBranchVersionResponse, error) {
	var modelVersions []modeldbmodels.ModelVersion
	err := ds.DB.Where("branch_uuid = ?", modelBranchUUID).Preload("Branch").Preload("CreatedByUser").Order("LENGTH(version) DESC").Order("version DESC").Find(&modelVersions).Error
	if err != nil {
		return nil, err
	}
	var modelVersionsResponse []modelmodels.ModelBranchVersionResponse
	for _, modelVersion := range modelVersions {
		modelBranchVersion := modelmodels.ModelBranchVersionResponse{
			UUID:    modelVersion.UUID,
			Hash:    modelVersion.Hash,
			Version: modelVersion.Version,
			Branch: modelmodels.ModelBranchNameResponse{
				UUID: modelVersion.Branch.UUID,
				Name: modelVersion.Branch.Name,
			},
			Path:    modelVersion.Path,
			IsEmpty: modelVersion.IsEmpty,
			CreatedBy: userorgmodels.UserHandleResponse{
				UUID:   modelVersion.CreatedByUser.UUID,
				Handle: modelVersion.CreatedByUser.Handle,
				Name:   modelVersion.CreatedByUser.Name,
				Avatar: modelVersion.CreatedByUser.Avatar,
				Email:  modelVersion.CreatedByUser.Email,
			},
			CreatedAt: modelVersion.CreatedAt,
		}
		if withLogs {
			logResponses, err := ds.GetLogForModelVersion(modelVersion.UUID)
			if err != nil {
				return nil, err
			}
			var logData []commonmodels.LogDataResponse
			for _, log := range logResponses {
				logData = append(logData, commonmodels.LogDataResponse{
					Key:  log.Key,
					Data: log.Data,
				})
			}
			modelBranchVersion.Logs = logData
		}
		modelVersionsResponse = append(modelVersionsResponse, modelBranchVersion)
	}
	return modelVersionsResponse, nil
}

func (ds *Datastore) GetModelBranchVersion(modelBranchUUID uuid.UUID, version string) (*modelmodels.ModelBranchVersionResponse, error) {
	var modelVersion modeldbmodels.ModelVersion
	var res *gorm.DB
	if strings.ToLower(version) == "latest" {
		res = ds.DB.Order("created_at desc").Where("branch_uuid = ?", modelBranchUUID).Preload("CreatedByUser").Preload("Branch").Limit(1).Find(&modelVersion)
	} else {
		res = ds.DB.Where("version = ?", version).Where("branch_uuid = ?", modelBranchUUID).Preload("CreatedByUser").Preload("Branch").Limit(1).Find(&modelVersion)
	}
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &modelmodels.ModelBranchVersionResponse{
		UUID:    modelVersion.UUID,
		Hash:    modelVersion.Hash,
		Version: modelVersion.Version,
		Branch: modelmodels.ModelBranchNameResponse{
			UUID: modelVersion.Branch.UUID,
			Name: modelVersion.Branch.Name,
		},
		Path:       modelVersion.Path,
		SourceType: modelVersion.SourceType,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   modelVersion.CreatedByUser.UUID,
			Handle: modelVersion.CreatedByUser.Handle,
			Name:   modelVersion.CreatedByUser.Name,
			Avatar: modelVersion.CreatedByUser.Avatar,
			Email:  modelVersion.CreatedByUser.Email,
		},
		CreatedAt: modelVersion.CreatedAt,
		IsEmpty:   modelVersion.IsEmpty,
	}, nil
}

/////////////////////////////// DATASET METHODS/////////////////////////////////

func (ds *Datastore) GetDatasetByName(orgId uuid.UUID, datasetName string) (*datasetmodels.DatasetResponse, error) {
	var dataset datasetdbmodels.Dataset
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Preload("Readme.ReadmeVersions", func(db *gorm.DB) *gorm.DB {
		return db.Order("LENGTH(readme_versions.version) DESC").Order("readme_versions.version DESC").Limit(1)
	}).Where("name = ?", datasetName).Where("organization_uuid = ?", orgId).Limit(1).Find(&dataset)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &datasetmodels.DatasetResponse{
		UUID: dataset.UUID,
		Name: dataset.Name,
		Wiki: dataset.Wiki,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   dataset.CreatedByUser.UUID,
			Handle: dataset.CreatedByUser.Handle,
			Name:   dataset.CreatedByUser.Name,
			Avatar: dataset.CreatedByUser.Avatar,
			Email:  dataset.CreatedByUser.Email,
		},
		UpdatedBy: userorgmodels.UserHandleResponse{
			UUID:   dataset.UpdatedByUser.UUID,
			Handle: dataset.UpdatedByUser.Handle,
			Name:   dataset.UpdatedByUser.Name,
			Avatar: dataset.UpdatedByUser.Avatar,
			Email:  dataset.UpdatedByUser.Email,
		},
		IsPublic: dataset.IsPublic,
		Readme: commonmodels.ReadmeResponse{
			UUID: dataset.Readme.UUID,
			LatestVersion: commonmodels.ReadmeVersionResponse{
				UUID:     dataset.Readme.ReadmeVersions[0].UUID,
				Version:  dataset.Readme.ReadmeVersions[0].Version,
				FileType: dataset.Readme.ReadmeVersions[0].FileType,
				Content:  dataset.Readme.ReadmeVersions[0].Content,
			},
		},
	}, nil
}

func (ds *Datastore) GetDatasetByUUID(datasetUUID uuid.UUID) (*datasetmodels.DatasetResponse, error) {
	var dataset datasetdbmodels.Dataset
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Preload("Readme.ReadmeVersions", func(db *gorm.DB) *gorm.DB {
		return db.Order("LENGTH(readme_versions.version) DESC").Order("readme_versions.version DESC").Limit(1)
	}).Where("uuid = ?", datasetUUID).Limit(1).Find(&dataset)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &datasetmodels.DatasetResponse{
		UUID: dataset.UUID,
		Name: dataset.Name,
		Wiki: dataset.Wiki,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   dataset.CreatedByUser.UUID,
			Handle: dataset.CreatedByUser.Handle,
			Name:   dataset.CreatedByUser.Name,
			Avatar: dataset.CreatedByUser.Avatar,
			Email:  dataset.CreatedByUser.Email,
		},
		UpdatedBy: userorgmodels.UserHandleResponse{
			UUID:   dataset.UpdatedByUser.UUID,
			Handle: dataset.UpdatedByUser.Handle,
			Name:   dataset.UpdatedByUser.Name,
			Avatar: dataset.UpdatedByUser.Avatar,
			Email:  dataset.UpdatedByUser.Email,
		},
		IsPublic: dataset.IsPublic,
		Readme: commonmodels.ReadmeResponse{
			UUID: dataset.Readme.UUID,
			LatestVersion: commonmodels.ReadmeVersionResponse{
				UUID:     dataset.Readme.ReadmeVersions[0].UUID,
				Version:  dataset.Readme.ReadmeVersions[0].Version,
				FileType: dataset.Readme.ReadmeVersions[0].FileType,
				Content:  dataset.Readme.ReadmeVersions[0].Content,
			},
		},
	}, nil
}

func (ds *Datastore) GetDatasetReadmeVersion(datasetUUID uuid.UUID, version string) (*commonmodels.ReadmeVersionResponse, error) {
	var dataset datasetdbmodels.Dataset
	result := ds.DB.Preload("Readme.ReadmeVersions", func(db *gorm.DB) *gorm.DB {
		return db.Where("version = ?", version).Limit(1)
	}).Where("uuid = ?", datasetUUID).Limit(1).Find(&dataset)
	if result.RowsAffected == 0 || len(dataset.Readme.ReadmeVersions) == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &commonmodels.ReadmeVersionResponse{
		UUID:     dataset.Readme.ReadmeVersions[0].UUID,
		Version:  dataset.Readme.ReadmeVersions[0].Version,
		FileType: dataset.Readme.ReadmeVersions[0].FileType,
		Content:  dataset.Readme.ReadmeVersions[0].Content,
	}, nil
}

func (ds *Datastore) GetDatasetReadmeAllVersions(datasetUUID uuid.UUID) ([]commonmodels.ReadmeVersionResponse, error) {
	var dataset datasetdbmodels.Dataset
	result := ds.DB.Preload("Readme.ReadmeVersions").Where("uuid = ?", datasetUUID).Limit(1).Find(&dataset)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	var versions []commonmodels.ReadmeVersionResponse
	for _, version := range dataset.Readme.ReadmeVersions {
		versions = append(versions, commonmodels.ReadmeVersionResponse{
			UUID:     version.UUID,
			Version:  version.Version,
			FileType: version.FileType,
			Content:  version.Content,
		})
	}
	return versions, nil
}

func (ds *Datastore) UpdateDatasetReadme(datasetUUID uuid.UUID, fileType string, content string) (*commonmodels.ReadmeVersionResponse, error) {
	var dataset datasetdbmodels.Dataset
	result := ds.DB.Preload("Readme.ReadmeVersions", func(db *gorm.DB) *gorm.DB {
		return db.Order("LENGTH(version) DESC").Order("version DESC").Limit(1)
	}).Where("uuid = ?", datasetUUID).Limit(1).Find(&dataset)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	var version string
	if len(dataset.Readme.ReadmeVersions) == 0 {
		version = "v1"
	} else {
		version = IncrementVersion(dataset.Readme.ReadmeVersions[0].Version)
	}
	readmeVersion := commondbmodels.ReadmeVersion{
		Version:  version,
		FileType: fileType,
		Content:  content,
		Readme: commondbmodels.Readme{
			BaseModel: commondbmodels.BaseModel{
				UUID: dataset.Readme.UUID,
			},
		},
	}
	result = ds.DB.Create(&readmeVersion)
	if result.Error != nil {
		return nil, result.Error
	}
	return &commonmodels.ReadmeVersionResponse{
		UUID:     readmeVersion.UUID,
		Version:  readmeVersion.Version,
		FileType: readmeVersion.FileType,
		Content:  readmeVersion.Content,
	}, nil
}

func (ds *Datastore) CreateDataset(orgId uuid.UUID, name string, wiki string, isPublic bool, readmeData *commonmodels.ReadmeRequest, createdByUser uuid.UUID) (*datasetmodels.DatasetResponse, error) {
	dataset := datasetdbmodels.Dataset{
		Name: name,
		Wiki: wiki,
		Org: userorgdbmodels.Organization{
			BaseModel: commondbmodels.BaseModel{
				UUID: orgId,
			},
		},
		CreatedByUser: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: createdByUser,
			},
		},
		UpdatedByUser: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: createdByUser,
			},
		},
		IsPublic: isPublic,
		Readme: commondbmodels.Readme{
			ReadmeVersions: []commondbmodels.ReadmeVersion{
				{
					Version:  "v1",
					FileType: readmeData.FileType,
					Content:  readmeData.Content,
				},
			},
		},
	}
	var user userorgdbmodels.User
	err := ds.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&dataset)
		if result.Error != nil {
			return result.Error
		}
		result = tx.Where("uuid = ?", createdByUser).First(&user)
		if result.Error != nil {
			return result.Error
		}
		datasetUser := datasetdbmodels.DatasetUser{
			UserUUID:    user.UUID,
			DatasetUUID: dataset.UUID,
			Role:        "owner",
		}
		result = tx.Create(&datasetUser)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if ds.SearchClient != nil {
		err := ds.SearchClient.AddDocument("datasets", map[string]interface{}{
			"uuid":              dataset.UUID,
			"name":              dataset.Name,
			"wiki":              dataset.Wiki,
			"organization_uuid": dataset.Org.UUID,
			"is_public":         dataset.IsPublic,
		})
		if err != nil {
			return nil, err
		}
	}
	return &datasetmodels.DatasetResponse{
		UUID: dataset.UUID,
		Name: dataset.Name,
		Wiki: dataset.Wiki,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   dataset.CreatedByUser.UUID,
			Handle: dataset.CreatedByUser.Handle,
			Name:   dataset.CreatedByUser.Name,
			Avatar: dataset.CreatedByUser.Avatar,
			Email:  dataset.CreatedByUser.Email,
		},
		UpdatedBy: userorgmodels.UserHandleResponse{
			UUID:   dataset.UpdatedByUser.UUID,
			Handle: dataset.UpdatedByUser.Handle,
			Name:   dataset.UpdatedByUser.Name,
			Avatar: dataset.UpdatedByUser.Avatar,
			Email:  dataset.UpdatedByUser.Email,
		},
		IsPublic: dataset.IsPublic,
		Readme: commonmodels.ReadmeResponse{
			UUID: dataset.Readme.UUID,
			LatestVersion: commonmodels.ReadmeVersionResponse{
				UUID:     dataset.Readme.ReadmeVersions[0].UUID,
				Version:  dataset.Readme.ReadmeVersions[0].Version,
				FileType: dataset.Readme.ReadmeVersions[0].FileType,
				Content:  dataset.Readme.ReadmeVersions[0].Content,
			},
		},
	}, nil
}

func (ds *Datastore) GetAllPublicDatasets() ([]datasetmodels.DatasetResponse, error) {
	var datasets []datasetdbmodels.Dataset
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Preload("Org").Where("is_public = ?", true).Find(&datasets)
	if result.Error != nil {
		return nil, result.Error
	}
	modelResponses := make([]datasetmodels.DatasetResponse, len(datasets))
	for i, dataset := range datasets {
		modelResponses[i] = datasetmodels.DatasetResponse{
			UUID: dataset.UUID,
			Name: dataset.Name,
			Wiki: dataset.Wiki,
			CreatedBy: userorgmodels.UserHandleResponse{
				UUID:   dataset.CreatedByUser.UUID,
				Handle: dataset.CreatedByUser.Handle,
				Avatar: dataset.CreatedByUser.Avatar,
				Name:   dataset.CreatedByUser.Name,
				Email:  dataset.CreatedByUser.Email,
			},
			UpdatedBy: userorgmodels.UserHandleResponse{
				UUID:   dataset.UpdatedByUser.UUID,
				Handle: dataset.UpdatedByUser.Handle,
				Avatar: dataset.UpdatedByUser.Avatar,
				Name:   dataset.UpdatedByUser.Name,
				Email:  dataset.UpdatedByUser.Email,
			},
			Org: userorgmodels.OrganizationHandleResponse{
				UUID:        dataset.Org.UUID,
				Name:        dataset.Org.Name,
				Handle:      dataset.Org.Handle,
				Avatar:      dataset.Org.Avatar,
				Description: dataset.Org.Description,
			},
			IsPublic: dataset.IsPublic,
		}
	}
	return modelResponses, nil
}

func (ds *Datastore) GetAllDatasets(orgId uuid.UUID, showPublic bool) ([]datasetmodels.DatasetResponse, error) {
	var datasets []datasetdbmodels.Dataset
	var result *gorm.DB
	if showPublic {
		result = ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Where("organization_uuid = ?", orgId).Where("is_public LIKE ?", showPublic).Find(&datasets)
	} else {
		result = ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Where("organization_uuid = ?", orgId).Find(&datasets)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	datasetResponses := make([]datasetmodels.DatasetResponse, len(datasets))
	for i, dataset := range datasets {
		datasetResponses[i] = datasetmodels.DatasetResponse{
			UUID: dataset.UUID,
			Name: dataset.Name,
			Wiki: dataset.Wiki,
			CreatedBy: userorgmodels.UserHandleResponse{
				UUID:   dataset.CreatedByUser.UUID,
				Handle: dataset.CreatedByUser.Handle,
				Avatar: dataset.CreatedByUser.Avatar,
				Name:   dataset.CreatedByUser.Name,
				Email:  dataset.CreatedByUser.Email,
			},
			UpdatedBy: userorgmodels.UserHandleResponse{
				UUID:   dataset.UpdatedByUser.UUID,
				Handle: dataset.UpdatedByUser.Handle,
				Avatar: dataset.UpdatedByUser.Avatar,
				Name:   dataset.UpdatedByUser.Name,
				Email:  dataset.UpdatedByUser.Email,
			},
			IsPublic: dataset.IsPublic,
		}
	}
	return datasetResponses, nil
}

func (ds *Datastore) GetDatasetAllBranches(datasetUUID uuid.UUID) ([]datasetmodels.DatasetBranchResponse, error) {
	var datasetBranches []datasetdbmodels.DatasetBranch
	result := ds.DB.Preload("Dataset").Where("dataset_uuid = ?", datasetUUID).Find(&datasetBranches)
	if result.Error != nil {
		return nil, result.Error
	}
	branches := make([]datasetmodels.DatasetBranchResponse, len(datasetBranches))
	for i, branch := range datasetBranches {
		branches[i] = datasetmodels.DatasetBranchResponse{
			UUID: branch.UUID,
			Name: branch.Name,
			Dataset: datasetmodels.DatasetNameResponse{
				UUID: branch.Dataset.UUID,
				Name: branch.Dataset.Name,
			},
			IsDefault: branch.IsDefault,
		}
	}
	return branches, nil
}

func (ds *Datastore) CreateDatasetBranch(datasetUUID uuid.UUID, datasetBranchName string) (*datasetmodels.DatasetBranchResponse, error) {
	datasetBranch := datasetdbmodels.DatasetBranch{
		Name: datasetBranchName,
		Dataset: datasetdbmodels.Dataset{
			BaseModel: commondbmodels.BaseModel{
				UUID: datasetUUID,
			},
		},
	}
	err := ds.DB.Create(&datasetBranch).Preload("Dataset").Error
	if err != nil {
		return nil, err
	}
	return &datasetmodels.DatasetBranchResponse{
		UUID: datasetBranch.UUID,
		Name: datasetBranch.Name,
		Dataset: datasetmodels.DatasetNameResponse{
			UUID: datasetBranch.Dataset.UUID,
			Name: datasetBranch.Dataset.Name,
		},
		IsDefault: datasetBranch.IsDefault,
	}, nil
}

func (ds *Datastore) RegisterDatasetFile(datasetBranchUUID uuid.UUID, sourceType string, sourcePublicURL string, filePath string, isEmpty bool, hash string, lineage string, userUUID uuid.UUID) (*datasetmodels.DatasetBranchVersionResponse, error) {
	latestDatasetVersion := datasetdbmodels.DatasetVersion{
		BranchUUID: datasetBranchUUID,
	}
	res := ds.DB.Where(&latestDatasetVersion).Order("created_at desc").Limit(1).Find(&latestDatasetVersion)

	var newVersion string
	if res.RowsAffected == 0 {
		newVersion = "v1"
	} else {
		latestVersion := latestDatasetVersion.Version
		newVersion = IncrementVersion(latestVersion)
	}

	datasetVersion := datasetdbmodels.DatasetVersion{
		Hash:    hash,
		Version: newVersion,
		Branch: datasetdbmodels.DatasetBranch{
			BaseModel: commondbmodels.BaseModel{
				UUID: datasetBranchUUID,
			},
		},
		CreatedByUser: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: userUUID,
			},
		},
		Lineage: datasetdbmodels.Lineage{
			Lineage: lineage,
		},
		Path:       fmt.Sprintf("%s/%s", sourcePublicURL, filePath),
		SourceType: sourceType,
		IsEmpty:    isEmpty,
	}
	err := ds.DB.Create(&datasetVersion).Error
	if err != nil {
		return nil, err
	}
	err = ds.DB.Preload("Lineage").Preload("Branch").Preload("CreatedByUser").Find(&datasetVersion).Error
	if err != nil {
		return nil, err
	}

	return &datasetmodels.DatasetBranchVersionResponse{
		UUID:    datasetVersion.UUID,
		Hash:    datasetVersion.Hash,
		Version: datasetVersion.Version,
		Branch: datasetmodels.DatasetBranchNameResponse{
			UUID: datasetVersion.Branch.UUID,
			Name: datasetVersion.Branch.Name,
		},
		Path:       datasetVersion.Path,
		SourceType: datasetVersion.SourceType,
		Lineage: datasetmodels.LineageResponse{
			UUID:    datasetVersion.Lineage.UUID,
			Lineage: datasetVersion.Lineage.Lineage,
		},
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   datasetVersion.CreatedByUser.UUID,
			Name:   datasetVersion.CreatedByUser.Name,
			Avatar: datasetVersion.CreatedByUser.Avatar,
			Email:  datasetVersion.CreatedByUser.Email,
			Handle: datasetVersion.CreatedByUser.Handle,
		},
		CreatedAt: datasetVersion.CreatedAt,
		IsEmpty:   datasetVersion.IsEmpty,
	}, nil
}

func (ds *Datastore) MigrateDatasetVersionBranch(datasetVersion uuid.UUID, toBranch uuid.UUID) (*datasetmodels.DatasetBranchVersionResponse, error) {
	var datasetVersionDB datasetdbmodels.DatasetVersion
	err := ds.DB.Preload("Branch").Preload("Lineage").Preload("Path").Preload("CreatedByUser").Where("uuid = ?", datasetVersion).First(&datasetVersionDB).Error
	if err != nil {
		return nil, err
	}
	//Update the branch of the dataset version
	datasetVersionDB.Branch.UUID = toBranch
	datasetVersionDB.BaseModel.UUID = uuid.Nil
	//Create a new lineage record
	datasetVersionDB.Lineage.UUID = uuid.Nil
	datasetVersionDB.CreatedAt = time.Now()
	datasetVersionDB.UpdatedAt = time.Now()
	err = ds.DB.Save(&datasetVersionDB).Error
	if err != nil {
		return nil, err
	}

	// Migrate logs
	var logs []dbmodels.Log
	err = ds.DB.Where("dataset_version_uuid = ?", datasetVersion).Preload("DatasetVersion").Find(&logs).Error
	if err != nil {
		return nil, err
	}
	err = ds.DB.Transaction(func(tx *gorm.DB) error {
		for _, log := range logs {
			log.BaseModel.UUID = uuid.Nil
			log.DatasetVersion.BaseModel.UUID = datasetVersionDB.UUID
			log.CreatedAt = time.Now()
			log.UpdatedAt = time.Now()
			err = tx.Save(&log).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &datasetmodels.DatasetBranchVersionResponse{
		UUID:    datasetVersionDB.UUID,
		Hash:    datasetVersionDB.Hash,
		Version: datasetVersionDB.Version,
		Branch: datasetmodels.DatasetBranchNameResponse{
			UUID: datasetVersionDB.BranchUUID,
			Name: datasetVersionDB.Branch.Name,
		},
		Path:       datasetVersionDB.Path,
		SourceType: datasetVersionDB.SourceType,
		Lineage: datasetmodels.LineageResponse{
			UUID:    datasetVersionDB.Lineage.UUID,
			Lineage: datasetVersionDB.Lineage.Lineage,
		},
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   datasetVersionDB.CreatedByUser.UUID,
			Name:   datasetVersionDB.CreatedByUser.Name,
			Avatar: datasetVersionDB.CreatedByUser.Avatar,
			Email:  datasetVersionDB.CreatedByUser.Email,
			Handle: datasetVersionDB.CreatedByUser.Handle,
		},
		CreatedAt: datasetVersionDB.CreatedAt,
		IsEmpty:   datasetVersionDB.IsEmpty,
	}, nil
}

func (ds *Datastore) GetDatasetAllVersions(datasetUUID uuid.UUID) ([]datasetmodels.DatasetBranchVersionResponse, error) {
	var datasetVersions []datasetdbmodels.DatasetVersion
	err := ds.DB.Select("dataset_versions.*").Joins("JOIN dataset_branches ON dataset_branches.uuid = dataset_versions.branch_uuid").Where("dataset_branches.dataset_uuid = ?", datasetUUID).Preload("Branch").Preload("CreatedByUser").Preload("Lineage").Preload("Path").Find(&datasetVersions).Error
	if err != nil {
		return nil, err
	}
	var datasetVersionsResponse []datasetmodels.DatasetBranchVersionResponse
	for _, datasetVersion := range datasetVersions {
		datasetVersionsResponse = append(datasetVersionsResponse, datasetmodels.DatasetBranchVersionResponse{
			UUID:    datasetVersion.UUID,
			Hash:    datasetVersion.Hash,
			Version: datasetVersion.Version,
			Branch: datasetmodels.DatasetBranchNameResponse{
				UUID: datasetVersion.Branch.UUID,
				Name: datasetVersion.Branch.Name,
			},
			Path:       datasetVersion.Path,
			SourceType: datasetVersion.SourceType,
			Lineage: datasetmodels.LineageResponse{
				UUID:    datasetVersion.Lineage.UUID,
				Lineage: datasetVersion.Lineage.Lineage,
			},
			CreatedBy: userorgmodels.UserHandleResponse{
				UUID:   datasetVersion.CreatedByUser.UUID,
				Name:   datasetVersion.CreatedByUser.Name,
				Avatar: datasetVersion.CreatedByUser.Avatar,
				Email:  datasetVersion.CreatedByUser.Email,
				Handle: datasetVersion.CreatedByUser.Handle,
			},
			CreatedAt: datasetVersion.CreatedAt,
			IsEmpty:   datasetVersion.IsEmpty,
		})
	}
	return datasetVersionsResponse, nil
}

func (ds *Datastore) GetDatasetBranchByName(orgId uuid.UUID, datasetName string, datasetBranchName string) (*datasetmodels.DatasetBranchResponse, error) {
	var datasetBranch datasetdbmodels.DatasetBranch
	dataset, err := ds.GetDatasetByName(orgId, datasetName)
	if err != nil {
		return nil, err
	}
	res := ds.DB.Where("name = ?", datasetBranchName).Where("dataset_uuid = ?", dataset.UUID).Preload("Dataset").Limit(1).Find(&datasetBranch)
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &datasetmodels.DatasetBranchResponse{
		UUID: datasetBranch.UUID,
		Name: datasetBranch.Name,
		Dataset: datasetmodels.DatasetNameResponse{
			UUID: datasetBranch.Dataset.UUID,
			Name: datasetBranch.Dataset.Name,
		},
		IsDefault: datasetBranch.IsDefault,
	}, nil
}

func (ds *Datastore) GetDatasetBranchByUUID(datasetBranchUUID uuid.UUID) (*datasetmodels.DatasetBranchResponse, error) {
	var datasetBranch datasetdbmodels.DatasetBranch
	err := ds.DB.Where("uuid = ?", datasetBranchUUID).Preload("Dataset").Find(&datasetBranch).Error
	if err != nil {
		return nil, err
	}
	return &datasetmodels.DatasetBranchResponse{
		UUID: datasetBranch.UUID,
		Name: datasetBranch.Name,
		Dataset: datasetmodels.DatasetNameResponse{
			UUID: datasetBranch.Dataset.UUID,
			Name: datasetBranch.Dataset.Name,
		},
		IsDefault: datasetBranch.IsDefault,
	}, nil
}

func (ds *Datastore) GetDatasetBranchAllVersions(datasetBranchUUID uuid.UUID) ([]datasetmodels.DatasetBranchVersionResponse, error) {
	var datasetVersions []datasetdbmodels.DatasetVersion
	err := ds.DB.Where("branch_uuid = ?", datasetBranchUUID).Preload("Lineage").Preload("Branch").Preload("CreatedByUser").Order("LENGTH(version) DESC").Order("version DESC").Find(&datasetVersions).Error
	if err != nil {
		return nil, err
	}
	var datasetVersionsResponse []datasetmodels.DatasetBranchVersionResponse
	for _, datasetVersion := range datasetVersions {
		datasetVersionsResponse = append(datasetVersionsResponse, datasetmodels.DatasetBranchVersionResponse{
			UUID:    datasetVersion.UUID,
			Hash:    datasetVersion.Hash,
			Version: datasetVersion.Version,
			Branch: datasetmodels.DatasetBranchNameResponse{
				UUID: datasetVersion.Branch.UUID,
				Name: datasetVersion.Branch.Name,
			},
			Path:       datasetVersion.Path,
			SourceType: datasetVersion.SourceType,
			Lineage: datasetmodels.LineageResponse{
				UUID:    datasetVersion.Lineage.UUID,
				Lineage: datasetVersion.Lineage.Lineage,
			},
			CreatedBy: userorgmodels.UserHandleResponse{
				UUID:   datasetVersion.CreatedByUser.UUID,
				Name:   datasetVersion.CreatedByUser.Name,
				Avatar: datasetVersion.CreatedByUser.Avatar,
				Email:  datasetVersion.CreatedByUser.Email,
				Handle: datasetVersion.CreatedByUser.Handle,
			},
			CreatedAt: datasetVersion.CreatedAt,
			IsEmpty:   datasetVersion.IsEmpty,
		})
	}
	return datasetVersionsResponse, nil
}

func (ds *Datastore) GetDatasetBranchVersion(datasetBranchUUID uuid.UUID, version string) (*datasetmodels.DatasetBranchVersionResponse, error) {
	var datasetVersion datasetdbmodels.DatasetVersion
	var res *gorm.DB
	if strings.ToLower(version) == "latest" {
		res = ds.DB.Preload("Branch").Where("branch_uuid = ?", datasetBranchUUID).Preload("Lineage").Preload("CreatedByUser").Order("created_at desc").Limit(1).Find(&datasetVersion)
	} else {
		res = ds.DB.Where("version = ?", version).Where("branch_uuid = ?", datasetBranchUUID).Preload("Branch").Preload("Lineage").Preload("CreatedByUser").Limit(1).Find(&datasetVersion)
	}
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &datasetmodels.DatasetBranchVersionResponse{
		UUID:    datasetVersion.UUID,
		Hash:    datasetVersion.Hash,
		Version: datasetVersion.Version,
		Branch: datasetmodels.DatasetBranchNameResponse{
			UUID: datasetVersion.Branch.UUID,
			Name: datasetVersion.Branch.Name,
		},
		Path:       datasetVersion.Path,
		SourceType: datasetVersion.SourceType,
		Lineage: datasetmodels.LineageResponse{
			UUID:    datasetVersion.Lineage.UUID,
			Lineage: datasetVersion.Lineage.Lineage,
		},
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   datasetVersion.CreatedByUser.UUID,
			Name:   datasetVersion.CreatedByUser.Name,
			Avatar: datasetVersion.CreatedByUser.Avatar,
			Email:  datasetVersion.CreatedByUser.Email,
			Handle: datasetVersion.CreatedByUser.Handle,
		},
		CreatedAt: datasetVersion.CreatedAt,
		IsEmpty:   datasetVersion.IsEmpty,
	}, nil
}

//////////////////////////////// LOG METHODS /////////////////////////////////

func (ds *Datastore) GetLogForModelVersion(modelVersionUUID uuid.UUID) ([]models.LogResponse, error) {
	var logs []dbmodels.Log
	err := ds.DB.Where("model_version_uuid = ?", modelVersionUUID).Preload("ModelVersion").Find(&logs).Error
	if err != nil {
		return nil, err
	}
	var logsResponse []models.LogResponse
	for _, log := range logs {
		logsResponse = append(logsResponse, models.LogResponse{
			Key:  log.Key,
			Data: log.Data,
			ModelVersion: modelmodels.ModelBranchVersionNameResponse{
				UUID:    log.ModelVersion.UUID,
				Version: log.ModelVersion.Version,
			},
		})
	}
	return logsResponse, nil
}

func (ds *Datastore) GetKeyLogForModelVersion(modelVersionUUID uuid.UUID, key string) ([]models.LogResponse, error) {
	var logs []dbmodels.Log
	err := ds.DB.Where("key = ?", key).Where("model_version_uuid = ?", modelVersionUUID).Preload("ModelVersion").Find(&logs).Error
	if err != nil {
		return nil, err
	}
	var logsResponse []models.LogResponse
	for _, log := range logs {
		logsResponse = append(logsResponse, models.LogResponse{
			Key:  log.Key,
			Data: log.Data,
			ModelVersion: modelmodels.ModelBranchVersionNameResponse{
				UUID:    log.ModelVersion.UUID,
				Version: log.ModelVersion.Version,
			},
		})
	}
	return logsResponse, nil
}

func (ds *Datastore) CreateLogForModelVersion(key string, data string, modelVersionUUID uuid.UUID) (*models.LogResponse, error) {
	var keyLog dbmodels.Log
	err := ds.DB.Where("key = ?", key).Where("model_version_uuid = ?", modelVersionUUID).First(&keyLog).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	log := dbmodels.Log{
		BaseModel: commondbmodels.BaseModel{
			UUID: keyLog.UUID,
		},
		Key:  key,
		Data: data,
		ModelVersion: modeldbmodels.ModelVersion{
			BaseModel: commondbmodels.BaseModel{
				UUID: modelVersionUUID,
			},
		},
	}
	err = ds.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&log).Preload("ModelVersion").Find(&log).Error
	if err != nil {
		return nil, err
	}
	return &models.LogResponse{
		Key:  log.Key,
		Data: log.Data,
		ModelVersion: modelmodels.ModelBranchVersionNameResponse{
			UUID:    log.ModelVersion.UUID,
			Version: log.ModelVersion.Version,
		},
	}, nil
}

func (ds *Datastore) GetLogForDatasetVersion(datasetVersion uuid.UUID) ([]models.LogResponse, error) {
	var logs []dbmodels.Log
	err := ds.DB.Where("dataset_version_uuid = ?", datasetVersion).Preload("DatasetVersion").Find(&logs).Error
	if err != nil {
		return nil, err
	}
	var logsResponse []models.LogResponse
	for _, log := range logs {
		logsResponse = append(logsResponse, models.LogResponse{
			Key:  log.Key,
			Data: log.Data,
			DatasetVersion: datasetmodels.DatasetBranchVersionNameResponse{
				UUID:    log.DatasetVersion.UUID,
				Version: log.DatasetVersion.Version,
			},
		})
	}
	return logsResponse, nil
}

func (ds *Datastore) GetKeyLogForDatasetVersion(datasetVersion uuid.UUID, key string) ([]models.LogResponse, error) {
	var logs []dbmodels.Log
	err := ds.DB.Where("key = ?", key).Where("dataset_version_uuid = ?", datasetVersion).Preload("DatasetVersion").Find(&logs).Error
	if err != nil {
		return nil, err
	}
	var logsResponse []models.LogResponse
	for _, log := range logs {
		logsResponse = append(logsResponse, models.LogResponse{
			Key:  log.Key,
			Data: log.Data,
			DatasetVersion: datasetmodels.DatasetBranchVersionNameResponse{
				UUID:    log.DatasetVersion.UUID,
				Version: log.DatasetVersion.Version,
			},
		})
	}
	return logsResponse, nil
}

func (ds *Datastore) CreateLogForDatasetVersion(key string, data string, datasetVersionUUID uuid.UUID) (*models.LogResponse, error) {
	var keyLog dbmodels.Log
	err := ds.DB.Where("key = ?", key).Where("dataset_version_uuid = ?", datasetVersionUUID).First(&keyLog).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	log := dbmodels.Log{
		BaseModel: commondbmodels.BaseModel{
			UUID: keyLog.UUID,
		},
		Key:  key,
		Data: data,
		DatasetVersion: datasetdbmodels.DatasetVersion{
			BaseModel: commondbmodels.BaseModel{
				UUID: datasetVersionUUID,
			},
		},
	}
	err = ds.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&log).Preload("DatasetVersion").Find(&log).Error
	if err != nil {
		return nil, err
	}
	return &models.LogResponse{
		Key:  log.Key,
		Data: log.Data,
		DatasetVersion: datasetmodels.DatasetBranchVersionNameResponse{
			UUID:    log.DatasetVersion.UUID,
			Version: log.DatasetVersion.Version,
		},
	}, nil
}

//////////////////////////////// ACTIVITY METHODS /////////////////////////////////

func (ds *Datastore) GetModelActivity(modelUUID uuid.UUID, category string) (*models.ActivityResponse, error) {
	var activity dbmodels.Activity
	res := ds.DB.Where("model_uuid = ?", modelUUID).Where("category = ?", category).Preload("Model").Preload("User").Limit(1).Find(&activity)
	if res.RowsAffected == 0 {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &models.ActivityResponse{
		UUID:     activity.UUID,
		Category: activity.Category,
		Activity: activity.Activity,
		Model: modelmodels.ModelNameResponse{
			UUID: activity.Model.UUID,
			Name: activity.Model.Name,
		},
		User: userorgmodels.UserHandleResponse{
			UUID: activity.User.UUID,
			Name: activity.User.Name,
		},
	}, nil
}

func (ds *Datastore) CreateModelActivity(modelUUID uuid.UUID, userUUID uuid.UUID, category string, activity string) (*models.ActivityResponse, error) {
	dbactivity := dbmodels.Activity{
		Category: category,
		Activity: activity,
		Model: modeldbmodels.Model{
			BaseModel: commondbmodels.BaseModel{
				UUID: modelUUID,
			},
		},
		User: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: userUUID,
			},
		},
	}
	err := ds.DB.Create(&dbactivity).Preload("Model").Preload("User").Find(&dbactivity).Error
	if err != nil {
		return nil, err
	}
	return &models.ActivityResponse{
		UUID:     dbactivity.UUID,
		Category: dbactivity.Category,
		Activity: dbactivity.Activity,
		Model: modelmodels.ModelNameResponse{
			UUID: dbactivity.Model.UUID,
			Name: dbactivity.Model.Name,
		},
		User: userorgmodels.UserHandleResponse{
			UUID: dbactivity.User.UUID,
			Name: dbactivity.User.Name,
		},
	}, nil
}

func (ds *Datastore) UpdateModelActivity(activityUUID uuid.UUID, updatedAttributes map[string]string) (*models.ActivityResponse, error) {
	var activity dbmodels.Activity
	err := ds.DB.Where("uuid = ?", activityUUID).Preload("Model").Preload("User").Limit(1).Find(&activity).Error
	if err != nil {
		return nil, err
	}
	err = ds.DB.Model(&activity).Updates(updatedAttributes).Error
	if err != nil {
		return nil, err
	}
	return &models.ActivityResponse{
		UUID:     activity.UUID,
		Category: activity.Category,
		Activity: activity.Activity,
		Model: modelmodels.ModelNameResponse{
			UUID: activity.Model.UUID,
			Name: activity.Model.Name,
		},
		User: userorgmodels.UserHandleResponse{
			UUID: activity.User.UUID,
			Name: activity.User.Name,
		},
	}, nil
}

func (ds *Datastore) DeleteModelActivity(activityUUID uuid.UUID) error {
	var activity dbmodels.Activity
	err := ds.DB.Where("uuid = ?", activityUUID).Limit(1).Find(&activity).Error
	if err != nil {
		return err
	}
	err = ds.DB.Delete(&activity).Error
	if err != nil {
		return err
	}
	return nil
}

func (ds *Datastore) GetDatasetActivity(datasetUUID uuid.UUID, category string) (*models.ActivityResponse, error) {
	var activity dbmodels.Activity
	res := ds.DB.Where("dataset_uuid = ?", datasetUUID).Where("category = ?", category).Preload("Dataset").Preload("User").Limit(1).Find(&activity)
	if res.RowsAffected == 0 {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &models.ActivityResponse{
		UUID:     activity.UUID,
		Category: activity.Category,
		Activity: activity.Activity,
		Dataset: datasetmodels.DatasetNameResponse{
			UUID: activity.Dataset.UUID,
			Name: activity.Dataset.Name,
		},
		User: userorgmodels.UserHandleResponse{
			UUID: activity.User.UUID,
			Name: activity.User.Name,
		},
	}, nil
}

func (ds *Datastore) CreateDatasetActivity(datasetUUID uuid.UUID, userUUID uuid.UUID, category string, activity string) (*models.ActivityResponse, error) {
	dbactivity := dbmodels.Activity{
		Category: category,
		Activity: activity,
		Dataset: datasetdbmodels.Dataset{
			BaseModel: commondbmodels.BaseModel{
				UUID: datasetUUID,
			},
		},
		User: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: userUUID,
			},
		},
	}
	err := ds.DB.Create(&dbactivity).Preload("Dataset").Preload("User").Find(&dbactivity).Error
	if err != nil {
		return nil, err
	}
	return &models.ActivityResponse{
		UUID:     dbactivity.UUID,
		Category: dbactivity.Category,
		Activity: dbactivity.Activity,
		Dataset: datasetmodels.DatasetNameResponse{
			UUID: dbactivity.Dataset.UUID,
			Name: dbactivity.Dataset.Name,
		},
		User: userorgmodels.UserHandleResponse{
			UUID: dbactivity.User.UUID,
			Name: dbactivity.User.Name,
		},
	}, nil
}

func (ds *Datastore) UpdateDatasetActivity(activityUUID uuid.UUID, updatedAttributes map[string]string) (*models.ActivityResponse, error) {
	var activity dbmodels.Activity
	err := ds.DB.Where("uuid = ?", activityUUID).Preload("Dataset").Preload("User").Limit(1).Find(&activity).Error
	if err != nil {
		return nil, err
	}
	err = ds.DB.Model(&activity).Updates(updatedAttributes).Error
	if err != nil {
		return nil, err
	}
	return &models.ActivityResponse{
		UUID:     activity.UUID,
		Category: activity.Category,
		Activity: activity.Activity,
		Dataset: datasetmodels.DatasetNameResponse{
			UUID: activity.Dataset.UUID,
			Name: activity.Dataset.Name,
		},
		User: userorgmodels.UserHandleResponse{
			UUID: activity.User.UUID,
			Name: activity.User.Name,
		},
	}, nil
}

func (ds *Datastore) DeleteDatasetActivity(activityUUID uuid.UUID) error {
	var activity dbmodels.Activity
	err := ds.DB.Where("uuid = ?", activityUUID).Limit(1).Find(&activity).Error
	if err != nil {
		return err
	}
	err = ds.DB.Delete(&activity).Error
	if err != nil {
		return err
	}
	return nil
}

/////////////////////////////// SECRET API METHODS ///////////////////////////////

func (ds *Datastore) GetOrganizationSecrets(orgId uuid.UUID) ([]string, error) {
	var secretNames []string
	res := ds.DB.Model(&userorgdbmodels.Secret{}).Where("org_uuid = ?", orgId).Distinct().Pluck("name", &secretNames)
	if res.RowsAffected == 0 {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return secretNames, nil
}

func (ds *Datastore) GetSecretByName(orgId uuid.UUID, secretName string) (*commonmodels.SourceSecrets, error) {
	var secrets []userorgdbmodels.Secret
	res := ds.DB.Where("org_uuid = ? and name = ?", orgId, secretName).Find(&secrets)
	if res.RowsAffected == 0 {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	if len(secrets) == 0 {
		return nil, fmt.Errorf("no secrets found for secret name %s", secretName)
	}
	var sourceSecret commonmodels.SourceSecrets
	secretType := strings.ToUpper(strings.Split(secrets[0].Key, "_")[0])
	switch secretType {
	case "S3":
		for _, secret := range secrets {
			switch strings.ToUpper(secret.Key) {
			case "S3_ACCESS_KEY_ID":
				sourceSecret.AccessKeyId = secret.Value
			case "S3_ACCESS_KEY_SECRET":
				sourceSecret.AccessKeySecret = secret.Value
			case "S3_BUCKET_NAME":
				sourceSecret.BucketName = secret.Value
			case "S3_BUCKET_LOCATION":
				sourceSecret.BucketLocation = secret.Value
			case "S3_PUBLIC_URL":
				sourceSecret.PublicURL = secret.Value
			}

		}
		if sourceSecret.AccessKeyId == "" || sourceSecret.AccessKeySecret == "" || sourceSecret.BucketName == "" || sourceSecret.BucketLocation == "" {
			return nil, fmt.Errorf("s3 secrets not found or invalid")
		}
	case "R2":
		for _, secret := range secrets {
			switch strings.ToUpper(secret.Key) {
			case "R2_ACCOUNT_ID":
				sourceSecret.AccountId = secret.Value
			case "R2_ACCESS_KEY_ID":
				sourceSecret.AccessKeyId = secret.Value
			case "R2_ACCESS_KEY_SECRET":
				sourceSecret.AccessKeySecret = secret.Value
			case "R2_BUCKET_NAME":
				sourceSecret.BucketName = secret.Value
			case "R2_PUBLIC_URL":
				sourceSecret.PublicURL = secret.Value
			}
		}
		if sourceSecret.AccountId == "" || sourceSecret.AccessKeyId == "" || sourceSecret.AccessKeySecret == "" || sourceSecret.BucketName == "" || sourceSecret.PublicURL == "" {
			return nil, fmt.Errorf("r2 secrets not found or invalid")
		}
	}
	sourceSecret.SourceType = secretType
	return &sourceSecret, nil
}

func (ds *Datastore) CreateR2Secrets(orgId uuid.UUID, secretName string, accountId string, accessKeyId string, accessKeySecret string, bucketName string, endpoint string) (*commonmodels.SourceSecrets, error) {
	secret := userorgdbmodels.Secret{
		Org: userorgdbmodels.Organization{
			BaseModel: commondbmodels.BaseModel{
				UUID: orgId,
			},
		},
		Name: secretName,
	}
	publicURL := fmt.Sprintf("https://%s", endpoint)
	err := ds.DB.Transaction(func(tx *gorm.DB) error {
		secret.Key = "R2_ACCOUNT_ID"
		secret.Value = accountId
		err := tx.Create(&secret).Error
		if err != nil {
			return err
		}
		secret.UUID = uuid.Nil
		secret.Key = "R2_ACCESS_KEY_ID"
		secret.Value = accessKeyId
		err = tx.Create(&secret).Error
		if err != nil {
			return err
		}
		secret.UUID = uuid.Nil
		secret.Key = "R2_ACCESS_KEY_SECRET"
		secret.Value = accessKeySecret
		err = tx.Create(&secret).Error
		if err != nil {
			return err
		}
		secret.UUID = uuid.Nil
		secret.Key = "R2_BUCKET_NAME"
		secret.Value = bucketName
		err = tx.Create(&secret).Error
		if err != nil {
			return err
		}
		secret.UUID = uuid.Nil
		secret.Key = "R2_PUBLIC_URL"
		secret.Value = publicURL
		err = tx.Create(&secret).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &commonmodels.SourceSecrets{
		AccountId: accountId,
		// AccessKeyId:     accessKeyId,
		// AccessKeySecret: accessKeySecret,
		BucketName: bucketName,
		PublicURL:  publicURL,
	}, nil
}

func (ds *Datastore) CreateS3Secrets(orgId uuid.UUID, secretName string, accessKeyId string, accessKeySecret string, bucketName string, bucketLocation string) (*commonmodels.SourceSecrets, error) {
	secret := userorgdbmodels.Secret{
		Org: userorgdbmodels.Organization{
			BaseModel: commondbmodels.BaseModel{
				UUID: orgId,
			},
		},
		Name: secretName,
	}
	publicURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", bucketName, bucketLocation)
	err := ds.DB.Transaction(func(tx *gorm.DB) error {
		secret.Key = "S3_ACCESS_KEY_ID"
		secret.Value = accessKeyId
		err := tx.Create(&secret).Error
		if err != nil {
			return err
		}
		secret.UUID = uuid.Nil
		secret.Key = "S3_ACCESS_KEY_SECRET"
		secret.Value = accessKeySecret
		err = tx.Create(&secret).Error
		if err != nil {
			return err
		}
		secret.UUID = uuid.Nil
		secret.Key = "S3_BUCKET_NAME"
		secret.Value = bucketName
		err = tx.Create(&secret).Error
		if err != nil {
			return err
		}
		secret.UUID = uuid.Nil
		secret.Key = "S3_BUCKET_LOCATION"
		secret.Value = bucketLocation
		err = tx.Create(&secret).Error
		if err != nil {
			return err
		}
		secret.UUID = uuid.Nil
		secret.Key = "S3_PUBLIC_URL"
		secret.Value = publicURL
		err = tx.Create(&secret).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &commonmodels.SourceSecrets{
		// AccessKeyId:     accessKeyId,
		// AccessKeySecret: accessKeySecret,
		BucketName:     bucketName,
		BucketLocation: bucketLocation,
		PublicURL:      publicURL,
	}, nil
}

func (ds *Datastore) DeleteSecrets(orgId uuid.UUID, secretName string) error {
	var secrets []userorgdbmodels.Secret
	err := ds.DB.Where("org_uuid = ?", orgId).Where("name = ?", secretName).Delete(&secrets).Error
	if err != nil {
		return err
	}
	return nil
}

/////////////////////////////// REVIEW API METHODS ///////////////////////////////

func (ds *Datastore) GetModelReview(reviewUUID uuid.UUID) (*modelmodels.ModelReviewResponse, error) {
	var review modeldbmodels.ModelReview
	err := ds.DB.Preload("Model").Preload("FromBranch").Preload("FromBranchVersion").Preload("ToBranch").Preload("CreatedByUser").Where("uuid = ?", reviewUUID).Find(&review).Error
	if err != nil {
		return nil, err
	}
	return &modelmodels.ModelReviewResponse{
		UUID: review.UUID,
		Model: modelmodels.ModelNameResponse{
			UUID: review.Model.UUID,
			Name: review.Model.Name,
		},
		FromBranch: modelmodels.ModelBranchNameResponse{
			UUID: review.FromBranch.UUID,
			Name: review.FromBranch.Name,
		},
		FromBranchVersion: modelmodels.ModelBranchVersionNameResponse{
			UUID:    review.FromBranchVersion.UUID,
			Version: review.FromBranchVersion.Version,
		},
		ToBranch: modelmodels.ModelBranchNameResponse{
			UUID: review.ToBranch.UUID,
			Name: review.ToBranch.Name,
		},
		Title:       review.Title,
		Description: review.Description,
		IsComplete:  review.IsComplete,
		IsAccepted:  review.IsAccepted,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   review.CreatedByUser.UUID,
			Handle: review.CreatedByUser.Handle,
			Name:   review.CreatedByUser.Name,
			Avatar: review.CreatedByUser.Avatar,
			Email:  review.CreatedByUser.Email,
		},
	}, nil
}

func (ds *Datastore) GetModelReviews(modelUUID uuid.UUID) ([]modelmodels.ModelReviewResponse, error) {
	var reviews []modeldbmodels.ModelReview
	err := ds.DB.Preload("Model").Preload("FromBranch").Preload("FromBranchVersion").Preload("ToBranch").Preload("CreatedByUser").Where("model_uuid = ?", modelUUID).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	var reviewResponses []modelmodels.ModelReviewResponse
	for _, review := range reviews {
		reviewResponses = append(reviewResponses, modelmodels.ModelReviewResponse{
			UUID: review.UUID,
			Model: modelmodels.ModelNameResponse{
				UUID: review.Model.UUID,
				Name: review.Model.Name,
			},
			FromBranch: modelmodels.ModelBranchNameResponse{
				UUID: review.FromBranch.UUID,
				Name: review.FromBranch.Name,
			},
			FromBranchVersion: modelmodels.ModelBranchVersionNameResponse{
				UUID:    review.FromBranchVersion.UUID,
				Version: review.FromBranchVersion.Version,
			},
			ToBranch: modelmodels.ModelBranchNameResponse{
				UUID: review.ToBranch.UUID,
				Name: review.ToBranch.Name,
			},
			Title:       review.Title,
			Description: review.Description,
			IsComplete:  review.IsComplete,
			IsAccepted:  review.IsAccepted,
			CreatedBy: userorgmodels.UserHandleResponse{
				UUID:   review.CreatedByUser.UUID,
				Handle: review.CreatedByUser.Handle,
				Name:   review.CreatedByUser.Name,
				Avatar: review.CreatedByUser.Avatar,
				Email:  review.CreatedByUser.Email,
			},
		})
	}
	return reviewResponses, nil
}

func (ds *Datastore) CreateModelReview(modelUUID uuid.UUID, userUUID uuid.UUID, fromBranch uuid.UUID, fromBranchVersion uuid.UUID, toBranch uuid.UUID, title string, desc string, isComplete bool, isAccepted bool) (*modelmodels.ModelReviewResponse, error) {
	review := modeldbmodels.ModelReview{
		Model: modeldbmodels.Model{
			BaseModel: commondbmodels.BaseModel{
				UUID: modelUUID,
			},
		},
		FromBranch: modeldbmodels.ModelBranch{
			BaseModel: commondbmodels.BaseModel{
				UUID: fromBranch,
			},
		},
		FromBranchVersion: modeldbmodels.ModelVersion{
			BaseModel: commondbmodels.BaseModel{
				UUID: fromBranchVersion,
			},
		},
		ToBranch: modeldbmodels.ModelBranch{
			BaseModel: commondbmodels.BaseModel{
				UUID: toBranch,
			},
		},
		CreatedByUser: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: userUUID,
			},
		},
		Title:       title,
		Description: desc,
		IsComplete:  isComplete,
		IsAccepted:  isAccepted,
	}
	err := ds.DB.Create(&review).Find(&review).Error
	if err != nil {
		return nil, err
	}
	return &modelmodels.ModelReviewResponse{
		UUID: review.UUID,
		Model: modelmodels.ModelNameResponse{
			UUID: review.Model.UUID,
			Name: review.Model.Name,
		},
		FromBranch: modelmodels.ModelBranchNameResponse{
			UUID: review.FromBranch.UUID,
			Name: review.FromBranch.Name,
		},
		FromBranchVersion: modelmodels.ModelBranchVersionNameResponse{
			UUID:    review.FromBranchVersion.UUID,
			Version: review.FromBranchVersion.Version,
		},
		ToBranch: modelmodels.ModelBranchNameResponse{
			UUID: review.ToBranch.UUID,
			Name: review.ToBranch.Name,
		},
		Title:       review.Title,
		Description: review.Description,
		IsComplete:  review.IsComplete,
		IsAccepted:  review.IsAccepted,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   review.CreatedByUser.UUID,
			Handle: review.CreatedByUser.Handle,
			Name:   review.CreatedByUser.Name,
			Avatar: review.CreatedByUser.Avatar,
			Email:  review.CreatedByUser.Email,
		},
	}, nil
}

func (ds *Datastore) UpdateModelReview(reviewUUID uuid.UUID, updatedAttributes map[string]any) (*modelmodels.ModelReviewResponse, error) {
	var review modeldbmodels.ModelReview
	err := ds.DB.Preload("Model").Preload("FromBranch").Preload("FromBranchVersion").Preload("ToBranch").Preload("CreatedByUser").Where("uuid = ?", reviewUUID).Find(&review).Error
	if err != nil {
		return nil, err
	}
	if review.IsComplete {
		return nil, errors.New("review already complete")
	}
	alreadyAccepted := review.IsAccepted
	for key, value := range updatedAttributes {
		switch key {
		case "title":
			review.Title = value.(string)
		case "description":
			review.Description = value.(string)
		case "is_complete":
			review.IsComplete = value.(bool)
		case "is_accepted":
			// CANNOT UNACCEPT ALREADY ACCEPTED REVIEW
			if !alreadyAccepted {
				review.IsAccepted = value.(bool)
			}
		}
	}
	err = ds.DB.Transaction(func(tx *gorm.DB) error {
		// CHECK IF REVIEW ACCEPTED
		if review.IsAccepted && !alreadyAccepted {
			// make a new version in to_branch from from_branch specified version
			_, err := ds.MigrateModelVersionBranch(review.FromBranchVersion.UUID, review.ToBranch.UUID)
			if err != nil {
				return err
			}
		}
		result := tx.Save(&review).Find(&review)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &modelmodels.ModelReviewResponse{
		UUID: review.UUID,
		Model: modelmodels.ModelNameResponse{
			UUID: review.Model.UUID,
			Name: review.Model.Name,
		},
		FromBranch: modelmodels.ModelBranchNameResponse{
			UUID: review.FromBranch.UUID,
			Name: review.FromBranch.Name,
		},
		ToBranch: modelmodels.ModelBranchNameResponse{
			UUID: review.ToBranch.UUID,
			Name: review.ToBranch.Name,
		},
		Title:       review.Title,
		Description: review.Description,
		IsComplete:  review.IsComplete,
		IsAccepted:  review.IsAccepted,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   review.CreatedByUser.UUID,
			Handle: review.CreatedByUser.Handle,
			Name:   review.CreatedByUser.Name,
			Avatar: review.CreatedByUser.Avatar,
			Email:  review.CreatedByUser.Email,
		},
	}, nil
}

func (ds *Datastore) GetDatasetReview(reviewUUID uuid.UUID) (*datasetmodels.DatasetReviewResponse, error) {
	var review datasetdbmodels.DatasetReview
	err := ds.DB.Preload("Dataset").Preload("FromBranch").Preload("FromBranchVersion").Preload("ToBranch").Preload("CreatedByUser").Where("uuid = ?", reviewUUID).Find(&review).Error
	if err != nil {
		return nil, err
	}
	return &datasetmodels.DatasetReviewResponse{
		UUID: review.UUID,
		Dataset: datasetmodels.DatasetNameResponse{
			UUID: review.Dataset.UUID,
			Name: review.Dataset.Name,
		},
		FromBranch: datasetmodels.DatasetBranchNameResponse{
			UUID: review.FromBranch.UUID,
			Name: review.FromBranch.Name,
		},
		FromBranchVersion: datasetmodels.DatasetBranchVersionNameResponse{
			UUID:    review.FromBranchVersion.UUID,
			Version: review.FromBranchVersion.Version,
		},
		ToBranch: datasetmodels.DatasetBranchNameResponse{
			UUID: review.ToBranch.UUID,
			Name: review.ToBranch.Name,
		},
		Title:       review.Title,
		Description: review.Description,
		IsComplete:  review.IsComplete,
		IsAccepted:  review.IsAccepted,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   review.CreatedByUser.UUID,
			Handle: review.CreatedByUser.Handle,
			Name:   review.CreatedByUser.Name,
			Avatar: review.CreatedByUser.Avatar,
			Email:  review.CreatedByUser.Email,
		},
	}, nil
}

func (ds *Datastore) GetDatasetReviews(datasetUUID uuid.UUID) ([]datasetmodels.DatasetReviewResponse, error) {
	var reviews []datasetdbmodels.DatasetReview
	err := ds.DB.Preload("Dataset").Preload("FromBranch").Preload("FromBranchVersion").Preload("ToBranch").Preload("CreatedByUser").Where("dataset_uuid = ?", datasetUUID).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	var reviewResponses []datasetmodels.DatasetReviewResponse
	for _, review := range reviews {
		reviewResponses = append(reviewResponses, datasetmodels.DatasetReviewResponse{
			UUID: review.UUID,
			Dataset: datasetmodels.DatasetNameResponse{
				UUID: review.Dataset.UUID,
				Name: review.Dataset.Name,
			},
			FromBranch: datasetmodels.DatasetBranchNameResponse{
				UUID: review.FromBranch.UUID,
				Name: review.FromBranch.Name,
			},
			FromBranchVersion: datasetmodels.DatasetBranchVersionNameResponse{
				UUID:    review.FromBranchVersion.UUID,
				Version: review.FromBranchVersion.Version,
			},
			ToBranch: datasetmodels.DatasetBranchNameResponse{
				UUID: review.ToBranch.UUID,
				Name: review.ToBranch.Name,
			},
			Title:       review.Title,
			Description: review.Description,
			IsComplete:  review.IsComplete,
			IsAccepted:  review.IsAccepted,
			CreatedBy: userorgmodels.UserHandleResponse{
				UUID:   review.CreatedByUser.UUID,
				Handle: review.CreatedByUser.Handle,
				Name:   review.CreatedByUser.Name,
				Avatar: review.CreatedByUser.Avatar,
				Email:  review.CreatedByUser.Email,
			},
		})
	}
	return reviewResponses, nil
}

func (ds *Datastore) CreateDatasetReview(datasetUUID uuid.UUID, userUUID uuid.UUID, fromBranch uuid.UUID, fromBranchVerison uuid.UUID, toBranch uuid.UUID, title string, desc string, isComplete bool, isAccepted bool) (*datasetmodels.DatasetReviewResponse, error) {
	review := datasetdbmodels.DatasetReview{
		Dataset: datasetdbmodels.Dataset{
			BaseModel: commondbmodels.BaseModel{
				UUID: datasetUUID,
			},
		},
		FromBranch: datasetdbmodels.DatasetBranch{
			BaseModel: commondbmodels.BaseModel{
				UUID: fromBranch,
			},
		},
		FromBranchVersion: datasetdbmodels.DatasetVersion{
			BaseModel: commondbmodels.BaseModel{
				UUID: fromBranchVerison,
			},
		},
		ToBranch: datasetdbmodels.DatasetBranch{
			BaseModel: commondbmodels.BaseModel{
				UUID: toBranch,
			},
		},
		CreatedByUser: userorgdbmodels.User{
			BaseModel: commondbmodels.BaseModel{
				UUID: userUUID,
			},
		},
		Title:       title,
		Description: desc,
		IsComplete:  isComplete,
		IsAccepted:  isAccepted,
	}
	err := ds.DB.Create(&review).Find(&review).Error
	if err != nil {
		return nil, err
	}
	return &datasetmodels.DatasetReviewResponse{
		UUID: review.UUID,
		Dataset: datasetmodels.DatasetNameResponse{
			UUID: review.Dataset.UUID,
			Name: review.Dataset.Name,
		},
		FromBranch: datasetmodels.DatasetBranchNameResponse{
			UUID: review.FromBranch.UUID,
			Name: review.FromBranch.Name,
		},
		FromBranchVersion: datasetmodels.DatasetBranchVersionNameResponse{
			UUID:    review.FromBranchVersion.UUID,
			Version: review.FromBranchVersion.Version,
		},
		ToBranch: datasetmodels.DatasetBranchNameResponse{
			UUID: review.ToBranch.UUID,
			Name: review.ToBranch.Name,
		},
		Title:       review.Title,
		Description: review.Description,
		IsComplete:  review.IsComplete,
		IsAccepted:  review.IsAccepted,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   review.CreatedByUser.UUID,
			Handle: review.CreatedByUser.Handle,
			Name:   review.CreatedByUser.Name,
			Avatar: review.CreatedByUser.Avatar,
			Email:  review.CreatedByUser.Email,
		},
	}, nil
}

func (ds *Datastore) UpdateDatasetReview(reviewUUID uuid.UUID, updatedAttributes map[string]any) (*datasetmodels.DatasetReviewResponse, error) {
	var review datasetdbmodels.DatasetReview
	err := ds.DB.Preload("Dataset").Preload("FromBranch").Preload("FromBranchVersion").Preload("ToBranch").Preload("CreatedByUser").Where("uuid = ?", reviewUUID).Find(&review).Error
	if err != nil {
		return nil, err
	}
	if review.IsComplete {
		return nil, errors.New("review already complete")
	}
	alreadyAccepted := review.IsAccepted
	for key, value := range updatedAttributes {
		switch key {
		case "title":
			review.Title = value.(string)
		case "description":
			review.Description = value.(string)
		case "is_complete":
			review.IsComplete = value.(bool)
		case "is_accepted":
			// CANNOT UNACCEPT ALREADY ACCEPTED REVIEW
			if !alreadyAccepted {
				review.IsAccepted = value.(bool)
			}
		}
	}
	err = ds.DB.Transaction(func(tx *gorm.DB) error {
		// CHECK IF REVIEW ACCEPTED
		if review.IsAccepted && !alreadyAccepted {
			// make a new version in "to_branch" from "from_branch" specified version
			// also migrate all corresponsing logs
			_, err := ds.MigrateDatasetVersionBranch(review.FromBranchVersion.UUID, review.ToBranch.UUID)
			if err != nil {
				return err
			}
		}
		result := tx.Save(&review).Find(&review)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &datasetmodels.DatasetReviewResponse{
		UUID: review.UUID,
		Dataset: datasetmodels.DatasetNameResponse{
			UUID: review.Dataset.UUID,
			Name: review.Dataset.Name,
		},
		FromBranch: datasetmodels.DatasetBranchNameResponse{
			UUID: review.FromBranch.UUID,
			Name: review.FromBranch.Name,
		},
		ToBranch: datasetmodels.DatasetBranchNameResponse{
			UUID: review.ToBranch.UUID,
			Name: review.ToBranch.Name,
		},
		Title:       review.Title,
		Description: review.Description,
		IsComplete:  review.IsComplete,
		IsAccepted:  review.IsAccepted,
		CreatedBy: userorgmodels.UserHandleResponse{
			UUID:   review.CreatedByUser.UUID,
			Handle: review.CreatedByUser.Handle,
			Name:   review.CreatedByUser.Name,
			Avatar: review.CreatedByUser.Avatar,
			Email:  review.CreatedByUser.Email,
		},
	}, nil
}
