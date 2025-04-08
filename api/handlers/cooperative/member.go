package cooperative_handler

import (
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/cooperative"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type CooperativeMemberHandler struct {
	ctx           *context.ERPContext
	coopertiveSrv *cooperative.CooperativeService
}

func NewCooperativeMemberHandler(ctx *context.ERPContext) *CooperativeMemberHandler {
	cooperativeSrv, ok := ctx.CooperativeService.(*cooperative.CooperativeService)
	if !ok {
		panic("CooperativeService is not found")
	}
	return &CooperativeMemberHandler{
		ctx:           ctx,
		coopertiveSrv: cooperativeSrv,
	}
}

func (h *CooperativeMemberHandler) GetMembersHandler(c *gin.Context) {
	members, err := h.coopertiveSrv.CooperativeMemberService.GetMembers(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": members})
}

func (h *CooperativeMemberHandler) GetMemberHandler(c *gin.Context) {
	id := c.Param("id")
	member, err := h.coopertiveSrv.CooperativeMemberService.GetMemberByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": member})
}

func (h *CooperativeMemberHandler) CreateMemberHandler(c *gin.Context) {
	var input models.CooperativeMemberModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = h.coopertiveSrv.CooperativeMemberService.CreateMember(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Member created successfully"})
}

func (h *CooperativeMemberHandler) UpdateMemberHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.CooperativeMemberModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = h.coopertiveSrv.CooperativeMemberService.UpdateMember(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Member updated successfully"})
}

func (h *CooperativeMemberHandler) DeleteMemberHandler(c *gin.Context) {
	id := c.Param("id")
	err := h.coopertiveSrv.CooperativeMemberService.DeleteMember(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Member deleted successfully"})
}

func (h *CooperativeMemberHandler) InviteMemberHandler(c *gin.Context) {
	var input models.MemberInvitationModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := h.coopertiveSrv.CooperativeMemberService.InviteMember(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Member invitation sent successfully", "token": token})
}

func (h *CooperativeMemberHandler) ApproveMemberHandler(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(string)
	err := h.coopertiveSrv.CooperativeMemberService.ApproveMemberByID(id, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Member approved successfully"})
}
