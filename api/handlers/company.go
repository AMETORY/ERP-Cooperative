package handlers

import (
	"ametory-cooperative/objects"
	"ametory-cooperative/services"

	"github.com/AMETORY/ametory-erp-modules/company"
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/AMETORY/ametory-erp-modules/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompanyHandler struct {
	ctx        *context.ERPContext
	companySrv *company.CompanyService
	appSrv     *services.AppService
}

func NewCompanyHandler(ctx *context.ERPContext) *CompanyHandler {
	companySrv, ok := ctx.CompanyService.(*company.CompanyService)
	if !ok {
		panic("CompanyService is not instance of company.CompanyService")
	}
	appSrv, ok := ctx.AppService.(*services.AppService)
	if !ok {
		panic("AppService is not instance of services.AppService")
	}
	return &CompanyHandler{
		ctx:        ctx,
		companySrv: companySrv,
		appSrv:     appSrv,
	}
}

func (c *CompanyHandler) GetCategories(ctx *gin.Context) {
	var sectorID *string
	if ctx.Query("sector_id") != "" {
		secId := ctx.Query("sector_id")
		sectorID = &secId
	}
	categories, err := c.companySrv.GetCategories(sectorID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"data": categories})
}

func (c *CompanyHandler) GetSectors(ctx *gin.Context) {
	sectors, err := c.companySrv.GetSectors()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"data": sectors})
}

func (c *CompanyHandler) CreateCompany(ctx *gin.Context) {
	var input objects.CompanyRequest
	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user := ctx.MustGet("user").(models.UserModel)

	company := models.CompanyModel{
		Name:              input.Name,
		Address:           input.Address,
		Email:             input.Email,
		Phone:             input.Phone,
		CompanyCategoryID: &input.CompanyCategoryID,
		ProvinceID:        &input.ProvinceID,
		RegencyID:         &input.RegencyID,
		DistrictID:        &input.DistrictID,
		VillageID:         &input.VillageID,
		ZipCode:           &input.ZipCode,
		IsCooperation:     input.IsCooperation,
	}
	err := c.ctx.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&company).Error; err != nil {
			return err
		}

		if input.IsCooperation {
			setting := models.CooperativeSettingModel{
				CompanyID:                  &company.ID,
				NetSurplusReserve:          10,
				NetSurplusMandatorySavings: 25,
				NetSurplusBusinessProfit:   25,
				NetSurplusSocialFund:       10,
				NetSurplusEducationFund:    10,
				NetSurplusManagement:       10,
				NetSurplusOtherFunds:       10,
				StaticCharacter:            "LOAN",
				NumberFormat:               "{static-character}-{auto-numeric}/{month-roman}/{year-yyyy}",
				AutoNumericLength:          5,
				RandomNumericLength:        3,
				RandomCharacterLength:      3,
				InterestRatePerMonth:       2,
				ExpectedProfitRatePerMonth: 3,
				IsIslamic:                  input.IsIslamic,
			}
			setting.ID = utils.Uuid()
			if err := tx.Create(&setting).Error; err != nil {
				return err
			}

		}

		for _, v := range input.Accounts {
			v.CompanyID = &company.ID
			v.UserID = &user.ID
			if err := tx.Create(&v).Error; err != nil {
				return err
			}
		}

		superAdmin := models.RoleModel{}
		roles := c.appSrv.GenerateDefaultRoles(company.ID)
		for _, v := range roles {
			if err := tx.Create(&v).Error; err != nil {
				return err
			}

			if v.IsSuperAdmin {
				superAdmin = v
			}
		}

		err := tx.Model(&user).Association("Roles").Append(&superAdmin)
		if err != nil {
			return err
		}
		err = tx.Model(&user).Association("Companies").Append(&company)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.ctx.DB.Preload("Companies").Find(&user)
	ctx.JSON(200, gin.H{"data": company, "companies": user.Companies})
}
