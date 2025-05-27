package sale

import (
	"errors"
	"math/rand"
	"parte3/internal/user" // <-- Importante
	"time"

	"github.com/google/uuid"
)

// Define errores específicos para ventas si es necesario
var ErrUserNotFound = errors.New("user not found for sale")
var ErrInvalidAmount = errors.New("sale amount must be positive")

// Service provides high-level user management operations on a LocalStorage backend.
type Service struct {
	// storage is the underlying persistence for User entities.
	salesStorage *LocalStorage // Para guardar ventas (¡usa el storage de ventas!)
	userService  *user.Service // Para validar usuarios
}

// NewService creates a new Service.
func NewService(salesStorage *LocalStorage, userService *user.Service) *Service {
	return &Service{
		salesStorage: salesStorage,
		userService:  userService,
	}
}

// Create adds a brand-new user to the system.
// It sets CreatedAt and UpdatedAt to the current time and initializes Version to 1.
// Returns ErrEmptyID if user.ID is empty.
func (s *Service) Create(userID string, amount float64) (*Sale, error) {
	// 1. Validar monto
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// 2. Validar que el user_id exista
	_, err := s.userService.Get(userID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) { // Comprueba si el error es 'user.ErrNotFound'
			return nil, ErrUserNotFound // Devuelve nuestro error específico
		}
		return nil, err // Devuelve otros errores (ej: problemas internos del servicio de usuario)
	}

	// 3. Asignar estado aleatorio
	statuses := []string{"pending", "approved", "rejected"}
	status := statuses[rand.Intn(len(statuses))] // Elige uno al azar

	// 4. Crear la venta
	now := time.Now()
	sale := &Sale{
		ID:        uuid.NewString(),
		UserID:    userID,
		Amount:    amount,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
		Version:   1,
		// No necesitas 'Estado' como en User, a menos que quieras un borrado lógico para ventas también.
	}

	// 5. Guardar la venta
	if err := s.salesStorage.Set(sale); err != nil {
		return nil, err // Devuelve error si falla el guardado
	}

	// 6. Devolver la venta creada
	return sale, nil
}
