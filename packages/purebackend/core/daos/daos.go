package daos

import (
	"errors"

	authmodels "github.com/PureMLHQ/PureML/packages/purebackend/auth/models"
	commonmodels "github.com/PureMLHQ/PureML/packages/purebackend/core/common/models"
	impl "github.com/PureMLHQ/PureML/packages/purebackend/core/daos/datastore"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/models"
	"github.com/PureMLHQ/PureML/packages/purebackend/core/tools/search"
	datasetmodels "github.com/PureMLHQ/PureML/packages/purebackend/dataset/models"
	modelmodels "github.com/PureMLHQ/PureML/packages/purebackend/model/models"
	userorgmodels "github.com/PureMLHQ/PureML/packages/purebackend/user_org/models"
	uuid "github.com/satori/go.uuid"
)

type Dao struct {
	datastore *impl.Datastore
}

// TODO: add function documentation descriptions
func InitDB(dataDir string, databaseType string, databaseUrl string, searchEnabled bool, searchClient *search.SearchClient) (*Dao, error) {
	if databaseType == "" {
		//default SQLite3 db
		databaseType = "sqlite3"
	}
	dao := &Dao{
		datastore: nil,
	}
	if databaseType == "sqlite3" {
		//SQLite3 db
		dao.datastore = impl.NewSQLiteDatastore(dataDir)
	} else if databaseType == "postgres" {
		//Postgres db
		if databaseUrl == "" {
			return nil, errors.New("databaseUrl is required for postgres")
		}
		dao.datastore = impl.NewPostgresDatastore(databaseUrl)
	}
	if searchEnabled {
		if searchClient == nil {
			return nil, errors.New("searchClient is required for searchEnabled")
		}
		dao.datastore.SearchClient = searchClient
	}
	return dao, nil
}

func (dao *Dao) Datastore() *impl.Datastore {
	ds := dao.datastore
	if ds == nil {
		panic("datastore not initialized")
	}
	return dao.datastore
}

func (dao *Dao) ExecuteSQL(sqlString string) error {
	return dao.Datastore().ExecuteSQL(sqlString)
}

func (dao *Dao) Close() error {
	return dao.Datastore().Close()
}

func (dao *Dao) GetAllAdminOrgs() ([]userorgmodels.OrganizationResponse, error) {
	return dao.Datastore().GetAllAdminOrgs()
}

func (dao *Dao) GetOrgById(orgId uuid.UUID) (*userorgmodels.OrganizationResponseWithMembers, error) {
	return dao.Datastore().GetOrgByID(orgId)
}

func (dao *Dao) GetOrgByHandle(orgHandle string) (*userorgmodels.OrganizationResponse, error) {
	return dao.Datastore().GetOrgByHandle(orgHandle)
}

func (dao *Dao) GetOrgByJoinCode(joinCode string) (*userorgmodels.OrganizationResponse, error) {
	return dao.Datastore().GetOrgByJoinCode(joinCode)
}

func (dao *Dao) CreateOrgFromEmail(email string, orgName string, orgDesc string, orgHandle string) (*userorgmodels.OrganizationResponse, error) {
	return dao.Datastore().CreateOrgFromEmail(email, orgName, orgDesc, orgHandle)
}

func (dao *Dao) GetUserOrganizationsByEmail(email string) ([]userorgmodels.UserOrganizationsResponse, error) {
	return dao.Datastore().GetUserOrganizationsByEmail(email)
}

func (dao *Dao) GetUserOrganizationByOrgIdAndUserUUID(orgId uuid.UUID, userUUID uuid.UUID) (*userorgmodels.UserOrganizationsRoleResponse, error) {
	return dao.Datastore().GetUserOrganizationByOrgIdAndUserUUID(orgId, userUUID)
}

func (dao *Dao) CreateUserOrganizationFromEmailAndOrgId(email string, orgId uuid.UUID) (*userorgmodels.UserOrganizationsResponse, error) {
	return dao.Datastore().CreateUserOrganizationFromEmailAndOrgId(email, orgId)
}

