package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"indasto1.com/unit3-mhs/broker"
	"indasto1.com/unit3-mhs/configuration"
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

	err = broker.SendMessage(sumModel, configuration.SUM_TOPIC)
	if err != nil {
		log.WithError(err).WithField("sumModel", sumModel).Error("Failed to push sum model to kafka")
		c.Error(UnknownErorr(err)).SetType(gin.ErrorTypePublic)
		return
	}

	log.WithField("sumModel", sumModel).Debug("Pushed message to kafka")

	c.Status(http.StatusOK)
}
