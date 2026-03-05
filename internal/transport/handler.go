package transport

import (
	"encoding/json"
	"net/http"

	"transfers-api/internal/models"
	"transfers-api/internal/services"
)

type Handler struct {
	svc *services.TransferService
}

func NewHandler(s *services.TransferService) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) HandleTransfers(w http.ResponseWriter, r *http.Request) {
	// Establecemos el Header para que Postman sepa que recibe JSON
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// 1. Verificamos si viene un ID en la URL: ?id=TX-001
		id := r.URL.Query().Get("id")
		if id != "" {
			t, err := h.svc.GetByID(id)
			if err != nil {
				http.Error(w, "Transfer not found", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(t)
			return
		}

		// 2. Si no hay ID, listamos todos los registros
		list, err := h.svc.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(list)

	case http.MethodPost:
		// Crear nueva transferencia
		var t models.Transfer
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			// Enviamos el error real para saber si falló el JSON o el Enum
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.svc.Create(&t); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)

	case http.MethodPut:
		// Actualizar transferencia existente
		var t models.Transfer
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if err := h.svc.Update(&t); err != nil {
			http.Error(w, "Could not update: "+err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(t)

	case http.MethodDelete:
		// Eliminar transferencia por ID: ?id=TX-001
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Query parameter 'id' is required", http.StatusBadRequest)
			return
		}
		if err := h.svc.Delete(id); err != nil {
			http.Error(w, "Could not delete: "+err.Error(), http.StatusNotFound)
			return
		}
		// 204 No Content es el estándar para deletes exitosos
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
