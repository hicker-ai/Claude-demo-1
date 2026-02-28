package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the standard API response format.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OK sends a successful response.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: data})
}

// Error sends an error response.
func Error(c *gin.Context, httpStatus int, msg string) {
	c.JSON(httpStatus, Response{Code: -1, Message: msg})
}
