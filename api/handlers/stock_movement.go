package handlers

import (
	"net/http"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/inventory"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type StockMovementHandler struct {
	ctx              *context.ERPContext
	inventoryService *inventory.InventoryService
}

func NewStockMovementHandler(ctx *context.ERPContext) *StockMovementHandler {
	var inventorySrv *inventory.InventoryService
	invSrv, ok := ctx.InventoryService.(*inventory.InventoryService)
	if ok {
		inventorySrv = invSrv
	}
	return &StockMovementHandler{ctx: ctx, inventoryService: inventorySrv}
}

func (h *StockMovementHandler) CreateStockMovementHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to create an stockMovement
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	var data models.StockMovementModel
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if data.Type == models.MovementTypeTransfer {
		_, err = h.inventoryService.StockMovementService.TransferStock(data.Date, data.SourceWarehouseID, data.WarehouseID, data.ProductID, data.VariantID, data.Quantity, data.Description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		err = h.inventoryService.StockMovementService.CreateStockMovement(&data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "StockMovement created successfully"})
}

func (h *StockMovementHandler) GetStockMovementHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to get an stockMovement
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	search, _ := c.GetQuery("search")
	data, err := h.inventoryService.StockMovementService.GetStockMovements(*c.Request, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "StockMovement retrieved successfully", "data": data})
}

func (h *StockMovementHandler) GetStockMovementByIdHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to get an stockMovement by ID
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	id := c.Param("id")
	data, err := h.inventoryService.StockMovementService.GetStockMovementByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "StockMovement retrieved successfully", "data": data})
}

func (h *StockMovementHandler) UpdateStockMovementHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to update an stockMovement
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	var data models.StockMovementModel
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	_, err = h.inventoryService.StockMovementService.GetStockMovementByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	h.inventoryService.StockMovementService.UpdateStockMovement(id, &data)
	c.JSON(http.StatusOK, gin.H{"message": "StockMovement updated successfully"})
}

func (h *StockMovementHandler) DeleteStockMovementHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to delete an stockMovement
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	id := c.Param("id")
	err := h.inventoryService.StockMovementService.DeleteStockMovement(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "StockMovement deleted successfully"})
}
