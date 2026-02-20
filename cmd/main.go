// @title           Ecommerce BE Test API
// @version         1.0
// @description     This is a sample ecommerce backend API server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Api-Key
// @description API Key Authentication

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Bearer Token Authentication. Format: "Bearer {token}"
package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/arthurhzna/ecommerce_be_test/common/response"
	"github.com/arthurhzna/ecommerce_be_test/config"
	"github.com/arthurhzna/ecommerce_be_test/constants"
	"github.com/arthurhzna/ecommerce_be_test/controllers"
	"github.com/arthurhzna/ecommerce_be_test/database/seeders"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/arthurhzna/ecommerce_be_test/middlewares"
	"github.com/arthurhzna/ecommerce_be_test/repositories"
	"github.com/arthurhzna/ecommerce_be_test/routes"
	"github.com/arthurhzna/ecommerce_be_test/services"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/arthurhzna/ecommerce_be_test/docs"
)

var command = &cobra.Command{
	Use:   "ecommerce-be",
	Short: "ecommerce be",
	Long:  "ecommerce be",
	Run: func(cmd *cobra.Command, args []string) {
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
			&models.Role{},
			&models.User{},
			&models.Product{},
			&models.Order{},
			&models.OrderItem{},
			&models.Payment{},
		)

		if err != nil {
			panic(err)
		}

		seeders.NewSeederRegistry(db).Run()

		repository := repositories.NewRepositoryRegistry(db)
		service := services.NewServiceRegistry(repository)
		controller := controllers.NewControllerRegistry(service)

		router := gin.Default()
		router.Use(middlewares.HandlePanic())

		// Swagger route - PUBLIC
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, response.Response{
				Status:  constants.Error,
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, response.Response{
				Status:  constants.Success,
				Message: "Welcome to Ecommerce BE Test",
			})
		})
		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Api-Key")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		})

		lmt := tollbooth.NewLimiter(
			config.Config.RateLimiterMaxRequest,
			&limiter.ExpirableOptions{
				DefaultExpirationTTL: time.Duration(config.Config.RateLimiterTimeSecond) * time.Second,
			})
		router.Use(middlewares.RateLimiter(lmt))

		group := router.Group("/api/v1")
		route := routes.NewRouteRegistry(controller, group)
		route.Serve()

		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, response.Response{
				Status:  constants.Error,
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})

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
