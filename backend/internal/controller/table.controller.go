package controller

import "github.com/gin-gonic/gin"

type TableController struct{}

func NewTableController() *TableController {
	return &TableController{}
}

func (t *TableController) AddNewTable(c *gin.Context) {

}
