package handler

import (
	apperror "memrizr/account/entity/apperrors"
	"memrizr/account/observability/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Updateorderstatus struct {
	Status string `json:"status"`
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handler update-order-status")
	defer span.End()
	
	var Input Updateorderstatus
	var orderid = c.Param("o_id")
	if err := c.ShouldBindJSON(&Input); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, apperror.CusBadRequest(err.Error(), "1040"))
		logger.LogError(apperror.CusBadRequest(err.Error(), "1040"),"", logrus.Fields{
			"At": "Service",
		})

		return
	}
	id, err := uuid.Parse(orderid)
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "invalid orderid"})
		c.JSON(http.StatusBadRequest, apperror.CusBadRequest("invlid orderid", "1140"))
		logger.LogError(apperror.CusBadRequest("invlid orderid", "1140"),"invlid orderid", logrus.Fields{
			"At": "Service",
		})
		return
	}
	res, err := h.OrderService.UpdateOrderStatus(ctx, id, Input.Status)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, apperror.CusInternal( "1250"))
		logger.LogError(apperror.CusInternal( "1250"),"internal update error", logrus.Fields{
			"At": "Service",
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"UpdateOrderStatus": "complete",
			"orderid":           res.OID,
			"status":            res.Status,
		})
		logger.LogInfo("UpdateOrderStatus successfully", logrus.Fields{
			"UpdateOrderStatus": "complete",
			"orderid":           res.OID,
			"status":            res.Status,
		})
	}
}
