## Assumptions
1. Meetings can be scheduled with or without an event
2. An event mean, 
   1. the detail entry of name of the event (like hiring for xyz), 
   2. offering of meeting slots against that event
   3. customisation of slot duration (30m, 60m)
   4. an automated message to be sent when a slot is booked
3. The user is available throughout the day, but slot timing are restricted from 9am to 5pm
   1. Which means meeting can be scheduled throughout the day. 
4. After meeting is booked, a predefined message is sent to the email with call link. 

## How to install
1. Clone the repo and cd inside it
2. Run `docker-compose up` from the base dir.
3. You should be good to go, fire your api's away

## How to play
1. Create a user
2. Create unavailable time if needed.
3. Create an event for ex Recruitment for xyz, Book my calender to make a sale. Define a message that can be sent to the
person booking it
4. Get available slots for the person. 
5. Book a meeting, specify the event ID.
6. You can also choose to reschedule the meeting or cancel it if plans change


## Future Work
1. Save the details of who booked the slot. Like name email etc.
2. Send email to booker with predefined message and a video call link.
This can be done using a async worker and queue.
3. Reminder for a booked meeting to both the booker and user.
4. Integration with multiple calendars so that unavailable slot / the calendar can be updated with meetings
5. 

## Api Spec

1. POST /user : This will create a user whose calendar is maintained
```shell
curl --location 'localhost:8080/users' \
--header 'Content-Type: application/json' \
--data '{
    "name": "harshil",
    "email": "hg"
}'
```

2. GET /user/{id} : Gets an user
```shell
curl --location 'localhost:8080/user/1'
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
5. POST /users/{id}/event create an event
```shell
curl --location 'localhost:8080/users/1/event' \
--header 'Content-Type: application/json' \
--data '{
    "name": "test",
    "slots":"30"
}'
```
6. GET /users/{id}/event list all event
```shell
curl --location --request GET 'localhost:8080/users/1/event' \
--header 'Content-Type: application/json'
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
7. DELETE /users/{id}/event/{event_id} delete an event
8. GET /users/{id}/overview get an overview of for a date
```shell
curl --location 'localhost:8080/users/1/overview?date=2024-12-13'
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
9. GET /meetings?event_id=1&user_id=1&date=2024-12-12 get available slots for an event on a day
```shell
curl --location 'localhost:8080/meetings?event_id=1&user_id=1&date=2024-12-12'
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
11. PUT /meetings/{id} update a meeting 
```shell
curl --location --request PUT 'localhost:8080/meetings/5' \
--header 'Content-Type: application/json' \
--data '{
    "end_time":"11:00:00"
}'
```
12. DELETE /meetings/{id}
13. GET /meetings/{id}
