package controller

import (
	"net/http"
	"restaurant-management-app/internal/service"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	PaymentService *service.PaymentService
}

func NewPaymentController(paymentService *service.PaymentService) *PaymentController {
	return &PaymentController{
		PaymentService: paymentService,
	}
}

type CreatePaymentRequest struct {
	Amount float32 `json:"amount" binding:"required,gt=0"`
}

// CreatePayment handles POST /bills/:id/payments.
func (pc *PaymentController) CreatePayment(c *gin.Context) {
	billID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req CreatePaymentRequest
	if !bindJSON(c, &req) {
		return
	}

	payment, err := pc.PaymentService.CreatePayment(billID, req.Amount)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, payment)
}

// GetPaymentsByBill handles GET /bills/:id/payments.
func (pc *PaymentController) GetPaymentsByBill(c *gin.Context) {
	billID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	payments, err := pc.PaymentService.GetPaymentsByBillID(billID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, payments)
}

// GetPaymentByID handles GET /payments/:id.
func (pc *PaymentController) GetPaymentByID(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	payment, err := pc.PaymentService.GetPaymentByID(id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, payment)
}
