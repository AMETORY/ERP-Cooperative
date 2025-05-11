package services

import (
	"ametory-cooperative/app_models"
	"ametory-cooperative/config"
	"encoding/json"
	"errors"
	"log"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/shared"
	"github.com/AMETORY/ametory-erp-modules/shared/constants"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/AMETORY/ametory-erp-modules/utils"
	"github.com/go-redis/redis/v8"
	"gopkg.in/olahol/melody.v1"
	"gorm.io/gorm"
)

type AppService struct {
	ctx       *context.ERPContext
	Config    *config.Config
	Redis     *redis.Client
	Websocket *melody.Melody
	App       *app_models.AppModel
}

func NewAppService(ctx *context.ERPContext, config *config.Config, redis *redis.Client, ws *melody.Melody) *AppService {
	if !ctx.SkipMigration {
		ctx.DB.AutoMigrate(
			&app_models.AppModel{},
			&app_models.CustomSettingModel{},
		)
	}
	var app *app_models.AppModel
	ctx.DB.First(&app)
	return &AppService{
		ctx:       ctx,
		Config:    config,
		Redis:     redis,
		Websocket: ws,
		App:       app,
	}
}

func (a AppService) GenerateDefaultPermissions() []models.PermissionModel {
	var (
		cruds    = []string{"create", "read", "update", "delete"}
		services = map[string][]map[string][]string{
			"auth": {{"user": cruds, "admin": cruds, "rbac": cruds}},
			"contact": {
				{"customer": cruds},
				{"vendor": cruds},
				{"supplier": cruds},
				{"all": cruds},
			},
			"company": {
				{"company": append(cruds, "approval")},
				{"user": []string{"invite"}},
			},
			"order": {
				{"sales": append(cruds, "approval", "publish")},
				{"merchant": append(cruds, "approval", "publish")},
			},
			"menu": {
				{"admin": []string{
					"dashboard",
					"preferences",
					"feature",
					"inventory",
				}},
			},
			"inventory": {
				{"purchase": cruds},
				{"purchase_return": cruds},
				{"product": cruds},
				{"product_category": cruds},
				{"price_category": cruds},
				{"product_attribute": cruds},
				{"warehouse": cruds},
				{"unit": cruds},
			},
			"finance": {
				{"account": cruds},
				{"transaction": cruds},
				{"journal": cruds},
				{"report": cruds},
				{"bank": cruds},
				{"tax": cruds},
				{"asset": append(cruds, "apply", "activate", "deactivate")},
				{"report": []string{
					"menu",
					"cash_flow",
					"balance_sheet",
					"income_statement",
					"general_ledger",
					"trial_balance",
					"inventory",
					"cogs",
					"profit_loss",
					"purchase_report",
					"sales_report",
					"inventory_journal",
					"general_journal",
					"reconcile_account_payable",
					"reconcile_account_receivable",
					"sales_per_customer",
					"purchase_per_supplier",
				}},
			},

			"cooperative": {
				{"cooperative_member": append(cruds, "approval", "invite")},
				{"cooperative_setting": cruds},
				{"loan_application": append(cruds, "request", "approve", "reject")},
				{"saving": cruds},
				{"net_surplus": cruds},
			},
		}
	)

	return a.generatePermissions(services)
}

func (a AppService) GenerateAdminPermissions() []models.PermissionModel {
	var (
		cruds    = []string{"create", "read", "update", "delete"}
		services = map[string][]map[string][]string{
			"auth": {{"user": cruds, "admin": cruds, "rbac": cruds}},
			"contact": {
				{"customer": cruds},
				{"vendor": cruds},
				{"supplier": cruds},
				{"all": cruds},
			},
			"company": {
				{"company": append(cruds, "approval")},
			},
			"order": {
				{"sales": append(cruds, "approval", "publish")},
			},
			"menu": {
				{"admin": []string{
					"dashboard",
					"preferences",
					"feature",
					"inventory",
				}},
			},
			"inventory": {
				{"purchase": cruds},
				{"purchase_return": cruds},
				{"product": cruds},
				{"product_category": cruds},
				{"price_category": cruds},
				{"product_attribute": cruds},
				{"warehouse": cruds},
				{"unit": cruds},
			},
			"finance": {
				{"account": cruds},
				{"asset": append(cruds, "apply", "activate", "deactivate")},
				{"transaction": cruds},
				{"journal": cruds},
				{"report": cruds},
				{"bank": cruds},
				{"tax": cruds},
				{"report": []string{
					"menu",
					"cash_flow",
					"balance_sheet",
					"income_statement",
					"general_ledger",
					"trial_balance",
					"inventory",
					"cogs",
					"profit_loss",
					"purchase_report",
					"sales_report",
					"inventory_journal",
					"general_journal",
					"reconcile_account_payable",
					"reconcile_account_receivable",
					"sales_per_customer",
					"purchase_per_supplier",
				}},
			},
			"cooperative": {
				{"cooperative_member": append(cruds, "approval", "invite")},
				{"cooperative_setting": cruds},
				{"loan_application": append(cruds, "request", "approve", "reject")},
				{"saving": cruds},
				{"net_surplus": cruds},
			},
		}
	)

	return a.generatePermissions(services)
}
func (a AppService) GenerateMemberPermissions() []models.PermissionModel {
	var (
		services = map[string][]map[string][]string{
			"cooperative": {
				{"loan_application": []string{"my", "request"}},
				{"saving": []string{"my", "request"}},
			},
		}
	)

	return a.generatePermissions(services)
}

