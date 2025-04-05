package handlers

import (
	"ametory-cooperative/objects"
	"encoding/json"
	"time"

	"github.com/AMETORY/ametory-erp-modules/contact"
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/order"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/AMETORY/ametory-erp-modules/utils"
	"github.com/gin-gonic/gin"
)

type SalesHandler struct {
	ctx        *context.ERPContext
	orderSrv   *order.OrderService
	contactSrv *contact.ContactService
}

func NewSalesHandler(ctx *context.ERPContext) *SalesHandler {
	orderSrv, ok := ctx.OrderService.(*order.OrderService)
	if !ok {
		panic("order service is not found")
	}
	contactSrv, ok := ctx.ContactService.(*contact.ContactService)
	if !ok {
		panic("contact service is not found")
	}
	return &SalesHandler{
		ctx:        ctx,
		orderSrv:   orderSrv,
		contactSrv: contactSrv,
	}
}

func (s *SalesHandler) GetSalesHandler(c *gin.Context) {
	sales, err := s.orderSrv.SalesService.GetSales(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": sales, "message": "Sales retrieved successfully"})
}

func (s *SalesHandler) GetSalesByIdHandler(c *gin.Context) {
	id := c.Param("id")
	sales, err := s.orderSrv.SalesService.GetSalesByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": sales, "message": "Sales retrieved successfully"})
}

func (s *SalesHandler) CreateSalesHandler(c *gin.Context) {
	var input objects.SalesRequest
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	contact, err := s.contactSrv.GetContactByID(*input.ContactID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	b, _ := json.Marshal(*contact)

	var data models.SalesModel = models.SalesModel{
		SalesNumber:      input.SalesNumber,
		Code:             utils.RandString(8, false),
		Description:      input.Description,
		Notes:            input.Notes,
		Status:           "DRAFT",
		SalesDate:        input.SalesDate,
		DueDate:          input.DueDate,
		PaymentTerms:     input.PaymentTerms,
		ContactID:        input.ContactID,
		Type:             input.Type,
		DocumentType:     input.DocumentType,
		ContactData:      string(b),
		DeliveryData:     "{}",
		TaxBreakdown:     "{}",
		RefID:            input.RefID,
		RefType:          input.RefType,
		SecondaryRefID:   input.SecondaryRefID,
		SecondaryRefType: input.SecondaryRefType,
		PaymentTermsCode: input.PaymentTermsCode,
		TermCondition:    input.TermCondition,
	}
	if input.DeliveryID != nil && input.DeliveryData != "" {
		data.DeliveryID = input.DeliveryID
		data.DeliveryData = input.DeliveryData
	}

	if input.Status != "" {
		data.Status = input.Status
	}
	companyID := c.MustGet("companyID").(string)
	userID := c.MustGet("userID").(string)
	data.CompanyID = &companyID
	data.UserID = &userID
	s.ctx.Request = c.Request
	err = s.ctx.DB.Create(&data).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(input.Items) > 0 {
		items := make([]models.SalesItemModel, len(input.Items))
		for _, item := range input.Items {
			item.SalesID = &data.ID
			item.ID = utils.Uuid()
			err = s.orderSrv.SalesService.AddItem(&data, &item)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			items = append(items, item)
		}
		data.Items = items
		s.orderSrv.SalesService.UpdateTotal(&data)
	}
	c.JSON(201, gin.H{"message": "Sales created successfully", "data": data})
}

func (s *SalesHandler) PostInvoiceHandler(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Sales           models.SalesModel `json:"sales"`
		TransactionDate time.Time         `json:"transaction_date"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(string)
	err = s.orderSrv.SalesService.PostInvoice(id, &input.Sales, userID, input.TransactionDate)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Invoice Posted successfully"})
}

func (s *SalesHandler) PublishSalesHandler(c *gin.Context) {
	id := c.Param("id")
	sales, err := s.orderSrv.SalesService.GetSalesByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	userID := c.MustGet("userID").(string)
	sales.PublishedAt = &now
	sales.PublishedByID = &userID
	sales.Status = "RELEASED"
	err = s.orderSrv.SalesService.UpdateSales(id, sales)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": sales, "message": "Purchase retrieved successfully"})
}

func (s *SalesHandler) UpdateSalesHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.SalesModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = s.orderSrv.SalesService.UpdateSales(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Sales updated successfully"})
}

func (s *SalesHandler) DeleteSalesHandler(c *gin.Context) {
	id := c.Param("id")
	err := s.orderSrv.SalesService.DeleteSales(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Sales deleted successfully"})
}

func (s *SalesHandler) GetItemsHandler(c *gin.Context) {
	id := c.Param("id")

	items, err := s.orderSrv.SalesService.GetItems(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": items, "message": "Items retrieved successfully"})
}
func (s *SalesHandler) AddItemHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.SalesItemModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	sales, err := s.orderSrv.SalesService.GetSalesByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = s.orderSrv.SalesService.AddItem(sales, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Item added successfully", "data": input})
}
func (s *SalesHandler) PaymentHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.SalesPaymentModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	sales, err := s.orderSrv.SalesService.GetSalesByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("userID").(string)
	companyID := c.MustGet("companyID").(string)
	sales.UserID = &userID
	sales.CompanyID = &companyID
	err = s.orderSrv.SalesService.CreateSalesPayment(sales, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Payment added successfully", "data": input})
}

func (s *SalesHandler) DeleteItemHandler(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")

	sales, err := s.orderSrv.SalesService.GetSalesByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = s.orderSrv.SalesService.DeleteItem(sales, itemId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Item delete successfully"})
}

func (s *SalesHandler) UpdateItemHandler(c *gin.Context) {
	id := c.Param("id")
	itemID := c.Param("itemID")
	var input models.SalesItemModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	sales, err := s.orderSrv.SalesService.GetSalesByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = s.orderSrv.SalesService.UpdateItem(sales, itemID, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Item updated successfully", "data": input})
}
