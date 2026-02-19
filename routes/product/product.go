package routes

import (
	"github.com/arthurhzna/ecommerce_be_test/constants"
	"github.com/arthurhzna/ecommerce_be_test/controllers"
	"github.com/arthurhzna/ecommerce_be_test/middlewares"
	"github.com/gin-gonic/gin"
)

type ProductRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IProductRoute interface {
	Run()
}

func NewProductRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IProductRoute {
	return &ProductRoute{controller: controller, group: group}
}

func (p *ProductRoute) Run() {
	group := p.group.Group("/products")
	group.GET("", middlewares.AuthenticateWithoutToken(), p.controller.GetProduct().GetProductsWithoutPagination)
	group.Use(middlewares.Authenticate())
	group.POST("/create", middlewares.CheckRole([]string{constants.AdminAuth}), p.controller.GetProduct().CreateProduct)
}
