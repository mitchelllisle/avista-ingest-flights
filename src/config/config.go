package config

import (
	"github.com/mitchelllisle/avista-ingest-flights/src/utils"
)

type AviationStackConfig struct {
	URL string
	AccessToken string
	ArrivalCode string
	PageLimit int
}

type AppConfig struct {
	AvistaConfig *AviationStackConfig
	GCPConfig *GCPConfig
}

type GCPConfig struct {
	Project string
	StorageBucket string
	StorageFolder string
}

func NewAppConfig() *AppConfig {
	avista := AviationStackConfig{
		URL:         utils.GetEnvOrString("AVIATION_STACK_URL", "http://api.aviationstack.com/v1"),
		AccessToken: utils.GetEnvOrPanic("AVIATION_STACK_API_KEY"),
		ArrivalCode: utils.GetEnvOrString("AVIATION_STACK_ARRIVAL_CODE", "SYD"),
		PageLimit: utils.GetEnvOrInt("AVIATION_STACK_PAGE_LIMIT", 100),
	}

	gcp := GCPConfig{
		Project: utils.GetEnvOrPanic("GOOGLE_PROJECT"),
		StorageBucket: utils.GetEnvOrPanic("GOOGLE_BUCKET"),
		StorageFolder: utils.GetEnvOrPanic("GOOGLE_FOLDER"),
	}
	return &AppConfig{
		AvistaConfig: &avista,
		GCPConfig: &gcp,
	}
}