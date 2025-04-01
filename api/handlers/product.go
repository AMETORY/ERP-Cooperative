package handlers

import (
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/inventory"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	ctx          *context.ERPContext
	inventorySrv *inventory.InventoryService
}

func NewProductHandler(ctx *context.ERPContext) *ProductHandler {
	inventorySrv, ok := ctx.InventoryService.(*inventory.InventoryService)
	if !ok {
		panic("product service is not found")
	}
	return &ProductHandler{
		ctx:          ctx,
		inventorySrv: inventorySrv,
	}
}

func (p *ProductHandler) GetProductHandler(c *gin.Context) {
	id := c.Param("id")
	product, err := p.inventorySrv.ProductService.GetProductByID(id, c.Request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": product, "message": "Product retrieved successfully"})
}

func (p *ProductHandler) ListProductsHandler(c *gin.Context) {
	products, err := p.inventorySrv.ProductService.GetProducts(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": products, "message": "Products retrieved successfully"})
}

func (p *ProductHandler) CreateProductHandler(c *gin.Context) {
	var input models.ProductModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	input.CompanyID = &companyID
	err = p.inventorySrv.ProductService.CreateProduct(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Product created successfully", "data": input})
}

func (p *ProductHandler) UpdateProductHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.ProductModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if id != input.ID {
		c.JSON(400, gin.H{"error": "ID mismatch"})
	}

	err = p.inventorySrv.ProductService.UpdateProduct(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Product updated successfully"})
}

func (p *ProductHandler) DeleteProductHandler(c *gin.Context) {
	id := c.Param("id")
	err := p.inventorySrv.ProductService.DeleteProduct(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Product deleted successfully"})
}
