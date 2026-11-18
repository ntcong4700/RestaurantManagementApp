package controller

import (
	"errors"
	"io"
	"net/http"
	"restaurant-management-app/internal/models"
	"restaurant-management-app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BillController struct {
	BillService *service.BillService
}

func NewBillController(billService *service.BillService) *BillController {
	return &BillController{
		BillService: billService,
	}
}

// BillResponse is a bill plus how much is still unpaid.
type BillResponse struct {
	*models.Bill
	PaidAmount      float32 `json:"paid_amount"`
	RemainingAmount float32 `json:"remaining_amount"`
}

func newBillResponse(bill *models.Bill) BillResponse {
	var paid float32
	for _, p := range bill.Payments {
		paid += p.Amount
	}
	return BillResponse{
		Bill:            bill,
		PaidAmount:      paid,
		RemainingAmount: bill.TotalAmount - paid,
	}
}

// CreateBillRequest: leave order_item_ids empty to bill every unbilled item of
// the order, or list specific items to split the bill.
type CreateBillRequest struct {
	OrderItemIDs []uuid.UUID `json:"order_item_ids"`
}

// CreateBill handles POST /orders/:id/bills. The body is optional.
func (bc *BillController) CreateBill(c *gin.Context) {
	orderID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req CreateBillRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bill, err := bc.BillService.CreateBill(orderID, req.OrderItemIDs)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, newBillResponse(bill))
}

// GetBillsByOrder handles GET /orders/:id/bills.
func (bc *BillController) GetBillsByOrder(c *gin.Context) {
	orderID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	bills, err := bc.BillService.GetBillsByOrderID(orderID)
	if err != nil {
		respondError(c, err)
		return
	}

	resp := make([]BillResponse, 0, len(bills))
	for i := range bills {
		resp = append(resp, newBillResponse(&bills[i]))
	}
	c.JSON(http.StatusOK, resp)
}

// GetBillByID handles GET /bills/:id.
func (bc *BillController) GetBillByID(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	bill, err := bc.BillService.GetBillByID(id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, newBillResponse(bill))
}

// DeleteBill handles DELETE /bills/:id (only bills without payments).
func (bc *BillController) DeleteBill(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := bc.BillService.DeleteBill(id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
