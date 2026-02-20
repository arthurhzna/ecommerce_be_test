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
	"github.com/google/uuid"
)

type PaymentController struct {
	service services.IServiceRegistry
}

type IPaymentController interface {
	ProcessPayment(*gin.Context)
}

func NewPaymentController(service services.IServiceRegistry) IPaymentController {
	return &PaymentController{service: service}
}

// ProcessPayment godoc
// @Summary      Process payment
// @Description  Process payment for an order (Customer only)
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Param        orderUUID  path  string  true  "Order UUID"
// @Param        request   body  dto.PaymentRequest  false  "Payment Request (optional)"
// @Success      200       {object}  response.Response{data=dto.PaymentResponse}
// @Failure      400       {object}  response.Response
// @Failure      401       {object}  response.Response
// @Router       /payments/{orderUUID}/pay [post]
func (p *PaymentController) ProcessPayment(ctx *gin.Context) {

	orderUUIDStr := ctx.Param("orderUUID")
	orderUUID, err := uuid.Parse(orderUUIDStr)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	request := &dto.PaymentRequest{}
	err = ctx.ShouldBindJSON(request)
	if err != nil {
		request = &dto.PaymentRequest{}
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

	userModel, err := p.service.GetUser().GetUserByEmail(ctx, user.Email)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusUnauthorized,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	payment, err := p.service.GetPayment().ProcessPayment(ctx, userModel.ID, orderUUID, request)
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
		Data: payment,
		Gin:  ctx,
	})
}
