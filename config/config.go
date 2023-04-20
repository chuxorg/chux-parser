package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	bo "github.com/chuxorg/chux-models/config"
	"github.com/spf13/viper"
)

type DataStoreConfig struct {
	Target         string `mapstructure:"target"`
	URI            string `mapstructure:"uri"`
	Timeout        int    `mapstructure:"timeout"`
	DatabaseName   string `mapstructure:"databaseName"`
	CollectionName string `mapstructure:"collectionName"`
}

type ParserConfig struct {
	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`
	AWS struct {
		BucketName    string `mapstructure:"bucketName"`
		DownloadPath  string `mapstructure:"downloadPath"`
		ArchiveBucket string `mapstructure:"archiveBucket"`
		Profile       string `mapstructure:"profile"`
		Region        string `mapstructure:"region"`
		AccessKey     string `mapstructure:"accessKey"`
		SecretKey     string `mapstructure:"secretKey"`
		SecretPath    string `mapstructure:"secretPath"`
	} `mapstructure:"aws"`
	Auth struct {
		IssuerURL string `mapstructure:"issuerUrl"`
		TokenURL  string `mapstructure:"tokenUrl"`
	} `mapstructure:"auth"`
	DataPath struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"data"`
	DataStores struct {
		// A map of data store configurations keyed by the data store name
		// e.g., "mongo" or "redis"
		DataStoreMap map[string]DataStoreConfig `mapstructure:"dataStore"`
	} `mapstructure:"dataStores"`
	Products []string `mapstructure:"productSources"`
}

func LoadConfig(env string) (*ParserConfig, error) {
	// Set the configuration file format
	viper.SetConfigType("yaml")

	// Set the configuration file name based on the environment
	// e.g., config.development.yaml or config.production.yaml
	viper.SetConfigName(fmt.Sprintf("config.%s", env))

	// Add a path where Viper should look for the configuration file
	// In this case, it will look for the file in the "../config" directory
	viper.AddConfigPath("../config")

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		// Return an error if the configuration file could not be read
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	// Declare a ParserConfig instance to store the loaded configuration values
	var cfg ParserConfig

	// Unmarshal the loaded configuration values into the ParserConfig instance
	err = viper.Unmarshal(&cfg)
	if err != nil {
		// Return an error if the unmarshalling process failed
		return nil, fmt.Errorf("failed to unmarshal configuration: %v", err)
	}

	// Initialize the DataStores.DataStoreMap field if it's not already set
	// This ensures that the map is not nil and can be used later in the code
	if cfg.DataStores.DataStoreMap == nil {
		cfg.DataStores.DataStoreMap = make(map[string]DataStoreConfig)
	}

	// Return the ParserConfig instance containing the loaded configuration values
	return &cfg, nil
}

// LoadConfig reads and parses the YAML configuration file based on the given environment
// and returns a ParserConfig instance containing the loaded configuration values.
func GetSecret(name string) (string, error) {
	// Create a new AWS session
	sess := session.Must(session.NewSession())

	// Create a Secrets Manager client
	svc := secretsmanager.New(sess)

	// Define the name of the secret to retrieve
	secretName := "dev/secrets"

	// Retrieve the secret value
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValueWithContext(context.Background(), input)
	if err != nil {
		fmt.Println("Error retrieving secret value:", err)
		return "", err
	}

	// Print the secret value
	fmt.Println("Secret value:", *result.SecretString) // Return the ParserConfig instance containing the loaded configuration values
	return *result.SecretString, nil
}

func NewBizObjConfig(parserConfig *ParserConfig) *bo.BizObjConfig {
	return &bo.BizObjConfig{
		Logging: struct {
			Level string `mapstructure:"level"`
		}{
			Level: parserConfig.Logging.Level,
		},
		DataStores: struct {
			DataStoreMap map[string]bo.DataStoreConfig `mapstructure:"dataStore"`
		}{
			DataStoreMap: ConvertDataStoreMap(parserConfig.DataStores.DataStoreMap),
		},
	}
}

func ConvertDataStoreMap(src map[string]DataStoreConfig) map[string]bo.DataStoreConfig {
	dst := make(map[string]bo.DataStoreConfig)
	for k, v := range src {
		dst[k] = bo.DataStoreConfig{
			Target:         v.Target,
			URI:            v.URI,
			Timeout:        v.Timeout,
			DatabaseName:   v.DatabaseName,
			CollectionName: v.CollectionName,
		}
	}
	return dst
}
