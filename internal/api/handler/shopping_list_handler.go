package handler

import (
	"NutriPlan/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShoppingListHandler struct {
	service service.ShoppingListService
}

func NewShoppingListHandler(service service.ShoppingListService) *ShoppingListHandler {
	return &ShoppingListHandler{service: service}
}

type CreateListRequest struct {
	RecipeIDs []uint `json:"recipeIds" binding:"required"`
	ListName  string `json:"listName"`
}

// CreateShoppingList 创建购物清单
func (h *ShoppingListHandler) CreateShoppingList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req CreateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ListName == "" {
		req.ListName = "购物清单"
	}

	list, err := h.service.GenerateFromRecipes(userID.(uint), req.RecipeIDs, req.ListName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建清单失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

// GetShoppingLists 获取用户购物清单列表
func (h *ShoppingListHandler) GetShoppingLists(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	lists, total, err := h.service.GetShoppingLists(userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"lists": lists,
		"total": total,
	})
}

// GetShoppingListDetail 获取清单详情
func (h *ShoppingListHandler) GetShoppingListDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}

	list, err := h.service.GetShoppingListDetail(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "清单不存在"})
		return
	}

	// 解析 Items JSON 以便前端使用
	var items []service.ShoppingItem
	json.Unmarshal([]byte(list.Items), &items)

	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"items": items,
	})
}

// UpdateShoppingList 更新清单（主要是勾选状态）
func (h *ShoppingListHandler) UpdateShoppingList(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}

	var items []service.ShoppingItem
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list, err := h.service.GetShoppingListDetail(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "清单不存在"})
		return
	}

	itemsJSON, _ := json.Marshal(items)
	list.Items = string(itemsJSON)

	if err := h.service.UpdateShoppingList(list); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteShoppingList 删除清单
func (h *ShoppingListHandler) DeleteShoppingList(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}

	if err := h.service.DeleteShoppingList(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// CompleteShoppingList 标记清单为完成
func (h *ShoppingListHandler) CompleteShoppingList(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}

	if err := h.service.CompleteShoppingList(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "已标记为完成"})
}
