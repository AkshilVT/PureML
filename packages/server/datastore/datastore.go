package datastore

import (
	"mime/multipart"
	"os"

	"github.com/PureML-Inc/PureML/server/datastore/impl"
	"github.com/PureML-Inc/PureML/server/models"
	uuid "github.com/satori/go.uuid"
)

var ds *impl.SQLiteDatastore = nil

func init() {
	stage := os.Getenv("STAGE")
	if stage == "Testing" {
		//For testing
		ds = impl.NewTestSQLiteDatastore()
	} else {
		//Real db
		ds = impl.NewSQLiteDatastore()
	}
}

func GetAllAdminOrgs() ([]models.OrganizationResponse, error) {
	return ds.GetAllAdminOrgs()
}

func GetOrgById(orgId uuid.UUID) (*models.OrganizationResponse, error) {
	return ds.GetOrgByID(orgId)
}

func GetOrgByJoinCode(joinCode string) (*models.OrganizationResponse, error) {
	return ds.GetOrgByJoinCode(joinCode)
}

func CreateOrgFromEmail(email string, orgName string, orgDesc string, orgHandle string) (*models.OrganizationResponse, error) {
	return ds.CreateOrgFromEmail(email, orgName, orgDesc, orgHandle)
}

func GetUserOrganizationsByEmail(email string) ([]models.UserOrganizationsResponse, error) {
	return ds.GetUserOrganizationsByEmail(email)
}

func GetUserOrganizationByOrgIdAndEmail(orgId uuid.UUID, email string) (*models.UserOrganizationsResponse, error) {
	return ds.GetUserOrganizationByOrgIdAndEmail(orgId, email)
}

func CreateUserOrganizationFromEmailAndOrgId(email string, orgId uuid.UUID) (*models.UserOrganizationsResponse, error) {
	return ds.CreateUserOrganizationFromEmailAndOrgId(email, orgId)
}

func DeleteUserOrganizationFromEmailAndOrgId(email string, orgId uuid.UUID) error {
	return ds.DeleteUserOrganizationFromEmailAndOrgId(email, orgId)
}

func CreateUserOrganizationFromEmailAndJoinCode(email string, joinCode string) (*models.UserOrganizationsResponse, error) {
	return ds.CreateUserOrganizationFromEmailAndJoinCode(email, joinCode)
}

func UpdateOrg(orgId uuid.UUID, orgName string, orgDesc string, orgAvatar string) (*models.OrganizationResponse, error) {
	return ds.UpdateOrg(orgId, orgName, orgDesc, orgAvatar)
}

func GetUserByEmail(email string) (*models.UserResponse, error) {
	return ds.GetUserByEmail(email)
}

func GetUserByHandle(email string) (*models.UserResponse, error) {
	return ds.GetUserByHandle(email)
}

func CreateUser(name string, email string, handle string, bio string, avatar string, hashedPassword string) (*models.UserResponse, error) {
	return ds.CreateUser(name, email, handle, bio, avatar, hashedPassword)
}

func UpdateUser(email string, name string, avatar string, bio string) (*models.UserResponse, error) {
	return ds.UpdateUser(email, name, avatar, bio)
}

func CreateLogForModelVersion(data string, modelVersionUUID uuid.UUID) (*models.LogResponse, error) {
	return ds.CreateLogForModelVersion(data, modelVersionUUID)
}

func CreateLogForDatasetVersion(data string, datasetVersionUUID uuid.UUID) (*models.LogResponse, error) {
	return ds.CreateLogForDatasetVersion(data, datasetVersionUUID)
}

func GetAllModels(orgId uuid.UUID) ([]models.ModelResponse, error) {
	return ds.GetAllModels(orgId)
}

func GetModelByName(orgId uuid.UUID, modelName string) (*models.ModelResponse, error) {
	return ds.GetModelByName(orgId, modelName)
}

func CreateModel(orgId uuid.UUID, name string, wiki string, userUUID uuid.UUID) (*models.ModelResponse, error) {
	return ds.CreateModel(orgId, name, wiki, userUUID)
}

func CreateModelBranch(modelUUID uuid.UUID, branchName string) (*models.ModelBranchResponse, error) {
	return ds.CreateModelBranch(modelUUID, branchName)
}

func CreateModelBranches(modelUUID uuid.UUID, branchNames []string) ([]models.ModelBranchResponse, error) {
	var branches []models.ModelBranchResponse

	for _, branchName := range branchNames {
		branch, err := CreateModelBranch(modelUUID, branchName)
		if err != nil {
			return nil, err
		}
		branches = append(branches, *branch)
	}

	return branches, nil
}

func UploadAndRegisterModelFile(modelBranchUUID uuid.UUID, file *multipart.FileHeader, hash string, source string) (*models.ModelVersionResponse, error) {
	return ds.UploadAndRegisterModelFile(modelBranchUUID, file, hash, source)
}

func GetModelAllBranches(modelUUID uuid.UUID) ([]models.ModelBranchResponse, error) {
	return ds.GetModelAllBranches(modelUUID)
}

func GetModelAllVersions(modelUUID uuid.UUID) ([]models.ModelVersionResponse, error) {
	return ds.GetModelAllVersions(modelUUID)
}

