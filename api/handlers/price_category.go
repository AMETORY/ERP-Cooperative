package handlers

import (
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/inventory"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type PriceCategoryHandler struct {
	ctx          *context.ERPContext
	inventorySrv *inventory.InventoryService
}

func NewPriceCategoryHandler(ctx *context.ERPContext) *PriceCategoryHandler {
	inventorySrv, ok := ctx.InventoryService.(*inventory.InventoryService)
	if !ok {
		panic("price service is not found")
	}
	return &PriceCategoryHandler{
		ctx:          ctx,
		inventorySrv: inventorySrv,
	}
}

func (p *PriceCategoryHandler) GetPriceCategoryHandler(c *gin.Context) {
	id := c.Param("id")
	price, err := p.inventorySrv.PriceCategoryService.GetPriceCategoryByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": price, "message": "PriceCategory retrieved successfully"})
}

func (p *PriceCategoryHandler) ListPriceCategoriesHandler(c *gin.Context) {
	prices, err := p.inventorySrv.PriceCategoryService.GetPriceCategories(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": prices, "message": "PriceCategories retrieved successfully"})
}

func (p *PriceCategoryHandler) CreatePriceCategoryHandler(c *gin.Context) {
	var input models.PriceCategoryModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	input.CompanyID = &companyID
	err = p.inventorySrv.PriceCategoryService.CreatePriceCategory(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "PriceCategory created successfully", "data": input})
}

func (p *PriceCategoryHandler) UpdatePriceCategoryHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.PriceCategoryModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.inventorySrv.PriceCategoryService.UpdatePriceCategory(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "PriceCategory updated successfully"})
}

func (p *PriceCategoryHandler) DeletePriceCategoryHandler(c *gin.Context) {
	id := c.Param("id")
	err := p.inventorySrv.PriceCategoryService.DeletePriceCategory(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "PriceCategory deleted successfully"})
}
