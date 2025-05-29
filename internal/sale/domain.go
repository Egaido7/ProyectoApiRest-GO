package sale

import (
	"time"
)

// User represents a system user with metadata for auditing and versioning.
type Sale struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
	Status    string    `json:"status"` // Estado de la venta (pending, approved, rejected)
}

type CreateSaleRequest struct {
	UserID string  `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"` // Monto de la venta
}

type GetSalesRequest struct {
}

// UpdateFields represents the optional fields for updating a User.
// A nil pointer means “no change” for that field.
type UpdateSale struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}

type Metadata struct {
	Quantity    int     `json:"quantity"`
	Approved    int     `json:"approved"`
	Rejected    int     `json:"rejected"`
	Pending     int     `json:"pending"`
	TotalAmount float64 `json:"total_amount"`
}
