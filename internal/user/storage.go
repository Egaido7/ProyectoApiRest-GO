package user

import "errors"

// ErrNotFound is returned when a user with the given ID is not found.
var ErrNotFound = errors.New("user not found")

// ErrEmptyID is returned when trying to store a user with an empty ID.
var ErrEmptyID = errors.New("empty user ID")

// LocalStorage provides an in-memory implementation for storing users.
type LocalStorage struct {
	m map[string]*User
}

// NewLocalStorage instantiates a new LocalStorage with an empty map.
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		m: map[string]*User{},
	}
}

// Set stores or updates a user in the local storage.
// Returns ErrEmptyID if the user has an empty ID.
func (l *LocalStorage) Set(user *User) error {
	if user.ID == "" {
		return ErrEmptyID
	}

	l.m[user.ID] = user
	return nil
}

// Read retrieves a user from the local storage by ID.
// Returns ErrNotFound if the user is not found.
func (l *LocalStorage) Get(id string) (*User, error) {
	//lado izquierdo tipo de mapa (tipo mapa), booleano si existe o no en el mapa
	u, ok := l.m[id]
	if !ok || !u.Estado {
		return nil, ErrNotFound
	}

	return u, nil
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

func (l *LocalStorage) ListActive() ([]*User, error) {
	var activeUsers []*User
	for _, user := range l.m {
		if user.Estado {
			activeUsers = append(activeUsers, user)
		}
	}

	return activeUsers, nil
}
