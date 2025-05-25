package pos

import (
	"ametory-cooperative/services"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/order"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/AMETORY/ametory-erp-modules/thirdparty/payment/xendit"
	"github.com/AMETORY/ametory-erp-modules/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/olahol/melody.v1"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PosHandler struct {
	ctx           *context.ERPContext
	financeSrv    *finance.FinanceService
	OrderSrv      *order.OrderService
	appService    *services.AppService
	xenditService *xendit.XenditService
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

	xenditService := xendit.NewXenditService()
	return &PosHandler{
		ctx:           ctx,
		financeSrv:    financeSrv,
		OrderSrv:      orderSrv,
		appService:    appService,
		xenditService: xenditService,
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
		Status       string `json:"status"`
		ContactName  string `json:"contact_name"`
		ContactPhone string `json:"contact_phone"`
		ContactID    string `json:"contact_id"`
	}{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// utils.LogJson(input)
	err := p.OrderSrv.MerchantService.UpdateTableStatus(id, tableId, input.Status)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if input.ContactName != "" || input.ContactPhone != "" || input.ContactID != "" {
		err = p.OrderSrv.MerchantService.UpdateTableContact(id, tableId, input.ContactName, input.ContactPhone, input.ContactID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
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
	userID := c.MustGet("userID").(string)
	input.CashierID = &userID
	companyID := c.GetHeader("ID-Company")

	if (input.ContactName != "" || input.ContactPhone != "") && input.ContactID == nil {
		var contact models.ContactModel
		phoneNumber := ""
		if input.ContactPhone != "" {
			phoneNumber = utils.ParsePhoneNumber(input.ContactPhone, "ID")
		}
		err := p.ctx.DB.Where("company_id = ?  AND phone = ?", companyID, input.ContactName, phoneNumber).First(&contact).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			contact.CompanyID = &companyID
			contact.Name = input.ContactName
			contact.Phone = &phoneNumber
			contact.IsCustomer = true
			err := p.ctx.DB.Create(&contact).Error
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			input.ContactID = &contact.ID

			b, _ := json.Marshal(contact)
			input.ContactData = b
		}

	} else if input.ContactID != nil {
		var contact models.ContactModel
		err := p.ctx.DB.First(&contact, "id = ?", *input.ContactID).Error
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		b, _ := json.Marshal(contact)
		input.ContactData = b
	}

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

	msg := gin.H{
		"command":           "ORDER_CREATED",
		"message":           "order  created",
		"sender_id":         userID,
		"merchant_id":       merchantID,
		"merchant_order_id": input.ID,
		"merchant_desk_id":  input.MerchantDeskID,
	}
	b, _ := json.Marshal(msg)
	p.ctx.AppService.(*services.AppService).Websocket.BroadcastFilter(b, func(q *melody.Session) bool {
		url := fmt.Sprintf("/api/v1/ws/%s", c.GetHeader("ID-Company"))
		// fmt.Println("ORDER_STATION_CREATED", url, q.Request.URL.Path)
		return q.Request.URL.Path == url
	})
	c.JSON(200, gin.H{"message": "Merchant order created successfully"})
}

func (p *PosHandler) PaymentCheckHandler(c *gin.Context) {
	merchantID := c.MustGet("merchantID").(string)
	orderID := c.Param("orderId")

	merchant, err := p.OrderSrv.MerchantService.GetMerchantByID(merchantID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	order, err := p.OrderSrv.MerchantService.GetOrderDetail(merchantID, orderID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("userID").(string)
	msg := gin.H{
		"command":     "PAYMENT_UPDATED",
		"message":     "Order payment updated",
		"sender_id":   userID,
		"merchant_id": merchantID,
		"order_id":    orderID,
	}
	msgByte, _ := json.Marshal(msg)
	for _, v := range order.Payments {
		if v.PaymentProvider == "QRIS" {
			p.xenditService.SetAPIKey(merchant.XenditApiKey)

			payments, err := p.xenditService.GetQRPayments(v.ExternalID)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			if len(payments) > 0 {
				order.OrderStatus = "PAID"
				p.ctx.DB.Save(&order)

				if order.ParentID != nil {
					err = p.ctx.DB.Model(&models.MerchantDesk{}).Where("merchant_id = ? AND id = ?", merchantID, order.MerchantDeskID).
						Updates(map[string]any{
							"contact_name":  "",
							"contact_phone": "",
							"contact_id":    nil,
							"status":        "AVAILABLE",
						}).Error
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
						return
					}
				}

				p.ctx.AppService.(*services.AppService).Websocket.BroadcastFilter(msgByte, func(q *melody.Session) bool {
					url := fmt.Sprintf("/api/v1/ws/%s", c.GetHeader("ID-Company"))
					// fmt.Println("ORDER_STATION_CREATED", url, q.Request.URL.Path)
					return q.Request.URL.Path == url
				})
				b, err := json.Marshal(payments[0])
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				v.PaymentData = b
				err = p.ctx.DB.Save(&v).Error
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				c.JSON(200, gin.H{"message": "Payment check completed"})
				return
			}

			resp, err := p.xenditService.GetQRByID(v.ExternalID)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			if resp.Status == "INACTIVE" {
				order.OrderStatus = "PAYMENT_FAILED"
				p.ctx.DB.Omit(clause.Associations).Save(&order)

				msg := gin.H{
					"command":     "PAYMENT_FAILED",
					"message":     "Order payment failed",
					"sender_id":   userID,
					"merchant_id": merchantID,
					"order_id":    orderID,
				}
				msgByte, _ := json.Marshal(msg)

				p.ctx.AppService.(*services.AppService).Websocket.BroadcastFilter(msgByte, func(q *melody.Session) bool {
					url := fmt.Sprintf("/api/v1/ws/%s", c.GetHeader("ID-Company"))
					// fmt.Println("ORDER_STATION_CREATED", url, q.Request.URL.Path)
					return q.Request.URL.Path == url
				})

				b, _ := json.Marshal(resp)
				v.PaymentData = b
				err = p.ctx.DB.Save(&v).Error
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

			}
		}
	}

	c.JSON(200, gin.H{"message": "Payment check completed"})
}
func (p *PosHandler) SplitBillHandler(c *gin.Context) {
	merchantID := c.Param("id")
	orderID := c.Param("orderId")
	companyID := c.GetHeader("ID-Company")

	input := struct {
		ContactName  string                     `json:"contact_name"`
		ContactPhone string                     `json:"contact_phone"`
		ContactID    *string                    `json:"contact_id"`
		Items        []models.MerchantOrderItem `json:"items"`
	}{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	order, err := p.OrderSrv.MerchantService.GetOrderDetail(merchantID, orderID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var contact *models.ContactModel
	if (input.ContactName != "" || input.ContactPhone != "") && input.ContactID == nil {
		phoneNumber := ""
		if input.ContactPhone != "" {
			phoneNumber = utils.ParsePhoneNumber(input.ContactPhone, "ID")
		}
		err := p.ctx.DB.Where("company_id = ?  AND phone = ?", companyID, input.ContactName, phoneNumber).First(&contact).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			contact.CompanyID = &companyID
			contact.Name = input.ContactName
			contact.Phone = &phoneNumber
			contact.IsCustomer = true
			err := p.ctx.DB.Create(contact).Error
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
		}

	} else if input.ContactID != nil {
		var existingContact models.ContactModel
		err := p.ctx.DB.First(&existingContact, "id = ?", *input.ContactID).Error
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		contact = &existingContact
	}
	// utils.LogJson(input)
	newOrder, err := p.OrderSrv.MerchantService.SplitBill(order, contact, input.Items)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(string)

	msg := gin.H{
		"command":     "SPLIT_BILL",
		"message":     "Split bill successfully",
		"sender_id":   userID,
		"merchant_id": merchantID,
		"order_id":    orderID,
	}
	b, _ := json.Marshal(msg)
	p.ctx.AppService.(*services.AppService).Websocket.BroadcastFilter(b, func(q *melody.Session) bool {
		url := fmt.Sprintf("/api/v1/ws/%s", c.GetHeader("ID-Company"))
		// fmt.Println("ORDER_STATION_CREATED", url, q.Request.URL.Path)
		return q.Request.URL.Path == url
	})

	c.JSON(200, gin.H{"data": newOrder, "message": "Order detail retrieved successfully"})

}
func (p *PosHandler) RepaymentOrderHandler(c *gin.Context) {
	merchantID := c.MustGet("merchantID").(string)
	orderID := c.Param("orderId")

	merchant, err := p.OrderSrv.MerchantService.GetMerchantByID(merchantID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	order, err := p.OrderSrv.MerchantService.GetOrderDetail(merchant.ID, orderID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	for _, v := range order.Payments {
		if v.PaymentProvider == "QRIS" {
			p.xenditService.SetAPIKey(merchant.XenditApiKey)
			resp, err := p.xenditService.CreateQR(xendit.XenditQRrequest{
				ReferenceID: v.OrderID,
				Amount:      v.Amount,
				Currency:    "IDR",
				Type:        "DYNAMIC",
				ExpiresAt:   time.Now().Add(time.Minute * 5).Format(time.RFC3339),
			})
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			fmt.Println("REPAYMENT QRIS DATA")
			utils.LogJson(resp)
			b, err := json.Marshal(resp)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			v.PaymentData = b
			v.ExternalID = resp.ID
			v.ExternalProvider = "xendit"
			err = p.ctx.DB.Save(&v).Error
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			order.OrderStatus = "PENDING_PAYMENT"
			err = p.ctx.DB.Save(&order).Error
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}
	}

	c.JSON(200, gin.H{"message": "Order repayment successfully"})
}
func (p *PosHandler) PaymentOrderHandler(c *gin.Context) {
	merchantID := c.MustGet("merchantID").(string)
	orderID := c.Param("orderId")

	var input []models.MerchantPayment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	merchant, err := p.OrderSrv.MerchantService.GetMerchantByID(merchantID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	order, err := p.OrderSrv.MerchantService.GetOrderDetail(merchantID, orderID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(input) == 0 {
		c.JSON(400, gin.H{"error": "payment data is required"})
		return
	}
	if len(input) == 1 {
		order.OrderStatus = "PAID"
		if input[0].PaymentProvider == "QRIS" && (!merchant.EnableXendit || !merchant.Xendit.EnableQRIS) {
			c.JSON(400, gin.H{"error": "QRIS is not enabled"})
			return
		}

		if input[0].PaymentProvider == "QRIS" {
			order.OrderStatus = "PENDING_PAYMENT"
		}
	}
	if len(input) > 2 {
		c.JSON(400, gin.H{"error": "only 2 payment method is allowed"})
		return
	}

	if len(input) > 1 {
		order.OrderStatus = "PENDING_PAYMENT"
		if input[0].PaymentProvider != "CASH" {
			c.JSON(400, gin.H{"error": "first payment method only cash is allowed"})
			return
		}
		if input[0].Amount+input[1].Amount != order.Total {
			c.JSON(400, gin.H{"error": "invalid payment amount"})
			return
		}

		if input[1].PaymentProvider == "QRIS" && (!merchant.EnableXendit || !merchant.Xendit.EnableQRIS) {
			c.JSON(400, gin.H{"error": "QRIS is not enabled"})
			return
		}
		// if input[0].PaymentMethod == "GOPAY" && (!merchant.EnableXendit || !merchant.Xendit.EnableGOPAY) {
		// 	c.JSON(400, gin.H{"error": "GOPAY is not enabled"})
		// 	return
		// }
		// if input[0].PaymentMethod == "DANA" && (!merchant.EnableXendit || !merchant.Xendit.EnableDANA) {
		// 	c.JSON(400, gin.H{"error": "DANA is not enabled"})
		// 	return
		// }
		// if input[0].PaymentMethod == "OVO" && (!merchant.EnableXendit || !merchant.Xendit.EnableOVO) {
		// 	c.JSON(400, gin.H{"error": "OVO is not enabled"})
		// 	return
		// }
		// if input[0].PaymentMethod == "BCA" && (!merchant.EnableXendit || !merchant.Xendit.EnableBCA) {
		// 	c.JSON(400, gin.H{"error": "BCA is not enabled"})
		// 	return
		// }
		// if input[0].PaymentMethod == "MANDIRI" && (!merchant.EnableXendit || !merchant.Xendit.EnableMANDIRI) {
		// 	c.JSON(400, gin.H{"error": "MANDIRI is not enabled"})
		// 	return
		// }
		// if input[0].PaymentMethod == "BNI" && (!merchant.EnableXendit || !merchant.Xendit.EnableBNI) {
		// 	c.JSON(400, gin.H{"error": "BNI is not enabled"})
		// 	return
		// }
		// if input[0].PaymentMethod == "BRI" && (!merchant.EnableXendit || !merchant.Xendit.EnableBRI) {
		// 	c.JSON(400, gin.H{"error": "BRI is not enabled"})
		// 	return
		// }
		if input[1].PaymentProvider == "EDC" {
			order.OrderStatus = "PAID"
		}
	}

	err = p.ctx.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Save(order).Error
		if err != nil {
			return err
		}

		if order.OrderStatus == "PAID" && order.ParentID == nil {
			err = tx.Model(&models.MerchantDesk{}).Where("merchant_id = ? AND id = ?", merchantID, order.MerchantDeskID).
				Updates(map[string]any{
					"contact_name":  "",
					"contact_phone": "",
					"contact_id":    nil,
					"status":        "AVAILABLE",
				}).Error
			if err != nil {
				return err
			}
		}

		for _, v := range input {
			v.ID = uuid.NewString()
			v.Order = order
			v.OrderID = orderID
			v.Date = time.Now()
			err := tx.Create(&v).Error
			if err != nil {
				return err
			}

			if v.PaymentProvider == "QRIS" && (merchant.EnableXendit && merchant.Xendit.EnableQRIS) {
				p.xenditService.SetAPIKey(merchant.XenditApiKey)
				resp, err := p.xenditService.CreateQR(xendit.XenditQRrequest{
					ReferenceID: v.OrderID,
					Amount:      v.Amount,
					Currency:    "IDR",
					Type:        "DYNAMIC",
					ExpiresAt:   time.Now().Add(time.Minute * 5).Format(time.RFC3339),
				})
				if err != nil {
					return err
				}
				utils.LogJson(resp)
				b, err := json.Marshal(resp)
				if err != nil {
					return err
				}
				v.PaymentData = b
				v.ExternalID = resp.ID
				v.ExternalProvider = "xendit"
				err = tx.Save(&v).Error
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("userID").(string)

	msg := gin.H{
		"command":     "PAYMENT_CREATED",
		"message":     "Order paid successfully",
		"sender_id":   userID,
		"merchant_id": merchantID,
		"order_id":    orderID,
	}
	b, _ := json.Marshal(msg)
	p.ctx.AppService.(*services.AppService).Websocket.BroadcastFilter(b, func(q *melody.Session) bool {
		url := fmt.Sprintf("/api/v1/ws/%s", c.GetHeader("ID-Company"))
		// fmt.Println("ORDER_STATION_CREATED", url, q.Request.URL.Path)
		return q.Request.URL.Path == url
	})

	c.JSON(200, gin.H{"message": "Order paid successfully"})
}

func (s *PosHandler) DownloadOrderDetailPdfHandler(c *gin.Context) {
	orderId := c.Param("orderId")
	merchantId := c.Param("id")

	order, err := s.OrderSrv.MerchantService.GetOrderDetail(merchantId, orderId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	pdf, err := s.OrderSrv.MerchantService.GetPrintReceipt(order, "templates/pdf/receipt.html", "")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "application/pdf", pdf)
	c.Writer.Header().Add("Content-Disposition", "attachment; filename="+order.Code+".pdf")
}
