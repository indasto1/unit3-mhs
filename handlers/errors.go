// Errors to return in HandleError middleware on any api call failure
package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type apiError struct {
	Status  int
	Message string
}

func (err apiError) Error() string {
	return err.Message
}

func (err apiError) SendEncodedResponse(c *gin.Context) {
	c.JSON(err.Status, gin.H{
		"error": err.Error(),
	})
}

func InvalidRequestBodyError(err error) *apiError {
	return &apiError{
		Status:  http.StatusBadRequest,
		Message: fmt.Sprintf("invalid request body: %v", err),
	}
}

func UnknownErorr(err error) *apiError {
	return &apiError{
		Status:  http.StatusInternalServerError,
		Message: fmt.Sprintf("unknow error: %v", err),
	}
}
