package models

type Event struct {
	ID      int     `json:"id"`
	UID     int     `json:"uid"`
	Name    string  `json:"name"`
	Message string  `json:"message"`
	Slots   *string `json:"slots,omitempty"`
}
