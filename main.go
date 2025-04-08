package main

import (
	"ametory-cooperative/api/routes"
	"ametory-cooperative/config"
	"ametory-cooperative/services"
	"ametory-cooperative/workers"
	"net/mail"
	"os"

	ctx "context"

	"github.com/AMETORY/ametory-erp-modules/auth"
	"github.com/AMETORY/ametory-erp-modules/company"
	"github.com/AMETORY/ametory-erp-modules/contact"
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/cooperative"
	"github.com/AMETORY/ametory-erp-modules/distribution"
	"github.com/AMETORY/ametory-erp-modules/file"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/inventory"
	"github.com/AMETORY/ametory-erp-modules/order"
	"github.com/AMETORY/ametory-erp-modules/shared/audit_trail"
	"github.com/AMETORY/ametory-erp-modules/shared/indonesia_regional"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/AMETORY/ametory-erp-modules/thirdparty"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := ctx.Background()
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	db, err := services.InitDB(cfg)
	if err != nil {
		panic(err)
	}
	redisClient := services.InitRedis()
	websocket := services.InitWS()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3035",
			"http://localhost:3035/",
		},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "DELETE", "HEAD"},
		AllowHeaders: []string{
			"Origin",
			"Authorization",
			"Content-Length",
			"Content-Type",
			"Access-Control-Allow-Origin",
			"API-KEY",
			"Currency-Code",
			"Cache-Control",
			"X-Requested-With",
			"Content-Disposition",
			"Content-Description",
			"ID-Company",
			"ID-Distributor",
			"timezone",
		},
		ExposeHeaders: []string{"Content-Length", "Content-Disposition", "Content-Description"},
	}))

	skipMigration := true

	if os.Getenv("MIGRATION") != "" {
		skipMigration = false
	}

	erpContext := context.NewERPContext(db, nil, &ctx, skipMigration)
	authService := auth.NewAuthService(erpContext)
	erpContext.AuthService = authService

	fileService := file.NewFileService(erpContext, cfg.Server.BaseURL)
	erpContext.FileService = fileService

	companyService := company.NewCompanyService(erpContext)
	erpContext.CompanyService = companyService

	rbacSrv := auth.NewRBACService(erpContext)
	erpContext.RBACService = rbacSrv

	financeSrv := finance.NewFinanceService(erpContext)
	erpContext.FinanceService = financeSrv

	cooperativeSrv := cooperative.NewCooperativeService(erpContext, companyService, financeSrv)
	erpContext.CooperativeService = cooperativeSrv

	orderSrv := order.NewOrderService(erpContext)
	erpContext.OrderService = orderSrv

	inventorySrv := inventory.NewInventoryService(erpContext)
	erpContext.InventoryService = inventorySrv

	auditTrailSrv := audit_trail.NewAuditTrailService(erpContext)

	distributionSrv := distribution.NewDistributionService(erpContext, auditTrailSrv, inventorySrv, orderSrv)
	erpContext.DistributionService = distributionSrv

	orderSrv = order.NewOrderService(erpContext)
	erpContext.OrderService = orderSrv

	inventorySrv = inventory.NewInventoryService(erpContext)
	erpContext.InventoryService = inventorySrv

	contactSrv := contact.NewContactService(erpContext, companyService)
	erpContext.ContactService = contactSrv

	appService := services.NewAppService(erpContext, cfg, redisClient, websocket)
	erpContext.AppService = appService

	indonesiaRegSrv := indonesia_regional.NewIndonesiaRegService(erpContext)
	erpContext.IndonesiaRegService = indonesiaRegSrv

	emailSender := thirdparty.NewSMTPSender(cfg.Email.Server, cfg.Email.Port, cfg.Email.Username, cfg.Email.Password, mail.Address{Name: cfg.Email.From, Address: cfg.Email.From})
	emailSender.SetTemplate("./templates/email/layout.html", "./templates/email/body.html")

	erpContext.EmailSender = emailSender

	v1 := r.Group("/api/v1")

	r.Static("/assets/files", "./assets/files")
	routes.SetupWSRoutes(v1, erpContext)
	routes.SetupAuthRoutes(v1, erpContext)
	routes.SetupCompanyRoutes(v1, erpContext)
	routes.SetupAccountRoutes(v1, erpContext)
	routes.SetupRegionalRoutes(v1, erpContext)
	routes.SetupCommonRoutes(v1, erpContext)
	routes.SetupTransactionRoutes(v1, erpContext)
	routes.SetupJournalRoutes(v1, erpContext)
	routes.SetupTaxRoutes(v1, erpContext)
	routes.SetContactRoutes(v1, erpContext)
	routes.SetupSalesRoutes(v1, erpContext)
	routes.SetupProductRoutes(v1, erpContext)
	routes.SetupProductCategoryRoutes(v1, erpContext)
	routes.SetupPriceCategoryRoutes(v1, erpContext)
	routes.SetupProductAttributeRoutes(v1, erpContext)
	routes.SetupWarehouseRoutes(v1, erpContext)
	routes.SetupUnitRoutes(v1, erpContext)
	routes.SetupPaymentTermRoutes(v1, erpContext)
	routes.SetupStockMovementRoutes(v1, erpContext)
	routes.SetupPurchaseRoutes(v1, erpContext)
	routes.SetupReportRoutes(v1, erpContext)
	routes.SetupPurchaseReturnRoutes(v1, erpContext)
	routes.SetupSalesReturnRoutes(v1, erpContext)
	routes.SetupCooperativeRoutes(v1, erpContext)

	go func() {
		workers.SendMail(erpContext)
	}()

	if os.Getenv("GEN_PERMISSIONS") != "" {
		for _, v := range appService.GenerateDefaultPermissions() {
			erpContext.DB.Create(&v)
		}

		var companies []models.CompanyModel
		erpContext.DB.Find(&companies)
		for _, company := range companies {
			roles := appService.GenerateDefaultRoles(company.ID)
			for _, v := range roles {
				err := erpContext.DB.Where("name = ?", v.Name).First(&v).Error
				if err == nil {
					erpContext.DB.Model(&v).Association("Permissions").Append(&v.Permissions)
				}
			}

			// for _, v := range services.GenerateCustomAccounts {
			// 	v.CompanyID = &company.ID
			// 	erpContext.DB.Create(&v)
			// }
		}
	}
	if os.Getenv("DEFAULT_CATEGORY") != "" {
		appService.GenerateDefaultCategories()

	}

	if os.Getenv("PAYMENT_TERMS") != "" {
		err := orderSrv.PaymentTermService.InitPaymentTerms()
		if err != nil {
			panic(err) // TODO: handle this properly
		}

	}

	r.Run(":" + config.App.Server.Port)
}
