package config

import (
	"fmt"
	"os"
	"strings"

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
	err := validate.Struct(conf)

	if err == nil {
		return conf, nil
	}

	if _, ok := err.(*validator.ValidationErrors); ok {
		return nil, fmt.Errorf("configuration validate failed:  %v", err)
	}

	var errMsgs []string
	for _, err := range err.(validator.ValidationErrors) {
		errMsgs = append(errMsgs, fmt.Sprintf("field `%s` failed on %s tag", err.Field(), err.Tag()))
	}

	return nil, fmt.Errorf("configuration validation failed: \n%v", strings.Join(errMsgs, "\n"))
}
