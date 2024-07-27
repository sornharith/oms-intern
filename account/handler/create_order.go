package handler

import (
	apperror "memrizr/account/entity/apperrors"
	"memrizr/account/observability/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CreateOrderInput struct {
	TID uuid.UUID `json:"t_id"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handler create-order")
	defer span.End()
	
	var input CreateOrderInput

	if err := c.ShouldBindJSON(&input); err != nil {
		if err.Error() == "invalid UUID length: 35" {
			// c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			c.JSON(http.StatusBadRequest, apperror.CusBadRequest("invalid id length", "1040"))
			logger.LogError(apperror.CusBadRequest("invalid id length ", "1040"),"invalid id length", logrus.Fields{
				"At": "Service",
			})
			return
		}
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, apperror.CusBadRequest(err.Error(), "1140"))
		logger.LogError(apperror.CusBadRequest(err.Error(), "1140"),"", logrus.Fields{
			"At": "Service",
		})
		return
	}
	
	order, err := h.OrderService.CreateOrder(ctx, input.TID)
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, apperror.CusBadRequest(err.Error(), "1240"))
		logger.LogError(apperror.CusBadRequest(err.Error(), "1240"),"", logrus.Fields{
			"At": "Service",
		})
		return
	}

	c.JSON(http.StatusCreated, order)
	logger.LogInfo("Create order successfully", logrus.Fields{
		"Order id": order.OID,
	})
}
