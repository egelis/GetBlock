package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
)

type errorResponse struct {
	Message string `json:"error"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	if !isCanceledCtxError(message) {
		logrus.Error(message)

		c.AbortWithStatusJSON(statusCode, errorResponse{message})
	}
}

func isCanceledCtxError(err string) bool {
	return strings.Contains(err, context.Canceled.Error())
}
