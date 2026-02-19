package cmd

import (
	"time"

	"github.com/arthurhzna/ecommerce_be_test/config"
	"github.com/arthurhzna/ecommerce_be_test/database/seeders"
	"github.com/arthurhzna/ecommerce_be_test/domain/models"
	"github.com/spf13/cobra"
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
		)
		if err != nil {
			panic(err)
		}

		seeders.NewSeederRegistry(db).Run()

		repository := repositories.NewRepositoryRegistry(db)
		service := services.NewServiceRegistry(repository)
		controller := controllers.NewControllerRegistry(service)
	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