func (dao *Dao) DeleteUserOrganizationFromEmailAndOrgId(email string, orgId uuid.UUID) error {
	return dao.Datastore().DeleteUserOrganizationFromEmailAndOrgId(email, orgId)
}

func (dao *Dao) CreateUserOrganizationFromEmailAndJoinCode(email string, joinCode string) (*userorgmodels.UserOrganizationsResponse, error) {
	return dao.Datastore().CreateUserOrganizationFromEmailAndJoinCode(email, joinCode)
}

func (dao *Dao) UpdateUserRoleByOrgIdAndUserUUID(orgId uuid.UUID, userUUID uuid.UUID, role string) error {
	return dao.Datastore().UpdateUserRoleByOrgIdAndUserUUID(orgId, userUUID, role)
}

func (dao *Dao) UpdateOrg(orgId uuid.UUID, updatedAttributes map[string]interface{}) (*userorgmodels.OrganizationResponse, error) {
	return dao.Datastore().UpdateOrg(orgId, updatedAttributes)
}

func (dao *Dao) GetOrgAllPublicModels(orgId uuid.UUID) ([]modelmodels.ModelResponse, error) {
	return dao.Datastore().GetOrgAllPublicModels(orgId)
}

func (dao *Dao) GetOrgAllPublicDatasets(orgId uuid.UUID) ([]datasetmodels.DatasetResponse, error) {
	return dao.Datastore().GetOrgAllPublicDatasets(orgId)
}

func (dao *Dao) GetUserByEmail(email string) (*userorgmodels.UserResponse, error) {
	return dao.Datastore().GetUserByEmail(email)
}

func (dao *Dao) GetUserByHandle(handle string) (*userorgmodels.UserProfileResponse, error) {
	return dao.Datastore().GetUserByHandle(handle)
}

func (dao *Dao) GetSecureUserByEmail(email string) (*userorgmodels.UserResponse, error) {
	return dao.Datastore().GetSecureUserByEmail(email)
}

func (dao *Dao) GetSecureUserByHandle(handle string) (*userorgmodels.UserResponse, error) {
	return dao.Datastore().GetSecureUserByHandle(handle)
}

func (dao *Dao) GetSecureUserByUUID(userUUID uuid.UUID) (*userorgmodels.UserResponse, error) {
	return dao.Datastore().GetSecureUserByUUID(userUUID)
}

func (dao *Dao) GetUserByUUID(userUUID uuid.UUID) (*userorgmodels.UserResponse, error) {
	return dao.Datastore().GetUserByUUID(userUUID)
}

func (dao *Dao) GetUserProfileByUUID(userUUID uuid.UUID) (*userorgmodels.UserProfileResponse, error) {
	return dao.Datastore().GetUserProfileByUUID(userUUID)
}

func (dao *Dao) CreateUser(name string, email string, handle string, bio string, avatar string, hashedPassword string, isVerified bool) (*userorgmodels.UserResponse, error) {
	return dao.Datastore().CreateUser(name, email, handle, bio, avatar, hashedPassword, isVerified)
}

func (dao *Dao) VerifyUserEmail(userUUID uuid.UUID) error {
	return dao.Datastore().VerifyUserEmail(userUUID)
}

func (dao *Dao) UpdateUser(email string, updatedAttributes map[string]interface{}) (*userorgmodels.UserResponse, error) {
	return dao.Datastore().UpdateUser(email, updatedAttributes)
}

func (dao *Dao) UpdateUserPassword(userUUID uuid.UUID, hashedPassword string) error {
	return dao.Datastore().UpdateUserPassword(userUUID, hashedPassword)
}

func (dao *Dao) GetSession(sessionUUID uuid.UUID) (*authmodels.SessionResponse, error) {
	return dao.Datastore().GetSession(sessionUUID)
}

