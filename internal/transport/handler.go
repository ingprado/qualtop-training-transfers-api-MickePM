package transport

import (
	"encoding/json"
	"fmt"
	"net/http"

	"transfers-api/internal/enums"
	"transfers-api/internal/models"
	"transfers-api/internal/services"
)

type Handler struct {
	svc *services.TransferService
}

func NewHandler(s *services.TransferService) *Handler {
	return &Handler{svc: s}
}

type transferResponse struct {
	ID         string  `json:"ID"`
	SenderID   string  `json:"SenderID"`
	ReceiverID string  `json:"ReceiverID"`
	Currency   string  `json:"Currency"`
	Amount     float64 `json:"Amount"`
	State      string  `json:"State"`
}

func (h *Handler) HandleTransfers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		senderID := r.URL.Query().Get("sender_id")

		// Esta línea usa el paquete fmt. Si la borras, borra también el import.
		fmt.Printf("DEBUG: Buscando por id='%s' o sender_id='%s'\n", id, senderID)

		if id != "" {
			t, err := h.svc.GetByID(id)
			if err != nil {
				http.Error(w, "Transfer not found", http.StatusNotFound)
				return
			}
			res := transferResponse{
				ID: t.ID, SenderID: t.SenderID, ReceiverID: t.ReceiverID,
				Currency: t.Currency.String(), Amount: t.Amount, State: t.State,
			}
			json.NewEncoder(w).Encode(res)
			return // Se detiene aquí si encontró por ID
		}

		if senderID != "" {
			list, err := h.svc.GetBySenderID(senderID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var finalResponse []transferResponse
			for _, t := range list {
				finalResponse = append(finalResponse, transferResponse{
					ID: t.ID, SenderID: t.SenderID, ReceiverID: t.ReceiverID,
					Currency: t.Currency.String(), Amount: t.Amount, State: t.State,
				})
			}
			json.NewEncoder(w).Encode(finalResponse)
			return // Se detiene aquí si filtró por SenderID
		}

		// Si llegamos aquí, es que no hay parámetros. Traemos todos.
		list, err := h.svc.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var finalResponse []transferResponse
		for _, t := range list {
			finalResponse = append(finalResponse, transferResponse{
				ID: t.ID, SenderID: t.SenderID, ReceiverID: t.ReceiverID,
				Currency: t.Currency.String(), Amount: t.Amount, State: t.State,
			})
		}
		json.NewEncoder(w).Encode(finalResponse)

	case http.MethodPost:
		var req transferResponse
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		t := models.Transfer{
			ID: req.ID, SenderID: req.SenderID, ReceiverID: req.ReceiverID,
			Currency: enums.ParseCurrency(req.Currency), Amount: req.Amount, State: req.State,
		}
		if err := h.svc.Create(&t); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.Currency = t.Currency.String()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(req)

	case http.MethodPut:
		var req transferResponse
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		t := models.Transfer{
			ID: req.ID, SenderID: req.SenderID, ReceiverID: req.ReceiverID,
			Currency: enums.ParseCurrency(req.Currency), Amount: req.Amount, State: req.State,
		}
		if err := h.svc.Update(&t); err != nil {
			http.Error(w, "Could not update: "+err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(req)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Query parameter 'id' is required", http.StatusBadRequest)
			return
		}
		if err := h.svc.Delete(id); err != nil {
			http.Error(w, "Could not delete: "+err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
