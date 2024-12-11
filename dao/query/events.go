package query

var InsertEvent = `
INSERT INTO event (uid, name, message, slots) values ($1, $2, $3, $4) returning id;`

var SelectEventWithUserID = `
SELECT id, uid, name, message, slots FROM event WHERE uid = $1`

var SelectEvent = `
SELECT uid, name, message, slots FROM event WHERE id = $1`

var DeleteEvent = `
DELETE FROM event WHERE id = $1 AND uid = $2`