func (dao *Dao) CreateSession(deviceData string, deviceId string, deviceLocation string) (*authmodels.CreateSessionResponse, error) {
	return dao.Datastore().CreateSession(deviceData, deviceId, deviceLocation)
}

func (dao *Dao) UpdateSession(sessionUUID uuid.UUID, userUUID uuid.UUID, updatedAttributes map[string]interface{}) (*authmodels.SessionResponse, error) {
	return dao.Datastore().UpdateSession(sessionUUID, userUUID, updatedAttributes)
}

func (dao *Dao) GetLogForModelVersion(modelVersionUUID uuid.UUID) ([]models.LogResponse, error) {
	return dao.Datastore().GetLogForModelVersion(modelVersionUUID)
}

func (dao *Dao) GetKeyLogForModelVersion(modelVersionUUID uuid.UUID, key string) ([]models.LogResponse, error) {
	return dao.Datastore().GetKeyLogForModelVersion(modelVersionUUID, key)
}

func (dao *Dao) CreateLogForModelVersion(key string, data string, modelVersionUUID uuid.UUID) (*models.LogResponse, error) {
	return dao.Datastore().CreateLogForModelVersion(key, data, modelVersionUUID)
}

func (dao *Dao) GetLogForDatasetVersion(datasetVersionUUID uuid.UUID) ([]models.LogResponse, error) {
	return dao.Datastore().GetLogForDatasetVersion(datasetVersionUUID)
}

func (dao *Dao) GetKeyLogForDatasetVersion(datasetVersionUUID uuid.UUID, key string) ([]models.LogResponse, error) {
	return dao.Datastore().GetKeyLogForDatasetVersion(datasetVersionUUID, key)
}

func (dao *Dao) CreateLogForDatasetVersion(key string, data string, datasetVersionUUID uuid.UUID) (*models.LogResponse, error) {
	return dao.Datastore().CreateLogForDatasetVersion(key, data, datasetVersionUUID)
}

func (dao *Dao) GetAllPublicModels() ([]modelmodels.ModelResponse, error) {
	return dao.Datastore().GetAllPublicModels()
}

func (dao *Dao) GetAllModels(orgId uuid.UUID) ([]modelmodels.ModelResponse, error) {
	return dao.Datastore().GetAllModels(orgId)
}

func (dao *Dao) GetModelByName(orgId uuid.UUID, modelName string) (*modelmodels.ModelResponse, error) {
	return dao.Datastore().GetModelByName(orgId, modelName)
}

func (dao *Dao) CreateModel(orgId uuid.UUID, name string, wiki string, isPublic bool, readmeData *commonmodels.ReadmeRequest, userUUID uuid.UUID) (*modelmodels.ModelResponse, error) {
	return dao.Datastore().CreateModel(orgId, name, wiki, isPublic, readmeData, userUUID)
}

func (dao *Dao) CreateModelBranch(modelUUID uuid.UUID, branchName string) (*modelmodels.ModelBranchResponse, error) {
	return dao.Datastore().CreateModelBranch(modelUUID, branchName)
}

func (dao *Dao) CreateModelBranches(modelUUID uuid.UUID, branchNames []string) ([]modelmodels.ModelBranchResponse, error) {
	var branches []modelmodels.ModelBranchResponse

	for _, branchName := range branchNames {
		branch, err := dao.CreateModelBranch(modelUUID, branchName)
		if err != nil {
			return nil, err
		}
		branches = append(branches, *branch)
	}

	return branches, nil
}

func (dao *Dao) RegisterModelFile(modelBranchUUID uuid.UUID, sourceType string, sourcePublicURL string, path string, isEmpty bool, hash string, userUUID uuid.UUID) (*modelmodels.ModelBranchVersionResponse, error) {
	return dao.Datastore().RegisterModelFile(modelBranchUUID, sourceType, sourcePublicURL, path, isEmpty, hash, userUUID)
}

func (dao *Dao) GetModelAllBranches(modelUUID uuid.UUID) ([]modelmodels.ModelBranchResponse, error) {
	return dao.Datastore().GetModelAllBranches(modelUUID)
}

