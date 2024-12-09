package models

type UnavailableSlot struct {
	ID              int     `json:"id"`
	UID             int     `json:"uid"`
	UnavailableDate string  `json:"unavailable_date"`
	StartTime       *string `json:"start_time"`
	EndTime         *string `json:"end_time"`
	CreatedAt       string  `json:"created_at"`
}

type UnavailableSlots struct {
	ID        int     `json:"id"`
	UID       int     `json:"uid"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
}
