package handlers

import (
	"fmt"

	"github.com/AMETORY/ametory-erp-modules/auth"
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	ctx     *context.ERPContext
	rbacSrv *auth.RBACService
}

func NewRoleHandler(ctx *context.ERPContext) *RoleHandler {
	rbacSrv, ok := ctx.RBACService.(*auth.RBACService)
	if !ok {
		panic("invalid rbac service")
	}
	return &RoleHandler{
		ctx:     ctx,
		rbacSrv: rbacSrv,
	}
}

func (h *RoleHandler) GetAllPermissionsHandler(c *gin.Context) {
	permissions, err := h.rbacSrv.GetAllPermissions()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": permissions})
}
func (h *RoleHandler) GetRolesHandler(c *gin.Context) {
	roles, err := h.rbacSrv.GetAllRoles(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": roles})
}

func (h *RoleHandler) DeleteRoleHandler(c *gin.Context) {
	id := c.Param("id")
	role, err := h.rbacSrv.GetRoleByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return

	}

	if err := h.rbacSrv.DeleteRole(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Role %s deleted successfully", role.Name)})
}

func (h *RoleHandler) CreateRoleHandler(c *gin.Context) {
	var role models.RoleModel
	if err := c.BindJSON(&role); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	if _, err := h.rbacSrv.CreateRole(role.Name, role.IsAdmin, role.IsSuperAdmin, role.IsMerchant, &companyID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": role})
}

func (h *RoleHandler) UpdateRoleHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.RoleModel
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	role, err := h.rbacSrv.GetRoleByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return

	}

	err = h.ctx.DB.Model(&role).Association("Permissions").Clear()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.rbacSrv.UpdateRole(id, input.Name, input.IsAdmin, input.IsSuperAdmin, input.IsMerchant, input.IsOwner); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	for _, v := range input.PermissionNames {
		fmt.Println("ASSIGN", id, v)
		h.rbacSrv.AssignPermissionToRoleID(id, v)
	}

	c.JSON(200, gin.H{"data": role})
}
