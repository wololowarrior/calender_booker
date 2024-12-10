package dao

import (
	"calendly_adventures/dao/query"
	"calendly_adventures/db"
	"calendly_adventures/models"
)

func InsertEvent(event *models.Event) error {
	err := db.DB.QueryRow(query.InsertEvent, event.UID, event.Name, event.Message, event.Slotted, event.Slots).Scan(&event.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetAllEvents(uid int) ([]*models.Event, error) {
	var events []*models.Event
	rows, err := db.DB.Query(query.SelectEvent, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events = make([]*models.Event, 0)
	for rows.Next() {
		var event models.Event
		err = rows.Scan(&event.ID, &event.UID, &event.Name, &event.Message, &event.Slotted, &event.Slots)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}
