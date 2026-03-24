package models

import "transfers-api/internal/enums"

type Transfer struct {
	IDTransaction int            `json:"IDTransaction,omitempty"`
	ID            string         `json:"id"`
	SenderID      string         `json:"sender_id"`
	ReceiverID    string         `json:"receiver_id"`
	Currency      enums.Currency `json:"currency"`
	Amount        float64        `json:"amount"`
	State         string         `json:"state"`
}
