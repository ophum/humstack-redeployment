package api

import "github.com/gin-gonic/gin"

type RedeploymentHandlerInterface interface {
	FindAll(ctx *gin.Context)
	Find(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type RedeploymentHandler struct {
	router *gin.RouterGroup
	rdhi   RedeploymentHandlerInterface
}

func NewRedeploymentHandler(router *gin.RouterGroup, rdhi RedeploymentHandlerInterface) *RedeploymentHandler {
	return &RedeploymentHandler{
		router: router,
		rdhi:   rdhi,
	}
}

func (h *RedeploymentHandler) RegisterHandlers() {
	h.router.GET("", h.rdhi.FindAll)
	h.router.GET("/:redeployment_id", h.rdhi.Find)
	h.router.POST("", h.rdhi.Create)
	h.router.PUT("/:redeployment_id", h.rdhi.Update)
	h.router.DELETE("/:redeployment_id", h.rdhi.Delete)
}
