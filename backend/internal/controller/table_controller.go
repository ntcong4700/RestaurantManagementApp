package controller

import (
	"net/http"
	"restaurant-management-app/internal/models"
	"restaurant-management-app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TableController struct {
	TableService *service.TableService
}

func NewTableController(tableService *service.TableService) *TableController {
	return &TableController{
		TableService: tableService,
	}
}

type CreateTableRequest struct {
	TableNumber int `json:"table_number" binding:"required"`
}

func (tc *TableController) CreateTable(c *gin.Context) {
	var req CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	table := models.Tables{
		Status:      models.TableStatusAvailable,
		TableNumber: req.TableNumber,
		ID:          uuid.New(),
	}
	createdTable, err := tc.TableService.CreateTable(table)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdTable)
}

func (tc *TableController) GetAllTables(c *gin.Context) {
	tables, err := tc.TableService.GetAllTable()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tables)
}

type UpdateTableRequest struct {
	TableNumber int    `json:"table_number" binding:"required"`
	Status      string `json:"status" binding:"required"`
}

func (tc *TableController) UpdateTable(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table id"})
		return
	}

	var req UpdateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	table := models.Tables{
		ID:          id,
		TableNumber: req.TableNumber,
		Status:      req.Status,
	}
	updatedTable, err := tc.TableService.UpdateTables(table)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedTable)
}

func (tc *TableController) DeleteTable(c *gin.Context) {
	id := c.Param("id")
	if err := tc.TableService.DeleteTable(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
