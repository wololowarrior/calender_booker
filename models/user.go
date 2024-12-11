package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

var StartOfDay = "09:00:00"
var EndOfDay = "17:00:00"

type Overview struct {
	UnavailableSlots []*UnavailableSlot `json:"unavailable_slots"`
	Meetings         []*Meeting         `json:"meetings"`
}
