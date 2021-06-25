package errex

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiError struct {
	Status  int    `json:"-"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err ApiError) Error() string {
	return err.Message
}

var InvalidParamError = ApiError{
	Status:  http.StatusOK,
	Code:    401,
	Message: "invalid param",
}
var NotFoundError = ApiError{
	Status:  http.StatusOK,
	Code:    402,
	Message: "cannot found record yet",
}

var MissingTimeError = ApiError{
	Status:  http.StatusOK,
	Code:    403,
	Message: "param time is missing",
}

var MissingTypeError = ApiError{
	Status:  http.StatusOK,
	Code:    404,
	Message: "param type is missing",
}

var MissingThresholdError = ApiError{
	Status:  http.StatusOK,
	Code:    405,
	Message: "param threshold is missing",
}

var MissingPMTypeError = ApiError{
	Status:  http.StatusOK,
	Code:    406,
	Message: "param pmType is missing",
}

var InternalServerError = ApiError{
	Status:  http.StatusOK,
	Code:    500,
	Message: "internal server error",
}

type WrapperHandle func(c *gin.Context) error

func ErrorWrapper(handle WrapperHandle) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := handle(c)
		if err != nil {
			apiError := err.(ApiError)
			c.JSON(apiError.Status, apiError)
			return
		}
	}
}
