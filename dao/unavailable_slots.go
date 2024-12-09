package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"calendly_adventures/dao/query"
	"calendly_adventures/db"
	"calendly_adventures/models"
)

func CreateUnavailableSlots(slot *models.UnavailableSlot) error {
	toRemove, toAdd, err2 := getOverlappingSlots(slot.UID, slot.UnavailableDate, *slot.StartTime, *slot.EndTime)
	if err2 != nil {
		return err2
	}

	for _, unavailableSlot := range toRemove {
		log.Printf("Removing unavailable slot: %s", unavailableSlot.ID)
		err := deleteUnavailableSlots(unavailableSlot.ID)
		if err != nil {
			return err
		}
	}
	for _, unavailableSlot := range toAdd {
		log.Printf("adding")
		log.Println(unavailableSlot.UID, unavailableSlot.UnavailableDate, unavailableSlot.StartTime, unavailableSlot.EndTime)
		err := db.DB.QueryRow(query.InsertUnavailableSlots, slot.UID, slot.UnavailableDate, unavailableSlot.StartTime, unavailableSlot.EndTime).Scan(&slot.ID, &slot.CreatedAt)
		if err != nil {
			return err
		}
	}
	if toAdd == nil {
		log.Println(slot.UID, slot.UnavailableDate, slot.StartTime, slot.EndTime)
		err := db.DB.QueryRow(query.InsertUnavailableSlots, slot.UID, slot.UnavailableDate, slot.StartTime, slot.EndTime).Scan(&slot.ID, &slot.CreatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetUnavailableSlots(uid int) ([]*models.UnavailableSlot, error) {
	rows, err := db.DB.Query(query.SelectUnavailableSlots, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var slots []*models.UnavailableSlot
	for rows.Next() {
		slot := new(models.UnavailableSlot)
		var startTime, endTime sql.NullString
		err = rows.Scan(&slot.ID, &slot.UID, &slot.UnavailableDate, &startTime, &endTime, &slot.CreatedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("unavailability with uid %d not found", uid)
			}
			return nil, fmt.Errorf("failed to retrieve unavailability slots: %w", err)
		}
		if slot.UnavailableDate != "" {
			t, err := time.Parse("2006-01-02T15:04:05Z", slot.UnavailableDate)
			if err != nil {
				return nil, err
			}
			t1 := t.Format("2006-01-02")
			slot.UnavailableDate = t1
		}
		if startTime.Valid {
			t, err := time.Parse("2006-01-02T15:04:05Z", startTime.String)
			if err != nil {
				return nil, err
			}
			t1 := t.Format("15:04:05")
			slot.StartTime = &t1
		}
		if endTime.Valid {
			t, err := time.Parse("2006-01-02T15:04:05Z", endTime.String)
			if err != nil {
				return nil, err
			}
			t1 := t.Format("15:04:05")
			slot.EndTime = &t1
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func deleteUnavailableSlots(id int) error {
	err := db.DB.QueryRow(query.DeleteUnavailableSlots, id)
	if err != nil {
		return err.Err()
	}
	return nil
}

func getOverlappingSlots(uid int, date, startTime, endTime string) ([]*models.UnavailableSlot, []*models.UnavailableSlot, error) {
	// Parse the input start and end times
	inputStartTime, err := time.Parse("15:04:05", startTime)
	if err != nil {
		return nil, nil, errors.New("invalid start time")
	}
	inputEndTime, err := time.Parse("15:04:05", endTime)
	if err != nil {
		return nil, nil, errors.New("invalid end time")
	}
	log.Println(inputStartTime, inputEndTime)
	rows, err := db.DB.Query(query.GetOverlappingUnavailableSlots, uid, date, startTime, endTime)
	if rows == nil {
		return nil, nil, err
	}
	defer rows.Close()

	var overlappingSlots []*models.UnavailableSlot
	earliestStartTime := inputStartTime
	latestEndTime := inputEndTime

	for rows.Next() {
		var id int
		var startTimeStr, endTimeStr sql.NullString
		if err := rows.Scan(&id, &startTimeStr, &endTimeStr); err != nil {
			return nil, nil, errors.New("Failed to parse overlapping slots")
		}
		// Parse times from the database
		var startTime, endTime time.Time
		if startTimeStr.Valid {
			startTime, err = time.Parse("2006-01-02T15:04:05Z", startTimeStr.String)
			if err != nil {
				log.Println(err)
				return nil, nil, errors.New("Failed to parse start_time from DB")
			}
		}
		if endTimeStr.Valid {
			endTime, err = time.Parse("2006-01-02T15:04:05Z", endTimeStr.String)
			if err != nil {
				return nil, nil, errors.New("Failed to parse end_time from DB")
			}
		}

		// Adjust earliest and latest times
		if startTime.Before(earliestStartTime) {
			earliestStartTime = startTime
		}
		if endTime.After(latestEndTime) {
			latestEndTime = endTime
		}

		overlappingSlots = append(overlappingSlots, &models.UnavailableSlot{
			ID: id,
		})
	}

	// Convert merged time back to string for JSON response
	earliestStartTimeStr := earliestStartTime.Format("15:04:05")
	latestEndTimeStr := latestEndTime.Format("15:04:05")

	// Return slots to remove and the merged slot to add
	toRemove := overlappingSlots
	toAdd := []*models.UnavailableSlot{
		{
			UID:             uid,
			UnavailableDate: date,
			StartTime:       &earliestStartTimeStr,
			EndTime:         &latestEndTimeStr,
		},
	}

	return toRemove, toAdd, nil
}
