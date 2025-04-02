package handlers

import (
	"net/http"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/inventory"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	ctx              *context.ERPContext
	inventoryService *inventory.InventoryService
}

func NewWarehouseHandler(ctx *context.ERPContext) *WarehouseHandler {
	var inventorySrv *inventory.InventoryService
	invSrv, ok := ctx.InventoryService.(*inventory.InventoryService)
	if ok {
		inventorySrv = invSrv
	}
	return &WarehouseHandler{ctx: ctx, inventoryService: inventorySrv}
}

func (h *WarehouseHandler) CreateWarehouseHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to create an Warehouse
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	var data models.WarehouseModel
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	data.CompanyID = &companyID
	err = h.inventoryService.WarehouseService.CreateWarehouse(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Warehouse created successfully"})
}

func (h *WarehouseHandler) ListWarehousesHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to get an Warehouse
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	search, _ := c.GetQuery("search")
	data, err := h.inventoryService.WarehouseService.GetWarehouses(*c.Request, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Warehouse retrieved successfully", "data": data})
}

func (h *WarehouseHandler) GetWarehouseByIdHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to get an Warehouse by ID
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	id := c.Param("id")
	data, err := h.inventoryService.WarehouseService.GetWarehouseByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Warehouse retrieved successfully", "data": data})
}

func (h *WarehouseHandler) UpdateWarehouseHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to update an Warehouse
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	var data models.WarehouseModel
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	_, err = h.inventoryService.WarehouseService.GetWarehouseByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	h.inventoryService.WarehouseService.UpdateWarehouse(id, &data)
	c.JSON(http.StatusOK, gin.H{"message": "Warehouse updated successfully"})
}

func (h *WarehouseHandler) DeleteWarehouseHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to delete an Warehouse
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	id := c.Param("id")
	err := h.inventoryService.WarehouseService.DeleteWarehouse(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Warehouse deleted successfully"})
}
