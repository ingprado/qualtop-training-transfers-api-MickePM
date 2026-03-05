package repositories

import "transfers-api/internal/models"

// Repository define operaciones CRUD para Transfer.
type Repository interface {
	Create(t *models.Transfer) error
	GetByID(id string) (*models.Transfer, error)
	Update(t *models.Transfer) error
	Delete(id string) error
	List() ([]*models.Transfer, error)
}
