package config

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

var lock = &sync.Mutex{}

type SecretData struct {
	IDC_OKTA_API_CLIENT_ID            string `json:"IDC_OKTA_API_CLIENT_ID"`
	IDC_OKTA_API_CLIENT_SECRET        string `json:"IDC_OKTA_API_CLIENT_SECRET"`
	KMS_KEY_ID                        string `json:KEY_ID"`
	IDC_OKTA_API_JKS_CLIENT_ID        string `json:"IDC_OKTA_API_JKS_CLIENT_ID"`
	IDC_OKTA_AUTHENTICATION_SERVER_ID string `json:"IDC_OKTA_AUTHENTICATION_SERVER_ID"`
	IDC_OKTA_ENVIRONMENT_FQDN         string `json:"IDC_OKTA_ENVIRONMENT_FQDN"`
}

var (
	OKTA_OAUTH2_ISSUER    = "https://%s/oauth2/%s"
	OKTA_GROUPS_URL       = "https://%s/api/v1/groups"
	OKTA_GROUP_USERS_URL  = "https://%s/api/v1/groups/%s/users"
	OKTA_OAUTH2_TOKEN_URL = "https://%s/oauth2/v1/token"
)

var (
	secretName   string = "idc-okta-api"
	region       string = "us-east-1"
	versionStage string = "AWSCURRENT"
)

var singleInstance *SecretData

// Returns the Client Id for the idc-okta-api-jks Okta Application
func GetOktaAPIClientId() string {

	return singleInstance.IDC_OKTA_API_CLIENT_ID
}

// Returns the Client Secret for the idc-okta-api-jks Okta Application
func GetOktaAPIClientSecret() string {

	return singleInstance.IDC_OKTA_API_CLIENT_SECRET
}

// Returns the Key of the key/cert in AWS KMS
func GetKeyId() string {
	return singleInstance.KMS_KEY_ID
}

func GetOktaApiJksClientId() string {
	return singleInstance.IDC_OKTA_API_JKS_CLIENT_ID
}

func GetOktaAuthenticationServerId() string {
	return singleInstance.IDC_OKTA_AUTHENTICATION_SERVER_ID
}

func GetOktaEnvironmentFQDN() string {
	return singleInstance.IDC_OKTA_ENVIRONMENT_FQDN
}

func GetOktaOAuth2Issuer() string {
	fqdn := GetOktaEnvironmentFQDN()
	authid := GetOktaAuthenticationServerId()

	return fmt.Sprintf("https://%s/oauth2/%s", fqdn, authid)
}

func GetOktaGroupsUrl() string {
	fqdn := GetOktaEnvironmentFQDN()
	return fmt.Sprintf(OKTA_GROUPS_URL, fqdn)
}
func GetOAuth2TokenUrl() string {
	fqdn := GetOktaEnvironmentFQDN()
	return fmt.Sprintf(OKTA_OAUTH2_TOKEN_URL, fqdn)
}

func GetOktaGroupsUsersUrl(groupId string) string {
	fqdn := GetOktaEnvironmentFQDN()
	return fmt.Sprintf(
		OKTA_GROUP_USERS_URL,
		fqdn,
		groupId,
	)
}

func fetchSecretData(instance SecretData) *SecretData {
	svc := secretsmanager.New(
		session.New(),
		aws.NewConfig().WithRegion(region),
	)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String(versionStage),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		fmt.Println(err.Error())
	}

	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
	}
	//var secretData SecretData
	err = json.Unmarshal([]byte(secretString), &instance)
	if err != nil {
		fmt.Println(err.Error())
	}

	return &instance
}

// Use a singleton for SecretData
func GetInstance() *SecretData {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating single instance now.")
			singleInstance = &SecretData{}
			singleInstance = fetchSecretData(*singleInstance)
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singleInstance
}
