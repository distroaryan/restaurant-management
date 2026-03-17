package handler

import (
	"net/http"

	"github.com/distroaryan/restaurant-management/internal/repository"
	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuRepo *repository.MenuRepository
}

func NewMenuHandler(menuRepo *repository.MenuRepository) *MenuHandler {
	return &MenuHandler{
		menuRepo: menuRepo,
	}
}

func (h *MenuHandler) GetAllMenus(c *gin.Context) {
	menus, err := h.menuRepo.GetAllMenu(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch menu",
		})
	}
	c.JSON(http.StatusOK, menus)
}