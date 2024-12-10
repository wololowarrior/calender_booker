package query

var ClashingWithUnavailableSlots = `
SELECT id
FROM unavailable_slots 
		WHERE uid = $1 AND
		unavailable_date = $2 AND
		(($3 > start_time AND $3 < end_time) OR
		($4 > start_time AND $4 < end_time) OR
		($3 < start_time AND $4 > end_time) OR
		($3 > start_time AND $4 < end_time))
`

var ClashingWithMeetings = `
SELECT id
FROM meetings 
		WHERE uid = $1 AND
		date = $2 AND
		(($3 > start_time AND $3 < end_time) OR
		($4 > start_time AND $4 < end_time) OR
		($3 < start_time AND $4 > end_time) OR
		($3 >= start_time AND $4 <= end_time))
`

var InsertMeeting = `
INSERT INTO meetings (uid, date, start_time, end_time, event_id) VALUES ($1, $2, $3, $4, $5) returning id`

var GetMeetings = `
SELECT id, start_time, end_time FROM meetings WHERE uid = $1 AND date=$2`

var GetMeeting = `
select uid, date,start_time, end_time, event_id FROM meetings WHERE id = $1`

var UpdateMeeting = `
UPDATE meetings
SET 
    start_time = COALESCE($2, start_time),
    end_time = COALESCE($3, end_time)
WHERE 
    id = $1 
    returning id`

var DeleteMeeting = `
DELETE FROM meetings WHERE id = $1 returning id`
