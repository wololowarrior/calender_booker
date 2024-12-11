package dao

import (
	"database/sql"
	"errors"
	"fmt"

	"calendly_adventures/dao/query"
	"calendly_adventures/db"
	"calendly_adventures/models"
	"github.com/lib/pq"
)

func GetUser(userID int) (*models.User, error) {
	var user models.User
	err := db.DB.QueryRow(query.GetUser, userID).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return &user, nil
}

func Create(user *models.User) (*models.User, error) {
	err := db.DB.QueryRow(query.InsertUser, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// PostgreSQL error code for unique violation
			return nil, fmt.Errorf("user with ID %d already exists", user.ID)
		}
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}
	return user, nil
}

func Overview(userID int, date string) (*models.Overview, error) {
	rows, err := db.DB.Query(query.GetUnavailableSlots, userID, date)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	unavailableSlots := make([]*models.UnavailableSlot, 0)
	for rows.Next() {
		unavailableSlot := &models.UnavailableSlot{UID: userID, UnavailableDate: date}
		err = rows.Scan(&unavailableSlot.ID, &unavailableSlot.StartTime, &unavailableSlot.EndTime)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			} else {
				break
			}
		}
		*unavailableSlot.StartTime = timeParse(*unavailableSlot.StartTime)
		*unavailableSlot.EndTime = timeParse(*unavailableSlot.EndTime)
		unavailableSlots = append(unavailableSlots, unavailableSlot)
	}
	rows, err = db.DB.Query(query.GetMeetings, userID, date)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	meetings := make([]*models.Meeting, 0)
	for rows.Next() {
		meeting := &models.Meeting{UID: userID, Date: date}
		err = rows.Scan(&meeting.ID, &meeting.StartTime, &meeting.EndTime)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			} else {
				break
			}
		}
		meeting.StartTime = timeParse(meeting.StartTime)
		meeting.EndTime = timeParse(meeting.EndTime)
		meetings = append(meetings, meeting)
	}
	overview := models.Overview{
		UnavailableSlots: unavailableSlots,
		Meetings:         meetings,
	}
	return &overview, nil
}
