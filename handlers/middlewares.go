package handlers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func HandleError(c *gin.Context) {
	c.Next()

	publicErr := c.Errors.ByType(gin.ErrorTypePublic).Last()
	if publicErr == nil {
		return
	}

	apiErr, ok := publicErr.Err.(*apiError)
	if !ok || apiErr == nil {
		apiErr = UnknownErorr(publicErr.Err)
	}

	log.WithFields(log.Fields{
		"status":   apiErr.Status,
		"response": apiErr.Message,
	}).Warn("Error response")

	apiErr.SendEncodedResponse(c)
}
