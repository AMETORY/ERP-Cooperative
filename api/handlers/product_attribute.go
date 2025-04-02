package handlers

import (
	"net/http"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/inventory"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type ProductAttributeHandler struct {
	ctx              *context.ERPContext
	inventoryService *inventory.InventoryService
}

func NewProductAttributeHandler(ctx *context.ERPContext) *ProductAttributeHandler {
	var inventorySrv *inventory.InventoryService
	invSrv, ok := ctx.InventoryService.(*inventory.InventoryService)
	if ok {
		inventorySrv = invSrv
	}
	return &ProductAttributeHandler{ctx: ctx, inventoryService: inventorySrv}
}

func (h *ProductAttributeHandler) CreateProductAttributeHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to create an ProductAttribute
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	var data models.ProductAttributeModel
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	data.CompanyID = &companyID
	err = h.inventoryService.ProductAttributeService.CreateProductAttribute(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ProductAttribute created successfully"})
}

func (h *ProductAttributeHandler) GetProductAttributeHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to get an ProductAttribute
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	search, _ := c.GetQuery("search")
	data, err := h.inventoryService.ProductAttributeService.GetProductAttributes(*c.Request, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ProductAttribute retrieved successfully", "data": data})
}

func (h *ProductAttributeHandler) GetProductAttributeByIdHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to get an ProductAttribute by ID
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	id := c.Param("id")
	data, err := h.inventoryService.ProductAttributeService.GetProductAttributeByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ProductAttribute retrieved successfully", "data": data})
}

func (h *ProductAttributeHandler) UpdateProductAttributeHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to update an ProductAttribute
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	var data models.ProductAttributeModel
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	_, err = h.inventoryService.ProductAttributeService.GetProductAttributeByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	h.inventoryService.ProductAttributeService.UpdateProductAttribute(id, &data)
	c.JSON(http.StatusOK, gin.H{"message": "ProductAttribute updated successfully"})
}

func (h *ProductAttributeHandler) DeleteProductAttributeHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to delete an ProductAttribute
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	id := c.Param("id")
	err := h.inventoryService.ProductAttributeService.DeleteProductAttribute(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ProductAttribute deleted successfully"})
}
