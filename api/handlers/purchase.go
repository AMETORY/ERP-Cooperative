package handlers

import (
	"ametory-cooperative/objects"
	"encoding/json"
	"time"

	"github.com/AMETORY/ametory-erp-modules/contact"
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/inventory"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/AMETORY/ametory-erp-modules/utils"
	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	ctx          *context.ERPContext
	inventorySrv *inventory.InventoryService
	contactSrv   *contact.ContactService
}

func NewPurchaseHandler(ctx *context.ERPContext) *PurchaseHandler {
	inventorySrv, ok := ctx.InventoryService.(*inventory.InventoryService)
	if !ok {
		panic("inventory service is not found")
	}
	contactSrv, ok := ctx.ContactService.(*contact.ContactService)
	if !ok {
		panic("contact service is not found")
	}
	return &PurchaseHandler{
		ctx:          ctx,
		inventorySrv: inventorySrv,
		contactSrv:   contactSrv,
	}
}

func (s *PurchaseHandler) GetPurchaseHandler(c *gin.Context) {
	purchase, err := s.inventorySrv.PurchaseService.GetPurchases(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": purchase, "message": "Purchase retrieved successfully"})
}

func (s *PurchaseHandler) GetPurchaseByIdHandler(c *gin.Context) {
	id := c.Param("id")
	purchase, err := s.inventorySrv.PurchaseService.GetPurchaseByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": purchase, "message": "Purchase retrieved successfully"})
}

func (s *PurchaseHandler) CreatePurchaseHandler(c *gin.Context) {
	var input objects.PurchaseRequest
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

	var data models.PurchaseOrderModel = models.PurchaseOrderModel{
		PurchaseNumber:   input.PurchaseNumber,
		Code:             utils.RandString(8, false),
		Description:      input.Description,
		Notes:            input.Notes,
		Status:           "DRAFT",
		PurchaseDate:     input.PurchaseDate,
		DueDate:          input.DueDate,
		PaymentTerms:     input.PaymentTerms,
		ContactID:        input.ContactID,
		Type:             input.Type,
		DocumentType:     input.DocumentType,
		ContactData:      string(b),
		RefID:            input.RefID,
		RefType:          input.RefType,
		SecondaryRefID:   input.SecondaryRefID,
		SecondaryRefType: input.SecondaryRefType,
		PaymentTermsCode: input.PaymentTermsCode,
		TermCondition:    input.TermCondition,
		TaxBreakdown:     "{}",
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
		items := make([]models.PurchaseOrderItemModel, len(input.Items))
		for _, item := range input.Items {
			item.PurchaseID = &data.ID
			item.ID = utils.Uuid()
			err = s.inventorySrv.PurchaseService.AddItem(&data, &item)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			items = append(items, item)
		}
		data.Items = items
		s.inventorySrv.PurchaseService.UpdateTotal(&data)
	}
	c.JSON(200, gin.H{"message": "Purchase created successfully", "data": data})
}

func (s *PurchaseHandler) PostPurchaseHandler(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Purchase        models.PurchaseOrderModel `json:"purchase"`
		TransactionDate time.Time                 `json:"transaction_date"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(string)
	err = s.inventorySrv.PurchaseService.PostPurchase(id, &input.Purchase, userID, input.TransactionDate)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Invoice Posted successfully"})
}
func (s *PurchaseHandler) PublishPurchaseHandler(c *gin.Context) {
	id := c.Param("id")
	purchase, err := s.inventorySrv.PurchaseService.GetPurchaseByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	userID := c.MustGet("userID").(string)
	purchase.PublishedAt = &now
	purchase.PublishedByID = &userID
	purchase.Status = "RELEASED"
	err = s.inventorySrv.PurchaseService.UpdatePurchase(id, purchase)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": purchase, "message": "Purchase retrieved successfully"})
}

func (s *PurchaseHandler) UpdatePurchaseHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.PurchaseOrderModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = s.inventorySrv.PurchaseService.UpdatePurchase(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Purchase updated successfully"})
}

func (s *PurchaseHandler) DeletePurchaseHandler(c *gin.Context) {
	id := c.Param("id")
	err := s.inventorySrv.PurchaseService.DeletePurchase(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Purchase deleted successfully"})
}

func (s *PurchaseHandler) GetItemsHandler(c *gin.Context) {
	id := c.Param("id")

	items, err := s.inventorySrv.PurchaseService.GetItems(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": items, "message": "Items retrieved successfully"})
}
func (s *PurchaseHandler) AddItemHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.PurchaseOrderItemModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	purchase, err := s.inventorySrv.PurchaseService.GetPurchaseByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = s.inventorySrv.PurchaseService.AddItem(purchase, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Item added successfully", "data": input})
}

func (s *PurchaseHandler) DeleteItemHandler(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")

	purchase, err := s.inventorySrv.PurchaseService.GetPurchaseByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = s.inventorySrv.PurchaseService.DeleteItem(purchase, itemId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Item delete successfully"})
}

func (s *PurchaseHandler) UpdateItemHandler(c *gin.Context) {
	id := c.Param("id")
	itemID := c.Param("itemID")
	var input models.PurchaseOrderItemModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	purchase, err := s.inventorySrv.PurchaseService.GetPurchaseByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = s.inventorySrv.PurchaseService.UpdateItem(purchase, itemID, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Item updated successfully", "data": input})
}
