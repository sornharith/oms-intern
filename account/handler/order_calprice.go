package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"memrizr/account/entity"
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid error input"})
		return
	}
	if len(input.UserSelect) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product"})
		return
	}
	// Convert UserSelect to JSON string
	userSelectJSON, err := json.Marshal(input.UserSelect)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert payload to domain model
	calPrice := entity.CalPrice{
		UserSelect: string(userSelectJSON),
		Address:    input.Address,
	}

	createdCalPrice, err := h.CalpriceService.CreateCalPrice(ctx, &calPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdCalPrice)
}
