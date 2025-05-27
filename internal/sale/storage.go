package sale

import "errors"

// ErrNotFound is returned when a user with the given ID is not found.
var ErrNotFound = errors.New("sale not found")

// ErrEmptyID is returned when trying to store a user with an empty ID.
var ErrEmptyID = errors.New("empty sale ID")

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

func (l *LocalStorage) ListActive() ([]*Sale, error) {
	var activeSales []*Sale
	for _, sale := range l.m {
		if sale.Status != "active" || sale.Status != "pending" || sale.Status != "approved" {
			activeSales = append(activeSales, sale)
		}
	}

	return activeSales, nil
}
