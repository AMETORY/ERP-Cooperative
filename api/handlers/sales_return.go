package handlers

import (
	"time"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/order"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type SalesReturnHandler struct {
	ctx      *context.ERPContext
	orderSrv *order.OrderService
}

func NewSalesReturnHandler(ctx *context.ERPContext) *SalesReturnHandler {
	orderSrv, ok := ctx.OrderService.(*order.OrderService)
	if !ok {
		panic("inventory service is not found")
	}
	return &SalesReturnHandler{
		ctx:      ctx,
		orderSrv: orderSrv,
	}
}

func (s *SalesReturnHandler) GetSalesReturnListHandler(c *gin.Context) {
	salesReturn, err := s.orderSrv.SalesReturnService.GetReturns(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": salesReturn, "message": "Sales return list retrieved successfully"})
}

func (s *SalesReturnHandler) GetSalesReturnHandler(c *gin.Context) {
	id := c.Param("id")

	salesReturn, err := s.orderSrv.SalesReturnService.GetReturnByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": salesReturn, "message": "Sales return retrieved successfully"})
}

func (s *SalesReturnHandler) UpdateSalesReturnHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.ReturnModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = s.orderSrv.SalesReturnService.UpdateReturn(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Sales return updated successfully"})
}

func (s *SalesReturnHandler) ReleaseSalesReturnHandler(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Date      time.Time `json:"date"`
		AccountID *string   `json:"account_id"`
		Notes     string    `json:"notes"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(string)
	err = s.orderSrv.SalesReturnService.ReleaseReturn(id, userID, input.Date, input.Notes, input.AccountID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Sales return released successfully"})
}
func (s *SalesReturnHandler) CreateSalesReturnHandler(c *gin.Context) {
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

	err = s.orderSrv.SalesReturnService.CreateReturn(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"data": input, "message": "Sales return created successfully"})
}

func (s *SalesReturnHandler) AddItemSalesReturnHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.ReturnItemModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	salesReturn, err := s.orderSrv.SalesReturnService.GetReturnByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = s.orderSrv.SalesReturnService.AddItem(salesReturn, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"data": input, "message": "Item added successfully"})
}

func (s *SalesReturnHandler) UpdateItemSalesReturnHandler(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")
	var input models.ReturnItemModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err = s.orderSrv.SalesReturnService.GetReturnByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	input.ID = itemId
	err = s.orderSrv.SalesReturnService.UpdateItem(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Item updated successfully", "data": input})
}
func (s *SalesReturnHandler) DeleteItemSalesReturnHandler(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")

	err := s.orderSrv.SalesReturnService.DeleteItem(id, itemId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Item deleted successfully"})
}

func (s *SalesReturnHandler) DeleteSalesReturnHandler(c *gin.Context) {
	id := c.Param("id")

	err := s.orderSrv.SalesReturnService.DeleteReturn(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Sales return deleted successfully"})
}
