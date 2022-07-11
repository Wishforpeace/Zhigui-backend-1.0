package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 请求响应
type Response struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
} //@name Response

func SendResponse(c *gin.Context, message interface{}, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

func SendBadRequest(c *gin.Context, message interface{}, data interface{}, err interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: message,
		Data:    data,
		Error:   err,
	})
}

func SendError(c *gin.Context, message interface{}, data interface{}, err interface{}) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: message,
		Data:    data,
		Error:   err,
	})
}
