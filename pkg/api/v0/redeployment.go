package v0

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack-redeployment/pkg/api"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/store"
)

type RedeploymentHandler struct {
	api.RedeploymentHandlerInterface

	store store.Store
}

func NewRedeploymentHandler(store store.Store) *RedeploymentHandler {
	return &RedeploymentHandler{
		store: store,
	}
}

func (h *RedeploymentHandler) FindAll(ctx *gin.Context) {
	rdList := []*api.Redeployment{}
	f := func(n int) []interface{} {
		m := []interface{}{}
		for i := 0; i < n; i++ {
			rd := &api.Redeployment{}
			rdList = append(rdList, rd)
			m = append(m, rd)
		}
		return m
	}
	h.store.List(getKey(""), f)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"redeployments": rdList,
	})

}

func (h *RedeploymentHandler) Find(ctx *gin.Context) {
	rdID := ctx.Param("redeployment_id")

	var rd api.Redeployment
	err := h.store.Get(getKey(rdID), &rd)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Redeployment `%s` is not found.", rdID), nil)
		return
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"redeployment": rd,
	})
}

func (h *RedeploymentHandler) Create(ctx *gin.Context) {

	var request api.Redeployment
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if request.ID == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: id is empty."), nil)
		return
	}

	key := getKey(request.ID)
	var rd api.Redeployment
	err = h.store.Get(key, &rd)
	if err == nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: Redeployment `%s` is already exists.", request.ID), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	request.APIType = api.APITypeRedeploymentV0
	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"redeployment": request,
	})
}

func (h *RedeploymentHandler) Update(ctx *gin.Context) {
	rdID := ctx.Param("redeployment_id")

	var request api.Redeployment
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if rdID != request.ID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: Can't change Redeployment ID."), nil)
		return
	}

	key := getKey(request.ID)
	var rd api.Redeployment
	err = h.store.Get(key, &rd)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: Redeployment `%s` is not found.", request.ID), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"redeployment": request,
	})
}

func (h *RedeploymentHandler) Delete(ctx *gin.Context) {
	rdID := ctx.Param("redeployment_id")

	key := getKey(rdID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"redeployment": nil,
	})
}

func getKey(rdID string) string {
	return filepath.Join("redeployment", rdID)
}
