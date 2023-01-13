package impl

import (
	"context"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/PureML-Inc/PureML/server/config"
	"github.com/PureML-Inc/PureML/server/datastore/dbmodels"
	"github.com/PureML-Inc/PureML/server/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	uuid "github.com/satori/go.uuid"
	"github.com/teris-io/shortid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLiteDatastore() *SQLiteDatastore {
	db, err := gorm.Open(sqlite.Open("db/pureml.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("Error connecting to database")
	}
	err = db.AutoMigrate(
		&dbmodels.Activity{},
		&dbmodels.Dataset{},
		&dbmodels.DatasetBranch{},
		&dbmodels.DatasetReview{},
		&dbmodels.DatasetUser{},
		&dbmodels.DatasetVersion{},
		&dbmodels.Lineage{},
		&dbmodels.Log{},
		&dbmodels.Model{},
		&dbmodels.ModelBranch{},
		&dbmodels.ModelReview{},
		&dbmodels.ModelUser{},
		&dbmodels.ModelVersion{},
		&dbmodels.Organization{},
		&dbmodels.Path{},
		&dbmodels.User{},
		&dbmodels.UserOrganizations{},
	)
	if err != nil {
		return &SQLiteDatastore{}
	}
	return &SQLiteDatastore{
		DB: db,
	}
}

type SQLiteDatastore struct {
	DB *gorm.DB
}

//////////////////////////// ORGANIZATION METHODS ////////////////////////////

func (ds *SQLiteDatastore) GetAllAdminOrgs() ([]models.OrganizationResponse, error) {
	var organizations []dbmodels.Organization
	ds.DB.Find(&organizations)
	var responseOrganizations []models.OrganizationResponse
	for _, org := range organizations {
		responseOrganizations = append(responseOrganizations, models.OrganizationResponse{
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

func (ds *SQLiteDatastore) GetOrgByID(orgId uuid.UUID) (*models.OrganizationResponse, error) {
	org := dbmodels.Organization{
		BaseModel: dbmodels.BaseModel{
			UUID: orgId,
		},
	}
	result := ds.DB.Limit(1).Find(&org)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.OrganizationResponse{
		UUID:        org.UUID,
		Name:        org.Name,
		Handle:      org.Handle,
		Avatar:      org.Avatar,
		Description: org.Description,
		JoinCode:    org.JoinCode,
	}, nil
}

func (ds *SQLiteDatastore) GetOrgByJoinCode(joinCode string) (*models.OrganizationResponse, error) {
	var org dbmodels.Organization
	result := ds.DB.Where("join_code = ?", joinCode).Limit(1).Find(&org)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.OrganizationResponse{
		UUID:        org.UUID,
		Name:        org.Name,
		Handle:      org.Handle,
		Avatar:      org.Avatar,
		Description: org.Description,
		JoinCode:    org.JoinCode,
	}, nil
}

func (ds *SQLiteDatastore) CreateOrgFromEmail(email string, orgName string, orgDesc string, orgHandle string) (*models.OrganizationResponse, error) {
	org := dbmodels.Organization{
		Name:         orgName,
		Handle:       orgHandle,
		Avatar:       "",
		Description:  orgDesc,
		JoinCode:     shortid.MustGenerate(),
		APITokenHash: "",
	}
	user := dbmodels.User{
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
		userOrg := dbmodels.UserOrganizations{
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
	return &models.OrganizationResponse{
		UUID:        org.UUID,
		Name:        org.Name,
		Handle:      org.Handle,
		Avatar:      org.Avatar,
		Description: org.Description,
		JoinCode:    org.JoinCode,
	}, nil
}

func (ds *SQLiteDatastore) GetOrgByHandle(handle string) (*models.OrganizationResponse, error) {
	var org dbmodels.Organization
	result := ds.DB.Where("handle = ?", handle).Limit(1).Find(&org)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.OrganizationResponse{
		UUID:        org.UUID,
		Name:        org.Name,
		Handle:      org.Handle,
		Avatar:      org.Avatar,
		Description: org.Description,
		JoinCode:    org.JoinCode,
	}, nil
}

func (ds *SQLiteDatastore) GetUserOrganizationsByEmail(email string) ([]models.UserOrganizationsResponse, error) {
	var orgs []models.UserOrganizationsResponse
	var tableOrgs []struct {
		UUID   uuid.UUID
		Handle string
		Name   string
		Avatar string
		Role   string
	}
	result := ds.DB.Table("organizations").Select("organizations.uuid, organizations.handle, organizations.name, organizations.avatar, user_organizations.role").Joins("JOIN user_organizations ON user_organizations.organization_uuid = organizations.uuid").Joins("JOIN users ON users.uuid = user_organizations.user_uuid").Where("users.email = ?", email).Scan(&tableOrgs)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, org := range tableOrgs {
		orgs = append(orgs, models.UserOrganizationsResponse{
			Org: models.OrganizationHandleResponse{
				UUID:   org.UUID,
				Handle: org.Handle,
				Name:   org.Name,
				Avatar: org.Avatar,
			},
			Role: org.Role,
		})
	}
	return orgs, nil
}

func (ds *SQLiteDatastore) GetUserOrganizationByOrgIdAndEmail(orgId uuid.UUID, email string) (*models.UserOrganizationsResponse, error) {
	var org models.UserOrganizationsResponse
	result := ds.DB.Table("organizations").Select("organizations.uuid, organizations.handle, organizations.name, organizations.avatar, user_organization.role").Joins("JOIN user_organizations ON user_organizations.organization_uuid = organizations.uuid").Joins("JOIN users ON users.uuid = user_organizations.user_uuid").Where("users.email = ?", email).Where("organizations.uuid = ?", orgId).Scan(&org)
	if result.Error != nil {
		return nil, result.Error
	}
	return &org, nil
}

func (ds *SQLiteDatastore) CreateUserOrganizationFromEmailAndOrgId(email string, orgId uuid.UUID) (*models.UserOrganizationsResponse, error) {
	var org dbmodels.Organization
	result := ds.DB.First(&org, orgId)
	if result.Error != nil {
		return nil, result.Error
	}
	var user dbmodels.User
	result = ds.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	userOrganization := dbmodels.UserOrganizations{
		OrganizationUUID: org.UUID,
		UserUUID:         user.UUID,
		Role:             "member",
	}
	result = ds.DB.Create(&userOrganization)
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.UserOrganizationsResponse{
		Org: models.OrganizationHandleResponse{
			UUID:   org.UUID,
			Name:   org.Name,
			Handle: org.Handle,
			Avatar: org.Avatar,
		},
		Role: userOrganization.Role,
	}, nil
}

func (ds *SQLiteDatastore) DeleteUserOrganizationFromEmailAndOrgId(email string, orgId uuid.UUID) error {
	var org dbmodels.Organization
	result := ds.DB.First(&org, orgId)
	if result.Error != nil {
		return result.Error
	}
	var user dbmodels.User
	result = ds.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return result.Error
	}
	result = ds.DB.Where("organization_uuid = ?", org.UUID).Where("user_uuid = ?", user.UUID).Delete(&dbmodels.UserOrganizations{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ds *SQLiteDatastore) CreateUserOrganizationFromEmailAndJoinCode(email string, joinCode string) (*models.UserOrganizationsResponse, error) {
	var org dbmodels.Organization
	result := ds.DB.Where("join_code = ?", joinCode).First(&org)
	if result.Error != nil {
		return nil, result.Error
	}
	var user dbmodels.User
	result = ds.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	userOrganization := dbmodels.UserOrganizations{
		OrganizationUUID: org.UUID,
		UserUUID:         user.UUID,
		Role:             "member",
	}
	result = ds.DB.Create(&userOrganization)
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.UserOrganizationsResponse{
		Org: models.OrganizationHandleResponse{
			UUID:   org.UUID,
			Name:   org.Name,
			Handle: org.Handle,
			Avatar: org.Avatar,
		},
		Role: userOrganization.Role,
	}, nil
}

func (ds *SQLiteDatastore) UpdateOrg(orgId uuid.UUID, orgName string, orgDesc string, orgAvatar string) (*models.OrganizationResponse, error) {
	var org dbmodels.Organization
	result := ds.DB.First(&org, orgId)
	if result.Error != nil {
		return nil, result.Error
	}
	org.Name = orgName
	org.Description = orgDesc
	org.Avatar = orgAvatar
	result = ds.DB.Save(&org)
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.OrganizationResponse{
		UUID:        org.UUID,
		Name:        org.Name,
		Handle:      org.Handle,
		Avatar:      org.Avatar,
		Description: org.Description,
		JoinCode:    org.JoinCode,
	}, nil
}

/////////////////////////////// USER METHODS /////////////////////////////////

func (ds *SQLiteDatastore) GetUserByEmail(email string) (*models.UserResponse, error) {
	var user dbmodels.User
	result := ds.DB.Where("email = ?", email).Limit(1).Find(&user)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.UserResponse{
		UUID:     user.UUID,
		Name:     user.Name,
		Email:    user.Email,
		Handle:   user.Handle,
		Bio:      user.Bio,
		Avatar:   user.Avatar,
		Password: user.Password,
	}, nil
}

func (ds *SQLiteDatastore) GetUserByHandle(handle string) (*models.UserResponse, error) {
	var user dbmodels.User
	result := ds.DB.Where("handle = ?", handle).Limit(1).Find(&user)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.UserResponse{
		UUID:     user.UUID,
		Name:     user.Name,
		Email:    user.Email,
		Handle:   user.Handle,
		Bio:      user.Bio,
		Avatar:   user.Avatar,
		Password: user.Password,
	}, nil
}

func (ds *SQLiteDatastore) CreateUser(name string, email string, handle string, bio string, avatar string, hashedPassword string) (*models.UserResponse, error) {
	user := dbmodels.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Handle:   handle,
		Bio:      bio,
		Avatar:   avatar,

		Orgs: []dbmodels.Organization{
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
	return &models.UserResponse{
		UUID:   user.UUID,
		Name:   user.Name,
		Email:  user.Email,
		Handle: user.Handle,
		Bio:    user.Bio,
		Avatar: user.Avatar,
	}, nil
}

func (ds *SQLiteDatastore) UpdateUser(email string, name string, bio string, avatar string) (*models.UserResponse, error) {
	var user dbmodels.User
	result := ds.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if name != "" {
		user.Name = name
	}
	if bio != "" {
		user.Bio = bio
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	result = ds.DB.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.UserResponse{
		UUID:   user.UUID,
		Name:   user.Name,
		Email:  user.Email,
		Handle: user.Handle,
		Bio:    user.Bio,
		Avatar: user.Avatar,
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

/////////////////////////////// MODEL METHODS/////////////////////////////////

func (ds *SQLiteDatastore) GetModelByName(orgId uuid.UUID, modelName string) (*models.ModelResponse, error) {
	var model dbmodels.Model
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Where("name = ?", modelName).Where("organization_uuid = ?", orgId).Limit(1).Find(&model)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.ModelResponse{
		UUID: model.UUID,
		Name: model.Name,
		Wiki: model.Wiki,
		CreatedBy: models.UserHandleResponse{
			UUID:   model.CreatedByUser.UUID,
			Handle: model.CreatedByUser.Handle,
		},
		UpdatedBy: models.UserHandleResponse{
			UUID:   model.UpdatedByUser.UUID,
			Handle: model.UpdatedByUser.Handle,
		},
		IsPublic: model.IsPublic,
	}, nil
}

func (ds *SQLiteDatastore) GetModelById(modelId string) (*models.ModelResponse, error) {
	var model dbmodels.Model
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Where("uuid = ?", modelId).Limit(1).Find(&model)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.ModelResponse{
		UUID: model.UUID,
		Name: model.Name,
		Wiki: model.Wiki,
		CreatedBy: models.UserHandleResponse{
			UUID:   model.CreatedByUser.UUID,
			Handle: model.CreatedByUser.Handle,
		},
		UpdatedBy: models.UserHandleResponse{
			UUID:   model.UpdatedByUser.UUID,
			Handle: model.UpdatedByUser.Handle,
		},
		IsPublic: model.IsPublic,
	}, nil
}

func (ds *SQLiteDatastore) CreateModel(orgId uuid.UUID, name string, wiki string, createdByUser uuid.UUID) (*models.ModelResponse, error) {
	model := dbmodels.Model{
		Name: name,
		Wiki: wiki,
		Org: dbmodels.Organization{
			BaseModel: dbmodels.BaseModel{
				UUID: orgId,
			},
		},
		CreatedByUser: dbmodels.User{
			BaseModel: dbmodels.BaseModel{
				UUID: createdByUser,
			},
		},
		UpdatedByUser: dbmodels.User{
			BaseModel: dbmodels.BaseModel{
				UUID: createdByUser,
			},
		},
		IsPublic: false,
	}
	err := ds.DB.Create(&model).Error
	if err != nil {
		return nil, err
	}
	return &models.ModelResponse{
		UUID: model.UUID,
		Name: model.Name,
		Wiki: model.Wiki,
		CreatedBy: models.UserHandleResponse{
			UUID:   model.CreatedByUser.UUID,
			Handle: model.CreatedByUser.Handle,
		},
		UpdatedBy: models.UserHandleResponse{
			UUID:   model.UpdatedByUser.UUID,
			Handle: model.UpdatedByUser.Handle,
		},
		IsPublic: model.IsPublic,
	}, nil
}

func (ds *SQLiteDatastore) GetAllModels(orgId uuid.UUID) ([]models.ModelResponse, error) {
	var mymodels []dbmodels.Model
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Where("organization_uuid = ?", orgId).Find(&mymodels)
	if result.Error != nil {
		return nil, result.Error
	}
	modelResponses := make([]models.ModelResponse, len(mymodels))
	for i, model := range mymodels {
		modelResponses[i] = models.ModelResponse{
			UUID: model.UUID,
			Name: model.Name,
			Wiki: model.Wiki,
			CreatedBy: models.UserHandleResponse{
				UUID:   model.CreatedByUser.UUID,
				Handle: model.CreatedByUser.Handle,
			},
			UpdatedBy: models.UserHandleResponse{
				UUID:   model.UpdatedByUser.UUID,
				Handle: model.UpdatedByUser.Handle,
			},
			IsPublic: model.IsPublic,
		}
	}
	return modelResponses, nil
}

func (ds *SQLiteDatastore) GetModelAllBranches(modelUUID uuid.UUID) ([]models.ModelBranchResponse, error) {
	var modelBranches []dbmodels.ModelBranch
	result := ds.DB.Preload("Model").Where("model_uuid = ?", modelUUID).Find(&modelBranches)
	if result.Error != nil {
		return nil, result.Error
	}
	branches := make([]models.ModelBranchResponse, len(modelBranches))
	for i, branch := range modelBranches {
		branches[i] = models.ModelBranchResponse{
			UUID: branch.UUID,
			Name: branch.Name,
			Model: models.ModelNameResponse{
				UUID: branch.Model.UUID,
				Name: branch.Model.Name,
			},
		}
	}
	return branches, nil
}

func (ds *SQLiteDatastore) CreateModelBranch(modelUUID uuid.UUID, modelBranchName string) (*models.ModelBranchResponse, error) {
	modelBranch := dbmodels.ModelBranch{
		Name: modelBranchName,
		Model: dbmodels.Model{
			BaseModel: dbmodels.BaseModel{
				UUID: modelUUID,
			},
		},
	}
	err := ds.DB.Create(&modelBranch).Preload("Model").Error
	if err != nil {
		return nil, err
	}
	return &models.ModelBranchResponse{
		UUID: modelBranch.UUID,
		Name: modelBranch.Name,
		Model: models.ModelNameResponse{
			UUID: modelBranch.Model.UUID,
			Name: modelBranch.Model.Name,
		},
	}, nil
}

func (ds *SQLiteDatastore) UploadAndRegisterModelFile(modelBranchUUID uuid.UUID, file *multipart.FileHeader, hash string, source string) (*models.ModelVersionResponse, error) {
	// For now source is R2 by default

	var sourceType dbmodels.SourceType
	var sourcePath dbmodels.Path
	if source == "R2" {
		sourceType.Name = "R2"
		err := ds.DB.Where(&sourceType).First(&sourceType).Error
		if err != nil {
			return nil, err
		}
		splitFile := strings.Split(file.Filename, ".")
		updatedFilename := fmt.Sprintf("%s-%s.%s", splitFile[0], shortid.MustGenerate(), splitFile[1])
		var uploadPath string

		fileData, err := file.Open()
		if err != nil {
			return nil, err
		}
		s3Client := GetR2Client()
		uploader := manager.NewUploader(s3Client)
		result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(config.R2BucketName()),
			Key:    aws.String(updatedFilename),
			Body:   fileData,
		})
		if err != nil {
			return nil, err
		}

		uploadPath = strings.Split(result.Location, "/")[3]

		sourcePath = dbmodels.Path{
			SourcePath: uploadPath,
			SourceType: sourceType,
		}
		err = ds.DB.Create(&sourcePath).Error
		if err != nil {
			return nil, err
		}
	}

	latestModelVersion := dbmodels.ModelVersion{
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

	modelVersion := dbmodels.ModelVersion{
		Hash:    hash,
		Version: newVersion,
		Branch: dbmodels.ModelBranch{
			BaseModel: dbmodels.BaseModel{
				UUID: modelBranchUUID,
			},
		},
		Path: sourcePath,
	}
	err := ds.DB.Create(&modelVersion).Preload("Branch").Preload("Path.SourceType").Error
	if err != nil {
		return nil, err
	}

	return &models.ModelVersionResponse{
		UUID:    modelVersion.UUID,
		Hash:    modelVersion.Hash,
		Version: modelVersion.Version,
		Branch: models.ModelBranchNameResponse{
			UUID: modelVersion.Branch.UUID,
			Name: modelVersion.Branch.Name,
		},
		Path: models.PathResponse{
			UUID:       modelVersion.Path.UUID,
			SourcePath: modelVersion.Path.SourcePath,
			SourceType: models.SourceTypeResponse{
				Name:      modelVersion.Path.SourceType.Name,
				PublicURL: modelVersion.Path.SourceType.PublicURL,
			},
		},
	}, nil
}

func (ds *SQLiteDatastore) GetModelAllVersions(modelUUID uuid.UUID) ([]models.ModelVersionResponse, error) {
	var modelVersions []dbmodels.ModelVersion
	err := ds.DB.Where("model_uuid = ?", modelUUID).Preload("Branch").Preload("Path.SourceType").Find(&modelVersions).Error
	if err != nil {
		return nil, err
	}
	var modelVersionsResponse []models.ModelVersionResponse
	for _, modelVersion := range modelVersions {
		modelVersionsResponse = append(modelVersionsResponse, models.ModelVersionResponse{
			UUID:    modelVersion.UUID,
			Hash:    modelVersion.Hash,
			Version: modelVersion.Version,
			Branch: models.ModelBranchNameResponse{
				UUID: modelVersion.Branch.UUID,
				Name: modelVersion.Branch.Name,
			},
			Path: models.PathResponse{
				UUID:       modelVersion.Path.UUID,
				SourcePath: modelVersion.Path.SourcePath,
				SourceType: models.SourceTypeResponse{
					Name:      modelVersion.Path.SourceType.Name,
					PublicURL: modelVersion.Path.SourceType.PublicURL,
				},
			},
		})
	}
	return modelVersionsResponse, nil
}

func (ds *SQLiteDatastore) GetModelBranchByName(orgId uuid.UUID, modelName string, modelBranchName string) (*models.ModelBranchResponse, error) {
	var modelBranch dbmodels.ModelBranch
	model, err := ds.GetModelByName(orgId, modelName)
	if err != nil {
		return nil, err
	}
	res := ds.DB.Where("name = ?", modelBranchName).Where("model_uuid = ?", model.UUID).Preload("Model").Limit(1).Find(&modelBranch)
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &models.ModelBranchResponse{
		UUID: modelBranch.UUID,
		Name: modelBranch.Name,
		Model: models.ModelNameResponse{
			UUID: modelBranch.Model.UUID,
			Name: modelBranch.Model.Name,
		},
		IsDefault: modelBranch.IsDefault,
	}, nil
}

func (ds *SQLiteDatastore) GetModelBranchByUUID(modelBranchUUID uuid.UUID) (*models.ModelBranchResponse, error) {
	var modelBranch dbmodels.ModelBranch
	err := ds.DB.Where("uuid = ?", modelBranchUUID).Preload("Model").Find(&modelBranch).Error
	if err != nil {
		return nil, err
	}
	return &models.ModelBranchResponse{
		UUID: modelBranch.UUID,
		Name: modelBranch.Name,
		Model: models.ModelNameResponse{
			UUID: modelBranch.Model.UUID,
			Name: modelBranch.Model.Name,
		},
		IsDefault: modelBranch.IsDefault,
	}, nil
}

func (ds *SQLiteDatastore) GetModelBranchAllVersions(modelBranchUUID uuid.UUID) ([]models.ModelVersionResponse, error) {
	var modelVersions []dbmodels.ModelVersion
	err := ds.DB.Where("branch_uuid = ?", modelBranchUUID).Preload("Branch").Preload("Path.SourceType").Find(&modelVersions).Error
	if err != nil {
		return nil, err
	}
	var modelVersionsResponse []models.ModelVersionResponse
	for _, modelVersion := range modelVersions {
		modelVersionsResponse = append(modelVersionsResponse, models.ModelVersionResponse{
			UUID:    modelVersion.UUID,
			Hash:    modelVersion.Hash,
			Version: modelVersion.Version,
			Branch: models.ModelBranchNameResponse{
				UUID: modelVersion.Branch.UUID,
				Name: modelVersion.Branch.Name,
			},
			Path: models.PathResponse{
				UUID:       modelVersion.Path.UUID,
				SourcePath: modelVersion.Path.SourcePath,
				SourceType: models.SourceTypeResponse{
					Name:      modelVersion.Path.SourceType.Name,
					PublicURL: modelVersion.Path.SourceType.PublicURL,
				},
			},
		})
	}
	return modelVersionsResponse, nil
}

func (ds *SQLiteDatastore) GetModelBranchVersion(modelBranchUUID uuid.UUID, version string) (*models.ModelVersionResponse, error) {
	var modelVersion dbmodels.ModelVersion
	res := ds.DB.Where("branch_uuid = ?", modelBranchUUID).Where("version = ?", version).Preload("Branch").Preload("Path.SourceType").Limit(1).Find(&modelVersion)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &models.ModelVersionResponse{
		UUID:    modelVersion.UUID,
		Hash:    modelVersion.Hash,
		Version: modelVersion.Version,
		Branch: models.ModelBranchNameResponse{
			UUID: modelVersion.Branch.UUID,
			Name: modelVersion.Branch.Name,
		},
		Path: models.PathResponse{
			UUID:       modelVersion.Path.UUID,
			SourcePath: modelVersion.Path.SourcePath,
			SourceType: models.SourceTypeResponse{
				Name:      modelVersion.Path.SourceType.Name,
				PublicURL: modelVersion.Path.SourceType.PublicURL,
			},
		},
	}, nil
}

/////////////////////////////// DATASET METHODS/////////////////////////////////

func (ds *SQLiteDatastore) GetDatasetByName(orgId uuid.UUID, datasetName string) (*models.DatasetResponse, error) {
	var dataset dbmodels.Dataset
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Where("name = ?", datasetName).Where("organization_uuid = ?", orgId).Limit(1).Find(&dataset)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.DatasetResponse{
		UUID: dataset.UUID,
		Name: dataset.Name,
		Wiki: dataset.Wiki,
		CreatedBy: models.UserHandleResponse{
			UUID:   dataset.CreatedByUser.UUID,
			Handle: dataset.CreatedByUser.Handle,
		},
		UpdatedBy: models.UserHandleResponse{
			UUID:   dataset.UpdatedByUser.UUID,
			Handle: dataset.UpdatedByUser.Handle,
		},
		IsPublic: dataset.IsPublic,
	}, nil
}

func (ds *SQLiteDatastore) GetDatasetById(modelId string) (*models.DatasetResponse, error) {
	var dataset dbmodels.Dataset
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Where("uuid = ?", modelId).Limit(1).Find(&dataset)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.DatasetResponse{
		UUID: dataset.UUID,
		Name: dataset.Name,
		Wiki: dataset.Wiki,
		CreatedBy: models.UserHandleResponse{
			UUID:   dataset.CreatedByUser.UUID,
			Handle: dataset.CreatedByUser.Handle,
		},
		UpdatedBy: models.UserHandleResponse{
			UUID:   dataset.UpdatedByUser.UUID,
			Handle: dataset.UpdatedByUser.Handle,
		},
		IsPublic: dataset.IsPublic,
	}, nil
}

func (ds *SQLiteDatastore) CreateDataset(orgId uuid.UUID, name string, wiki string, createdByUser uuid.UUID) (*models.DatasetResponse, error) {
	dataset := dbmodels.Dataset{
		Name: name,
		Wiki: wiki,
		Org: dbmodels.Organization{
			BaseModel: dbmodels.BaseModel{
				UUID: orgId,
			},
		},
		CreatedByUser: dbmodels.User{
			BaseModel: dbmodels.BaseModel{
				UUID: createdByUser,
			},
		},
		UpdatedByUser: dbmodels.User{
			BaseModel: dbmodels.BaseModel{
				UUID: createdByUser,
			},
		},
		IsPublic: false,
	}
	err := ds.DB.Create(&dataset).Error
	if err != nil {
		return nil, err
	}
	return &models.DatasetResponse{
		UUID: dataset.UUID,
		Name: dataset.Name,
		Wiki: dataset.Wiki,
		CreatedBy: models.UserHandleResponse{
			UUID:   dataset.CreatedByUser.UUID,
			Handle: dataset.CreatedByUser.Handle,
		},
		UpdatedBy: models.UserHandleResponse{
			UUID:   dataset.UpdatedByUser.UUID,
			Handle: dataset.UpdatedByUser.Handle,
		},
		IsPublic: dataset.IsPublic,
	}, nil
}

func (ds *SQLiteDatastore) GetAllDatasets(orgId uuid.UUID) ([]models.DatasetResponse, error) {
	var datasets []dbmodels.Dataset
	result := ds.DB.Preload("CreatedByUser").Preload("UpdatedByUser").Where("organization_uuid = ?", orgId).Find(&datasets)
	if result.Error != nil {
		return nil, result.Error
	}
	datasetResponses := make([]models.DatasetResponse, len(datasets))
	for i, dataset := range datasets {
		datasetResponses[i] = models.DatasetResponse{
			UUID: dataset.UUID,
			Name: dataset.Name,
			Wiki: dataset.Wiki,
			CreatedBy: models.UserHandleResponse{
				UUID:   dataset.CreatedByUser.UUID,
				Handle: dataset.CreatedByUser.Handle,
			},
			UpdatedBy: models.UserHandleResponse{
				UUID:   dataset.UpdatedByUser.UUID,
				Handle: dataset.UpdatedByUser.Handle,
			},
			IsPublic: dataset.IsPublic,
		}
	}
	return datasetResponses, nil
}

func (ds *SQLiteDatastore) GetDatasetAllBranches(datasetUUID uuid.UUID) ([]models.DatasetBranchResponse, error) {
	var datasetBranches []dbmodels.DatasetBranch
	result := ds.DB.Preload("Dataset").Where("dataset_uuid = ?", datasetUUID).Find(&datasetBranches)
	if result.Error != nil {
		return nil, result.Error
	}
	branches := make([]models.DatasetBranchResponse, len(datasetBranches))
	for i, branch := range datasetBranches {
		branches[i] = models.DatasetBranchResponse{
			UUID: branch.UUID,
			Name: branch.Name,
			Dataset: models.DatasetNameResponse{
				UUID: branch.Dataset.UUID,
				Name: branch.Dataset.Name,
			},
		}
	}
	return branches, nil
}

func (ds *SQLiteDatastore) CreateDatasetBranch(datasetUUID uuid.UUID, datasetBranchName string) (*models.DatasetBranchResponse, error) {
	datasetBranch := dbmodels.DatasetBranch{
		Name: datasetBranchName,
		Dataset: dbmodels.Dataset{
			BaseModel: dbmodels.BaseModel{
				UUID: datasetUUID,
			},
		},
	}
	err := ds.DB.Create(&datasetBranch).Preload("Dataset").Error
	if err != nil {
		return nil, err
	}
	return &models.DatasetBranchResponse{
		UUID: datasetBranch.UUID,
		Name: datasetBranch.Name,
		Dataset: models.DatasetNameResponse{
			UUID: datasetBranch.Dataset.UUID,
			Name: datasetBranch.Dataset.Name,
		},
	}, nil
}

func (ds *SQLiteDatastore) UploadAndRegisterDatasetFile(datasetBranchUUID uuid.UUID, file *multipart.FileHeader, hash string, source string, lineage string) (*models.DatasetVersionResponse, error) {
	// For now source is R2 by default

	var sourceType dbmodels.SourceType
	var sourcePath dbmodels.Path
	if source == "R2" {
		sourceType.Name = "R2"
		err := ds.DB.Where(&sourceType).First(&sourceType).Error
		if err != nil {
			return nil, err
		}
		splitFile := strings.Split(file.Filename, ".")
		updatedFilename := fmt.Sprintf("%s-%s.%s", splitFile[0], shortid.MustGenerate(), splitFile[1])
		var uploadPath string

		fileData, err := file.Open()
		if err != nil {
			return nil, err
		}
		s3Client := GetR2Client()
		uploader := manager.NewUploader(s3Client)
		result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(config.R2BucketName()),
			Key:    aws.String(updatedFilename),
			Body:   fileData,
		})
		if err != nil {
			return nil, err
		}

		uploadPath = strings.Split(result.Location, "/")[3]

		sourcePath = dbmodels.Path{
			SourcePath: uploadPath,
			SourceType: sourceType,
		}
		err = ds.DB.Create(&sourcePath).Error
		if err != nil {
			return nil, err
		}
	}

	latestDatasetVersion := dbmodels.DatasetVersion{
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

	datasetVersion := dbmodels.DatasetVersion{
		Hash:    hash,
		Version: newVersion,
		Branch: dbmodels.DatasetBranch{
			BaseModel: dbmodels.BaseModel{
				UUID: datasetBranchUUID,
			},
		},
		Lineage: dbmodels.Lineage{
			Lineage: lineage,
		},
		Path: sourcePath,
	}
	err := ds.DB.Create(&datasetVersion).Preload("Lineage").Preload("Branch").Preload("Path.SourceType").Error
	if err != nil {
		return nil, err
	}

	return &models.DatasetVersionResponse{
		UUID:    datasetVersion.UUID,
		Hash:    datasetVersion.Hash,
		Version: datasetVersion.Version,
		Branch: models.DatasetBranchNameResponse{
			UUID: datasetVersion.Branch.UUID,
			Name: datasetVersion.Branch.Name,
		},
		Path: models.PathResponse{
			UUID:       datasetVersion.Path.UUID,
			SourcePath: datasetVersion.Path.SourcePath,
			SourceType: models.SourceTypeResponse{
				Name:      datasetVersion.Path.SourceType.Name,
				PublicURL: datasetVersion.Path.SourceType.PublicURL,
			},
		},
		Lineage: models.LineageResponse{
			UUID:    datasetVersion.Lineage.UUID,
			Lineage: datasetVersion.Lineage.Lineage,
		},
	}, nil
}

func (ds *SQLiteDatastore) GetDatasetAllVersions(datasetUUID uuid.UUID) ([]models.DatasetVersionResponse, error) {
	var datasetVersions []dbmodels.DatasetVersion
	err := ds.DB.Where("dataset_uuid = ?", datasetUUID).Preload("Lineage").Preload("Branch").Preload("Path.SourceType").Find(&datasetVersions).Error
	if err != nil {
		return nil, err
	}
	var datasetVersionsResponse []models.DatasetVersionResponse
	for _, datasetVersion := range datasetVersions {
		datasetVersionsResponse = append(datasetVersionsResponse, models.DatasetVersionResponse{
			UUID:    datasetVersion.UUID,
			Hash:    datasetVersion.Hash,
			Version: datasetVersion.Version,
			Branch: models.DatasetBranchNameResponse{
				UUID: datasetVersion.Branch.UUID,
				Name: datasetVersion.Branch.Name,
			},
			Path: models.PathResponse{
				UUID:       datasetVersion.Path.UUID,
				SourcePath: datasetVersion.Path.SourcePath,
				SourceType: models.SourceTypeResponse{
					Name:      datasetVersion.Path.SourceType.Name,
					PublicURL: datasetVersion.Path.SourceType.PublicURL,
				},
			},
			Lineage: models.LineageResponse{
				UUID:    datasetVersion.Lineage.UUID,
				Lineage: datasetVersion.Lineage.Lineage,
			},
		})
	}
	return datasetVersionsResponse, nil
}

func (ds *SQLiteDatastore) GetDatasetBranchByName(orgId uuid.UUID, datasetName string, datasetBranchName string) (*models.DatasetBranchResponse, error) {
	var datasetBranch dbmodels.DatasetBranch
	model, err := ds.GetDatasetByName(orgId, datasetName)
	if err != nil {
		return nil, err
	}
	res := ds.DB.Where("name = ?", datasetBranchName).Where("dataset_uuid = ?", model.UUID).Preload("Dataset").Limit(1).Find(&datasetBranch)
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &models.DatasetBranchResponse{
		UUID: datasetBranch.UUID,
		Name: datasetBranch.Name,
		Dataset: models.DatasetNameResponse{
			UUID: datasetBranch.Dataset.UUID,
			Name: datasetBranch.Dataset.Name,
		},
		IsDefault: datasetBranch.IsDefault,
	}, nil
}

func (ds *SQLiteDatastore) GetDatasetBranchByUUID(datasetBranchUUID uuid.UUID) (*models.DatasetBranchResponse, error) {
	var datasetBranch dbmodels.DatasetBranch
	err := ds.DB.Where("uuid = ?", datasetBranchUUID).Preload("Dataset").Find(&datasetBranch).Error
	if err != nil {
		return nil, err
	}
	return &models.DatasetBranchResponse{
		UUID: datasetBranch.UUID,
		Name: datasetBranch.Name,
		Dataset: models.DatasetNameResponse{
			UUID: datasetBranch.Dataset.UUID,
			Name: datasetBranch.Dataset.Name,
		},
		IsDefault: datasetBranch.IsDefault,
	}, nil
}

func (ds *SQLiteDatastore) GetDatasetBranchAllVersions(datasetBranchUUID uuid.UUID) ([]models.DatasetVersionResponse, error) {
	var datasetVersions []dbmodels.DatasetVersion
	err := ds.DB.Where("branch_uuid = ?", datasetBranchUUID).Preload("Lineage").Preload("Branch").Preload("Path.SourceType").Find(&datasetVersions).Error
	if err != nil {
		return nil, err
	}
	var datasetVersionsResponse []models.DatasetVersionResponse
	for _, datasetVersion := range datasetVersions {
		datasetVersionsResponse = append(datasetVersionsResponse, models.DatasetVersionResponse{
			UUID:    datasetVersion.UUID,
			Hash:    datasetVersion.Hash,
			Version: datasetVersion.Version,
			Branch: models.DatasetBranchNameResponse{
				UUID: datasetVersion.Branch.UUID,
				Name: datasetVersion.Branch.Name,
			},
			Path: models.PathResponse{
				UUID:       datasetVersion.Path.UUID,
				SourcePath: datasetVersion.Path.SourcePath,
				SourceType: models.SourceTypeResponse{
					Name:      datasetVersion.Path.SourceType.Name,
					PublicURL: datasetVersion.Path.SourceType.PublicURL,
				},
			},
			Lineage: models.LineageResponse{
				UUID:    datasetVersion.Lineage.UUID,
				Lineage: datasetVersion.Lineage.Lineage,
			},
		})
	}
	return datasetVersionsResponse, nil
}

func (ds *SQLiteDatastore) GetDatasetBranchVersion(datasetBranchUUID uuid.UUID, version string) (*models.DatasetVersionResponse, error) {
	var datasetVersion dbmodels.DatasetVersion
	res := ds.DB.Where("branch_uuid = ?", datasetBranchUUID).Where("version = ?", version).Preload("Lineage").Preload("Branch").Preload("Path.SourceType").Limit(1).Find(&datasetVersion)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &models.DatasetVersionResponse{
		UUID:    datasetVersion.UUID,
		Hash:    datasetVersion.Hash,
		Version: datasetVersion.Version,
		Branch: models.DatasetBranchNameResponse{
			UUID: datasetVersion.Branch.UUID,
			Name: datasetVersion.Branch.Name,
		},
		Path: models.PathResponse{
			UUID:       datasetVersion.Path.UUID,
			SourcePath: datasetVersion.Path.SourcePath,
			SourceType: models.SourceTypeResponse{
				Name:      datasetVersion.Path.SourceType.Name,
				PublicURL: datasetVersion.Path.SourceType.PublicURL,
			},
		},
		Lineage: models.LineageResponse{
			UUID:    datasetVersion.Lineage.UUID,
			Lineage: datasetVersion.Lineage.Lineage,
		},
	}, nil
}

//////////////////////////////// LOG METHODS /////////////////////////////////

func (ds *SQLiteDatastore) CreateLogForModelVersion(data string, modelVersionUUID uuid.UUID) (*models.LogResponse, error) {
	log := dbmodels.Log{
		Data: data,
		ModelVersion: dbmodels.ModelVersion{
			BaseModel: dbmodels.BaseModel{
				UUID: modelVersionUUID,
			},
		},
	}
	err := ds.DB.Create(&log).Association("ModelVersion").Find(&log.ModelVersion)
	if err != nil {
		return nil, err
	}
	return &models.LogResponse{
		Data: log.Data,
		ModelVersion: models.ModelVersionNameResponse{
			UUID:    log.ModelVersion.UUID,
			Version: log.ModelVersion.Version,
		},
	}, nil
}

func (ds *SQLiteDatastore) CreateLogForDatasetVersion(data string, datasetVersionUUID uuid.UUID) (*models.LogResponse, error) {
	log := dbmodels.Log{
		Data: data,
		DatasetVersion: dbmodels.DatasetVersion{
			BaseModel: dbmodels.BaseModel{
				UUID: datasetVersionUUID,
			},
		},
	}
	err := ds.DB.Create(&log).Association("DatasetVersion").Find(&log.DatasetVersion)
	if err != nil {
		return nil, err
	}
	return &models.LogResponse{
		Data: log.Data,
		DatasetVersion: models.DatasetVersionNameResponse{
			UUID:    log.DatasetVersion.UUID,
			Version: log.DatasetVersion.Version,
		},
	}, nil
}
