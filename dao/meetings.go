package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"calendly_adventures/dao/query"
	"calendly_adventures/db"
	"calendly_adventures/models"
)

func CreateMeeting(meeting *models.Meeting) error {
	// is clashing with unavailable time
	err, valid := ValidMeeting(meeting)
	if !valid {
		return err
	}

	err = db.DB.QueryRow(query.InsertMeeting, meeting.UID, meeting.Date, meeting.StartTime, meeting.EndTime, meeting.EventID).Scan(&meeting.ID)
	if err != nil {
		return err
	}
	return nil
}

func ValidMeeting(meeting *models.Meeting) (error, bool) {
	var unavailableID int

	err := db.DB.QueryRow(query.ClashingWithUnavailableSlots, meeting.UID, meeting.Date, meeting.StartTime, meeting.EndTime).Scan(&unavailableID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err, false
		}
	}
	if unavailableID != 0 {
		return errors.New("unavailable during selected time"), true
	}

	// is clasing with other meetings
	var meetingID int

	err = db.DB.QueryRow(query.ClashingWithMeetings, meeting.UID, meeting.Date, meeting.StartTime, meeting.EndTime).Scan(&meetingID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err, false
		}
	}
	if meetingID != 0 {
		return errors.New("clashing with another meeting"), true
	}
	return nil, true
}

func GetSlottedMeetings(userID int, slotDuration, date string) (*[]models.Meeting, error) {
	//query unavailable
	rows, err := db.DB.Query(query.GetUnavailableSlots, userID, date)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	slots := make(map[time.Time]time.Time)
	err = parseSlots(rows, slots)
	if err != nil {
		return nil, err
	}

	rows, err = db.DB.Query(query.GetMeetings, userID, date)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	err = parseSlots(rows, slots)
	if err != nil {
		return nil, err
	}

	// query meetings
	// merge them and create available slots
	startOfDay, err := time.Parse("15:04:05", models.StartOfDay)
	if err != nil {
		log.Println(err)
	}
	endofDay, err := time.Parse("15:04:05", models.EndOfDay)
	duration, _ := strconv.Atoi(slotDuration)
	unusedSlots := generateUnusedSlots(startOfDay, endofDay, slots, time.Duration(duration)*time.Minute)
	fmt.Println(unusedSlots)
	meetingSlots := make([]models.Meeting, 0)
	for _, timeSlot := range unusedSlots {
		meeting := models.Meeting{
			UID:       userID,
			Date:      date,
			StartTime: timeSlot[0].Format("15:04:05"),
			EndTime:   timeSlot[1].Format("15:04:05"),
		}
		meetingSlots = append(meetingSlots, meeting)
	}
	return &meetingSlots, nil
}

func parseSlots(rows *sql.Rows, slots map[time.Time]time.Time) error {
	var err error
	for rows.Next() {
		var id int
		var startTimeStr, endTimeStr sql.NullString
		if err = rows.Scan(&id, &startTimeStr, &endTimeStr); err != nil {
			log.Println(err.Error())
			return errors.New("Failed to parse overlapping slots")
		}
		// Parse times from the database
		var startTime, endTime time.Time
		if startTimeStr.Valid {
			startTime, err = time.Parse("2006-01-02T15:04:05Z", startTimeStr.String)
			if err != nil {
				log.Println(err)
				return errors.New("Failed to parse start_time from DB")
			}
		}
		if endTimeStr.Valid {
			endTime, err = time.Parse("2006-01-02T15:04:05Z", endTimeStr.String)
			if err != nil {
				return errors.New("Failed to parse end_time from DB")
			}
		}
		if v, ok := slots[startTime]; ok {
			if v.After(endTime) {
				continue
			}
		}
		slots[startTime] = endTime
	}
	return nil
}

func generateUnusedSlots(startOfDay, endOfDay time.Time, usage map[time.Time]time.Time, slotDuration time.Duration) [][2]time.Time {
	var unusedSlots [][2]time.Time

	// Convert the map into a slice of start-end pairs and sort them by start time
	var usedPeriods [][2]time.Time
	for start, end := range usage {
		usedPeriods = append(usedPeriods, [2]time.Time{start, end})
	}

	sort.Slice(usedPeriods, func(i, j int) bool {
		return usedPeriods[i][0].Before(usedPeriods[j][0])
	})

	// Iterate through the day to find unused slots
	current := startOfDay
	for _, period := range usedPeriods {
		startTime := period[0]
		endTime := period[1]

		// Fill unused slots before the next used period starts
		for current.Add(slotDuration).Before(startTime) {
			slotEnd := current.Add(slotDuration)
			if slotEnd.After(startTime) {
				slotEnd = startTime
			}
			unusedSlots = append(unusedSlots, [2]time.Time{current, slotEnd})
			current = slotEnd
		}

		// Move current to the end of the used period
		if current.Before(endTime) {
			current = endTime
		}
	}

	// Fill unused slots after the last used period until the end of the day
	for current.Add(slotDuration).Before(endOfDay) {
		slotEnd := current.Add(slotDuration)
		unusedSlots = append(unusedSlots, [2]time.Time{current, slotEnd})
		current = slotEnd
	}
	log.Println(current, endOfDay, endOfDay.Sub(current), slotDuration)
	if endOfDay.Sub(current) >= slotDuration {
		unusedSlots = append(unusedSlots, [2]time.Time{current, current.Add(slotDuration)})
	}

	return unusedSlots
}

func GetMeeting(meeting *models.Meeting) error {
	err := db.DB.QueryRow(query.GetMeeting, meeting.ID).Scan(&meeting.UID, &meeting.Date, &meeting.StartTime, &meeting.EndTime, &meeting.EventID)
	meeting.StartTime = timeParse(meeting.StartTime)
	meeting.EndTime = timeParse(meeting.EndTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("meeting not found")
		}
		return err
	}
	return nil
}

func UpdateMeeting(meeting *models.Meeting) error {
	err, valid := ValidMeeting(meeting)
	if !valid {
		return errors.New("invalid Meeting")
	}
	var id int
	err = db.DB.QueryRow(query.UpdateMeeting, meeting.ID, meeting.StartTime, meeting.EndTime).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteMeeting(meeting *models.Meeting) error {
	err := db.DB.QueryRow(query.DeleteMeeting, meeting.ID).Scan(&meeting.ID)

	if err != nil {
		return err
	}
	return nil
}

func GetBookedMeetings(userID int, date string) ([]*models.Meeting, error) {
	getQuery := `SELECT id,start_time, end_time, event_id, created_at FROM meetings WHERE uid = $1 AND date=$2`
	rows, err := db.DB.Query(getQuery, userID, date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	meetings := make([]*models.Meeting, 0)
	for rows.Next() {
		meeting := models.Meeting{
			UID:  userID,
			Date: date,
		}
		err := rows.Scan(&meeting.ID, &meeting.StartTime, &meeting.EndTime, &meeting.EventID, &meeting.CreatedAt)

		if err != nil {
			return nil, err
		}
		meeting.StartTime = timeParse(meeting.StartTime)
		meeting.EndTime = timeParse(meeting.EndTime)
		meetings = append(meetings, &meeting)
	}
	return meetings, nil
}

func timeParse(s string) string {
	parse, err := time.Parse("2006-01-02T15:04:05Z", s)
	if err != nil {
		return ""
	}
	return parse.Format("15:04:05")
}