func GetModelBranchByName(orgId uuid.UUID, modelName string, branchName string) (*models.ModelBranchResponse, error) {
	return ds.GetModelBranchByName(orgId, modelName, branchName)
}

func GetModelBranchByUUID(modelBranchUUID uuid.UUID) (*models.ModelBranchResponse, error) {
	return ds.GetModelBranchByUUID(modelBranchUUID)
}

func GetModelBranchAllVersions(modelBranchUUID uuid.UUID) ([]models.ModelVersionResponse, error) {
	return ds.GetModelBranchAllVersions(modelBranchUUID)
}

func GetModelBranchVersion(modelBranchUUID uuid.UUID, version string) (*models.ModelVersionResponse, error) {
	return ds.GetModelBranchVersion(modelBranchUUID, version)
}

func GetAllDatasets(orgId uuid.UUID) ([]models.DatasetResponse, error) {
	return ds.GetAllDatasets(orgId)
}

func GetDatasetByName(orgId uuid.UUID, datasetName string) (*models.DatasetResponse, error) {
	return ds.GetDatasetByName(orgId, datasetName)
}

func CreateDataset(orgId uuid.UUID, name string, wiki string, userUUID uuid.UUID) (*models.DatasetResponse, error) {
	return ds.CreateDataset(orgId, name, wiki, userUUID)
}

func CreateDatasetBranch(datasetUUID uuid.UUID, branchName string) (*models.DatasetBranchResponse, error) {
	return ds.CreateDatasetBranch(datasetUUID, branchName)
}

func CreateDatasetBranches(datasetUUID uuid.UUID, branchNames []string) ([]models.DatasetBranchResponse, error) {
	var branches []models.DatasetBranchResponse

	for _, branchName := range branchNames {
		branch, err := CreateDatasetBranch(datasetUUID, branchName)
		if err != nil {
			return nil, err
		}
		branches = append(branches, *branch)
	}

	return branches, nil
}

func UploadAndRegisterDatasetFile(datasetBranchUUID uuid.UUID, file *multipart.FileHeader, hash string, source string, lineage string) (*models.DatasetVersionResponse, error) {
	return ds.UploadAndRegisterDatasetFile(datasetBranchUUID, file, hash, source, lineage)
}

func GetDatasetAllBranches(datasetUUID uuid.UUID) ([]models.DatasetBranchResponse, error) {
	return ds.GetDatasetAllBranches(datasetUUID)
}

func GetDatasetAllVersions(datasetUUID uuid.UUID) ([]models.DatasetVersionResponse, error) {
	return ds.GetDatasetAllVersions(datasetUUID)
}

func GetDatasetBranchByName(orgId uuid.UUID, datasetName string, branchName string) (*models.DatasetBranchResponse, error) {
	return ds.GetDatasetBranchByName(orgId, datasetName, branchName)
}

func GetDatasetBranchByUUID(datasetBranchUUID uuid.UUID) (*models.DatasetBranchResponse, error) {
	return ds.GetDatasetBranchByUUID(datasetBranchUUID)
}

func GetDatasetBranchAllVersions(datasetBranchUUID uuid.UUID) ([]models.DatasetVersionResponse, error) {
	return ds.GetDatasetBranchAllVersions(datasetBranchUUID)
}

func GetDatasetBranchVersion(datasetBranchUUID uuid.UUID, version string) (*models.DatasetVersionResponse, error) {
	return ds.GetDatasetBranchVersion(datasetBranchUUID, version)
}

type Datastore interface {
	GetAllAdminOrgs() ([]models.OrganizationResponse, error)
	GetOrgByID(orgId uuid.UUID) (*models.OrganizationResponse, error)
	GetOrgByJoinCode(joinCode string) (*models.OrganizationResponse, error)
	CreateOrgFromEmail(email string, orgName string, orgDesc string, orgHandle string) (*models.OrganizationResponse, error)
	GetUserOrganizationsByEmail(email string) ([]models.UserOrganizationsResponse, error)
	GetUserOrganizationByOrgIdAndEmail(orgId uuid.UUID, email string) (*models.UserOrganizationsResponse, error)
	CreateUserOrganizationFromEmailAndOrgId(email string, orgId uuid.UUID) (*models.UserOrganizationsResponse, error)
	DeleteUserOrganizationFromEmailAndOrgId(email string, orgId uuid.UUID) error
	CreateUserOrganizationFromEmailAndJoinCode(email string, joinCode string) (*models.UserOrganizationsResponse, error)
	UpdateOrg(orgId uuid.UUID, orgName string, orgDesc string, orgAvatar string) (*models.OrganizationResponse, error)
	GetUserByEmail(email string) (*models.UserResponse, error)
	GetUserByHandle(email string) (*models.UserResponse, error)
	CreateUser(name string, email string, handle string, bio string, avatar string, hashedPassword string) (*models.UserResponse, error)
	UpdateUser(email string, name string, avatar string, bio string) (*models.UserResponse, error)
	CreateLogForModelVersion(data string, modelVersionUUID uuid.UUID) (*models.LogResponse, error)
	CreateLogForDatasetVersion(data string, datasetVersionUUID uuid.UUID) (*models.LogResponse, error)
}
