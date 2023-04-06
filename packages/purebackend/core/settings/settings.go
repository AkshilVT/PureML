package settings

import (
	"github.com/PureMLHQ/PureML/packages/purebackend/core/tools/security"
)

// Settings defines common app configuration options.
type Settings struct {
	S3     S3Config     `form:"s3" json:"s3"`
	R2     R2Config     `form:"r2" json:"r2"`
	Search SearchConfig `form:"search" json:"search"`

	AdminAuthToken              TokenConfig       `form:"adminAuthToken" json:"adminAuthToken"`
	MailVerifificationAuthToken TokenConfig       `form:"mailVerifificationAuthToken" json:"mailVerifificationAuthToken"`
	PasswordResetAuthToken      TokenConfig       `form:"passwordResetAuthToken" json:"passwordResetAuthToken"`
	MailService                 MailServiceConfig `form:"mailService" json:"mailService"`

	Site SiteConfig `form:"site" json:"site"`
}

// New creates and returns a new default Settings instance.
func New() *Settings {
	return &Settings{
		AdminAuthToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 1209600, // 14 days,
		},
		MailVerifificationAuthToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 86400, // 1 day
		},
		PasswordResetAuthToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 86400, // 1 day
		},
	}
}

type TokenConfig struct {
	Secret   string `form:"secret" json:"secret"`
	Duration int64  `form:"duration" json:"duration"`
}

type S3Config struct {
	Enabled        bool   `form:"enabled" json:"enabled"`
	Bucket         string `form:"bucket" json:"bucket"`
	Region         string `form:"region" json:"region"`
	Endpoint       string `form:"endpoint" json:"endpoint"`
	AccessKey      string `form:"accessKey" json:"accessKey"`
	Secret         string `form:"secret" json:"secret"`
	ForcePathStyle bool   `form:"forcePathStyle" json:"forcePathStyle"`
}

type R2Config struct {
	Enabled        bool   `form:"enabled" json:"enabled"`
	AccountId      string `form:"accountId" json:"accountId"`
	Bucket         string `form:"bucket" json:"bucket"`
	Endpoint       string `form:"endpoint" json:"endpoint"`
	AccessKey      string `form:"accessKey" json:"accessKey"`
	Secret         string `form:"secret" json:"secret"`
	ForcePathStyle bool   `form:"forcePathStyle" json:"forcePathStyle"`
}

type SearchConfig struct {
	Enabled     bool   `form:"enabled" json:"enabled"`
	Host        string `form:"host" json:"host"`
	AdminAPIKey string `form:"admin_api_key" json:"admin_api_key"`
}

type MailServiceConfig struct {
	Enabled  bool   `form:"enabled" json:"enabled"`
	Host     string `form:"host" json:"host"`
	Port     int    `form:"port" json:"port"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type SiteConfig struct {
	BaseURL string `form:"baseUrl" json:"baseUrl"`
}

// func (s *Settings) LoadFromDB(dao *daos.Dao, source string) error {
// 	source = strings.ToUpper(source)
// 	defaultUUID := uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))
// 	sourceSecrets, err := dao.GetSecretByName(defaultUUID, source)
// 	if err != nil {
// 		return err
// 	}
// 	switch source {
// 	case "S3":
// 		s.S3.AccessKey = sourceSecrets.AccessKeyId
// 		s.S3.Secret = sourceSecrets.AccessKeySecret
// 		s.S3.Bucket = sourceSecrets.BucketName
// 		s.S3.Region = sourceSecrets.BucketLocation
// 	}
// 	return nil
// }
