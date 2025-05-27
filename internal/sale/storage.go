package sale

import "errors"

// ErrNotFound is returned when a user with the given ID is not found.
var ErrNotFound = errors.New("sale not found")

// ErrEmptyID is returned when trying to store a user with an empty ID.
var ErrEmptyID = errors.New("empty sale ID")

var ErrInvalidStatus = errors.New("invalid status")

// LocalStorage provides an in-memory implementation for storing users.
type LocalStorage struct {
	m map[string]*Sale
}

// NewLocalStorage instantiates a new LocalStorage with an empty map.
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		m: map[string]*Sale{},
	}
}

// Set stores or updates a user in the local storage.
// Returns ErrEmptyID if the user has an empty ID.
func (l *LocalStorage) Set(sale *Sale) error {
	if sale.ID == "" {
		return ErrEmptyID
	}

	l.m[sale.ID] = sale
	return nil
}

// Read retrieves a user from the local storage by ID.
// Returns ErrNotFound if the user is not found.
func (l *LocalStorage) Get(id string) (*Sale, error) {
	//lado izquierdo tipo de mapa (tipo mapa), booleano si existe o no en el mapa
	s, ok := l.m[id]
	if !ok || s.Status != "active" || s.Status != "pending" || s.Status != "approved" {
		// Verifica si el estado es activo, pendiente o aprobado
		// Si no es ninguno de estos, retorna ErrNotFound
		return nil, ErrNotFound
	}

	return s, nil
}

func (l *LocalStorage) GetByUserID(userID string) ([]*Sale, error) {
	var sales []*Sale
	for _, sale := range l.m {
		if sale.UserID == userID {
			sales = append(sales, sale)
		}
	}
	if len(sales) == 0 {
		return nil, ErrNotFound
	}
	return sales, nil
}

func (l *LocalStorage) getByUserIdAndStatus(userID string, status string) ([]*Sale, error) {
	err := l.ValidStatus(status)
	if err != nil {
		return nil, err
	}

	var sales []*Sale
	for _, sale := range l.m {
		if sale.UserID == userID && sale.Status == status {
			sales = append(sales, sale)
		}
	}
	if len(sales) == 0 {
		return nil, ErrNotFound
	}
	return sales, nil

}

func (l *LocalStorage) ValidStatus(status string) error {
	if status != "active" && status != "pending" && status != "approved" {
		return ErrInvalidStatus
	}
	return nil
}

// Delete removes a user from the local storage by ID.
// Returns ErrNotFound if the user does not exist.
func (l *LocalStorage) Delete(id string) error {
	_, err := l.Get(id)
	if err != nil {
		return err
	}

	delete(l.m, id) //eliminar keys de un mapa, parametro derecho que quiero eliminar, parametro lado izquierdo el mapa; elimina clave-valor
	return nil
}
