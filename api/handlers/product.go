package handlers

import (
	"net/http"

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
func (p *ProductHandler) GetProductVariantHandler(c *gin.Context) {
	p.ctx.Request = c.Request
	// Implement logic to get a variant of a product

	id := c.Param("id")
	data, err := p.inventorySrv.ProductService.GetProductVariants(id, *c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product variant retrieved successfully", "data": data})
}
func (p *ProductHandler) CreateProductVariantHandler(c *gin.Context) {
	p.ctx.Request = c.Request
	// Implement logic to create a variant of a product

	var data models.VariantModel
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	product, err := p.inventorySrv.ProductService.GetProductByID(id, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	data.ProductID = id
	if data.Height == 0 {
		data.Height = product.Height
	}
	if data.Width == 0 {
		data.Width = product.Width
	}
	if data.Length == 0 {
		data.Length = product.Length
	}
	if data.Weight == 0 {
		data.Weight = product.Weight
	}
	err = p.inventorySrv.ProductService.CreateProductVariant(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product variant created successfully"})
}
func (p *ProductHandler) GetProductDiscountHandler(c *gin.Context) {
	p.ctx.Request = c.Request

	productId := c.Param("id")
	data, err := p.inventorySrv.ProductService.GetAllDiscountByProductID(productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Active discounts retrieved successfully", "data": data})
}

func (h *ProductHandler) AddDiscountHandler(c *gin.Context) {
	h.ctx.Request = c.Request

	var data models.DiscountModel
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	productId := c.Param("id")
	discount, err := h.inventorySrv.ProductService.AddDiscount(productId, data.Type, data.Value, data.StartDate, data.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "discount added successfully", "data": discount})
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

func (h *ProductHandler) AddPriceProductHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to update an product

	var data models.PriceModel
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	_, err = h.inventorySrv.ProductService.GetProductByID(id, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	h.inventorySrv.ProductService.AddPriceToProduct(id, &data)
	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}
