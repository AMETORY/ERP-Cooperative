package handlers

import (
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"

	"github.com/AMETORY/ametory-erp-modules/inventory"
)

type PurchaseReturnHandler struct {
	ctx          *context.ERPContext
	inventorySrv *inventory.InventoryService
}

func NewPurchaseReturnHandler(ctx *context.ERPContext) *PurchaseReturnHandler {
	inventorySrv, ok := ctx.InventoryService.(*inventory.InventoryService)
	if !ok {
		panic("inventory service is not found")
	}
	return &PurchaseReturnHandler{
		ctx:          ctx,
		inventorySrv: inventorySrv,
	}
}

func (s *PurchaseReturnHandler) GetPurchaseReturnListHandler(c *gin.Context) {
	purchaseReturn, err := s.inventorySrv.PurchaseReturnService.GetReturns(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": purchaseReturn, "message": "Purchase return list retrieved successfully"})
}

func (s *PurchaseReturnHandler) GetPurchaseReturnHandler(c *gin.Context) {
	id := c.Param("id")

	purchaseReturn, err := s.inventorySrv.PurchaseReturnService.GetReturnByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": purchaseReturn, "message": "Purchase return retrieved successfully"})
}

func (s *PurchaseReturnHandler) CreatePurchaseReturnHandler(c *gin.Context) {
	var input models.ReturnModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(string)
	companyID := c.MustGet("companyID").(string)
	input.UserID = &userID
	input.CompanyID = &companyID

	err = s.inventorySrv.PurchaseReturnService.CreateReturn(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"data": input, "message": "Purchase return created successfully"})
}

func (s *PurchaseReturnHandler) AddItemPurchaseReturnHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.ReturnItemModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	purchaseReturn, err := s.inventorySrv.PurchaseReturnService.GetReturnByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = s.inventorySrv.PurchaseReturnService.AddItem(purchaseReturn, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"data": input, "message": "Item added successfully"})
}

func (s *PurchaseReturnHandler) UpdateItemPurchaseReturnHandler(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")
	var input models.ReturnItemModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err = s.inventorySrv.PurchaseReturnService.GetReturnByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	input.ID = itemId
	err = s.inventorySrv.PurchaseReturnService.UpdateItem(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Item updated successfully", "data": input})
}
func (s *PurchaseReturnHandler) DeleteItemPurchaseReturnHandler(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")

	err := s.inventorySrv.PurchaseReturnService.DeleteItem(id, itemId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Item deleted successfully"})
}

func (s *PurchaseReturnHandler) DeletePurchaseReturnHandler(c *gin.Context) {
	id := c.Param("id")

	err := s.inventorySrv.PurchaseReturnService.DeleteReturn(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Purchase return deleted successfully"})
}
