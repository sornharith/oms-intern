package handler

import (
	"memrizr/account/entity"
	apperror "memrizr/account/entity/apperrors"
	"memrizr/account/observability/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StockInput struct {
	Quantity int `json:"quantity"`
}

type StockOutput struct {
	SID      int `json:"s_id"`
	Quantity int `json:"quantity"`
}

func (h *Handler) updateStock(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "handler update-stock")
	defer span.End()

	var input StockInput
	sid := c.Param("p_id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "invalid productid"})
		c.JSON(http.StatusBadRequest, apperror.CusBadRequest("invalid productid", "1040"))
		logger.LogError(apperror.CusBadRequest("invalid productid", "1040"),"invalid productid", logrus.Fields{
			"At": "Service",
		})
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		c.JSON(http.StatusBadRequest, apperror.CusBadRequest("invalid input", "1140"))
		logger.LogError(apperror.CusBadRequest("invalid input", "1140"),"invalid input", logrus.Fields{
			"At": "Service",
		})
		return
	}
	if input.Quantity < 0 {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be positive"})
		c.JSON(http.StatusBadRequest, apperror.CusBadRequest("quantity must be positive", "1240"))
		logger.LogError(apperror.CusBadRequest("quantity must be positive", "1240"),"quantity must be positive", logrus.Fields{
			"At": "Service",
		})
		return
	}

	stock := entity.Stock{
		SID:      id,
		Quantity: input.Quantity,
	}

	res, err := h.StockService.UpdateStockById(ctx, &stock)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, apperror.CusInternal( "1250"))
		logger.LogError(apperror.CusInternal( "1250"),"internal update error", logrus.Fields{
			"At": "Service",
		})
		return
	}

	output := StockOutput{
		SID:      res.SID,
		Quantity: res.Quantity,
	}

	c.JSON(http.StatusOK, gin.H{"status": "update complete", "data": output})

	logger.LogInfo("update stock successfully", logrus.Fields{
		"PID": res.SID,
		"QTY": res.Quantity,
	})
}
