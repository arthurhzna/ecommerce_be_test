package controllers

import (
	"net/http"

	"github.com/arthurhzna/ecommerce_be_test/common/response"
	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/arthurhzna/ecommerce_be_test/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductController struct {
	service services.IServiceRegistry
}

type IProductController interface {
	GetProductsWithoutPagination(*gin.Context)
	CreateProduct(*gin.Context)
}

func NewProductController(service services.IServiceRegistry) IProductController {
	return &ProductController{service: service}
}

func (p *ProductController) GetProductsWithoutPagination(ctx *gin.Context) {
	products, err := p.service.GetProduct().GetProductsWithoutPagination(ctx)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusInternalServerError,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: products,
		Gin:  ctx,
	})
}

func (p *ProductController) CreateProduct(ctx *gin.Context) {
	request := &dto.CreateProductRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusUnprocessableEntity,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	product, err := p.service.GetProduct().CreateProduct(ctx, request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusInternalServerError,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusCreated,
		Data: product,
		Gin:  ctx,
	})
}
