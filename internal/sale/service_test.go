// go
package sale

import (
	"testing"

	"parte3/internal/user"

	"github.com/stretchr/testify/require"
)

// Mock de user.Service que siempre devuelve ErrNotFound
type mockUserService struct{}

func (m *mockUserService) Get(id string) (*user.User, error) {
	return nil, user.ErrNotFound
}

func TestService_Create_SaleWithNonExistentUser(t *testing.T) {
	// Arrange
	salesStorage := NewLocalStorage()
	userSvc := &mockUserService{}
	saleSvc := NewService(salesStorage, userSvc, nil)

	// Act
	sale, err := saleSvc.Create("non-existent-user", 100.0)

	// Assert
	require.Nil(t, sale)                     //no devuelve ninguna venta si el user no exite
	require.ErrorIs(t, err, ErrUserNotFound) //se verifica que verifica que el error devuelto por saleSvc.Create debe ser ErrUserNotFound (de ventas).
}