func (a AppService) GenerateDefaultRoles(companyID string) []models.RoleModel {
	return []models.RoleModel{
		{
			Name:         "Super Admin",
			IsSuperAdmin: true,
			IsOwner:      true,
			CompanyID:    &companyID,
			Permissions:  []models.PermissionModel{},
		},
		{
			Name:        "Admin",
			Permissions: a.GenerateAdminPermissions(),
			CompanyID:   &companyID,
		},
		{
			Name:        "Member",
			Permissions: a.GenerateMemberPermissions(),
			CompanyID:   &companyID,
		},
	}

}
func (a AppService) generatePermissions(services map[string][]map[string][]string) []models.PermissionModel {

	var permissions []models.PermissionModel

	for service, modules := range services {
		for _, module := range modules {
			for key, actions := range module {
				for _, action := range actions {
					var permission models.PermissionModel
					err := a.ctx.DB.First(&permission, "name = ?", service+":"+key+":"+action).Error
					if errors.Is(err, gorm.ErrRecordNotFound) {
						permission.Name = service + ":" + key + ":" + action
						a.ctx.DB.Create(&permission)
					}
					permissions = append(permissions, permission)
				}
			}
		}
	}
	return permissions
}

// GenerateDefaultCategories creates default company sectors and categories.
func (a AppService) GenerateDefaultCategories() {
	var categories map[string]any
	err := json.Unmarshal([]byte(companyCategoriesStr), &categories)
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range categories["sme"].([]interface{}) {
		var sectorID = utils.Uuid()
		companyCats := []models.CompanyCategory{}
		for _, c := range v.(map[string]interface{})["category"].([]interface{}) {

			companyCats = append(companyCats, models.CompanyCategory{
				BaseModel: shared.BaseModel{
					ID: utils.Uuid(),
				},
				Name:          c.(string),
				IsCooperative: false,
				SectorID:      &sectorID,
			})
		}
		var sector = models.CompanySector{
			BaseModel: shared.BaseModel{
				ID: sectorID,
			},
			Name:       v.(map[string]interface{})["sector"].(string),
			Categories: companyCats,
		}

		a.ctx.DB.Create(&sector)
	}
	for _, v := range categories["cooperative"].([]interface{}) {
		var sectorID = utils.Uuid()
		companyCats := []models.CompanyCategory{}
		for _, c := range v.(map[string]interface{})["category"].([]interface{}) {

			companyCats = append(companyCats, models.CompanyCategory{
				BaseModel: shared.BaseModel{
					ID: utils.Uuid(),
				},
				Name:          c.(string),
				IsCooperative: true,
				SectorID:      &sectorID,
			})
		}
		var sector = models.CompanySector{
			BaseModel: shared.BaseModel{
				ID: sectorID,
			},
			Name:          v.(map[string]interface{})["sector"].(string),
			Categories:    companyCats,
			IsCooperative: true,
		}

		a.ctx.DB.Create(&sector)
	}

}

