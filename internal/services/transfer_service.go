package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"transfers-api/internal/models"
	"transfers-api/internal/queue"
	"transfers-api/internal/repositories"
)

var ErrInvalidAmount = errors.New("monto inválido")

type TransferService struct {
	repo     repositories.Repository
	producer *queue.RabbitMQProducer
}

func NewTransferService(repo repositories.Repository, producer *queue.RabbitMQProducer) *TransferService {
	return &TransferService{
		repo:     repo,
		producer: producer,
	}
}
func (s *TransferService) Create(t *models.Transfer) error {
	err := s.repo.Create(t)
	if err != nil {
		return err
	}

	if s.producer != nil {
		body, _ := json.Marshal(t)
		if pubErr := s.producer.Publish(body); pubErr != nil {
			fmt.Printf("Advertencia: No se pudo enviar a RabbitMQ: %v\n", pubErr)
		}
	} else {
		fmt.Println("Advertencia: El servicio no tiene un productor de RabbitMQ configurado")
	}

	return nil
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

func (s *TransferService) GetBySenderID(senderID string) ([]*models.Transfer, error) {

	return s.repo.GetBySenderID(senderID)
}
