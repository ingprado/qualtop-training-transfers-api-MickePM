package services

import (
	"errors"

	"transfers-api/internal/models"
	"transfers-api/internal/repositories"
)

var ErrInvalidAmount = errors.New("monto inválido")

type TransferService struct {
	repo repositories.Repository
}

func NewTransferService(r repositories.Repository) *TransferService {
	return &TransferService{repo: r}
}

func (s *TransferService) Create(t *models.Transfer) error {
	if t == nil {
		return errors.New("transfer es nil")
	}
	if t.Amount <= 0 {
		return ErrInvalidAmount
	}
	return s.repo.Create(t)
}

func (s *TransferService) GetByID(id string) (*models.Transfer, error) {
	return s.repo.GetByID(id)
}

func (s *TransferService) Update(t *models.Transfer) error {
	if t == nil {
		return errors.New("transfer es nil")
	}
	return s.repo.Update(t)
}

func (s *TransferService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *TransferService) List() ([]*models.Transfer, error) {
	return s.repo.List()
}
