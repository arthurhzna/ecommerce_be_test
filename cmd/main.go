package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "ecommerce-be",
	Short: "ecommerce be",
	Long:  "ecommerce be",
	Run: func(cmd *cobra.Command, args []string) {
		gin.Default().Run(":8001")
	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
