package repositories

import (
	"database/sql" // Nuevo: para manejar la BD
	"fmt"
	"transfers-api/internal/models"

	_ "github.com/go-sql-driver/mysql" // Driver de MySQL
)

// MySQLRepository ahora contiene *sql.DB para conectarse de verdad.
type MySQLRepository struct {
	db *sql.DB // Cambiamos mu y store por la conexión real
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository { // Recibe la conexión desde el main
	return &MySQLRepository{
		db: db,
	}
}

func (r *MySQLRepository) Create(t *models.Transfer) error { //funcion para crear registro
	if t == nil || t.ID == "" { //valiudamos el id que no sea null
		return ErrInvalid // envia mensaje generico ya precargado
	}

	// Usamos Exec para insertar en la tabla real
	query := "INSERT INTO transfers (id, sender_id, receiver_id, currency, amount, state) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query, t.ID, t.SenderID, t.ReceiverID, int(t.Currency), t.Amount, t.State)
	if err != nil {
		return fmt.Errorf("error al insertar: %v", err)
	}

	return nil
}

func (r *MySQLRepository) GetByID(id string) (*models.Transfer, error) {
	query := "SELECT id, sender_id, receiver_id, currency, amount, state FROM transfers WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var t models.Transfer
	err := row.Scan(&t.ID, &t.SenderID, &t.ReceiverID, &t.Currency, &t.Amount, &t.State)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *MySQLRepository) Update(t *models.Transfer) error {
	query := "UPDATE transfers SET sender_id=?, receiver_id=?, currency=?, amount=?, state=? WHERE id=?"
	res, err := r.db.Exec(query, t.SenderID, t.ReceiverID, int(t.Currency), t.Amount, t.State, t.ID)
	if err != nil {
		return err
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *MySQLRepository) Delete(id string) error {
	query := "DELETE FROM transfers WHERE id = ?"
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *MySQLRepository) List() ([]*models.Transfer, error) {
	query := "SELECT id, sender_id, receiver_id, currency, amount, state FROM transfers"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Transfer
	for rows.Next() {
		var t models.Transfer
		if err := rows.Scan(&t.ID, &t.SenderID, &t.ReceiverID, &t.Currency, &t.Amount, &t.State); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	return out, nil
}

// METODO DE TRANSACCIONES POR USUARIO (Ahora con SQL)
func (r *MySQLRepository) GetBySenderID(senderID string) ([]*models.Transfer, error) {
	query := "SELECT id, sender_id, receiver_id, currency, amount, state FROM transfers WHERE sender_id = ?"
	rows, err := r.db.Query(query, senderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*models.Transfer, 0)
	for rows.Next() {
		var t models.Transfer
		if err := rows.Scan(&t.ID, &t.SenderID, &t.ReceiverID, &t.Currency, &t.Amount, &t.State); err != nil {
			return nil, err
		}
		results = append(results, &t)
	}

	return results, nil
}

// package repositories

// import (
// 	"sync"

// 	"transfers-api/internal/models"
// )

// // MySQLRepository implementa Repository. En el futuro contendrá *sql.DB.
// // Ahora usa un store en memoria para pruebas.
// type MySQLRepository struct {
// 	mu    sync.RWMutex                //para manejar muchas peticiones evitando duplicidad
// 	store map[string]*models.Transfer // diccionario del tipo models.Transfer donde esta definidos los tipos
// }

// func NewMySQLRepository() *MySQLRepository { //crea el repositorio como el main
// 	return &MySQLRepository{
// 		store: make(map[string]*models.Transfer), //inicialisa el diccionario vacio para usarlo
// 	}
// }

// func (r *MySQLRepository) Create(t *models.Transfer) error { //funcion para crear registro
// 	if t == nil || t.ID == "" { //valiudamos el id que no sea null
// 		return ErrInvalid // envia mensaje generico ya precargado
// 	}
// 	r.mu.Lock()         // cierra el candado para no leer ni escribir hasta que termine la operacion
// 	defer r.mu.Unlock() // en cuanto termine sea lo que sea abre de nuevo el candado

// 	if _, exists := r.store[t.ID]; exists { // checa que el id no exista
// 		return ErrInvalid // o crear ErrAlreadyExists si prefieres
// 	}

// 	r.store[t.ID] = &models.Transfer{
// 		ID:         t.ID,
// 		SenderID:   t.SenderID,
// 		ReceiverID: t.ReceiverID,
// 		Currency:   t.Currency,
// 		Amount:     t.Amount,
// 		State:      t.State,
// 	}

// 	return nil
// }

// func (r *MySQLRepository) GetByID(id string) (*models.Transfer, error) {
// 	r.mu.RLock()
// 	defer r.mu.RUnlock()
// 	if t, ok := r.store[id]; ok { // busca en el mapa temporal
// 		copy := *t //hace una copia del registro si encontro
// 		return &copy, nil
// 	}
// 	return nil, ErrNotFound
// }

// func (r *MySQLRepository) Update(t *models.Transfer) error {
// 	if t == nil || t.ID == "" {
// 		return ErrInvalid
// 	}
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	if _, ok := r.store[t.ID]; !ok { // busca el registro si no regresa error generico
// 		return ErrNotFound
// 	}
// 	// actualizar campos (mantén lógicas de negocio que necesites)
// 	r.store[t.ID] = &models.Transfer{
// 		ID:         t.ID,
// 		SenderID:   t.SenderID,
// 		ReceiverID: t.ReceiverID,
// 		Currency:   t.Currency,
// 		Amount:     t.Amount,
// 		State:      t.State,
// 	}
// 	return nil
// }

// func (r *MySQLRepository) Delete(id string) error {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	if _, ok := r.store[id]; !ok {
// 		return ErrNotFound
// 	}
// 	delete(r.store, id)
// 	return nil
// }

// func (r *MySQLRepository) List() ([]*models.Transfer, error) {
// 	r.mu.RLock()
// 	defer r.mu.RUnlock()
// 	out := make([]*models.Transfer, 0, len(r.store)) // crea una lista vacia con el tamaño exacto de la informqacion que hay
// 	for _, v := range r.store {                      //llena la tabla con la info que existe
// 		copy := *v
// 		out = append(out, &copy)
// 	}
// 	return out, nil
// }

// // METODO DE TRANSACCIONES POR USUARIO
// func (r *MySQLRepository) GetBySenderID(senderID string) ([]*models.Transfer, error) {
// 	r.mu.RLock()
// 	defer r.mu.RUnlock()

// 	results := make([]*models.Transfer, 0)

// 	for _, t := range r.store {
// 		if t.SenderID == senderID {
// 			copy := *t
// 			results = append(results, &copy)
// 		}
// 	}

// 	return results, nil
// }
