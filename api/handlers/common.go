package handlers

import (
	"ametory-cooperative/app_models"
	"ametory-cooperative/config"
	"ametory-cooperative/objects"
	"ametory-cooperative/services"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AMETORY/ametory-erp-modules/auth"
	"github.com/AMETORY/ametory-erp-modules/company"
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/cooperative"
	"github.com/AMETORY/ametory-erp-modules/file"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/AMETORY/ametory-erp-modules/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommonHandler struct {
	ctx            *context.ERPContext
	companyService *company.CompanyService
	appService     *services.AppService
	rbacService    *auth.RBACService
	authService    *auth.AuthService
	fileService    *file.FileService
	cooperativeSrv *cooperative.CooperativeService
}

func NewCommonHandler(ctx *context.ERPContext) *CommonHandler {
	companyService, ok := ctx.CompanyService.(*company.CompanyService)
	if !ok {
		panic("CompanyService is not instance of company.CompanyService")
	}
	appService, ok := ctx.AppService.(*services.AppService)
	if !ok {
		panic("AppService is not instance of app.AppService")
	}
	rbacService, ok := ctx.RBACService.(*auth.RBACService)
	if !ok {
		panic("RBACService is not instance of auth.RBACService")
	}
	authService, ok := ctx.AuthService.(*auth.AuthService)
	if !ok {
		panic("AuthService is not instance of auth.AuthService")
	}
	fileService, ok := ctx.FileService.(*file.FileService)
	if !ok {
		panic("FileService is not instance of file.FileService")
	}
	cooperativeSrv, ok := ctx.CooperativeService.(*cooperative.CooperativeService)
	if !ok {
		panic("FileService is not instance of cooperative.CooperativeService")
	}

	return &CommonHandler{
		ctx:            ctx,
		companyService: companyService,
		appService:     appService,
		rbacService:    rbacService,
		authService:    authService,
		fileService:    fileService,
		cooperativeSrv: cooperativeSrv,
	}
}

func (h *CommonHandler) GetMembersHandler(c *gin.Context) {
	members, err := h.cooperativeSrv.CooperativeMemberService.GetMembers(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": members})
}

func (h *CommonHandler) GetRolesHandler(c *gin.Context) {
	roles, err := h.rbacService.GetAllRoles(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	items := roles.Items.(*[]models.RoleModel)
	newItems := make([]models.RoleModel, 0)
	for _, v := range *items {
		if !v.IsSuperAdmin {
			v.Permissions = nil
			newItems = append(newItems, v)
		}
	}
	roles.Items = &newItems
	c.JSON(200, gin.H{"data": roles})
}

func (h *CommonHandler) InvitedHandler(c *gin.Context) {
	members, err := h.cooperativeSrv.CooperativeMemberService.GetInvitedMembers(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": members})
}

func (h *CommonHandler) DeleteInvitedHandler(c *gin.Context) {
	invitationID := c.Param("id")
	err := h.cooperativeSrv.CooperativeMemberService.DeleteInvitation(invitationID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Invitation deleted successfully"})
}

func (h *CommonHandler) DeleteUserHandler(c *gin.Context) {
	userID := c.Param("id")
	companyID := c.MustGet("companyID").(string)
	var user models.UserModel
	err := h.ctx.DB.Preload("Roles", "company_id = ?", companyID).First(&user, "id = ? ", userID).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var company models.CompanyModel
	err = h.ctx.DB.First(&company, "id = ? ", companyID).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	for _, role := range user.Roles {
		h.ctx.DB.Model(&user).Association("Roles").Delete(role)
	}

	err = h.ctx.DB.Model(&company).Association("Users").Delete(user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted successfully"})
}
func (h *CommonHandler) UpdateRoleHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.RoleModel
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	var user models.UserModel
	err = h.ctx.DB.Preload("Roles", "company_id = ?", companyID).Where("id = ?", id).First(&user).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	for _, role := range user.Roles {
		h.ctx.DB.Model(&user).Association("Roles").Delete(role)
	}
	h.ctx.DB.Model(&user).Association("Roles").Append(&input)
	c.JSON(200, gin.H{"message": "Role updated successfully"})

}
func (h *CommonHandler) GetCompanyUsersHandler(c *gin.Context) {
	companyID := c.MustGet("companyID").(string)
	users, err := h.companyService.GetCompanyUsers(companyID, *c.Request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": users})

}
func (h *CommonHandler) InviteMemberHandler(c *gin.Context) {
	var input struct {
		FullName  string  `json:"full_name"`
		RoleID    *string `json:"role_id"`
		Email     string  `json:"email"`
		ProjectID *string `json:"project_id"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var data models.MemberInvitationModel
	data.FullName = input.FullName
	data.RoleID = input.RoleID
	data.ProjectID = input.ProjectID
	data.Email = input.Email

	var user models.UserModel
	var link = ""
	var password = ""

	err = h.ctx.DB.Where("email = ?", input.Email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new user if not exists
		username := utils.CreateUsernameFromFullName(input.FullName)
		// fmt.Println("username", username)
		password = utils.RandString(8, false)
		u, err := h.authService.Register(input.FullName, username, input.Email, password, "")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user = *u

	} else {
		user.Password = ""
	}

	var company models.CompanyModel
	h.ctx.DB.Where("id = ?", c.GetHeader("ID-Company")).First(&company)

	err = h.ctx.DB.Model(&user).Association("Companies").Append(&company)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data.UserID = user.ID
	if user.VerificationToken != "" {
		data.Token = user.VerificationToken
	}

	data.InviterID = c.MustGet("userID").(string)
	data.CompanyID = &company.ID
	token, err := h.cooperativeSrv.CooperativeMemberService.InviteMember(&data)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var role models.RoleModel
	h.ctx.DB.Where("id = ?", input.RoleID).First(&role)
	link = fmt.Sprintf("%s/invitation/verify/%s", h.appService.Config.Server.FrontendURL, token)
	notif := fmt.Sprintf("Anda telah diundang untuk bergabung di %s ", company.Name)
	if input.ProjectID != nil {
		var project models.ProjectModel
		h.ctx.DB.Where("id = ?", *input.ProjectID).First(&project)
		notif += fmt.Sprintf("dalam proyek %s", project.Name)
	}
	var emailData objects.EmailData = objects.EmailData{
		FullName: user.FullName,
		Email:    user.Email,
		Subject:  "Selamat datang di " + h.appService.Config.Server.AppName,
		Notif:    notif,
		Link:     link,
		Password: password,
		RoleName: role.Name,
	}

	b, err := json.Marshal(emailData)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	// fmt.Println("SEND MAIL", string(b))
	err = h.appService.Redis.Publish(*h.ctx.Ctx, "SEND:MAIL", string(b)).Err()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Member invited successfully", "token": token})
}

func (h *CommonHandler) AcceptMemberInvitationHandler(c *gin.Context) {
	token := c.Param("token")

	var invitation models.MemberInvitationModel
	h.ctx.DB.Where("token = ?", token).First(&invitation)

	var user models.UserModel
	h.ctx.DB.Where("id = ?", invitation.UserID).First(&user)
	now := time.Now()
	if user.VerifiedAt == nil {
		user.VerifiedAt = &now
		user.VerificationToken = ""
		user.VerificationTokenExpiredAt = nil
		h.ctx.DB.Save(&user)
	}

	if invitation.IsCooperativeMember {
		err := h.cooperativeSrv.CooperativeMemberService.AcceptMemberInvitation(token, invitation.UserID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}
	if (invitation.RoleID) != nil {
		var role models.RoleModel
		h.ctx.DB.Where("id = ?", invitation.RoleID).First(&role)
		h.ctx.DB.Model(&user).Association("Roles").Append(&role)
	}
	c.JSON(200, gin.H{"message": "Member invitation accepted successfully"})
}

func (h *CommonHandler) UploadFileFromBase64Handler(c *gin.Context) {
	var input struct {
		FileName string `json:"fileName"`
		Base64   string `json:"base64"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.ctx.Request = c.Request

	fileObject := models.FileModel{}
	refID, _ := c.GetPostForm("ref_id")
	refType, _ := c.GetPostForm("ref_type")
	skipSave := false
	skipSaveStr, _ := c.GetPostForm("skip_save")
	if skipSaveStr == "true" || skipSaveStr == "1" {
		skipSave = true
	}

	filename := utils.GenerateRandomNumber(6)

	fileObject.FileName = utils.FilenameTrimSpace(filename)
	fileObject.RefID = refID
	fileObject.RefType = refType
	fileObject.SkipSave = skipSave
	err := h.fileService.UploadFileFromBase64(input.Base64, config.App.Server.StorageProvider, "files", &fileObject)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "File uploaded successfully", "data": fileObject})
}

func (h *CommonHandler) UploadFileHandler(c *gin.Context) {
	h.ctx.Request = c.Request

	fileObject := models.FileModel{}
	refID, _ := c.GetPostForm("ref_id")
	refType, _ := c.GetPostForm("ref_type")
	skipSave := false
	skipSaveStr, _ := c.GetPostForm("skip_save")
	if skipSaveStr == "true" || skipSaveStr == "1" {
		skipSave = true
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fileByte, err := utils.FileHeaderToBytes(file)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	filename := file.Filename

	fileObject.FileName = utils.FilenameTrimSpace(filename)
	fileObject.RefID = refID
	fileObject.RefType = refType
	fileObject.SkipSave = skipSave

	if err := h.fileService.UploadFile(fileByte, config.App.Server.StorageProvider, "files", &fileObject); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "File uploaded successfully", "data": fileObject})
}

func (h *CommonHandler) CompanySettingHandler(c *gin.Context) {
	// h.ctx.Request = c.Request
	// data, err := h.companyService.GetCompanyByID(c.GetHeader("ID-Company"))
	// if err != nil {
	// 	c.JSON(400, gin.H{"error": err.Error()})
	// 	return
	// }
	// var companySetting com.CompanySetting
	// err = h.ctx.DB.First(&companySetting, "company_id = ?", data.ID).Error
	// if err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		companySetting.ID = utils.Uuid()
	// 		companySetting.CompanyID = &data.ID
	// 		if err := h.ctx.DB.Create(&companySetting).Error; err != nil {
	// 			c.JSON(500, gin.H{"error": "Failed to create company setting"})
	// 			return
	// 		}

	// 	}
	// }
	// var response = struct {
	// 	models.CompanyModel
	// 	Setting com.CompanySetting `json:"setting"`
	// }{
	// 	CompanyModel: *data,
	// 	Setting:      companySetting,
	// }

	userID := c.MustGet("userID").(string)
	companyID := c.GetHeader("ID-Company")

	var setting app_models.CustomSettingModel
	err := h.ctx.DB.Where("id = ?", c.GetHeader("ID-Company")).First(&setting).Error
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	if setting.IsCooperation {
		var cooperationSetting models.CooperativeSettingModel
		err := h.ctx.DB.Preload(clause.Associations).Where("company_id = ?", c.GetHeader("ID-Company")).First(&cooperationSetting).Error
		if err == nil {
			setting.CooperativeSetting = &cooperationSetting
		}

	}
	isSuperAdmin, _ := checkSuperAdmin(h.ctx, userID, companyID)
	if !isSuperAdmin {
		setting.GeminiAPIKey = nil
		setting.WhatsappWebHost = nil
		setting.WhatsappWebMockNumber = nil
		setting.WhatsappWebIsMocked = nil
	}
	setting.CashflowGroupSettingData = nil

	c.JSON(200, gin.H{"message": "Get company setting successfully", "data": setting})

}
func (h *CommonHandler) UpdateCompanySettingHandler(c *gin.Context) {
	var input app_models.CustomSettingModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.ctx.Request = c.Request
	err = h.ctx.DB.Where("id = ?", input.ID).Updates(&input).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// utils.LogJson(input)
	// err = h.ctx.DB.Model(&input.Setting).Where("company_id = ?", input.CompanyModel.ID).Updates(map[string]any{
	// 	"gemini_api_key":           input.Setting.GeminiAPIKey,
	// 	"whatsapp_web_host":        input.Setting.WhatsappWebHost,
	// 	"whatsapp_web_mock_number": input.Setting.WhatsappWebMockNumber,
	// 	"whatsapp_web_is_mocked":   input.Setting.WhatsappWebIsMocked,
	// }).Error
	// if err != nil {
	// 	c.JSON(500, gin.H{"error": "Failed to update company setting"})
	// 	return
	// }
	c.JSON(200, gin.H{"message": " company setting update successfully"})
}
func (h *CommonHandler) UpdateCooperativeSettingHandler(c *gin.Context) {
	var input models.CooperativeSettingModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.ctx.Request = c.Request
	err = h.ctx.DB.Where("id = ?", input.ID).Updates(&input).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// err = h.ctx.DB.Model(&input.Setting).Where("company_id = ?", input.CompanyModel.ID).Updates(map[string]any{
	// 	"gemini_api_key":           input.Setting.GeminiAPIKey,
	// 	"whatsapp_web_host":        input.Setting.WhatsappWebHost,
	// 	"whatsapp_web_mock_number": input.Setting.WhatsappWebMockNumber,
	// 	"whatsapp_web_is_mocked":   input.Setting.WhatsappWebIsMocked,
	// }).Error
	// if err != nil {
	// 	c.JSON(500, gin.H{"error": "Failed to update company setting"})
	// 	return
	// }
	c.JSON(200, gin.H{"message": " company setting update successfully"})
}

func checkSuperAdmin(erpContext *context.ERPContext, userID string, companyID string) (bool, error) {
	var admin models.UserModel

	// Cari pengguna beserta peran dan izin
	if err := erpContext.DB.Preload("Roles", func(db *gorm.DB) *gorm.DB {
		return db.Where("company_id = ?", companyID)
	}).First(&admin, "id = ?", userID).Error; err != nil {
		return false, errors.New("admin not found")
	}

	// Periksa izin
	for _, role := range admin.Roles {
		if role.IsSuperAdmin {
			return true, nil
		}
	}

	return false, nil
}

func (h *CommonHandler) InviteUserHandler(c *gin.Context) {
	var input struct {
		FullName            string  `json:"full_name"`
		RoleID              *string `json:"role_id"`
		Email               string  `json:"email"`
		IsCooperativeMember bool    `json:"is_cooperative_member"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var data models.MemberInvitationModel
	data.FullName = input.FullName
	data.RoleID = input.RoleID
	data.Email = input.Email

	var user models.UserModel
	var link = ""
	var password = ""

	err = h.ctx.DB.Where("email = ?", input.Email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new user if not exists
		username := utils.CreateUsernameFromFullName(input.FullName)
		// fmt.Println("username", username)
		password = utils.RandString(8, false)
		u, err := h.authService.Register(input.FullName, username, input.Email, password, "")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user = *u

	} else {
		user.Password = ""
	}

	var company models.CompanyModel
	h.ctx.DB.Where("id = ?", c.GetHeader("ID-Company")).First(&company)

	err = h.ctx.DB.Model(&user).Association("Companies").Append(&company)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data.UserID = user.ID
	if user.VerificationToken != "" {
		data.Token = user.VerificationToken
	}

	data.InviterID = c.MustGet("userID").(string)
	data.CompanyID = &company.ID
	token, err := h.cooperativeSrv.CooperativeMemberService.InviteMember(&data)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var role models.RoleModel
	h.ctx.DB.Where("id = ?", input.RoleID).First(&role)
	link = fmt.Sprintf("%s/invitation/verify/%s", h.appService.Config.Server.FrontendURL, token)
	notif := fmt.Sprintf("Anda telah diundang untuk bergabung di %s ", company.Name)

	var emailData objects.EmailData = objects.EmailData{
		FullName: user.FullName,
		Email:    user.Email,
		Subject:  "Selamat datang di " + h.appService.Config.Server.AppName,
		Notif:    notif,
		Link:     link,
		Password: password,
		RoleName: role.Name,
	}

	b, err := json.Marshal(emailData)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	// fmt.Println("SEND MAIL", string(b))
	err = h.appService.Redis.Publish(*h.ctx.Ctx, "SEND:MAIL", string(b)).Err()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Member invited successfully", "token": token})
}
