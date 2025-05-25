package config

import (
	"field-service/common/util"
	"os"

	"github.com/sirupsen/logrus"
	_ "github.com/spf13/viper/remote"
)

var Config AppConfig

type AppConfig struct {
	Port                   int             `json:"port"`
	AppName                string          `json:"appName"`
	AppEnv                 string          `json:"appEnv"`
	SignatureKey           string          `json:"signatureKey"`
	Database               Database        `json:"database"`
	EnableRateLimiter      bool            `json:"enableRateLimiter"`
	RateLimiterMaxRequests float64         `json:"rateLimiterMaxRequests"`
	RateLimiterTimeSeconds int             `json:"rateLimiterTimeSeconds"`
	InternalService        InternalService `json:"internalService"`
	// GCSType                    string          `json:"gcsType"`
	// GCSProjectID               string          `json:"gcsProjectID"`
	// GCSPrivateKeyID            string          `json:"gcsPrivateKeyID"`
	// GCSPrivateKey              string          `json:"gcsPrivateKey"`
	// GCSClientEmail             string          `json:"gcsClientEmail"`
	// GCSClientId                string          `json:"gcsClientID"`
	// GCSAuthURI                 string          `json:"gcsAuthURI"`
	// GCSTokenURI                string          `json:"gcsTokenURI"`
	// GCSAuthProviderX509CertUrl string          `json:"gcsAuthProviderX509CertURL"`
	// GCSClientX509CertUrl       string          `json:"gcsClientX509CertURL"`
	// GCSBucketName              string          `json:"gcsBucketName"`
	// GCSUniverseDomain          string          `json:"gcsUniverseDomain"`

	//GCS config yang tadi sudah di definisikan di config.json
	GCSCredentialPath string `json:"gcsCredentialPath"`
	GCSBucketName     string `json:"gcsBucketName"`
}

type Database struct {
	Host                   string `json:"host"`
	Port                   int    `json:"port"`
	Name                   string `json:"name"`
	Username               string `json:"username"`
	Password               string `json:"password"`
	MaxOpenConnections     int    `json:"maxOpenConnections"`
	MaxIdleConnections     int    `json:"maxIdleConnections"`
	MaxLifetimeConnections int    `json:"maxLifetimeConnections"`
	MaxIdleTime            int    `json:"maxIdleTime"`
}

type InternalService struct {
	User User `json:"user"`
}

type User struct {
	Host         string `json:"host"`
	SignatureKey string `json:"signatureKey"`
}

func Init() {
	err := util.BindFromJSON(&Config, "config.json", ".")
	if err != nil {
		logrus.Infof("failed to bind config: %v", err)
		err = util.BindFromConsul(&Config, os.Getenv("CONSUL_HTTP_URL"), os.Getenv("CONSUL_HTTP_KEY"))
		if err != nil {
			panic(err)
		}
	}
}
