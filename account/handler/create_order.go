package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateOrderInput struct {
	TID int `json:"t_id"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var input CreateOrderInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.OrderService.CreateOrder(input.TID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
