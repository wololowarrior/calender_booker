package query

var InsertUnavailableSlots = `
INSERT INTO unavailable_slots (uid, unavailable_date, start_time, end_time) VALUES ($1, $2, $3, $4) returning id, created_at
`

var SelectUnavailableSlots = `
		SELECT id, uid , unavailable_date, start_time, end_time, created_at 
		FROM unavailable_slots 
		WHERE uid = $1 AND unavailable_date >= CURRENT_DATE`

var GetOverlappingUnavailableSlots = `
		SELECT id, start_time, end_time
		FROM unavailable_slots
		WHERE uid = $1
		  AND unavailable_date = $2
		  AND (start_time < $4 AND end_time > $3)`

var DeleteUnavailableSlots = `
DELETE FROM unavailable_slots
WHERE id = $1`

var GetUnavailableSlots = `
SELECT id, start_time, end_time FROM unavailable_slots 
WHERE uid = $1 AND unavailable_date = $2`
