package handlers

import (
	"ametory-cooperative/services"
	"net/http"

	"github.com/AMETORY/ametory-erp-modules/contact"
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	ctx            *context.ERPContext
	contactService *contact.ContactService
	appService     *services.AppService
}

func NewContactHandler(ctx *context.ERPContext) *ContactHandler {
	contactService, ok := ctx.ContactService.(*contact.ContactService)
	if !ok {
		panic("invalid contact service")
	}

	var appService *services.AppService
	appSrv, ok := ctx.AppService.(*services.AppService)
	if ok {
		appService = appSrv
	}

	return &ContactHandler{
		ctx:            ctx,
		contactService: contactService,
		appService:     appService,
	}
}

func (h *ContactHandler) CreateContactHandler(c *gin.Context) {
	var contact models.ContactModel
	if err := c.BindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.GetHeader("ID-Company")
	contact.CompanyID = &companyID

	if err := h.contactService.CreateContact(&contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact created successfully", "data": contact})
}

func (h *ContactHandler) GetContactHandler(c *gin.Context) {
	id := c.Param("id")

	contact, err := h.contactService.GetContactByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": contact, "message": "Contact created successfully"})
}

func (h *ContactHandler) UpdateContactHandler(c *gin.Context) {
	id := c.Param("id")
	var contact models.ContactModel
	if err := c.BindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.contactService.UpdateContact(id, &contact)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.ctx.DB.Model(&contact).Where("id = ?", id).Updates(map[string]any{
		"is_customer": contact.IsCustomer,
		"is_vendor":   contact.IsVendor,
		"is_supplier": contact.IsSupplier,
	})

	if contact.DebtLimit == 0 {
		h.ctx.DB.Model(&contact).Where("id = ?", id).Updates(map[string]any{
			"debt_limit": 0,
		})
	}
	if contact.ReceivablesLimit == 0 {
		h.ctx.DB.Model(&contact).Where("id = ?", id).Updates(map[string]any{
			"receivables_limit": 0,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact created successfully"})
}

func (h *ContactHandler) DeleteContactHandler(c *gin.Context) {
	id := c.Param("id")
	if err := h.contactService.DeleteContact(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted"})
}

func (h *ContactHandler) GetContactsHandler(c *gin.Context) {

	var isCustomer, isVendor, isSupplier bool

	if c.Query("is_customer") == "true" || c.Query("is_customer") == "1" {
		isCustomer = true
	}

	if c.Query("is_vendor") == "true" || c.Query("is_vendor") == "1" {
		isVendor = true
	}

	if c.Query("is_supplier") == "true" || c.Query("is_supplier") == "1" {
		isSupplier = true
	}
	contacts, err := h.contactService.GetContacts(*c.Request, c.Query("search"), &isCustomer, &isVendor, &isSupplier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": contacts, "message": "Contact created successfully"})
}