var GenerateCustomAccounts = []models.AccountModel{
	{
		CashflowSubGroup: constants.EQUITY_CAPITAL,
		CashflowGroup:    constants.CASHFLOW_GROUP_FINANCING,
		Category:         constants.CATEGORY_EQUITY,
		Name:             "Jasa Modal",
		Type:             models.EQUITY,
		Code:             "31006",
	},
	{
		CashflowSubGroup: constants.EQUITY_CAPITAL,
		CashflowGroup:    constants.CASHFLOW_GROUP_FINANCING,
		Category:         constants.CATEGORY_EQUITY,
		Name:             "Jasa Usaha",
		Type:             models.EQUITY,
		Code:             "31007",
	},
	{
		CashflowSubGroup: constants.EQUITY_CAPITAL,
		CashflowGroup:    constants.CASHFLOW_GROUP_FINANCING,
		Category:         constants.CATEGORY_EQUITY,
		Name:             "Dana Cadangan",
		Type:             models.EQUITY,
		Code:             "31008",
	},
	{
		CashflowSubGroup: constants.EQUITY_CAPITAL,
		CashflowGroup:    constants.CASHFLOW_GROUP_FINANCING,
		Category:         constants.CATEGORY_EQUITY,
		Name:             "Dana Sosial",
		Type:             models.EQUITY,
		Code:             "31009",
	},
	{
		CashflowSubGroup: constants.EQUITY_CAPITAL,
		CashflowGroup:    constants.CASHFLOW_GROUP_FINANCING,
		Category:         constants.CATEGORY_EQUITY,
		Name:             "Dana Pendidikan",
		Type:             models.EQUITY,
		Code:             "31010",
	},
	{
		CashflowSubGroup: constants.EQUITY_CAPITAL,
		CashflowGroup:    constants.CASHFLOW_GROUP_FINANCING,
		Category:         constants.CATEGORY_EQUITY,
		Name:             "Dana Pengurus",
		Type:             models.EQUITY,
		Code:             "31011",
	},
	{
		CashflowSubGroup: constants.EQUITY_CAPITAL,
		CashflowGroup:    constants.CASHFLOW_GROUP_FINANCING,
		Category:         constants.CATEGORY_EQUITY,
		Name:             "Dana Lainnya",
		Type:             models.EQUITY,
		Code:             "31012",
	},

	{
		CashflowSubGroup: constants.CASH_BANK,
		CashflowGroup:    constants.CASHFLOW_GROUP_CURRENT_ASSET,
		Category:         constants.CATEGORY_CURRENT_ASSET,
		Name:             "Kas Alokasi SHU Jasa Modal",
		Type:             models.ASSET,
		Code:             "11010",
	},
	{
		CashflowSubGroup: constants.CASH_BANK,
		CashflowGroup:    constants.CASHFLOW_GROUP_CURRENT_ASSET,
		Category:         constants.CATEGORY_CURRENT_ASSET,
		Name:             "Kas Alokasi SHU Jasa Usaha",
		Type:             models.ASSET,
		Code:             "11011",
	},
	{
		CashflowSubGroup: constants.CASH_BANK,
		CashflowGroup:    constants.CASHFLOW_GROUP_CURRENT_ASSET,
		Category:         constants.CATEGORY_CURRENT_ASSET,
		Name:             "Kas Alokasi SHU Dana Cadangan",
		Type:             models.ASSET,
		Code:             "11012",
	},
	{
		CashflowSubGroup: constants.CASH_BANK,
		CashflowGroup:    constants.CASHFLOW_GROUP_CURRENT_ASSET,
		Category:         constants.CATEGORY_CURRENT_ASSET,
		Name:             "Kas Alokasi SHU Dana Sosial",
		Type:             models.ASSET,
		Code:             "11013",
	},
	{
		CashflowSubGroup: constants.CASH_BANK,
		CashflowGroup:    constants.CASHFLOW_GROUP_CURRENT_ASSET,
		Category:         constants.CATEGORY_CURRENT_ASSET,
		Name:             "Kas Alokasi SHU Dana Pendidikan",
		Type:             models.ASSET,
		Code:             "11014",
	},
	{
		CashflowSubGroup: constants.CASH_BANK,
		CashflowGroup:    constants.CASHFLOW_GROUP_CURRENT_ASSET,
		Category:         constants.CATEGORY_CURRENT_ASSET,
		Name:             "Kas Alokasi SHU Dana Pengurus",
		Type:             models.ASSET,
		Code:             "11015",
	},
	{
		CashflowSubGroup: constants.CASH_BANK,
		CashflowGroup:    constants.CASHFLOW_GROUP_CURRENT_ASSET,
		Category:         constants.CATEGORY_CURRENT_ASSET,
		Name:             "Kas Alokasi SHU Dana Lainnya",
		Type:             models.ASSET,
		Code:             "11016",
	},
}

