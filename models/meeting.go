package models

type Meeting struct {
	ID        int    `json:"id"`
	UID       int    `json:"uid"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	EventID   int    `json:"eventID,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}
