package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	DB  *DBConfig  `validate:"required"`
	AWS *AWSConfig `validate:"required"`
}

type DBConfig struct {
	DBUser    string `validate:"required"`
	DBPasword string `validate:"required"`
	DBHost    string `validate:"required,hostname|ipv4"`
	DBPort    string `validate:"required,numeric"`
	DBName    string `validate:"required"`
}

type AWSConfig struct {
	AWSRegion                 string `validate:"required,oneof=ap-southeast-1 ap-southeast-2 ap-southeast-3"`
	AWSEbsSnapshotDescription string `validate:"required"`
}

func LoadConfig() (*Config, error) {
	dbConf := &DBConfig{
		DBUser:    os.Getenv("DB_USER"),
		DBPasword: os.Getenv("DB_PASSWORD"),
		DBHost:    os.Getenv("DB_HOST"),
		DBPort:    os.Getenv("DB_PORT"),
		DBName:    os.Getenv("DB_NAME"),
	}

	awsConf := &AWSConfig{
		AWSRegion:                 os.Getenv("AWS_REGION"),
		AWSEbsSnapshotDescription: os.Getenv("AWS_EBS_SNAPSHOT_DESC"),
	}

	conf := &Config{
		DB:  dbConf,
		AWS: awsConf,
	}

	validate := validator.New()
	if err := validate.Struct(conf); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %v", err)
	}

	return conf, nil
}
