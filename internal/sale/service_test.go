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

func TestService_Crear_ConUser_Inexistente(t *testing.T) {
	// fase de preparación (arrange)
	salesStorage := NewLocalStorage()
	userService := &mockUserService{}
	saleService := NewService(salesStorage, userService, nil)

	//paso donde se ejecuta la lógica (act)
	//se intenta crear una venta con un usuario que no existe y se espera que falle
	sale, err := saleService.Create("non-existent-user", 150.0)

	// validar que el código se comporta como debería (assert)
	require.Nil(t, sale)                     //no devuelve ninguna venta si el user no existe
	require.ErrorIs(t, err, ErrUserNotFound) //se verifica que verifica que el error devuelto por saleSvc.Create debe ser ErrUserNotFound (de ventas).
}
