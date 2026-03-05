package models

import "transfers-api/internal/enums"

// type Transfer struct {
// 	ID         string
// 	SenderID   string
// 	ReceiverID string
// 	Currency   enums.Currency
// 	Amount     float64
// 	State      string // TODO: replace with enums.State
// }

type Transfer struct {
	ID         string         `json:"id"`
	SenderID   string         `json:"sender_id"`
	ReceiverID string         `json:"receiver_id"`
	Currency   enums.Currency `json:"currency"`
	Amount     float64        `json:"amount"`
	State      string         `json:"state"`
}
