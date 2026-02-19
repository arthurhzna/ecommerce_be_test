package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/arthurhzna/ecommerce_be_test/common/response"
	"github.com/arthurhzna/ecommerce_be_test/config"
	"github.com/arthurhzna/ecommerce_be_test/constants"
	errConstant "github.com/arthurhzna/ecommerce_be_test/constants/error"
	"github.com/arthurhzna/ecommerce_be_test/domain/dto"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recovered from panic: %v", r)
				c.JSON(http.StatusInternalServerError, response.Response{
					Status:  constants.Error,
					Message: errConstant.ErrInternalServerError.Error(),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, response.Response{
				Status:  constants.Error,
				Message: errConstant.ErrTooManyRequests.Error(),
			})
			c.Abort()
		}
		c.Next()
	}
}

func extractBearerToken(token string) string {
	arrayToken := strings.Split(token, " ")
	if len(arrayToken) == 2 {
		return arrayToken[1]
	}
	return ""
}

func responseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, response.Response{
		Status:  constants.Error,
		Message: message,
	})
	c.Abort()
}

func validateAPIKey(c *gin.Context) error {
	apiKey := c.GetHeader(constants.XApiKey)
	apiKeyConfig := config.Config.ApiKey

	if apiKey != apiKeyConfig {
		return errConstant.ErrUnauthorized
	}
	return nil
}

func contains(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func CheckRole(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims := c.Request.Context().Value(constants.UserLogin)
		if userClaims == nil {
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		user, ok := userClaims.(*dto.Claims)
		if !ok {
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		if !contains(roles, user.User.Role) {
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}
		c.Next()
	}
}

func validateBearerToken(c *gin.Context, token string) error {
	if !strings.Contains(token, "Bearer") {
		return errConstant.ErrUnauthorized
	}

	tokenString := extractBearerToken(token)
	if tokenString == "" {
		return errConstant.ErrUnauthorized
	}

	claims := &dto.Claims{}
	tokenJwt, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errConstant.ErrInvalidToken
		}

		jwtSecret := []byte(config.Config.JwtSecretKey)
		return jwtSecret, nil
	})

	if err != nil || !tokenJwt.Valid {
		return errConstant.ErrUnauthorized
	}

	userLogin := c.Request.WithContext(context.WithValue(c.Request.Context(), constants.UserLogin, claims.User))
	c.Request = userLogin
	c.Set(constants.Token, token)
	return nil
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		token := c.GetHeader(constants.Authorization)
		if token == "" {
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		err = validateBearerToken(c, token)
		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}

		err = validateAPIKey(c)
		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}

		c.Next()
	}
}
