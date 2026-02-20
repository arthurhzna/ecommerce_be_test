package routes

import (
	"github.com/arthurhzna/ecommerce_be_test/constants"
	"github.com/arthurhzna/ecommerce_be_test/controllers"
	"github.com/arthurhzna/ecommerce_be_test/middlewares"
	"github.com/gin-gonic/gin"
)

type OrderRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IOrderRoute interface {
	Run()
}

func NewOrderRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IOrderRoute {
	return &OrderRoute{controller: controller, group: group}
}

func (o *OrderRoute) Run() {
	group := o.group.Group("/orders")
	group.Use(middlewares.Authenticate())
	group.Use(middlewares.CheckRole([]string{constants.CustomerAuth}))
	group.POST("", o.controller.GetOrder().CreateOrder)
	group.GET("/my", o.controller.GetOrder().GetMyOrders)
}
