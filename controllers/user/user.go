package controllers

import (
	"net/http"

	errWrap "github.com/arthurhzna/ecommerce_be_test/common/error"
	"github.com/arthurhzna/ecommerce_be_test/common/response"
	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/arthurhzna/ecommerce_be_test/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	service services.IServiceRegistry
}

type IUserController interface {
	Login(*gin.Context)
	Register(*gin.Context)
}

func NewUserController(service services.IServiceRegistry) IUserController {
	return &UserController{service: service}
}

// Login godoc
// @Summary      Login user
// @Description  Login with email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body  dto.LoginRequest  true  "Login Request"
// @Success      200      {object}  response.Response{data=dto.LoginResponse}
// @Failure      400      {object}  response.Response
// @Failure      422      {object}  response.Response
// @Router       /auth/login [post]
func (u *UserController) Login(ctx *gin.Context) {
	request := &dto.LoginRequest{}

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

	user, err := u.service.GetUser().Login(ctx, request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code:  http.StatusOK,
		Data:  user.User,
		Token: &user.Token,
		Gin:   ctx,
	})
}

// Register godoc
// @Summary      Register new user
// @Description  Register a new customer account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body  dto.RegisterRequest  true  "Register Request"
// @Success      200      {object}  response.Response{data=dto.RegisterResponse}
// @Failure      400      {object}  response.Response
// @Failure      422      {object}  response.Response
// @Router       /auth/register [post]
func (u *UserController) Register(ctx *gin.Context) {
	request := &dto.RegisterRequest{}

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

	user, err := u.service.GetUser().Register(ctx, request)
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
		Data: user.User,
		Gin:  ctx,
	})
}
