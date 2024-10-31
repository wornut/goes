package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Success(t *testing.T) {
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpassword")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "3456")
	t.Setenv("DB_NAME", "testdb")
	t.Setenv("AWS_REGION", "ap-southeast-1")
	t.Setenv("AWS_EBS_SNAPSHOT_DESC", "snapshot description")

	config, err := LoadConfig()
	assert.NoError(t, err, "expected no error for valid config")
	assert.NotNil(t, config, "expected config to be loaded")

	assert.Equal(t, "testuser", config.DB.DBUser)
	assert.Equal(t, "testpassword", config.DB.DBPasword)
	assert.Equal(t, "localhost", config.DB.DBHost)
	assert.Equal(t, "3456", config.DB.DBPort)
	assert.Equal(t, "testdb", config.DB.DBName)

	assert.Equal(t, "ap-southeast-1", config.AWS.AWSRegion)
	assert.Equal(t, "snapshot description", config.AWS.AWSEbsSnapshotDescription)
}

func TestLoadConfig_MissingEnvVars(t *testing.T) {
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpassword")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "3456")
	// do not mock DB_NAME
	t.Setenv("AWS_REGION", "ap-southeast-1")
	// do not mock AWS_EBS_SNAPSHOT_DESC

	config, err := LoadConfig()
	// assert.Error(t, err, "expected error for missing DB_NAME")
	assert.ErrorContains(t, err, "DBName")
	assert.ErrorContains(t, err, "AWSEbsSnapshotDescription")
	assert.Nil(t, config, "expected config to be nil")
}

func TestLoadConfig_InvalidAWSRegion(t *testing.T) {
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpassword")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "3456")
	t.Setenv("DB_NAME", "testdb")
	t.Setenv("AWS_REGION", "us-east-1")
	t.Setenv("AWS_EBS_SNAPSHOT_DESC", "snapshot description")

	config, err := LoadConfig()
	assert.ErrorContains(t, err, "AWSRegion")
	assert.Nil(t, config, "expected config to be nil")
}
