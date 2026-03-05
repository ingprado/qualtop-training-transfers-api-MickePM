package repositories

import (
	"sync"

	"transfers-api/internal/models"
)

// MySQLRepository implementa Repository. En el futuro contendrá *sql.DB.
// Ahora usa un store en memoria para pruebas.
type MySQLRepository struct {
	mu    sync.RWMutex                //para manejar muchas peticiones evitando duplicidad
	store map[string]*models.Transfer // diccionario del tipo models.Transfer donde esta definidos los tipos
}

func NewMySQLRepository() *MySQLRepository { //crea el repositorio como el main
	return &MySQLRepository{
		store: make(map[string]*models.Transfer), //inicialisa el diccionario vacio para usarlo
	}
}

func (r *MySQLRepository) Create(t *models.Transfer) error { //funcion para crear registro
	if t == nil || t.ID == "" { //valiudamos el id que no sea null
		return ErrInvalid // envia mensaje generico ya precargado
	}
	r.mu.Lock()         // cierra el candado para no leer ni escribir hasta que termine la operacion
	defer r.mu.Unlock() // en cuanto termine sea lo que sea abre de nuevo el candado

	if _, exists := r.store[t.ID]; exists { // checa que el id no exista
		return ErrInvalid // o crear ErrAlreadyExists si prefieres
	}

	r.store[t.ID] = &models.Transfer{
		ID:         t.ID,
		SenderID:   t.SenderID,
		ReceiverID: t.ReceiverID,
		Currency:   t.Currency,
		Amount:     t.Amount,
		State:      t.State,
	}

	return nil
}

func (r *MySQLRepository) GetByID(id string) (*models.Transfer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if t, ok := r.store[id]; ok { // busca en el mapa temporal
		copy := *t //hace una copia del registro si encontro
		return &copy, nil
	}
	return nil, ErrNotFound
}

func (r *MySQLRepository) Update(t *models.Transfer) error {
	if t == nil || t.ID == "" {
		return ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[t.ID]; !ok { // busca el registro si no regresa error generico
		return ErrNotFound
	}
	// actualizar campos (mantén lógicas de negocio que necesites)
	r.store[t.ID] = &models.Transfer{
		ID:         t.ID,
		SenderID:   t.SenderID,
		ReceiverID: t.ReceiverID,
		Currency:   t.Currency,
		Amount:     t.Amount,
		State:      t.State,
	}
	return nil
}

func (r *MySQLRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return ErrNotFound
	}
	delete(r.store, id)
	return nil
}

func (r *MySQLRepository) List() ([]*models.Transfer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*models.Transfer, 0, len(r.store)) // crea una lista vacia con el tamaño exacto de la informqacion que hay
	for _, v := range r.store {                      //llena la tabla con la info que existe
		copy := *v
		out = append(out, &copy)
	}
	return out, nil
}
