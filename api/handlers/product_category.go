package handlers

import (
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/inventory"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type ProductCategoryHandler struct {
	ctx          *context.ERPContext
	inventorySrv *inventory.InventoryService
}

func NewProductCategoryHandler(ctx *context.ERPContext) *ProductCategoryHandler {
	inventorySrv, ok := ctx.InventoryService.(*inventory.InventoryService)
	if !ok {
		panic("product service is not found")
	}
	return &ProductCategoryHandler{
		ctx:          ctx,
		inventorySrv: inventorySrv,
	}
}

func (p *ProductCategoryHandler) GetProductCategoryHandler(c *gin.Context) {
	id := c.Param("id")
	product, err := p.inventorySrv.ProductCategoryService.GetProductCategoryByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": product, "message": "ProductCategory retrieved successfully"})
}

func (p *ProductCategoryHandler) ListProductCategoriesHandler(c *gin.Context) {
	products, err := p.inventorySrv.ProductCategoryService.GetProductCategories(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": products, "message": "ProductCategories retrieved successfully"})
}

func (p *ProductCategoryHandler) CreateProductCategoryHandler(c *gin.Context) {
	var input models.ProductCategoryModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	input.CompanyID = &companyID
	err = p.inventorySrv.ProductCategoryService.CreateProductCategory(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "ProductCategory created successfully", "data": input})
}

func (p *ProductCategoryHandler) UpdateProductCategoryHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.ProductCategoryModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.inventorySrv.ProductCategoryService.UpdateProductCategory(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "ProductCategory updated successfully"})
}

func (p *ProductCategoryHandler) DeleteProductCategoryHandler(c *gin.Context) {
	id := c.Param("id")
	err := p.inventorySrv.ProductCategoryService.DeleteProductCategory(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "ProductCategory deleted successfully"})
}
