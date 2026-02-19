package routes

import (
	"github.com/arthurhzna/ecommerce_be_test/controllers"
	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IUserRoute interface {
	Run()
}

func NewUserRoute(controller controllers.IControllerRegistry, group *gin.RouterGroup) IUserRoute {
	return &UserRoute{controller: controller, group: group}
}

func (u *UserRoute) Run() {
	group := u.group.Group("/auth")
	group.POST("/login", u.controller.GetUser().Login)
	group.POST("/register", u.controller.GetUser().Register)
}
