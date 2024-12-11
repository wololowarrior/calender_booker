## Assumptions
1. Meetings can be scheduled with or without an event
2. An event mean, 
   1. the detail entry of name of the event (like hiring for xyz), 
   2. offering of meeting slots against that event
   3. customisation of slot duration (30m, 60m)
   4. an automated message to be sent when a slot is booked
3. Currently Only supporting one timezone. Or rather no timezone just working with plain clock numbers. 
4. The user is available throughout the day, but slot timing are restricted from 9am to 5pm
   1. Which means meeting can be scheduled throughout the day. 
5. After meeting is booked, a predefined message is sent to the email with call link. 

## MVP
A calendly like MVP should allow 
1. user creation
2. unavailability manager
3. slot creation, meeting management and customised links
4. Calender integration
5. Reminders
6. Event creation and event duration

My MVP mostly satisfies all of these, some of the complex features like customised link, calender integrations, reminders
require more engineering and explorations into 3rd party api's which isnt possible in a week.

## Trade OFFs and Hacks
1. Doing many things via code and not utilising innate RDS features like joins, views.
   1. To rapidly prototype
2. Not using date range type in postgres, that would've reduced the number of rows created for unavailability range.
   1. Instead hacked it by going in a loop and inserting a row for each day.
3. HACK, some ugliness in code due to not managing date column nicely in postgres. Do to which time.parse is needed in code
But this doesn't affect things as we mostly don't care about what day it is, just the time is of concern.

## How to install

### Start the service
1. Install Docker if needed
2. https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository
3. `sudo service docker status`. Start docker if needed
4. Clone the repo and `cd` inside it
5. Run `sudo docker-compose up` from the base dir. It should build the app image for the first time and start the container with the postgres
   1. Give `--build` tag if you want to force build.

### Install Goose to support postgres migration

1. https://pressly.github.io/goose/installation/
```shell
curl -fsSL     https://raw.githubusercontent.com/pressly/goose/master/install.sh |    sudo sh
```
2. Run migrations
```shell
EXPORT DB_URL=postgres://postgres:password@localhost:5432/calendly?sslmode=disable
goose -dir ./db/migrations postgres "$DB_URL" up
```
3. Nuke the db
```shell
goose -dir ./db/migrations postgres "$DB_URL" reset
```

You should be good to go, fire your api's away

## How to play
1. Create a user 
2. Get the list of unavailable time. Create unavailable time if needed. See #3 in [api spec](#api-spec). 
3. Create an event (for ex Recruitment xyz position, Book my calender to make a sale). See #5 in [api spec](#api-spec)
   1. Define a message that can be sent to the
   person booking it
   2. Maybe, get the event list that you've created. See #6 in [api spec](#api-spec)
   3. Or an overview of the day. See #8 in [api spec](#api-spec)
4. Get available slots for the user. See #9 in [api spec](#api-spec)
5. Book a meeting, specify the event ID. See #10 in [api spec](#api-spec)
6. You can also choose to reschedule the meeting or cancel it if plans change. See #12 in [api spec](#api-spec)


## Future Work
1. Save the details of who booked the slot in db. Like name email etc.
2. Send email to booker with predefined message and a video call link.
This can be done using a async worker and queue.
3. Reminder for a booked meeting to both the booker and user.
4. Integration with multiple calendars so that unavailable slot / the calendar can be updated with meetings
5. Support multiple timezones. Display values in local tz.
 

## Api Spec

1. POST /user : This will create a user whose calendar is maintained
```shell
curl --location 'localhost:8080/user' \
--header 'Content-Type: application/json' \
--data '{
    "name": "harshil",
    "email": "hg"
}'
```

2. GET /user/{id} : Gets an user
```shell
curl --location 'localhost:8080/user/1'

Response
{
    "id": 1,
    "name": "harshil",
    "email": "hg",
    "created_at": "2024-12-10T12:23:42.490824Z"
}
```

3. POST /user/{id}/unavailable create unavailable time.
This api works on a range of date. For ex you can create unavailable time on series of day 
or a single day by giving same start/end date.
```shell
curl --location --request GET 'localhost:8080/user/1/unavailable' \
--header 'Content-Type: application/json' \
--data '{
    "start_date": "2024-12-13",
    "end_date":"2024-12-13",
    "start_time": "11:00:00",
    "end_time": "15:00:00"
}'
```
4. GET /user/{id}/unavailable get unavailable time, return days >= current day
```shell
curl --location 'localhost:8080/user/1/unavailable'

Response
[
    {
        "id": 1,
        "uid": 1,
        "unavailable_date": "2024-12-13",
        "start_time": "11:00:00",
        "end_time": "15:00:00",
        "created_at": "2024-12-10T17:27:08.023987Z"
    }
]
```
5. POST /user/{id}/event create an event
```shell
curl --location 'localhost:8080/user/1/event' \
--header 'Content-Type: application/json' \
--data '{
    "name": "test",
    "slots":"30"
}'
```
6. GET /user/{id}/event list all event
```shell
curl --location --request GET 'localhost:8080/user/1/event' \
--header 'Content-Type: application/json'

Response
[
    {
        "id": 1,
        "uid": 1,
        "name": "test",
        "message": "",
        "slots": "30"
    }
]
```
7. DELETE /user/{id}/event/{event_id} delete an event
8. GET /user/{id}/overview get an overview of for a date
```shell
curl --location 'localhost:8080/user/1/overview?date=2024-12-13'

Response
{
    "unavailable_slots": [
        {
            "id": 1,
            "uid": 1,
            "unavailable_date": "2024-12-13",
            "start_time": "11:00:00",
            "end_time": "15:00:00"
        }
    ],
    "meetings": []
}
```
GET /user/{id}/meetings?date=<optional> to just get the list of meetings 
```shell
curl --location --request GET 'localhost:8080/user/1/meetings?date=YYYY-mm-dd'
```
9. GET `/meetings` get available slots for an event on a day
```shell
curl --location --request GET 'localhost:8080/meetings' \
--header 'Content-Type: application/json' \
--data '{
    "uid":1,
    "date": "2024-12-13",
    "event_id":2
}'

Response
[
    {
        "id": 0,
        "uid": 1,
        "date": "2024-12-12",
        "start_time": "09:00:00",
        "end_time": "09:30:00"
    },
    .
    .
    .
    .
    {
        "id": 0,
        "uid": 1,
        "date": "2024-12-12",
        "start_time": "16:30:00",
        "end_time": "17:00:00"
    }
]
```
10. POST /meetings create a meeting tied to an event. This will in future send an email with the 
video call link etc.
```shell
curl --location 'localhost:8080/meetings' \
--header 'Content-Type: application/json' \
--data '{
    "uid":1,
    "date": "2024-12-13",
    "start_time": "09:30:00",
    "end_time":"10:30:00",
    "event_id":1
}'
```
11. GET /meetings/{id} Get the meeting
12. PUT /meetings/{id} update a meeting. Update end_time or start_time. This validates if the slot is possible or not
```shell
curl --location --request PUT 'localhost:8080/meetings/5' \
--header 'Content-Type: application/json' \
--data '{
    "end_time":"11:00:00"
}'
```
13. DELETE /meetings/{id}

