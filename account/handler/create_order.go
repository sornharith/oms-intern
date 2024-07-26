package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	order, err := h.OrderService.CreateOrder(ctx, input.TID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}
