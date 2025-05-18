package handlers

import (
	"net/http"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/inventory"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type BrandHandler struct {
	ctx              *context.ERPContext
	inventoryService *inventory.InventoryService
}

func NewBrandHandler(ctx *context.ERPContext) *BrandHandler {
	var inventorySrv *inventory.InventoryService
	invSrv, ok := ctx.InventoryService.(*inventory.InventoryService)
	if ok {
		inventorySrv = invSrv
	}
	return &BrandHandler{ctx: ctx, inventoryService: inventorySrv}
}

func (h *BrandHandler) CreateBrandHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to create an Brand
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	var data models.BrandModel

	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	data.CompanyID = &companyID
	err = h.inventoryService.BrandService.CreateBrand(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Brand created successfully"})
}

func (h *BrandHandler) GetBrandHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to get an Brand
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	search, _ := c.GetQuery("search")
	data, err := h.inventoryService.BrandService.GetBrands(*c.Request, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Brand retrieved successfully", "data": data})
}

func (h *BrandHandler) GetBrandByIdHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to get an Brand by ID
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	id := c.Param("id")
	data, err := h.inventoryService.BrandService.GetBrandByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Brand retrieved successfully", "data": data})
}

func (h *BrandHandler) UpdateBrandHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to update an Brand
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	var data models.BrandModel
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	_, err = h.inventoryService.BrandService.GetBrandByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	h.inventoryService.BrandService.UpdateBrand(id, &data)
	c.JSON(http.StatusOK, gin.H{"message": "Brand updated successfully"})
}

func (h *BrandHandler) DeleteBrandHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to delete an Brand
	if h.inventoryService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	id := c.Param("id")
	err := h.inventoryService.BrandService.DeleteBrand(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Brand deleted successfully"})
}

// func (h *BrandHandler) GetTopBrandHandler(c *gin.Context) {
// 	h.ctx.Request = c.Request
// 	// Implement logic to get TopBrands
// 	if h.ctx.InternalService == nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal service is not initialized"})
// 	}
// 	limitStr := c.Query("limit")
// 	if limitStr == "" {
// 		limitStr = "20"
// 	}
// 	limit, _ := strconv.Atoi(limitStr)
// 	data, err := h.ctx.InternalService.(*additional.InternalService).ProductService.GetTopBrands(limit)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "TopBrand retrieved successfully", "data": data})
// }

// func (h *BrandHandler) AddTopBrandHandler(c *gin.Context) {
// 	h.ctx.Request = c.Request
// 	// Implement logic to add an TopBrand
// 	if h.ctx.InternalService == nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal service is not initialized"})
// 	}
// 	var data app_models.TopBrand
// 	err := c.ShouldBindBodyWithJSON(&data)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	err = h.ctx.InternalService.(*additional.InternalService).ProductService.AddTopBrand(data.BrandID, data.Value)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "TopBrand added successfully"})
// }

// func (h *BrandHandler) UpdateTopBrandHandler(c *gin.Context) {
// 	h.ctx.Request = c.Request
// 	// Implement logic to update an TopBrand
// 	if h.ctx.InternalService == nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal service is not initialized"})
// 	}
// 	brandId := c.Param("brandId")
// 	var data app_models.TopBrand
// 	err := c.ShouldBindBodyWithJSON(&data)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	err = h.ctx.InternalService.(*additional.InternalService).ProductService.UpdateTopBrand(brandId, data.Value)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "TopBrand updated successfully"})
// }

// func (h *BrandHandler) DeleteTopBrandHandler(c *gin.Context) {
// 	h.ctx.Request = c.Request
// 	// Implement logic to delete a TopBrand
// 	if h.ctx.InternalService == nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal service is not initialized"})
// 		return
// 	}
// 	brandId := c.Param("brandId")
// 	err := h.ctx.InternalService.(*additional.InternalService).ProductService.DeleteTopBrand(brandId)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "TopBrand deleted successfully"})
// }
