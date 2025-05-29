package user

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Getter interface {
	Get(id string) (*User, error)
}

// Service provides high-level user management operations on a LocalStorage backend.
type Service struct {
	// storage is the underlying persistence for User entities.
	storage *LocalStorage
	logger  *zap.Logger
}

func NewService(storage *LocalStorage, logger *zap.Logger) *Service {
	if logger == nil {
		logger, _ = zap.NewProduction()
		defer logger.Sync() // flushes buffer, if any
	}

	return &Service{
		storage: storage,
		logger:  logger,
	}
}

// Create adds a brand-new user to the system.
// It sets CreatedAt and UpdatedAt to the current time and initializes Version to 1.
// Returns ErrEmptyID if user.ID is empty.
func (s *Service) Create(user *User) error {
	user.ID = uuid.NewString()
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.Version = 1
	user.Estado = true

	if err := s.storage.Set(user); err != nil {
		s.logger.Error("failed to set user", zap.Error(err), zap.Any("user", user))
		return err
	}

	return s.storage.Set(user)
}

// Get retrieves a user by its ID.
// Returns ErrNotFound if no user exists with the given ID.
func (s *Service) Get(id string) (*User, error) {
	return s.storage.Get(id)
}

// Update modifies an existing user's data.
// It updates Name, Address, NickName, sets UpdatedAt to now and increments Version.
// Returns ErrNotFound if the user does not exist, or ErrEmptyID if user.ID is empty.
func (s *Service) Update(id string, user *UpdateFields, user2 User) (*User, error) {
	existing, err := s.storage.Get(id)
	if err != nil {
		return nil, err
	}
	if !user2.Estado {
		if user.Name != nil {
			existing.Name = *user.Name
		}

		if user.Address != nil {
			existing.Address = *user.Address
		}

		if user.NickName != nil {
			existing.NickName = *user.NickName
		}

		existing.UpdatedAt = time.Now()
		existing.Version++

		if err := s.storage.Set(existing); err != nil {
			return nil, err
		}
	}
	return existing, nil
}

// Delete removes a user from the system by its ID.
// Returns ErrNotFound if the user does not exist.
func (s *Service) Delete(id string) error {
	// Obtener el usuario existente
	existing, err := s.storage.Get(id)
	if err != nil {
		return err // Retorna ErrNotFound si no existe
	}

	// Cambiar el estado del usuario a false (borrado lógico)
	existing.Estado = false
	existing.UpdatedAt = time.Now() // Actualizar la fecha de modificación

	// Guardar los cambios en el almacenamiento
	return s.storage.Set(existing)
}
func (s *Service) ListActive() ([]*User, error) {
	return s.storage.ListActive()
}