func (dao *Dao) GetModelAllVersions(modelUUID uuid.UUID) ([]modelmodels.ModelBranchVersionResponse, error) {
	return dao.Datastore().GetModelAllVersions(modelUUID)
}

func (dao *Dao) GetModelBranchByName(orgId uuid.UUID, modelName string, branchName string) (*modelmodels.ModelBranchResponse, error) {
	return dao.Datastore().GetModelBranchByName(orgId, modelName, branchName)
}

func (dao *Dao) GetModelBranchByUUID(modelBranchUUID uuid.UUID) (*modelmodels.ModelBranchResponse, error) {
	return dao.Datastore().GetModelBranchByUUID(modelBranchUUID)
}

func (dao *Dao) GetModelBranchAllVersions(modelBranchUUID uuid.UUID, withLogs bool) ([]modelmodels.ModelBranchVersionResponse, error) {
	return dao.Datastore().GetModelBranchAllVersions(modelBranchUUID, withLogs)
}

func (dao *Dao) GetModelBranchVersion(modelBranchUUID uuid.UUID, version string) (*modelmodels.ModelBranchVersionResponse, error) {
	return dao.Datastore().GetModelBranchVersion(modelBranchUUID, version)
}

func (dao *Dao) GetAllPublicDatasets() ([]datasetmodels.DatasetResponse, error) {
	return dao.Datastore().GetAllPublicDatasets()
}

func (dao *Dao) GetAllDatasets(orgId uuid.UUID, showPublic bool) ([]datasetmodels.DatasetResponse, error) {
	return dao.Datastore().GetAllDatasets(orgId, showPublic)
}

func (dao *Dao) GetDatasetByName(orgId uuid.UUID, datasetName string) (*datasetmodels.DatasetResponse, error) {
	return dao.Datastore().GetDatasetByName(orgId, datasetName)
}

func (dao *Dao) CreateDataset(orgId uuid.UUID, name string, wiki string, isPublic bool, readmeData *commonmodels.ReadmeRequest, userUUID uuid.UUID) (*datasetmodels.DatasetResponse, error) {
	return dao.Datastore().CreateDataset(orgId, name, wiki, isPublic, readmeData, userUUID)
}

func (dao *Dao) CreateDatasetBranch(datasetUUID uuid.UUID, branchName string) (*datasetmodels.DatasetBranchResponse, error) {
	return dao.Datastore().CreateDatasetBranch(datasetUUID, branchName)
}

func (dao *Dao) CreateDatasetBranches(datasetUUID uuid.UUID, branchNames []string) ([]datasetmodels.DatasetBranchResponse, error) {
	var branches []datasetmodels.DatasetBranchResponse

	for _, branchName := range branchNames {
		branch, err := dao.CreateDatasetBranch(datasetUUID, branchName)
		if err != nil {
			return nil, err
		}
		branches = append(branches, *branch)
	}

	return branches, nil
}

func (dao *Dao) RegisterDatasetFile(datasetBranchUUID uuid.UUID, sourceType string, sourcePublicURL string, path string, isEmpty bool, hash string, lineage string, userUUID uuid.UUID) (*datasetmodels.DatasetBranchVersionResponse, error) {
	return dao.Datastore().RegisterDatasetFile(datasetBranchUUID, sourceType, sourcePublicURL, path, isEmpty, hash, lineage, userUUID)
}

func (dao *Dao) GetDatasetAllBranches(datasetUUID uuid.UUID) ([]datasetmodels.DatasetBranchResponse, error) {
	return dao.Datastore().GetDatasetAllBranches(datasetUUID)
}

func (dao *Dao) GetDatasetAllVersions(datasetUUID uuid.UUID) ([]datasetmodels.DatasetBranchVersionResponse, error) {
	return dao.Datastore().GetDatasetAllVersions(datasetUUID)
}

