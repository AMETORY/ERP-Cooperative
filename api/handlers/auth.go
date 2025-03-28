package handlers

import (
	"ametory-cooperative/objects"
	"ametory-cooperative/services"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AMETORY/ametory-erp-modules/auth"
	"github.com/AMETORY/ametory-erp-modules/company"
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/AMETORY/ametory-erp-modules/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	ctx            *context.ERPContext
	authService    *auth.AuthService
	companyService *company.CompanyService
	rbacService    *auth.RBACService
	appService     *services.AppService
}

func NewAuthHandler(ctx *context.ERPContext) *AuthHandler {
	authService, ok := ctx.AuthService.(*auth.AuthService)
	if !ok {
		panic("AuthService is not instance of auth.AuthService")
	}

	appService, ok := ctx.AppService.(*services.AppService)
	if !ok {
		panic("AppService is not instance of app.AppService")
	}
	companyService, ok := ctx.CompanyService.(*company.CompanyService)
	if !ok {
		panic("CompanyService is not instance of company.CompanyService")
	}
	rbacService, ok := ctx.RBACService.(*auth.RBACService)
	if !ok {
		panic("RBACService is not instance of auth.RBACService")
	}
	return &AuthHandler{
		ctx:            ctx,
		authService:    authService,
		appService:     appService,
		companyService: companyService,
		rbacService:    rbacService,
	}
}

func (a *AuthHandler) LoginHandler(c *gin.Context) {
	req := &objects.LoginRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := a.authService.Login(req.Email, req.Password, true)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	token, err := utils.GenerateJWT(user.ID, time.Now().AddDate(0, 0, a.appService.Config.Server.TokenExpiredDay).Unix(), a.appService.Config.Server.SecretKey)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"token": token})
}

func (a *AuthHandler) RegisterHandler(c *gin.Context) {
	req := &objects.RegisterRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	username := utils.CreateUsernameFromFullName(req.FullName)
	if req.Password == "" {
		req.Password = utils.RandString(8, false)
	}
	user, err := a.authService.Register(req.FullName, username, req.Email, req.Password, req.PhoneNumber)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var emailData objects.EmailData = objects.EmailData{
		FullName: user.FullName,
		Email:    user.Email,
		Subject:  "Selamat datang di " + a.appService.Config.Server.AppName,
		Notif:    "Silakan verifikasi akun Anda, dengan mengklik link berikut",
		Link:     fmt.Sprintf("%s/verify/%s", a.appService.Config.Server.FrontendURL, user.VerificationToken),
		Password: req.Password,
	}

	b, err := json.Marshal(emailData)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("SEND MAIL", string(b))
	err = a.appService.Redis.Publish(*a.ctx.Ctx, "SEND:MAIL", string(b)).Err()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Register success, please check your email to verify your account"})
}

func (h *AuthHandler) ForgotPasswordHandler(c *gin.Context) {
	var input struct {
		EmailOrPhoneNumber string `json:"email_or_phone_number"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}
	authSrv, ok := h.ctx.AuthService.(*auth.AuthService)
	if !ok {
		c.JSON(500, gin.H{"message": "Auth service is not available"})
		return
	}
	user, err := authSrv.GetUserByEmailOrPhone(input.EmailOrPhoneNumber)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	newPassword := utils.RandString(8, false)

	hashedPassword, err := models.HashPassword(newPassword)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	user.Password = hashedPassword

	if err := h.ctx.DB.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	var emailData objects.EmailData = objects.EmailData{
		FullName: user.FullName,
		Email:    user.Email,
		Subject:  "Permintaan Penggatian PASSWORD",
		Notif:    "Berikut ini adalah PASSWORD baru Anda",
		Link:     "",
		Password: newPassword,
	}

	b, err := json.Marshal(emailData)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("SEND MAIL", string(b))
	err = h.appService.Redis.Publish(*h.ctx.Ctx, "SEND:MAIL", string(b)).Err()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "New password sent"})
}

func (h *AuthHandler) ChangePasswordHandler(c *gin.Context) {
	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.GetUserByID(c.MustGet("userID").(string))
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	err = h.authService.ChangePassword(user.ID, input.OldPassword, input.NewPassword)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Password changed successfully"})
}

func (h *AuthHandler) VerificationEmailHandler(c *gin.Context) {
	token := c.Param("token")
	if h.authService == nil {
		c.JSON(500, gin.H{"message": "Auth service is not available"})
		return
	}
	err := h.authService.VerificationEmail(token)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Email verified"})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	user := c.MustGet("user").(models.UserModel)
	h.ctx.DB.Preload("Companies").Find(&user)

	c.JSON(200, gin.H{"user": user, "companies": user.Companies, "member": c.MustGet("member").(models.CooperativeMemberModel)})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var input struct {
		FullName       string            `json:"full_name"`
		Address        string            `json:"address"`
		ProfilePicture *models.FileModel `json:"profile_picture,omitempty" gorm:"-"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(models.UserModel)
	if input.FullName != "" {
		user.FullName = input.FullName
	}
	if input.Address != "" {
		user.Address = input.Address
	}

	if input.ProfilePicture != nil {
		input.ProfilePicture.RefID = user.ID
		input.ProfilePicture.RefType = "user"
		err = h.ctx.DB.Save(&input.ProfilePicture).Error
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	err = h.ctx.DB.Save(&user).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Profile updated", "user": user})
}
