package handlers

import (
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/inventory"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type UnitHandler struct {
	ctx          *context.ERPContext
	inventorySrv *inventory.InventoryService
}

func NewUnitHandler(ctx *context.ERPContext) *UnitHandler {
	inventorySrv, ok := ctx.InventoryService.(*inventory.InventoryService)
	if !ok {
		panic("inventory service is not found")
	}
	return &UnitHandler{
		ctx:          ctx,
		inventorySrv: inventorySrv,
	}
}

func (p *UnitHandler) GetUnitHandler(c *gin.Context) {
	id := c.Param("id")
	unit, err := p.inventorySrv.UnitService.GetUnitByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": unit, "message": "Unit retrieved successfully"})
}

func (p *UnitHandler) ListUnitsHandler(c *gin.Context) {
	units, err := p.inventorySrv.UnitService.GetUnits(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": units, "message": "Units retrieved successfully"})
}

func (p *UnitHandler) CreateUnitHandler(c *gin.Context) {
	var input models.UnitModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	input.CompanyID = &companyID
	err = p.inventorySrv.UnitService.CreateUnit(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Unit created successfully", "data": input})
}

func (p *UnitHandler) UpdateUnitHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.UnitModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = p.inventorySrv.UnitService.UpdateUnit(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Unit updated successfully"})
}

func (p *UnitHandler) DeleteUnitHandler(c *gin.Context) {
	id := c.Param("id")
	unit, err := p.inventorySrv.UnitService.GetUnitByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if unit.CompanyID == nil {
		c.JSON(403, gin.H{"error": "You do not have permission to delete this unit"})
	}
	err = p.inventorySrv.UnitService.DeleteUnit(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Unit deleted successfully"})
}