var companyCategoriesStr = `{
	"sme": [
		{
		  "sector": "Pertanian, Perkebunan, dan Perikanan",
		  "category": [
			"Usaha tani padi, jagung, kedelai",
			"Budidaya sayuran dan buah-buahan",
			"Perkebunan kelapa sawit, karet, kopi",
			"Peternakan ayam, sapi, kambing",
			"Perikanan tangkap dan budidaya",
			"Pengolahan hasil pertanian (minyak kelapa, tepung singkong)"
		  ]
		},
		{
		  "sector": "Makanan dan Minuman",
		  "category": [
			"Warung makan, katering, makanan ringan",
			"Produksi makanan kemasan (keripik, sambal, abon)",
			"Minuman tradisional (jamu, sirup, kopi bubuk)",
			"Roti, kue, dan camilan",
			"Produk olahan susu (yogurt, keju)"
		  ]
		},
		{
		  "sector": "Kerajinan Tangan dan Seni",
		  "category": [
			"Kerajinan kayu (furnitur, ukiran)",
			"Kerajinan bambu (anyaman, perabot)",
			"Batik, tenun, dan bordir",
			"Kerajinan logam (perhiasan, kuningan)",
			"Kerajinan kulit (tas, sepatu, dompet)",
			"Kerajinan keramik dan gerabah"
		  ]
		},
		{
		  "sector": "Fashion dan Tekstil",
		  "category": [
			"Konveksi dan garment",
			"Batik dan tenun tradisional",
			"Aksesoris fashion (tas, topi, syal)",
			"Produksi sepatu dan sandal",
			"Sablon dan printing kaos"
		  ]
		},
		{
		  "sector": "Jasa",
		  "category": [
			"Fotografi dan videografi",
			"Servis elektronik dan gadget",
			"Bengkel kendaraan",
			"Laundry dan cleaning service",
			"Tour & travel",
			"Bimbingan belajar dan kursus",
			"Salon, spa, barbershop",
			"Percetakan dan digital printing"
		  ]
		},
		{
		  "sector": "Teknologi dan Digital",
		  "category": [
			"Pengembangan software & aplikasi",
			"Pembuatan website & digital marketing",
			"E-commerce (online shop, dropshipping)",
			"Desain grafis & video editing"
		  ]
		},
		{
		  "sector": "Perdagangan (Retail & Grosir)",
		  "category": [
			"Toko kelontong & sembako",
			"Warung sembako & minimarket",
			"Grosir pakaian & aksesoris",
			"Toko elektronik",
			"Jual-beli produk pertanian"
		  ]
		},
		{
		  "sector": "Kesehatan & Produk Herbal",
		  "category": [
			"Produksi jamu tradisional",
			"Kosmetik alami & skincare",
			"Alat kesehatan sederhana",
			"Tanaman obat (herbal)"
		  ]
		},
		{
		  "sector": "Konstruksi & Properti",
		  "category": [
			"Material bangunan (batu bata, genteng)",
			"Kontraktor kecil & renovasi rumah",
			"Kerajinan batu alam & marmer"
		  ]
		},
		{
		  "sector": "Transportasi & Logistik",
		  "category": [
			"Rental kendaraan",
			"Jasa angkutan barang",
			"Ojek online & delivery"
		  ]
		},
		{
		  "sector": "Industri Kreatif",
		  "category": [
			"Musik & produksi alat musik",
			"Film indie & konten kreatif",
			"Craft beer & produk inovatif"
		  ]
		},
		{
		  "sector": "Pendidikan & Pelatihan",
		  "category": [
			"Bimbingan belajar & kursus bahasa",
			"Pelatihan keterampilan (menjahit, komputer)"
		  ]
		},
		{
		  "sector": "Lingkungan & Energi Terbarukan",
		  "category": [
			"Pengolahan sampah & daur ulang",
			"Produksi biogas & kompos",
			"Produk ramah lingkungan"
		  ]
		}
	  ],
	  "cooperative": [
		{
		  "sector": "Koperasi Simpan Pinjam",
		  "category": [
			"Koperasi kredit (KSP)",
			"Unit Simpan Pinjam (USP)"
		  ]
		},
		{
		  "sector": "Koperasi Konsumen",
		  "category": [
			"Koperasi karyawan (Kopkar)",
			"Koperasi sekolah/mahasiswa (Kopsis/Kopma)",
			"Koperasi konsumsi umum"
		  ]
		},
		{
		  "sector": "Koperasi Produsen",
		  "category": [
			"Koperasi petani (Koptan)",
			"Koperasi nelayan (Kopnel)",
			"Koperasi pengrajin (Kopra)",
			"Koperasi industri kecil"
		  ]
		},
		{
		  "sector": "Koperasi Pemasaran",
		  "category": [
			"Koperasi pemasaran hasil pertanian",
			"Koperasi pemasaran produk UMKM"
		  ]
		},
		{
		  "sector": "Koperasi Jasa",
		  "category": [
			"Koperasi transportasi (Koptrans)",
			"Koperasi jasa keuangan",
			"Koperasi jasa kesehatan"
		  ]
		},
		{
		  "sector": "Koperasi Serba Usaha (KSU)",
		  "category": [
			"Gabungan simpan pinjam & usaha produktif",
			"Koperasi dengan multi-bidang usaha"
		  ]
		}
	  ]
  }`
