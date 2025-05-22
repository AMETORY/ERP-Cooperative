package pos

import (
	"ametory-cooperative/services"
	"encoding/json"
	"fmt"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/order"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/AMETORY/ametory-erp-modules/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type PosHandler struct {
	ctx        *context.ERPContext
	financeSrv *finance.FinanceService
	OrderSrv   *order.OrderService
	appService *services.AppService
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
	appService, ok := ctx.AppService.(*services.AppService)
	if !ok {
		panic("AppService is not instance of app.AppService")
	}
	return &PosHandler{
		ctx:        ctx,
		financeSrv: financeSrv,
		OrderSrv:   orderSrv,
		appService: appService,
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

func (p *PosHandler) GetOrdersHandler(c *gin.Context) {
	merchantID := c.MustGet("merchantID").(string)
	orders, err := p.OrderSrv.MerchantService.GetOrders(*c.Request, merchantID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": orders, "message": "Orders retrieved successfully"})
}

func (p *PosHandler) GetOrderDetailHandler(c *gin.Context) {
	merchantID := c.MustGet("merchantID").(string)
	orderID := c.Param("orderId")

	order, err := p.OrderSrv.MerchantService.GetOrderDetail(merchantID, orderID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": order, "message": "Order detail retrieved successfully"})
}

func (p *PosHandler) GetStationsHandler(c *gin.Context) {
	merchantID := c.MustGet("merchantID").(string)
	stations, err := p.OrderSrv.MerchantService.GetMerchantStations(*c.Request, merchantID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": stations, "message": "Stations retrieved successfully"})
}

func (p *PosHandler) GetStationDetailHandler(c *gin.Context) {
	merchantID := c.MustGet("merchantID").(string)
	stationID := c.Param("stationId")

	station, err := p.OrderSrv.MerchantService.GetMerchantStationDetail(merchantID, stationID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": station, "message": "Station detail retrieved successfully"})
}

func (p *PosHandler) GetStationOrdersHandler(c *gin.Context) {
	merchantID := c.MustGet("merchantID").(string)
	stationID := c.Param("stationId")

	orders, err := p.OrderSrv.MerchantService.GetOrdersFromStation(*c.Request, merchantID, stationID, []string{"PENDING", "DISTRIBUTING", "PROCESSING", "COMPLETED"})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": orders, "message": "Orders from station retrieved successfully"})
}

func (p *PosHandler) UpdateStationOrderHandler(c *gin.Context) {
	merchantID := c.MustGet("merchantID").(string)
	stationID := c.Param("stationId")
	orderID := c.Param("orderId")

	var input struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	stationOrder, err := p.OrderSrv.MerchantService.GetMerchantOrderStation(orderID, stationID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = p.OrderSrv.MerchantService.UpdateStationOrderStatus(stationID, orderID, input.Status)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("userID").(string)

	msg := gin.H{
		"command":             "ORDER_STATION_UPDATED",
		"message":             "Order station updated",
		"sender_id":           userID,
		"merchant_id":         merchantID,
		"merchant_station_id": stationID,
		"merchant_order_id":   stationOrder.OrderID,
		"merchant_desk_id":    stationOrder.MerchantDeskID,
		"status":              input.Status,
	}
	b, _ := json.Marshal(msg)
	p.ctx.AppService.(*services.AppService).Websocket.BroadcastFilter(b, func(q *melody.Session) bool {
		url := fmt.Sprintf("/api/v1/ws/%s", c.GetHeader("ID-Company"))
		// fmt.Println("ORDER_STATION_CREATED", url, q.Request.URL.Path)
		return q.Request.URL.Path == url
	})
	c.JSON(200, gin.H{"message": "Station order updated successfully"})
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
	userID := c.MustGet("userID").(string)

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

		// utils.LogJson(orderStationMap)

		for stationID, v := range orderStationMap {
			msg := gin.H{
				"command":             "ORDER_STATION_CREATED",
				"message":             "order station created",
				"sender_id":           userID,
				"merchant_id":         merchantID,
				"merchant_station_id": stationID,
				"merchant_order_id":   input.ID,
				"merchant_desk_id":    input.MerchantDeskID,
				"items":               v,
			}
			b, _ := json.Marshal(msg)
			p.ctx.AppService.(*services.AppService).Websocket.BroadcastFilter(b, func(q *melody.Session) bool {
				url := fmt.Sprintf("/api/v1/ws/%s", c.GetHeader("ID-Company"))
				// fmt.Println("ORDER_STATION_CREATED", url, q.Request.URL.Path)
				return q.Request.URL.Path == url
			})
		}

	}
	c.JSON(200, gin.H{"message": "Merchant order created successfully"})
}
