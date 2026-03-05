package repositories

import "errors"

var (
	ErrNotFound = errors.New("registro no encontrado")
	ErrInvalid  = errors.New("datos inválidos")
)
