package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestParserConfig(t *testing.T) {
	
	_, err := LoadConfig("development")
	if err != nil {
		t.Fatalf("Error reading config file: %v", err)
	}

	var config ParserConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		t.Fatalf("Unable to decode config into struct: %v", err)
	}

	// Test Logging.Level
	assert.Equal(t, "info", config.Logging.Level)

	// Test AWS fields
	assert.Equal(t, "chux-crawler", config.AWS.BucketName)
	assert.Equal(t, "~/projects/chux/chux-parser", config.AWS.DownloadPath)

	// Test Auth fields
	assert.Equal(t, "https://dev-29752729.okta.com/oauth2/default", config.Auth.IssuerURL)
	assert.Equal(t, "", config.Auth.TokenURL)

	// Test BizConfig fields
	assert.Len(t, config.BizConfig.DataStores, 1)
	dataStore := config.BizConfig.DataStores[0].DataStore
	assert.Equal(t, "mongo", dataStore.Mongo.Target)
	assert.Equal(t, "mongodb://localhost:27017", dataStore.Mongo.URI)
	assert.Equal(t, 10, dataStore.Mongo.Timeout) // Assuming timeout is in seconds
	assert.Equal(t, "chux-cprs", dataStore.Mongo.DatabaseName)
	assert.Equal(t, "", dataStore.Mongo.CollectionName)
}
