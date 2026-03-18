package handler

import (
	"net/http"

	"github.com/distroaryan/restaurant-management/internal/errs"
	"github.com/distroaryan/restaurant-management/internal/repository"
	"github.com/gin-gonic/gin"
)

type TableHandler struct {
	tableRepo *repository.TableRepository
}

func NewTableRepositroy(repository *repository.Repository) *TableHandler {
	return &TableHandler{
		tableRepo: repository.Table,
	}
}



func (h *TableHandler) GetAllTables(c *gin.Context) {
	tables, err := h.tableRepo.GetAllTables(c.Request.Context())
	if err != nil {
		errs.InternalServerError(c, "Failed to fetch tables")
		return
	}
	c.JSON(http.StatusOK, tables)
}

func (h *TableHandler) GetTableById(c *gin.Context) {
	id := c.Param("tableId")

	table, err := h.tableRepo.GetTableById(c.Request.Context(), id)
	if err != nil {
		errs.InternalServerError(c, "Failed to fetch tables")
		return
	}
	c.JSON(http.StatusOK, table)
}

func (h *TableHandler) BookTable(c *gin.Context) {
	tableId := c.Param("tableId")
	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Valid user authentication required"})
		return
	}

	err := h.tableRepo.BookTable(c.Request.Context(), tableId, userId)
	if err != nil {
		errs.InternalServerError(c, "Failed to book table")
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *TableHandler) ReleaseTable(c *gin.Context) {
	tableId := c.Param("tableId")

	err := h.tableRepo.ReleaseTable(c.Request.Context(), tableId)
	if err != nil {
		errs.InternalServerError(c, "Failed to release table")
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
