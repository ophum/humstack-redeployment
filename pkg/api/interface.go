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
	r := h.router.Group("/redeployments")
	{
		r.GET("", h.rdhi.FindAll)
		r.GET("/:redeployment_id", h.rdhi.Find)
		r.POST("", h.rdhi.Create)
		r.PUT("/:redeployment_id", h.rdhi.Update)
		r.DELETE("/:redeployment_id", h.rdhi.Delete)
	}
}
