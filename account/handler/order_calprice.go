package handler

import (
	"encoding/json"
	"memrizr/account/entity"
	apperror "memrizr/account/entity/apperrors"
	"memrizr/account/observability/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CreateCalPriceInput struct {
	UserSelect []map[string]interface{} `json:"user_select"`
	Address    string                   `json:"address"`
}

func (h *Handler) CreateCalPrice(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handler create-calprice")
	defer span.End()
	
	var input CreateCalPriceInput

	if err := c.ShouldBindJSON(&input); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "invalid error input"})
		c.JSON(http.StatusBadRequest, apperror.CusBadRequest(err.Error(), "1040"))
		logger.LogError(apperror.CusBadRequest(err.Error(), "1040"),"", logrus.Fields{
			"At": "Service",
		})
		return
	}
	if len(input.UserSelect) == 0 {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product"})
		c.JSON(http.StatusBadRequest, apperror.CusBadRequest("invalid productid", "1140"))
		logger.LogError(apperror.CusBadRequest("invalid productid", "1140"),"invalid productid", logrus.Fields{
			"At": "Service",
		})
		return
	}
	// Convert UserSelect to JSON string
	userSelectJSON, err := json.Marshal(input.UserSelect)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, apperror.CusBadRequest(err.Error(), "1240"))
		logger.LogError(apperror.CusBadRequest(err.Error(), "1240"),"", logrus.Fields{
			"At": "Service",
		})
		return
	}

	// Convert payload to domain model
	calPrice := entity.CalPrice{
		UserSelect: string(userSelectJSON),
		Address:    input.Address,
	}

	createdCalPrice, err := h.CalpriceService.CreateCalPrice(ctx, &calPrice)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, apperror.CusBadRequest(err.Error(), "1340"))
		logger.LogError(apperror.CusBadRequest(err.Error(), "1340"),"", logrus.Fields{
			"At": "Service",
		})
		return
	}
	c.JSON(http.StatusCreated, createdCalPrice)
	logger.LogInfo("Create calprice successfully", logrus.Fields{
		"Transaction id": createdCalPrice.TID,
	})
}
