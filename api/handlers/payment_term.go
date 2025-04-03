package handlers

import (
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/order"
	"github.com/gin-gonic/gin"
)

type PaymentTermHandler struct {
	ctx      *context.ERPContext
	orderSrv *order.OrderService
}

func NewPaymentTermHandler(ctx *context.ERPContext) *PaymentTermHandler {
	orderSrv, ok := ctx.OrderService.(*order.OrderService)
	if !ok {
		panic("order service is not found")
	}
	return &PaymentTermHandler{
		ctx:      ctx,
		orderSrv: orderSrv,
	}
}

func (p *PaymentTermHandler) GetPaymentTermsHandler(c *gin.Context) {
	data := p.orderSrv.PaymentTermService.GetPaymentTerms()
	c.JSON(200, gin.H{"message": "success", "data": data})

}
func (p *PaymentTermHandler) GetPaymentTermsGroupHandler(c *gin.Context) {
	data := p.orderSrv.PaymentTermService.GroupPaymentTermsByCategory()
	c.JSON(200, gin.H{"message": "success", "data": data})

}
