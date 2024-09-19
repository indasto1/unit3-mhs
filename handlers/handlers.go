package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"indasto1.com/unit3-mhs/models"
)

func ProduceData(c *gin.Context) {
	var sumModel models.Sum
	err := c.ShouldBindJSON(&sumModel)
	if err != nil {
		log.WithError(err).Error("Failed to bind Sum model")
		c.Error(InvalidRequestBodyError(err)).SetType(gin.ErrorTypePublic)
		return
	}

	// TODO: put to sum topic

	c.Status(http.StatusOK)
}
