package models

type Meeting struct {
	ID            int           `json:"id"`
	UID           int           `json:"uid"`
	Date          string        `json:"date"`
	StartTime     string        `json:"start_time"`
	EndTime       string        `json:"end_time"`
	EventID       int           `json:"event_id,omitempty"`
	CreatedAt     string        `json:"createdAt,omitempty"`
	BookerDetails BookerDetails `json:"bookerDetails,omitempty"`
}

type SlottedMeeting struct {
	ID      int       `json:"id"`
	UID     int       `json:"uid"`
	Date    string    `json:"date"`
	Time    [2]string `json:"time"`
	EventID int       `json:"event_id,omitempty"`
}

type BookerDetails struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
