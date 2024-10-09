## Google Calendar Synchronisation Implementation

### OverView

This backend application lets you connect your google accounts and sync all events from your google calendar. 

### Setup

* Clone the repo and navigate into the directoru.
* Copy `.env.sample` to `.env` & fill in the values.  
* Do `go get ./...`
* To start the server, do run  `go run ./.../cmd`
* Run the cron job using `go run ./.../crons`
* Run tests using `go test ./...`

Auth middleware is bypassed using query params.
    
Navigate to `HOST:PORT/index?user_id=u1`
// change query param user_id=u2 for second user

Visit the /connect endpoint to connect google accounts, with the same authorization.

    - You can connect multiple google accounts to same user.
    - Same google accounts can be mapped to multiple users.
    - The application will maintain this as seperate entities.

