package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/chadgrant/go-http-infra/infra"
	"github.com/spf13/viper"
)

type (
	Configuration struct {
		Service      Service
		AWS          AWS
		Dynamo       Dynamo
		HealthChecks HealthConfiguration
		UseMocks     bool
	}

	Service struct {
		Name    string
		Address string
	}

	AWS struct {
		Region          string
		AccessKeyID     string
		SecretAccessKey string
	}

	Dynamo struct {
		Endpoint string
		Tables   Tables
	}

	HealthConfiguration struct {
		Dynamo HealthCheckItem
		DNS    HealthCheckItemValue
		TCP    HealthCheckItemValue
		HTTP   HealthCheckItemValue
	}

	HealthCheckItemValue struct {
		Enabled  bool
		Interval time.Duration
		Timeout  time.Duration
		Value    string
	}

	HealthCheckItem struct {
		Enabled  bool
		Interval time.Duration
		Timeout  time.Duration
	}

	Tables struct {
		Products   string
		Categories string
	}
)

func Load(cfgFile string) (*Configuration, error) {
	viper.SetConfigName("config")
	if len(cfgFile) > 0 {
		viper.SetConfigFile(cfgFile)
	}
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Printf("%v, using defaults\n", err)
	}

	cfg := &Configuration{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func init() {

	viper.SetDefault("service", map[string]interface{}{
		"name":    "dynamodb-go-sample",
		"address": "0.0.0.0:8080",
	})

	viper.SetDefault("aws", map[string]interface{}{
		"region":          "us-east-1",
		"accesskeyid":     infra.GetEnvVar("AWS_ACCESS_KEY_ID", "thisisfake-local-aws-wants-key"),
		"secretaccesskey": infra.GetEnvVar("AWS_SECRET_ACCESS_KEY", "thisisalsofake"),
	})

	viper.SetDefault("dynamo", map[string]interface{}{
		"endpoint": "http://localhost:8000",
		"tables": map[string]interface{}{
			"products":   "products",
			"categories": "categories",
		},
	})

	viper.SetDefault("healthchecks", map[string]interface{}{
		"dynamo": map[string]interface{}{
			"enabled":  true,
			"interval": "10s",
			"timeout":  "3s",
		},
		"dns": map[string]interface{}{
			"enabled":  true,
			"value":    "google.com",
			"interval": "10s",
			"timeout":  "3s",
		},
		"tcp": map[string]interface{}{
			"enabled":  true,
			"value":    "google.com:80",
			"interval": "10s",
			"timeout":  "3s",
		},
		"http": map[string]interface{}{
			"enabled":  true,
			"value":    "https://golang.org",
			"interval": "10s",
			"timeout":  "3s",
		},
	})

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath("./")
	viper.AutomaticEnv()
}
