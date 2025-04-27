package handlers

import (
	"ametory-cooperative/objects"
	"ametory-cooperative/services"
	"encoding/json"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/AMETORY/ametory-erp-modules/thirdparty/google"
	"github.com/gin-gonic/gin"
)

type GeminiHandler struct {
	ctx           *context.ERPContext
	geminiService *google.GeminiService
	appService    *services.AppService
}

func NewGeminiHandler(ctx *context.ERPContext) *GeminiHandler {
	geminiService, ok := ctx.ThirdPartyServices["GEMINI"].(*google.GeminiService)
	if !ok {
		panic("GeminiService is not found")
	}
	appService, ok := ctx.AppService.(*services.AppService)
	if !ok {
		panic("AppService is not instance of app.AppService")
	}

	return &GeminiHandler{
		ctx:           ctx,
		geminiService: geminiService,
		appService:    appService,
	}
}

func (h *GeminiHandler) GenerateContentHandler(c *gin.Context) {

	var input struct {
		Content string
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var companyID *string
	if c.GetHeader("ID-Company") != "" {
		compID := c.GetHeader("ID-Company")
		companyID = &compID
	}
	var histories []models.GeminiHistoryModel = h.geminiService.GetHistories(nil, companyID)

	chatHistories := []map[string]any{}
	for _, v := range histories {
		chatHistories = append(chatHistories, map[string]any{
			"role":    "user",
			"content": v.Input,
		})
		chatHistories = append(chatHistories, map[string]any{
			"role":    "model",
			"content": v.Output,
		})
	}

	// h.geminiService.SetupModel(companySetting.GeminiAPIKey)
	// utils.LogJson(chatHistories)
	output, err := h.geminiService.GenerateContent(*h.ctx.Ctx, input.Content, chatHistories, "", "")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	outputResp := objects.GeminiResponse{}
	err = json.Unmarshal([]byte(output), &outputResp)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if outputResp.Type == "command" {
		userID := c.MustGet("userID").(string)
		companyID := c.MustGet("companyID").(string)
		outputResp.UserID = userID
		outputResp.CompanyID = companyID
		b, _ := json.Marshal(outputResp)
		h.appService.Redis.Publish(*h.ctx.Ctx, "GEMINI:COMMAND", string(b))
	}

	c.JSON(200, gin.H{"data": outputResp})
}

func (h *GeminiHandler) DeleteHistoryHandler(c *gin.Context) {
	historyId := c.Param("historyId")

	err := h.ctx.DB.Delete(&models.GeminiHistoryModel{}, "id = ?", historyId).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "History deleted successfully"})
}

func (h *GeminiHandler) UpdateHistoryHandler(c *gin.Context) {
	// id := c.Param("id")
	historyId := c.Param("historyId")

	var input models.GeminiHistoryModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = h.ctx.DB.Model(&models.GeminiHistoryModel{}).Where("id = ?", historyId).Updates(map[string]any{
		"input":    input.Input,
		"output":   input.Output,
		"agent_id": input.AgentID,
	}).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "History updated successfully"})
}
