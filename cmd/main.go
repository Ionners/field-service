package cmd

import (
	"field-service/clients"
	"field-service/common/gcs"
	"field-service/common/response"
	"field-service/config"
	"field-service/constants"
	"field-service/controllers"
	"field-service/domain/models"
	"field-service/middlewares"
	"field-service/repositories"
	"field-service/routes"
	"field-service/services"
	"fmt"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "serve",
	Short: "start the server",
	Run: func(c *cobra.Command, args []string) {
		_ = godotenv.Load()
		// Load the environment variables from the .env file
		// and start the server
		config.Init()
		db, err := config.InitDatabase()
		if err != nil {
			panic(err)
		}

		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			panic(err)
		}

		time.Local = loc

		err = db.AutoMigrate(
			&models.Field{},
			&models.FieldSchedule{},
			&models.Time{},
		)
		if err != nil {
			panic(err)
		}

		gcs := gcs.NewGCSClient(config.Config.GCSCredentialPath, config.Config.GCSBucketName)
		client := clients.NewClientRegistry()

		repository := repositories.NewRepositoryRegistry(db)
		service := services.NewServiceRegistry(repository, gcs)
		controller := controllers.NewControllerRegistry(service)

		router := gin.Default()
		router.Use(middlewares.HandlePanic())
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, response.Response{
				Status:  constants.Error,
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, response.Response{
				Status:  constants.Success,
				Message: "Welcome to Field Service",
			})
		})
		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-COntrol-Allow_Methods", "GET, POST, PUT")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-service-name, x-api-key, x-request-id")
			c.Next()
		})

		if config.Config.EnableRateLimiter {
			lmt := tollbooth.NewLimiter(
				config.Config.RateLimiterMaxRequests/float64(config.Config.RateLimiterTimeSeconds),
				&limiter.ExpirableOptions{
					DefaultExpirationTTL: time.Duration(config.Config.RateLimiterTimeSeconds) * time.Second,
				})

			router.Use(middlewares.RateLimiter(lmt))
		}

		group := router.Group("/api/v1")
		route := routes.NewRouteRegistry(controller, group, client)
		route.Serve()

		port := fmt.Sprintf(":%d", config.Config.Port)
		router.Run(port)
	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}

// func initGCS() gcs.IGCSClient {
// 	decode, err := base64.StdEncoding.DecodeString(config.Config.SignatureKey)
// 	if err != nil {
// 		panic(err)
// 	}

// 	stringPrivateKey := string(decode)
// 	gcsServiceAccount := gcs.ServiceAccountKeyJSON{
// 		Type:                    config.Config.GCSType,
// 		ProjectID:               config.Config.GCSProjectID,
// 		PrivateKeyID:            config.Config.GCSPrivateKeyID,
// 		PrivateKey:              stringPrivateKey,
// 		ClientEmail:             config.Config.GCSClientEmail,
// 		ClientID:                config.Config.GCSClientId,
// 		AuthURI:                 config.Config.GCSAuthURI,
// 		TokenURI:                config.Config.GCSTokenURI,
// 		AuthProviderX509CertURL: config.Config.GCSAuthProviderX509CertUrl,
// 		ClientX509CertURL:       config.Config.GCSClientX509CertUrl,
// 		UniverseDomain:          config.Config.GCSUniverseDomain,
// 	}

// 	gcsClient := gcs.NewGCSClient(
// 		gcsServiceAccount,
// 		config.Config.GCSBucketName)

// 	return gcsClient
// }
