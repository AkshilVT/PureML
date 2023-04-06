package models

import uuid "github.com/satori/go.uuid"

// Request models

type R2SecretRequest struct {
	SecretName      string `json:"secret_name"`
	AccountId       string `json:"account_id"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
	PublicURL       string `json:"public_url"`
}

type S3SecretRequest struct {
	SecretName      string `json:"secret_name"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
	BucketLocation  string `json:"bucket_location"`
}

// Response models

type PathResponse struct {
	UUID       uuid.UUID          `json:"uuid"`
	SourcePath string             `json:"source_path"`
	SourceType SourceTypeResponse `json:"source_type"`
}

type SourceTypeResponse struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	PublicURL string    `json:"public_url"`
}

type SourceSecretResponse struct {
	Name  string `json:"name"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SourceSecrets struct {
	AccountId       string `json:"account_id"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
	BucketLocation  string `json:"bucket_location"`
	PublicURL       string `json:"public_url"`
	SourceType      string `json:"source_type"`
}

var SupportedSources = []string{"S3", "R2", "LOCAL", "PUREML-STORAGE"}
