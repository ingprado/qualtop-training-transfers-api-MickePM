package repositories

import (
	"log"
	"time"
	"transfers-api/internal/models"

	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

type CachedTransferRepository struct {
	realRepo Repository
	cache    *cache.Cache
	logger   *zap.SugaredLogger
}

func NewCachedRepository(repo Repository, logger *zap.SugaredLogger) *CachedTransferRepository {
	return &CachedTransferRepository{
		realRepo: repo,
		cache:    cache.New(1*time.Minute, 2*time.Minute),
		logger:   logger,
	}
}

func (r *CachedTransferRepository) GetByID(id string) (*models.Transfer, error) {
	// busca en el Caché
	if val, found := r.cache.Get(id); found {
		// Mensaje cuando los datos están en RAM
		log.Printf("obteniendo transferencia desde la memoria")
		return val.(*models.Transfer), nil
	}

	//si el metodo anterior no lo encuentra, ahora vamos a la BD
	log.Printf("transferencia no encontrada en caché, consultando MySQL...")

	transfer, err := r.realRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Guardar en caché para la próxima vez
	r.cache.Set(id, transfer, cache.DefaultExpiration)
	return transfer, nil
}

// Guarda en DB y limpia el caché para evitar datos viejos
func (r *CachedTransferRepository) Create(t *models.Transfer) error {
	err := r.realRepo.Create(t)
	if err == nil {
		r.cache.Set(t.ID, t, cache.DefaultExpiration) // lo guardaMOS  en cache
	}
	return err
}

func (r *CachedTransferRepository) Update(t *models.Transfer) error   { return r.realRepo.Update(t) }
func (r *CachedTransferRepository) Delete(id string) error            { return r.realRepo.Delete(id) }
func (r *CachedTransferRepository) List() ([]*models.Transfer, error) { return r.realRepo.List() }
func (r *CachedTransferRepository) GetBySenderID(id string) ([]*models.Transfer, error) {
	return r.realRepo.GetBySenderID(id)
}
