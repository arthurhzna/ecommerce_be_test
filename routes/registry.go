package routes

import (
	"github.com/arthurhzna/ecommerce_be_test/controllers"
	productRoutes "github.com/arthurhzna/ecommerce_be_test/routes/product"
	userRoutes "github.com/arthurhzna/ecommerce_be_test/routes/user"
	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IRouteRegister interface {
	Serve()
}

func NewRouteRegistry(controller controllers.IControllerRegistry, group *gin.RouterGroup) IRouteRegister {
	return &Registry{controller: controller, group: group}
}

func (r *Registry) Serve() {
	r.userRoute().Run()
	r.productRoute().Run()
}

func (r *Registry) userRoute() userRoutes.IUserRoute {
	return userRoutes.NewUserRoute(r.controller, r.group)
}

func (r *Registry) productRoute() productRoutes.IProductRoute {
	return productRoutes.NewProductRoute(r.controller, r.group)
}
