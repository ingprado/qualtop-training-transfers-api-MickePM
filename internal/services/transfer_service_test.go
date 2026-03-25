package services

import (
	"testing"
	"transfers-api/internal/models"
	repoMocks "transfers-api/internal/repositories/mocks" // Asegúrate de que la ruta coincida

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransferService_Create_Success(t *testing.T) {
	// 1. Setup: Instanciamos el mock del repositorio
	mockRepo := new(repoMocks.Repository)

	// 2. Definimos qué esperamos que pase:
	// "Esperamos que se llame a Create con cualquier puntero a Transfer y devuelva nil"
	mockRepo.On("Create", mock.AnythingOfType("*models.Transfer")).Return(nil)

	// 3. Creamos el servicio inyectando el mock
	// Nota: Si aún no tienes interface para el producer, pásale nil por ahora
	svc := NewTransferService(mockRepo, nil)

	// 4. Ejecutamos la acción
	transfer := &models.Transfer{
		ID:         "TX-0014",
		SenderID:   "MIGUEL-123",
		ReceiverID: "LUIS-456",
		Amount:     50.50,
		State:      "PENDING",
	}
	err := svc.Create(transfer)

	// 5. Verificaciones (Asserts)
	assert.NoError(t, err)         // Verificamos que no haya error
	mockRepo.AssertExpectations(t) // Verificamos que Create se llamó como dijimos
}
func TestTransferService_Create_RepositoryError(t *testing.T) {
	mockRepo := new(repoMocks.Repository)

	// Simulamos que el repositorio falla
	mockRepo.On("Create", mock.Anything).Return(assert.AnError)

	svc := NewTransferService(mockRepo, nil)
	err := svc.Create(&models.Transfer{Amount: 100})

	// Verificamos que el servicio propaga el error
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
func TestTransferService_Create_InvalidAmount(t *testing.T) {
	// 1. Setup
	mockRepo := new(repoMocks.Repository)
	svc := NewTransferService(mockRepo, nil)

	// 2. Ejecución con un monto inválido
	transfer := &models.Transfer{
		ID:     "TX-FAIL",
		Amount: -500.0, // <--- Monto negativo
	}
	err := svc.Create(transfer)

	// 3. Verificaciones
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "monto")

	// CRÍTICO: Verificamos que el repositorio NUNCA fue llamado.
	// Si el monto es malo, no debemos tocar la base de datos.
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}
