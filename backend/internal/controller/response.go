package controller

import (
	"errors"
	"net/http"
	"restaurant-management-app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var notFoundErrors = []error{
	service.ErrTableNotFound,
	service.ErrMenuItemNotFound,
	service.ErrOrderNotFound,
	service.ErrOrderItemNotFound,
	service.ErrBillNotFound,
	service.ErrPaymentNotFound,
}

var badRequestErrors = []error{
	service.ErrInvalidQuantity,
	service.ErrInvalidAmount,
	service.ErrInvalidTableStatus,
	service.ErrDuplicateOrderItems,
	service.ErrItemNotInOrder,
}

var conflictErrors = []error{
	service.ErrTableHasOpenOrder,
	service.ErrOrderNotOpen,
	service.ErrOrderItemBilled,
	service.ErrOrderHasItems,
	service.ErrBillHasPayments,
	service.ErrBillAlreadyPaid,
	service.ErrOverpayment,
	service.ErrNoItemsToBill,
	service.ErrTableStatusManaged,
}

// respondError maps a service error to an HTTP status. Unknown errors become a
// generic 500; the real error is attached to the gin context so the logger prints it.
func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrDuplicatedKey):
		c.JSON(http.StatusConflict, gin.H{"error": "resource already exists"})
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		c.JSON(http.StatusConflict, gin.H{"error": "resource is still referenced by other records"})
	case isAny(err, notFoundErrors):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case isAny(err, badRequestErrors):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case isAny(err, conflictErrors):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func isAny(err error, targets []error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// bindJSON binds the request body and writes a 400 on failure.
func bindJSON(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}

// parseUUIDParam reads a UUID path parameter and writes a 400 if it is invalid.
func parseUUIDParam(c *gin.Context, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return uuid.Nil, false
	}
	return id, true
}
