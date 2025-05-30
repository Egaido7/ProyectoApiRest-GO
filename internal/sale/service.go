package sale

import (
	"errors"
	"math/rand"
	"parte3/internal/user" // <-- Importante
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Define errores específicos para ventas si es necesario
var ErrUserNotFound = errors.New("user not found for sale")
var ErrInvalidAmount = errors.New("sale amount must be positive")
var ErrSaleNotActive = errors.New("sale is not active and cannot be updated")
var ErrInvalidSaleStateTransition = errors.New("invalid state transition for sale")
var ErrSaleMustBePending = errors.New("sale status must be pending to be updated")

// Service provides high-level user management operations on a LocalStorage backend.
type Service struct {
	// storage is the underlying persistence for User entities.
	salesStorage *LocalStorage // Para guardar ventas (¡usa el storage de ventas!)
	userService  user.Getter   // Para validar usuarios
	logger       *zap.Logger
}

// NewService creates a new Service.
func NewService(salesStorage *LocalStorage, userService user.Getter, logger *zap.Logger) *Service {
	if logger == nil {
		logger, _ = zap.NewProduction()
		defer logger.Sync() // flushes buffer, if any
	}
	return &Service{
		salesStorage: salesStorage,
		userService:  userService,
		logger:       logger,
	}
}

// Create adds a brand-new user to the system.
// It sets CreatedAt and UpdatedAt to the current time and initializes Version to 1.
// Returns ErrEmptyID if user.ID is empty.
func (s *Service) Create(userID string, amount float64) (*Sale, error) {

	// 1. Validar que el user_id exista
	_, err := s.userService.Get(userID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) { // Comprueba si el error es 'user.ErrNotFound'
			s.logger.Warn("user not found for sale", zap.String("userID", userID))
			return nil, ErrUserNotFound // Devuelve nuestro error específico
		}
		return nil, err // Devuelve otros errores (ej: problemas internos del servicio de usuario)
	}

	// 2. Validar monto
	if amount <= 0 {
		s.logger.Warn("invalid sale amount", zap.Float64("amount", amount))
		return nil, ErrInvalidAmount
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
	}

	// 5. Guardar la venta
	if err := s.salesStorage.Set(sale); err != nil {
		s.logger.Error("failed to save sale", zap.Error(err), zap.Any("sale", sale))
		return nil, err // Devuelve error si falla el guardado
	}

	// 6. Devolver la venta creada
	return sale, nil
}

func (s *Service) Get(userID string) ([]*Sale, *Metadata, error) {
	sales, err := s.salesStorage.GetByUserID(userID)
	if err != nil {
		return nil, nil, err
	}
	meta, err := s.salesStorage.FillMetadata(sales)
	if err != nil {
		return nil, nil, err
	}
	return sales, meta, nil
}

func (s *Service) GetByStatus(userID string, status *string) ([]*Sale, *Metadata, error) {
	err := s.salesStorage.ValidStatus(*status)
	if err != nil {
		return nil, nil, err
	}

	sales, err := s.salesStorage.getByUserIdAndStatus(userID, *status)
	if err != nil {
		return nil, nil, err
	}
	meta, err := s.salesStorage.FillMetadata(sales)
	if err != nil {
		return nil, nil, err
	}
	return sales, meta, nil
}

func (s *Service) Update(saleID string, status string) (*Sale, error) {
	// 1. Validar que la venta exista
	sale, err := s.salesStorage.GetForUpdate(saleID) // Asumiendo que tienes GetForUpdate como discutimos
	if err != nil {
		if errors.Is(err, ErrNotFound) { // ErrNotFound de sales.storage
			s.logger.Warn("sale not found for update", zap.String("saleID", saleID))
			return nil, ErrNotFound
		}
		return nil, err // Otro error del storage
	}

	if sale.Version == 0 {
		s.logger.Warn("sale not active for update", zap.String("saleID", saleID))
		return nil, ErrSaleNotActive // devuelve error si la venta no está activa
	}

	if sale.Status != "pending" {
		s.logger.Warn("sale must be pending for status update", zap.String("saleID", saleID), zap.String("current_status", sale.Status))
		return nil, ErrSaleMustBePending // Devuelve error si el estado no es válido
	}

	if status != "approved" && status != "rejected" {
		s.logger.Warn("invalid sale state transition", zap.String("saleID", saleID), zap.String("new_status", status))
		return nil, ErrInvalidSaleStateTransition // Devuelve error si el estado no es válido
	}

	// 2. Actualizar estado
	sale.Status = status
	sale.UpdatedAt = time.Now()
	sale.Version++

	// 5. Guardar la venta
	if err := s.salesStorage.Set(sale); err != nil {
		s.logger.Error("failed to update sale status", zap.Error(err), zap.Any("sale", sale))
		return nil, err // Devuelve error si falla el guardado
	}

	// 6. Devolver la venta creada
	return sale, nil

}
