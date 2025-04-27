package workers

import (
	"ametory-cooperative/objects"
	"ametory-cooperative/services"
	"encoding/json"
	"log"

	"github.com/AMETORY/ametory-erp-modules/context"
)

func GeminiCommand(erpContext *context.ERPContext) {
	appService, ok := erpContext.AppService.(*services.AppService)
	if ok {
		geminiCmdSrv := services.NewGeminiService(erpContext)
		dataSub := appService.Redis.Subscribe(*erpContext.Ctx, "GEMINI:COMMAND")
		for {
			msg, err := dataSub.ReceiveMessage(*erpContext.Ctx)
			if err != nil {
				log.Println(err)
			}
			var geminiData objects.GeminiResponse
			err = json.Unmarshal([]byte(msg.Payload), &geminiData)
			if err != nil {
				log.Println(err)
				continue
			}
			switch geminiData.Command {

			case "calculate_cashflow_health_score":
				geminiCmdSrv.CalculateCashflowHealthScore(&geminiData)
			case "identify_unusual_transactions":
				geminiCmdSrv.IdentifyUnusualTransactions(&geminiData)
			case "calculate_financial_health_ratios":
				geminiCmdSrv.CalculateFinancialHealthRatios(&geminiData)
			case "analyze_asset_efficiency":
				geminiCmdSrv.AnalyzeAssetEfficiency(&geminiData)
			case "analyze_debt_risk":
				geminiCmdSrv.AnalyzeDebtRisk(&geminiData)
			case "analyze_financial_growth":
				geminiCmdSrv.AnalyzeFinancialGrowth(&geminiData)
			case "forecast_sales":
				geminiCmdSrv.ForecastSales(&geminiData)
			case "analyze_business_opportunity":
				geminiCmdSrv.AnalyzeBusinessOpportunity(&geminiData)
			case "get_sales_chart":
				geminiCmdSrv.GetSalesChart(&geminiData)
			case "get_purchase_sales_chart":
				geminiCmdSrv.GetPurchaseSalesChart(&geminiData)
			case "get_purchase_chart":
				geminiCmdSrv.GetPurchaseChart(&geminiData)
			case "get_profit_loss_chart":
				geminiCmdSrv.GetProfitLossChart(&geminiData)
			case "get_best_selling_products":
				geminiCmdSrv.GetBestSellingProducts(&geminiData)
			case "change_password":
				geminiCmdSrv.ChangePassword(&geminiData)
			case "reset_password":
				geminiCmdSrv.ResetPassword(&geminiData)
			case "check_subscription_status":
				geminiCmdSrv.CheckSubscriptionStatus(&geminiData)
			case "get_quarterly_report":
				geminiCmdSrv.GetQuarterlyReport(&geminiData)
			case "get_balance_sheet_report":
				geminiCmdSrv.GetBalanceSheetReport(&geminiData)
			case "get_cash_flow_report":
				geminiCmdSrv.GetCashFlowReport(&geminiData)
			case "get_profit_loss_report":
				geminiCmdSrv.GetProfitLossReport(&geminiData)
			case "get_cash_balance":
				geminiCmdSrv.GetCashBalance(&geminiData)
			case "get_account_report":
				geminiCmdSrv.GetAccountReport(&geminiData)
			case "list_vendors":
				geminiCmdSrv.ListVendors(&geminiData)
			case "list_customers":
				geminiCmdSrv.ListCustomers(&geminiData)
			case "get_outstanding_balance":
				geminiCmdSrv.GetOutstandingBalance(&geminiData)
			case "get_total_outstanding":
				geminiCmdSrv.GetTotalOutstanding(&geminiData)
			case "get_total_outstanding_balance":
				geminiCmdSrv.GetTotalOutstandingBalance(&geminiData)
			case "get_unpaid_invoice_total":
				geminiCmdSrv.GetUnpaidInvoiceTotal(&geminiData)
			case "get_total_outstanding_invoices":
				geminiCmdSrv.GetTotalOutstandingInvoices(&geminiData)
			case "get_invoice_by_id":
				geminiCmdSrv.GetInvoiceByID(&geminiData)
			case "get_purchase_by_id":
				geminiCmdSrv.GetPurchaseByID(&geminiData)
			case "get_quotation_by_id":
				geminiCmdSrv.GetQuotationByID(&geminiData)
			case "get_purchase_order_by_id":
				geminiCmdSrv.GetPurchaseOrderByID(&geminiData)
			case "get_invoices":
				geminiCmdSrv.GetInvoices(&geminiData)
			case "get_purchases":
				geminiCmdSrv.GetPurchases(&geminiData)
			case "get_sales_report":
				geminiCmdSrv.GetSalesReport(&geminiData)
			case "get_purchase_report":
				geminiCmdSrv.GetPurchaseReport(&geminiData)
			case "get_transactions":
				geminiCmdSrv.GetTransactions(&geminiData)
			case "create_product":
				geminiCmdSrv.CreateProduct(&geminiData)
			case "create_transaction":
				geminiCmdSrv.CreateTransaction(&geminiData)
			case "create_purchase":
				geminiCmdSrv.CreatePurchase(&geminiData)
			case "create_purchase_order":
				geminiCmdSrv.CreatePurchaseOrder(&geminiData)
			case "create_invoice":
				geminiCmdSrv.CreateInvoice(&geminiData)
			case "create_quotation":
				geminiCmdSrv.CreateQuotation(&geminiData)
			case "get_customer_by_name":
				geminiCmdSrv.GetCustomerByName(&geminiData)
			case "get_vendor_by_name":
				geminiCmdSrv.GetVendorByName(&geminiData)
			case "get_supplier_by_name":
				geminiCmdSrv.GetSupplierByName(&geminiData)
			}
		}
	}
}
