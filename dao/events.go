package dao

import (
	"database/sql"
	"errors"

	"calendly_adventures/dao/query"
	"calendly_adventures/db"
	"calendly_adventures/models"
)

func InsertEvent(event *models.Event) error {
	err := db.DB.QueryRow(query.InsertEvent, event.UID, event.Name, event.Message, event.Slots).Scan(&event.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetAllEvents(uid int) ([]*models.Event, error) {
	var events []*models.Event
	rows, err := db.DB.Query(query.SelectEventWithUserID, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events = make([]*models.Event, 0)
	for rows.Next() {
		var event models.Event
		err = rows.Scan(&event.ID, &event.UID, &event.Name, &event.Message, &event.Slots)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

func GetEvent(id int) (*models.Event, error) {
	var event models.Event
	err := db.DB.QueryRow(query.SelectEvent, id).Scan(&event.UID, &event.Name, &event.Message, &event.Slots)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("sql: no rows in result set")
		}
		return nil, err
	}
	return &event, nil
}

func DeleteEvent(id int, userID int) error {
	_, err := db.DB.Exec(query.DeleteEvent, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("event doesn't exist")
		}
		return err
	}
	return nil
}
