package handlers

import (
	"time"

	"github.com/AMETORY/ametory-erp-modules/company"
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/file"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/AMETORY/ametory-erp-modules/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	ctx            *context.ERPContext
	userService    *user.UserService
	fileService    *file.FileService
	companyService *company.CompanyService
}

func NewUserHandler(ctx *context.ERPContext) *UserHandler {
	userService, ok := ctx.UserService.(*user.UserService)
	if !ok {
		panic("UserService is not instance of user.UserService")
	}
	fileService, ok := ctx.FileService.(*file.FileService)
	if !ok {
		panic("FileService is not instance of file.FileService")
	}

	companyService, ok := ctx.CompanyService.(*company.CompanyService)
	if !ok {
		panic("CompanyService is not instance of company.CompanyService")
	}

	return &UserHandler{
		ctx:            ctx,
		userService:    userService,
		fileService:    fileService,
		companyService: companyService,
	}
}

func (h *UserHandler) GetUserDetailHandler(c *gin.Context) {
	userID := c.Param("id")
	var user models.UserModel
	err := h.ctx.DB.First(&user, "id = ?", userID).Error
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, gin.H{"data": user})
}

func (h *UserHandler) GetUserListHandler(c *gin.Context) {
	companyID := c.MustGet("companyID").(string)
	users, err := h.companyService.GetCompanyUsers(companyID, *c.Request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": users})
}

func (h *UserHandler) FinishActivityHandler(c *gin.Context) {
	id := c.Param("id")
	input := struct {
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
		Notes     *string  `json:"notes"`
		FileID    *string  `json:"file_id"`
		RefType   *string  `json:"ref_type"`
	}{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userActivity := models.UserActivityModel{}
	err := h.ctx.DB.Where("id = ?", id).First(&userActivity).Error
	if err != nil {
		c.JSON(404, gin.H{"error": "UserActivity not found"})
		return
	}
	_, err = h.userService.FinishActivityByUser(c.MustGet("userID").(string), userActivity.ActivityType, input.Latitude, input.Longitude, input.Notes)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if input.FileID != nil {
		if input.RefType != nil {
			h.fileService.UpdateFileRefByID(*input.FileID, userActivity.ID, *input.RefType)
		} else {
			h.fileService.UpdateFileRefByID(*input.FileID, userActivity.ID, "user_activity")
		}
	}
	c.JSON(200, gin.H{"message": "Activity updated successfully"})

}
func (h *UserHandler) CreateActivityHandler(c *gin.Context) {
	input := struct {
		Latitude    *float64                `json:"latitude"`
		Longitude   *float64                `json:"longitude"`
		Notes       *string                 `json:"notes"`
		FileID      *string                 `json:"file_id"`
		RefID       *string                 `json:"ref_id"`
		RefType     *string                 `json:"ref_type"`
		ActivityTpe models.UserActivityType `json:"activity_type"`
	}{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()

	userID := c.MustGet("userID").(string)
	companyID := c.MustGet("companyID").(string)
	userActivity := models.UserActivityModel{
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
		Notes:        input.Notes,
		ActivityType: input.ActivityTpe,
		StartedAt:    &now,
		RefID:        input.RefID,
		RefType:      input.RefType,
		UserID:       userID,
		CompanyID:    &companyID,
	}

	err := h.userService.CreateActivity(c.MustGet("userID").(string), &userActivity)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to create activity", "error": err.Error()})
		return
	}

	if input.FileID != nil {
		if input.ActivityTpe == models.UserActivityCheckPoint {
			h.fileService.UpdateFileRefByID(*input.FileID, userActivity.ID, "check_point")
		} else {
			h.fileService.UpdateFileRefByID(*input.FileID, userActivity.ID, "user_activity")
		}
	}
	c.JSON(200, gin.H{"message": "Activity created successfully"})

}

func (h *UserHandler) GetLastClockinHandler(c *gin.Context) {
	var input struct {
		ThresholdDuration int    `json:"threshold_duration" binding:"required"`
		ThresholdUnit     string `json:"threshold_unit" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var thresholdDuration time.Duration
	switch input.ThresholdUnit {
	case "minutes":
		thresholdDuration = time.Minute * time.Duration(input.ThresholdDuration)
	case "hours":
		thresholdDuration = time.Hour * time.Duration(input.ThresholdDuration)
	case "days":
		thresholdDuration = time.Hour * 24 * time.Duration(input.ThresholdDuration)
	}
	companyID := c.MustGet("companyID").(string)
	userID := c.MustGet("userID").(string)
	userActivity, _ := h.userService.GetLastClockinByUser(userID, companyID, thresholdDuration)

	c.JSON(200, gin.H{"last_clockin": userActivity})
}
func (h *UserHandler) ClockInHandler(c *gin.Context) {
	input := struct {
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
		Notes     *string  `json:"notes"`
		FileID    *string  `json:"file_id"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}
	now := time.Now()
	var merchantID *string
	var refType *string
	merchant, ok := c.MustGet("merchant").(*models.MerchantModel)
	if ok && merchant != nil {
		refTypeStr := "MERCHANT"
		merchantID = &merchant.ID
		refType = &refTypeStr
	}
	userID := c.MustGet("userID").(string)
	companyID := c.MustGet("companyID").(string)
	userActivity := models.UserActivityModel{
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
		Notes:        input.Notes,
		ActivityType: models.UserActivityClockIn,
		StartedAt:    &now,
		RefID:        merchantID,
		RefType:      refType,
		UserID:       userID,
		CompanyID:    &companyID,
	}

	err := h.userService.CreateActivity(c.MustGet("userID").(string), &userActivity)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to create activity", "error": err.Error()})
		return
	}

	if input.FileID != nil {
		h.fileService.UpdateFileRefByID(*input.FileID, userActivity.ID, "clock_in")
	}
	c.JSON(200, gin.H{"message": "Activity created successfully"})
}

func (h *UserHandler) ClockOutHandler(c *gin.Context) {
	input := struct {
		ClockInID *string  `json:"clock_in_id"`
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
		Notes     *string  `json:"notes"`
		FileID    *string  `json:"file_id"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	userActivity, err := h.userService.FinishActivityByActivityID(c.MustGet("userID").(string), *input.ClockInID, input.Latitude, input.Longitude, input.Notes)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to create activity", "error": err.Error()})
		return
	}

	if input.FileID != nil {
		h.fileService.UpdateFileRefByID(*input.FileID, userActivity.ID, "clock_out")
	}
	c.JSON(200, gin.H{"message": "Activity updated successfully"})
}

func (h *UserHandler) GetActivityHandler(c *gin.Context) {

	data, err := h.userService.GetUserActivitiesByUserID(c.MustGet("userID").(string), *c.Request)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Activity created successfully", "data": data})
}

func (h *UserHandler) BreakHandler(c *gin.Context) {
	input := struct {
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
		Notes     *string  `json:"notes"`
		FileID    *string  `json:"file_id"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}
	now := time.Now()

	var merchantID *string
	var refType *string
	merchant, ok := c.MustGet("merchant").(*models.MerchantModel)
	if ok && merchant != nil {
		refTypeStr := "MERCHANT"
		merchantID = &merchant.ID
		refType = &refTypeStr
	}
	userID := c.MustGet("userID").(string)
	companyID := c.MustGet("companyID").(string)

	userActivity := models.UserActivityModel{
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
		Notes:        input.Notes,
		ActivityType: models.UserActivityBreak,
		StartedAt:    &now,
		RefID:        merchantID,
		RefType:      refType,
		UserID:       userID,
		CompanyID:    &companyID,
	}

	err := h.userService.CreateActivity(c.MustGet("userID").(string), &userActivity)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to create activity", "error": err.Error()})
		return
	}

	if input.FileID != nil {
		h.fileService.UpdateFileRefByID(*input.FileID, userActivity.ID, "user_activity")
	}
	c.JSON(200, gin.H{"message": "Activity created successfully"})
}

func (h *UserHandler) BreakOffHandler(c *gin.Context) {
	input := struct {
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
		Notes     *string  `json:"notes"`
		FileID    *string  `json:"file_id"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	userActivity, err := h.userService.FinishActivityByUser(c.MustGet("userID").(string), models.UserActivityBreak, input.Latitude, input.Longitude, input.Notes)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to create activity", "error": err.Error()})
		return
	}

	if input.FileID != nil {
		h.fileService.UpdateFileRefByID(*input.FileID, userActivity.ID, "user_activity")
	}
	c.JSON(200, gin.H{"message": "Activity updated successfully"})
}
