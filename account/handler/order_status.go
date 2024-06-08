package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Updateorderstatus struct {
	Status string `json:"status"`
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	var Input Updateorderstatus
	var orderid = c.Param("o_id")
	if err := c.ShouldBindJSON(&Input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := strconv.Atoi(orderid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid orderid"})
		return
	}
	if err := h.OrderService.UpdateOrderStatus(id, Input.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"UpdateOrderStatus": "complete",
			"orderid":           orderid,
			"status":            Input.Status,
		})
	}
}
