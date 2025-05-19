package pos

import (
	"fmt"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/order"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/AMETORY/ametory-erp-modules/utils"
	"github.com/gin-gonic/gin"
)

type PosHandler struct {
	ctx        *context.ERPContext
	financeSrv *finance.FinanceService
	OrderSrv   *order.OrderService
}

func NewPosHandler(ctx *context.ERPContext) *PosHandler {
	financeSrv, ok := ctx.FinanceService.(*finance.FinanceService)
	if !ok {
		panic("finance service is not found")
	}

	orderSrv, ok := ctx.OrderService.(*order.OrderService)
	if !ok {
		panic("order service is not found")
	}
	return &PosHandler{
		ctx:        ctx,
		financeSrv: financeSrv,
		OrderSrv:   orderSrv,
	}
}

func (p *PosHandler) GetMerchantsHandler(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	companyID := c.MustGet("companyID").(string)
	merchants, err := p.OrderSrv.MerchantService.GetMerchantsByUserID(*c.Request, userID, companyID, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": merchants, "message": "Merchants retrieved successfully"})
}

func (p *PosHandler) GetLayoutsMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	desks, err := p.OrderSrv.MerchantService.GetLayoutsFromID(*c.Request, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": desks, "message": "Merchant layout retrieved successfully"})
}

func (p *PosHandler) UpdateStatusTableHandler(c *gin.Context) {
	id := c.Param("id")
	tableId := c.Param("tableId")

	input := struct {
		Status string `json:"status"`
	}{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	utils.LogJson(input)
	err := p.OrderSrv.MerchantService.UpdateTableStatus(id, tableId, input.Status)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Table status updated successfully"})

}
func (p *PosHandler) GetLayoutDetailMerchantHandler(c *gin.Context) {
	id := c.Param("id")
	layoutId := c.Param("layoutId")

	layout, err := p.OrderSrv.MerchantService.GetLayoutDetailFromID(id, layoutId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": layout, "message": "Merchant layout retrieved successfully"})
}

func (p *PosHandler) CreateOrderHandler(c *gin.Context) {
	var input models.MerchantOrder
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	merchantID := c.MustGet("merchantID").(string)

	err = p.OrderSrv.MerchantService.CreateOrder(merchantID, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if input.NextStep == "distribute" {
		fmt.Println("DISTRIBUTE ORDER")
		orderStations, err := p.OrderSrv.MerchantService.DistributeOrder(merchantID, &input)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		orderStationMap := make(map[string][]models.MerchantStationOrder)
		for _, v := range orderStations {
			orderStationMap[*v.MerchantStationID] = append(orderStationMap[*v.MerchantStationID], v)
		}

		utils.LogJson(orderStationMap)

	}
	c.JSON(200, gin.H{"message": "Merchant order created successfully"})
}
