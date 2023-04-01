package config

import (
	"fmt"

	bmc "github.com/chuxorg/chux-models/config"
	"github.com/spf13/viper"
)

// Config struct for webapp config
type ParserConfig struct {
	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`
	AWS struct {
		S3BucketName   string `mapstructure:"bucketName", envconfig: "S3_BUCKET_NAME"`
		S3DownloadPath string `mapstructure:"downloadPath", envconfig: "S3_DOWNLOAD_PATH"`
	} `mapstructure:"aws"`
	Auth struct {
		//--This is the Okta Issuer Url for oAuth2
		OktaOuth2Issuer string `mapstructure:"issuerUrl"`
		//--The url where a token is requested
		OktaOauth2TokenUrl string `mapstructure:"tokenUrl"`
	} `mapstructure:"auth"`

	BizConfig bmc.BizObjConfig `mapstructure:"bizObjConfig"`
	
}

func LoadConfig(env string) (*ParserConfig, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName(fmt.Sprintf("config.%s.yaml", env)) // e.g., config.development.yaml or config.production.yaml
	viper.AddConfigPath(".")                           // Look for config files in the current directory
	viper.AddConfigPath("./config")                    // Look for config files in the config directory
	viper.AddConfigPath("../../config")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	var cfg ParserConfig
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %v", err)
	}

	// Make sure the DataStores and DataStoreMap are initialized
	if cfg.BizConfig.DataStores.DataStoreMap == nil {
		cfg.BizConfig.DataStores.DataStoreMap = make(map[string]bmc.DataStoreConfig)
	}

	return &cfg, nil
}

