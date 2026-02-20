package controllers

import (
	"net/http"

	errWrap "github.com/arthurhzna/ecommerce_be_test/common/error"
	"github.com/arthurhzna/ecommerce_be_test/common/response"
	"github.com/arthurhzna/ecommerce_be_test/constants"
	errConstant "github.com/arthurhzna/ecommerce_be_test/constants/error"
	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/arthurhzna/ecommerce_be_test/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrderController struct {
	service services.IServiceRegistry
}

type IOrderController interface {
	CreateOrder(*gin.Context)
	GetMyOrders(*gin.Context)
}

func NewOrderController(service services.IServiceRegistry) IOrderController {
	return &OrderController{service: service}
}

// CreateOrder godoc
// @Summary      Create new order
// @Description  Create a new order with product items (Customer only)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Param        request  body  dto.CreateOrderRequest  true  "Create Order Request"
// @Success      201      {object}  response.Response{data=dto.OrderResponse}
// @Failure      400      {object}  response.Response
// @Failure      401      {object}  response.Response
// @Failure      422      {object}  response.Response
// @Router       /orders [post]
func (o *OrderController) CreateOrder(ctx *gin.Context) {
	request := &dto.CreateOrderRequest{}

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
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	userClaims := ctx.Request.Context().Value(constants.UserLogin)
	if userClaims == nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusUnauthorized,
			Err:  errWrap.WrapError(errConstant.ErrUnauthorized),
			Gin:  ctx,
		})
		return
	}

	user, ok := userClaims.(*dto.UserResponse)
	if !ok {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusUnauthorized,
			Err:  errWrap.WrapError(errConstant.ErrUnauthorized),
			Gin:  ctx,
		})
		return
	}

	userModel, err := o.service.GetUser().GetUserByEmail(ctx, user.Email)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusUnauthorized,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	order, err := o.service.GetOrder().CreateOrder(ctx, userModel.ID, request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusCreated,
		Data: order,
		Gin:  ctx,
	})
}

// GetMyOrders godoc
// @Summary      Get my orders
// @Description  Get list of orders for the authenticated user (Customer only)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Success      200  {object}  response.Response{data=[]dto.OrderResponse}
// @Failure      401  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /orders/my [get]
func (o *OrderController) GetMyOrders(ctx *gin.Context) {
	userClaims := ctx.Request.Context().Value(constants.UserLogin)
	if userClaims == nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusUnauthorized,
			Err:  errWrap.WrapError(errConstant.ErrUnauthorized),
			Gin:  ctx,
		})
		return
	}

	user, ok := userClaims.(*dto.UserResponse)
	if !ok {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusUnauthorized,
			Err:  errWrap.WrapError(errConstant.ErrUnauthorized),
			Gin:  ctx,
		})
		return
	}

	userModel, err := o.service.GetUser().GetUserByEmail(ctx, user.Email)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusUnauthorized,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	orders, err := o.service.GetOrder().GetMyOrders(ctx, userModel.ID)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: orders,
		Gin:  ctx,
	})
}