func (dao *Dao) GetDatasetBranchByName(orgId uuid.UUID, datasetName string, branchName string) (*datasetmodels.DatasetBranchResponse, error) {
	return dao.Datastore().GetDatasetBranchByName(orgId, datasetName, branchName)
}

func (dao *Dao) GetDatasetBranchByUUID(datasetBranchUUID uuid.UUID) (*datasetmodels.DatasetBranchResponse, error) {
	return dao.Datastore().GetDatasetBranchByUUID(datasetBranchUUID)
}

func (dao *Dao) GetDatasetBranchAllVersions(datasetBranchUUID uuid.UUID) ([]datasetmodels.DatasetBranchVersionResponse, error) {
	return dao.Datastore().GetDatasetBranchAllVersions(datasetBranchUUID)
}

func (dao *Dao) GetDatasetBranchVersion(datasetBranchUUID uuid.UUID, version string) (*datasetmodels.DatasetBranchVersionResponse, error) {
	return dao.Datastore().GetDatasetBranchVersion(datasetBranchUUID, version)
}

func (dao *Dao) GetModelActivity(modelUUID uuid.UUID, category string) (*models.ActivityResponse, error) {
	return dao.Datastore().GetModelActivity(modelUUID, category)
}

func (dao *Dao) CreateModelActivity(modelUUID uuid.UUID, userUUID uuid.UUID, category string, activity string) (*models.ActivityResponse, error) {
	return dao.Datastore().CreateModelActivity(modelUUID, userUUID, category, activity)
}

func (dao *Dao) UpdateModelActivity(activityUUID uuid.UUID, updatedAttributes map[string]string) (*models.ActivityResponse, error) {
	return dao.Datastore().UpdateModelActivity(activityUUID, updatedAttributes)
}

func (dao *Dao) DeleteModelActivity(activityUUID uuid.UUID) error {
	return dao.Datastore().DeleteModelActivity(activityUUID)
}

func (dao *Dao) GetDatasetActivity(datasetUUID uuid.UUID, category string) (*models.ActivityResponse, error) {
	return dao.Datastore().GetDatasetActivity(datasetUUID, category)
}

func (dao *Dao) CreateDatasetActivity(datasetUUID uuid.UUID, userUUID uuid.UUID, category string, activity string) (*models.ActivityResponse, error) {
	return dao.Datastore().CreateDatasetActivity(datasetUUID, userUUID, category, activity)
}

func (dao *Dao) UpdateDatasetActivity(activityUUID uuid.UUID, updatedAttributes map[string]string) (*models.ActivityResponse, error) {
	return dao.Datastore().UpdateDatasetActivity(activityUUID, updatedAttributes)
}

func (dao *Dao) DeleteDatasetActivity(activityUUID uuid.UUID) error {
	return dao.Datastore().DeleteDatasetActivity(activityUUID)
}

func (dao *Dao) GetOrganizationSecrets(orgId uuid.UUID) ([]string, error) {
	return dao.Datastore().GetOrganizationSecrets(orgId)
}

func (dao *Dao) GetSecretByName(orgId uuid.UUID, secretName string) (*commonmodels.SourceSecrets, error) {
	return dao.Datastore().GetSecretByName(orgId, secretName)
}

func (dao *Dao) CreateR2Secrets(orgId uuid.UUID, secretName string, accountId string, accessKeyId string, accessKeySecret string, bucketName string, publicURL string) (*commonmodels.SourceSecrets, error) {
	return dao.Datastore().CreateR2Secrets(orgId, secretName, accountId, accessKeyId, accessKeySecret, bucketName, publicURL)
}

func (dao *Dao) CreateS3Secrets(orgId uuid.UUID, secretName string, accessKeyId string, accessKeySecret string, bucketName string, bucketLocation string) (*commonmodels.SourceSecrets, error) {
	return dao.Datastore().CreateS3Secrets(orgId, secretName, accessKeyId, accessKeySecret, bucketName, bucketLocation)
}

