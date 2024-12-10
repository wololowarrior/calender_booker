package query

var InsertEvent = `
INSERT INTO event (uid, name, message, slotted, slots) values ($1, $2, $3, $4, $5) returning id;`

var SelectEventWithUserID = `
SELECT id, uid, name, message, slotted, slots FROM event WHERE uid = $1`

var SelectEvent = `
SELECT uid, name, message, slotted, slots FROM event WHERE id = $1`
