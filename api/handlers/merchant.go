package handlers

import (
	"fmt"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/order"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	ctx      *context.ERPContext
	orderSrc *order.OrderService
}

func NewMerchantHandler(ctx *context.ERPContext) *MerchantHandler {
	orderSrc, ok := ctx.OrderService.(*order.OrderService)
	if !ok {
		panic("order service is not found")
	}
	return &MerchantHandler{
		ctx:      ctx,
		orderSrc: orderSrc,
	}
}

func (p *MerchantHandler) GetMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	merchant, err := p.orderSrc.MerchantService.GetMerchantByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": merchant, "message": "Merchant retrieved successfully"})
}

func (p *MerchantHandler) ListMerchantsHandler(c *gin.Context) {
	merchants, err := p.orderSrc.MerchantService.GetMerchants(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": merchants, "message": "Merchants retrieved successfully"})
}

func (p *MerchantHandler) CreateMerchantHandler(c *gin.Context) {
	var input models.MerchantModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	input.CompanyID = &companyID
	err = p.orderSrc.MerchantService.CreateMerchant(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Merchant created successfully", "data": input})
}

func (p *MerchantHandler) UpdateMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.MerchantModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.orderSrc.MerchantService.UpdateMerchant(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if input.Picture != nil {
		input.Picture.RefID = input.ID
		input.Picture.RefType = "merchant"
		err = p.ctx.DB.Save(&input.Picture).Error
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(200, gin.H{"message": "Merchant updated successfully"})
}

func (p *MerchantHandler) DeleteMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	merchant, err := p.orderSrc.MerchantService.GetMerchantByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if merchant.CompanyID == nil {
		c.JSON(403, gin.H{"error": "You do not have permission to delete this merchant"})
	}
	err = p.orderSrc.MerchantService.DeleteMerchant(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Merchant deleted successfully"})
}

func (p *MerchantHandler) AddProductMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	input := struct {
		ProductIDs []string `json:"product_ids"`
	}{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.orderSrc.MerchantService.AddProductsToMerchant(id, input.ProductIDs)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Product added to merchant successfully"})
}

func (p *MerchantHandler) GetMerchantProductsHandler(c *gin.Context) {
	id := c.Param("id")
	wareID := c.Query("warehouse_id")
	var warehouseID *string
	if wareID != "" {
		warehouseID = &wareID
	}
	products, err := p.orderSrc.MerchantService.GetMerchantProducts(*c.Request, c.Query("search"), id, warehouseID, nil, false)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": products, "message": "Merchant products retrieved successfully"})
}

func (p *MerchantHandler) DeleteProductsFromMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	input := struct {
		ProductIDs []string `json:"product_ids"`
	}{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.orderSrc.MerchantService.DeleteProductsFromMerchant(id, input.ProductIDs)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Product deleted from merchant successfully"})
}

func (p *MerchantHandler) GetMerchantUsersHandler(c *gin.Context) {
	id := c.Param("id")
	companyID := c.MustGet("companyID").(string)
	users, err := p.orderSrc.MerchantService.GetMerchantUsers(*c.Request, c.Query("search"), id, companyID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": users, "message": "Merchant users retrieved successfully"})
}

func (p *MerchantHandler) AddUserMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	input := struct {
		UserIDs []string `json:"user_ids"`
	}{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.orderSrc.MerchantService.AddMerchantUser(id, input.UserIDs)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User added to merchant successfully"})
}

func (p *MerchantHandler) DeleteUserFromMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	input := struct {
		UserIDs []string `json:"user_ids"`
	}{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.orderSrc.MerchantService.DeleteUserFromMerchant(id, input.UserIDs)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted from merchant successfully"})
}

func (p *MerchantHandler) AddDeskMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	input := models.MerchantDesk{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.orderSrc.MerchantService.AddDeskToMerchant(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Desk added to merchant successfully"})
}

func (p *MerchantHandler) GetDeskMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("GET DESK FROM MERCHANT", id)
	desks, err := p.orderSrc.MerchantService.GetDesksFromID(*c.Request, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": desks, "message": "Merchant desk retrieved successfully"})
}

func (p *MerchantHandler) UpdateDeskMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	deskId := c.Param("deskId")
	input := models.MerchantDesk{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.orderSrc.MerchantService.UpdateMerchantDesk(id, deskId, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Merchant desk updated successfully"})
}

func (p *MerchantHandler) DeleteDeskMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	deskId := c.Param("deskId")

	err := p.orderSrc.MerchantService.DeleteDeskFromMerchant(id, deskId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Desk deleted from merchant successfully"})
}

func (p *MerchantHandler) GetLayoutDetailMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	layoutId := c.Param("layoutId")

	layout, err := p.orderSrc.MerchantService.GetLayoutDetailFromID(id, layoutId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": layout, "message": "Merchant layout retrieved successfully"})
}

func (p *MerchantHandler) GetLayoutsMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	desks, err := p.orderSrc.MerchantService.GetLayoutsFromID(*c.Request, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": desks, "message": "Merchant layout retrieved successfully"})
}

func (p *MerchantHandler) AddLayoutMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	input := models.MerchantDeskLayout{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.orderSrc.MerchantService.AddLayoutToMerchant(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Layout added to merchant successfully"})
}

func (p *MerchantHandler) UpdateLayoutMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	layoutId := c.Param("layoutId")
	input := models.MerchantDeskLayout{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("UPDATE MERCHANT LAYOUT", id, layoutId, input)
	err = p.orderSrc.MerchantService.UpdateLayoutMerchant(id, layoutId, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Layout updated to merchant successfully"})
}

func (p *MerchantHandler) DeleteLayoutMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	layoutId := c.Param("layoutId")

	err := p.orderSrc.MerchantService.DeleteLayoutMerchant(id, layoutId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Layout deleted from merchant successfully"})
}

func (p *MerchantHandler) GetMerchantStationsHandler(c *gin.Context) {
	id := c.Param("id")
	stations, err := p.orderSrc.MerchantService.GetMerchantStations(*c.Request, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": stations, "message": "Merchant stations retrieved successfully"})
}

func (p *MerchantHandler) GetMerchantStationDetailHandler(c *gin.Context) {
	id := c.Param("id")
	stationId := c.Param("stationId")

	station, err := p.orderSrc.MerchantService.GetMerchantStationDetail(id, stationId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": station, "message": "Merchant station detail retrieved successfully"})
}

func (p *MerchantHandler) CreateMerchantStationHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.MerchantStation
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.orderSrc.MerchantService.CreateMerchantStation(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Merchant station created successfully"})
}

func (p *MerchantHandler) UpdateMerchantStationHandler(c *gin.Context) {
	id := c.Param("id")
	stationId := c.Param("stationId")
	var input models.MerchantStation
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.orderSrc.MerchantService.UpdateMerchantStation(id, stationId, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Merchant station updated successfully"})
}

func (p *MerchantHandler) DeleteMerchantStationHandler(c *gin.Context) {
	id := c.Param("id")
	stationId := c.Param("stationId")

	err := p.orderSrc.MerchantService.DeleteMerchantStation(id, stationId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Merchant station deleted successfully"})
}

func (p *MerchantHandler) AddProductMerchantStationHandler(c *gin.Context) {
	id := c.Param("id")
	stationId := c.Param("stationId")
	var input struct {
		ProductIDs []string `json:"product_ids"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.orderSrc.MerchantService.AddProductsToMerchantStation(id, stationId, input.ProductIDs)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Product added to merchant station successfully"})
}
func (p *MerchantHandler) DeleteProductMerchantStationHandler(c *gin.Context) {
	id := c.Param("id")
	stationId := c.Param("stationId")
	var input struct {
		ProductIDs []string `json:"product_ids"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.orderSrc.MerchantService.DeleteProductFromMerchantStation(id, stationId, input.ProductIDs)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Product deleted from merchant station successfully"})
}

func (p *MerchantHandler) GetProductsMerchantStationHandler(c *gin.Context) {
	id := c.Param("id")
	stationId := c.Param("stationId")

	products, err := p.orderSrc.MerchantService.GetProductsFromMerchantStation(id, stationId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("PRODUCT MERCHANT", products)
	c.JSON(200, gin.H{"data": products, "message": "Products from merchant station retrieved successfully"})
}
