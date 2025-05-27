package sale

import (
	"time"
)

// User represents a system user with metadata for auditing and versioning.
type Sale struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userid"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
	Status    string    `json:"status"` // Estado de la venta (pending, approved, rejected)
}

type CreateSaleRequest struct {
	UserID string  `json:"userid" binding:"required"`
	Amount float64 `json:"amount" binding:"required, gt=0"` // Monto de la venta
}

// UpdateFields represents the optional fields for updating a User.
// A nil pointer means “no change” for that field.
type UpdateSale struct {
	Name     *string `json:"name" binding:"required,regexp"` // Solo letras si se proporciona
	Address  *string `json:"address" binding:"required"`     // Opcional
	NickName *string `json:"nickname" binding:"regexp"`      // Solo letras si se
}
