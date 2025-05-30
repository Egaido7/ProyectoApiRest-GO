package user

import (
	"time"
)

// User represents a system user with metadata for auditing and versioning.
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" binding:"required,regexp`
	Address   string    `json:"address" binding:"required"` // Opcional, pero requerido si se proporciona
	NickName  string    `json:"nickname" binding:"omitempty,regexp`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
	Estado    bool      `json:"estado"` // Estado del usuario (activo/inactivo)
}

// UpdateFields represents the optional fields for updating a User.
// A nil pointer means “no change” for that field.
type UpdateFields struct {
	Name     *string `json:"name" binding:"required,regexp`      // Solo letras si se proporciona
	Address  *string `json:"address" binding:"required"`         // Opcional
	NickName *string `json:"nickname" binding:"omitempty,regexp` // Solo letras si se
}
