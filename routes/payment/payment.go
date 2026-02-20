package routes

import (
	"github.com/arthurhzna/ecommerce_be_test/constants"
	"github.com/arthurhzna/ecommerce_be_test/controllers"
	"github.com/arthurhzna/ecommerce_be_test/middlewares"
	"github.com/gin-gonic/gin"
)

type PaymentRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IPaymentRoute interface {
	Run()
}

func NewPaymentRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IPaymentRoute {
	return &PaymentRoute{controller: controller, group: group}
}

func (p *PaymentRoute) Run() {
	group := p.group.Group("/payments")
	group.Use(middlewares.Authenticate())
	group.Use(middlewares.CheckRole([]string{constants.CustomerAuth}))
	group.POST("/:orderUUID/pay", p.controller.GetPayment().ProcessPayment)
}
