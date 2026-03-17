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

func NewTableRepositroy(tableRepo *repository.TableRepository) *TableHandler {
	return &TableHandler{
		tableRepo: tableRepo,
	}
}

type seatRequest struct {
	Seats int `json:"seats" binding:"required,min=1"`
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

func (h *TableHandler) BookSeats(c *gin.Context) {
	tableId := c.Param("tableId")

	var req seatRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "seats must be a positive integer"})
		return
	}

	err := h.tableRepo.BookSeats(c.Request.Context(), tableId, req.Seats)
	if err != nil {
		errs.InternalServerError(c, "Failed to fetch tables")
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *TableHandler) ReleaseSeats(c *gin.Context) {
	tableId := c.Param("tableId")

	var req seatRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "seats must be a positive integer"})
		return
	}

	err := h.tableRepo.ReleaseSeats(c.Request.Context(), tableId, req.Seats)
	if err != nil {
		errs.InternalServerError(c, "Failed to fetch tables")
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
