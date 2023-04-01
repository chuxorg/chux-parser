package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config struct for webapp config
type Config struct {
	AWS struct {
		S3BucketName   string `yaml:"bucketName", envconfig: "S3_BUCKET_NAME"`
		S3DownloadPath string `yaml:"downloadPath", envconfig: "S3_DOWNLOAD_PATH"`
	} `yaml:"aws"`
	Auth struct {
		//--This is the Okta Issuer Url for oAuth2
		OktaOuth2Issuer string `yaml:"issuerUrl"`
		//--The url where a token is requested
		OktaOauth2TokenUrl string `yaml:"tokenUrl"`
	} `yaml:"okta"`
	Mongo struct {
		ConnectionString string `yaml:"connectionString", envconfig: "MONGO_CONNECTION_STRING"`
		Database         string `yaml:"database", envconfig: "MONGO_DB"`
	} `yaml:"mongo"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	//--Create config structure
	config := &Config{}

	//--Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//--Init new YAML decode
	d := yaml.NewDecoder(file)

	//--Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