func (dao *Dao) DeleteSecrets(orgId uuid.UUID, secretName string) error {
	return dao.Datastore().DeleteSecrets(orgId, secretName)
}

func (dao *Dao) GetModelReadmeVersion(modelUUID uuid.UUID, version string) (*commonmodels.ReadmeVersionResponse, error) {
	return dao.Datastore().GetModelReadmeVersion(modelUUID, version)
}

func (dao *Dao) GetModelReadmeAllVersions(modelUUID uuid.UUID) ([]commonmodels.ReadmeVersionResponse, error) {
	return dao.Datastore().GetModelReadmeAllVersions(modelUUID)
}

func (dao *Dao) UpdateModelReadme(modelUUID uuid.UUID, fileType string, content string) (*commonmodels.ReadmeVersionResponse, error) {
	return dao.Datastore().UpdateModelReadme(modelUUID, fileType, content)
}

func (dao *Dao) GetDatasetReadmeVersion(modelUUID uuid.UUID, version string) (*commonmodels.ReadmeVersionResponse, error) {
	return dao.Datastore().GetDatasetReadmeVersion(modelUUID, version)
}

func (dao *Dao) GetDatasetReadmeAllVersions(modelUUID uuid.UUID) ([]commonmodels.ReadmeVersionResponse, error) {
	return dao.Datastore().GetDatasetReadmeAllVersions(modelUUID)
}

func (dao *Dao) UpdateDatasetReadme(modelUUID uuid.UUID, fileType string, content string) (*commonmodels.ReadmeVersionResponse, error) {
	return dao.Datastore().UpdateDatasetReadme(modelUUID, fileType, content)
}

func (dao *Dao) GetModelReview(reviewUUID uuid.UUID) (*modelmodels.ModelReviewResponse, error) {
	return dao.Datastore().GetModelReview(reviewUUID)
}

func (dao *Dao) GetModelReviews(modelUUID uuid.UUID) ([]modelmodels.ModelReviewResponse, error) {
	return dao.Datastore().GetModelReviews(modelUUID)
}

func (dao *Dao) CreateModelReview(modelUUID uuid.UUID, userUUID uuid.UUID, fromBranch uuid.UUID, fromBranchVersion uuid.UUID, toBranch uuid.UUID, title string, desc string, isComplete bool, isAccepted bool) (*modelmodels.ModelReviewResponse, error) {
	return dao.Datastore().CreateModelReview(modelUUID, userUUID, fromBranch, fromBranchVersion, toBranch, title, desc, isComplete, isAccepted)
}

func (dao *Dao) UpdateModelReview(reviewUUID uuid.UUID, updatedAttributes map[string]any) (*modelmodels.ModelReviewResponse, error) {
	return dao.Datastore().UpdateModelReview(reviewUUID, updatedAttributes)
}

func (dao *Dao) GetDatasetReview(reviewUUID uuid.UUID) (*datasetmodels.DatasetReviewResponse, error) {
	return dao.Datastore().GetDatasetReview(reviewUUID)
}

func (dao *Dao) GetDatasetReviews(datasetUUID uuid.UUID) ([]datasetmodels.DatasetReviewResponse, error) {
	return dao.Datastore().GetDatasetReviews(datasetUUID)
}

func (dao *Dao) CreateDatasetReview(datasetUUID uuid.UUID, userUUID uuid.UUID, fromBranch uuid.UUID, fromBranchVerison uuid.UUID, toBranch uuid.UUID, title string, desc string, isComplete bool, isAccepted bool) (*datasetmodels.DatasetReviewResponse, error) {
	return dao.Datastore().CreateDatasetReview(datasetUUID, userUUID, fromBranch, fromBranchVerison, toBranch, title, desc, isComplete, isAccepted)
}

func (dao *Dao) UpdateDatasetReview(reviewUUID uuid.UUID, updatedAttributes map[string]any) (*datasetmodels.DatasetReviewResponse, error) {
	return dao.Datastore().UpdateDatasetReview(reviewUUID, updatedAttributes)
}
