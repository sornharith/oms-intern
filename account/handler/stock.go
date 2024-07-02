package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"memrizr/account/entity"
	"net/http"
	"strconv"
)

type StockInput struct {
	Quantity int `json:"quantity"`
}

type StockOutput struct {
	SID      int `json:"s_id"`
	Quantity int `json:"quantity"`
}

func (h *Handler) updateStock(c *gin.Context) {
	var input StockInput
	sid := c.Param("p_id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid productid"})
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if input.Quantity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be positive"})
		return
	}

	log.Println("input data " + strconv.Itoa(input.Quantity))

	stock := entity.Stock{
		SID:      id,
		Quantity: input.Quantity,
	}
	ctx := c.Request.Context()
	if err := h.StockService.UpdateStockById(ctx, &stock); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	output := StockOutput{
		SID:      id,
		Quantity: input.Quantity,
	}

	c.JSON(http.StatusOK, gin.H{"status": "update complete", "data": output})
}
