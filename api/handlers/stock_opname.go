package handlers

import (
	"ametory-cooperative/objects"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/inventory"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type StockOpnameHandler struct {
	ctx          *context.ERPContext
	inventorySrv *inventory.InventoryService
}

func NewStockOpnameHandler(ctx *context.ERPContext) *StockOpnameHandler {
	inventorySrv, ok := ctx.InventoryService.(*inventory.InventoryService)
	if !ok {
		panic("inventory service is not found")
	}
	return &StockOpnameHandler{
		ctx:          ctx,
		inventorySrv: inventorySrv,
	}
}

func (p *StockOpnameHandler) GetStockOpnameHandler(c *gin.Context) {
	id := c.Param("id")
	stockOpname, err := p.inventorySrv.StockOpnameService.GetStockOpnameByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": stockOpname, "message": "StockOpname retrieved successfully"})
}

func (p *StockOpnameHandler) ListStockOpnamesHandler(c *gin.Context) {
	stockOpnames, err := p.inventorySrv.StockOpnameService.GetStockOpnames(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": stockOpnames, "message": "StockOpnames retrieved successfully"})
}

func (p *StockOpnameHandler) CreateStockOpnameHandler(c *gin.Context) {
	var input models.StockOpnameHeader
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("userID").(string)
	companyID := c.MustGet("companyID").(string)
	input.CompanyID = &companyID
	input.CreatedByID = &userID

	err = p.inventorySrv.StockOpnameService.CreateStockOpnameFromHeader(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "StockOpname created successfully", "data": input})
}

func (p *StockOpnameHandler) UpdateStockOpnameHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.StockOpnameHeader
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.inventorySrv.StockOpnameService.UpdateStockOpname(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "StockOpname updated successfully"})
}

func (p *StockOpnameHandler) DeleteStockOpnameHandler(c *gin.Context) {
	id := c.Param("id")
	stockOpname, err := p.inventorySrv.StockOpnameService.GetStockOpnameByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// if stockOpname.CompanyID == nil {
	// 	c.JSON(403, gin.H{"error": "You do not have permission to delete this stockOpname"})
	// }
	err = p.inventorySrv.StockOpnameService.DeleteStockOpname(stockOpname.ID, false)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "StockOpname deleted successfully"})
}

func (p *StockOpnameHandler) CompleteStockOpnameHandler(c *gin.Context) {
	id := c.Param("id")
	input := objects.StockOpnameRequest{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(string)
	err = p.inventorySrv.StockOpnameService.CompleteStockOpname(id, input.Date, userID, input.InventoryID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "StockOpname completed successfully"})
}

func (p *StockOpnameHandler) AddItemHandler(c *gin.Context) {
	id := c.Param("id")
	input := models.StockOpnameDetail{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.inventorySrv.StockOpnameService.AddItem(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Item added to stock opname successfully"})
}

func (p *StockOpnameHandler) UpdateItemHandler(c *gin.Context) {
	id := c.Param("id")
	_, err := p.inventorySrv.StockOpnameService.GetStockOpnameByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	detailID := c.Param("detailId")
	input := models.StockOpnameDetail{}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.inventorySrv.StockOpnameService.UpdateItem(detailID, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Item updated in stock opname successfully"})
}

func (p *StockOpnameHandler) DeleteItemHandler(c *gin.Context) {
	id := c.Param("id")
	_, err := p.inventorySrv.StockOpnameService.GetStockOpnameByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	detailID := c.Param("detailId")
	err = p.inventorySrv.StockOpnameService.DeleteItem(detailID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Item deleted from stock opname successfully"})
}
