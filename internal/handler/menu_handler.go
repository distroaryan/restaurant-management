package handler

import (
	"net/http"

	"github.com/distroaryan/restaurant-management/internal/errs"
	"github.com/distroaryan/restaurant-management/internal/repository"
	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuRepo *repository.MenuRepository
}

func NewMenuHandler(repository *repository.Repository) *MenuHandler {
	return &MenuHandler{
		menuRepo: repository.Menu,
	}
}

func (h *MenuHandler) GetAllMenus(c *gin.Context) {
	menus, err := h.menuRepo.GetAllMenu(c.Request.Context())
	if err != nil {
		errs.InternalServerError(c, "Failed to fetch menu")
		return
	}
	c.JSON(http.StatusOK, menus)
}