## Flags

- dev - dev or prod
- configFilepath - filepath to config.toml file
- logLevel - trace / debug / info

## ENV variables

- $PORT - just listen port, defaults to 8080
- $DATABASE_URL - url to database in format of "postgres://127.0.0.1:5432"
- $SECRET_KEY - it's key for cookie store, server using this store to store login data.
- $SENDGRID_KEY - key for sendgrid.com, I'm using it as backend for sending emails.

## Common models

- Message
```json
{"user": "user", "chat": "chat", "id": "id", "text": "text", "time": "time"}
```

- User
```json
{"user": "user"}
```

- Chat
```json
{"chat": "chat"}
```

- Events
* LoginEvent, LogoutEvent, ChatJoinEvent, ChatLeaveEvent
```json
{User, Chat, Time}
```
* MessageEvent
```json
{Message}
```

* ChatCreatedEvent, ChatRemovedEvent
```json
{Chat, Time}
```

- ChatReport: contains info about chat for concrete user. If not joined - fields "messages" and "users" are empty.
```json
{
    User,
    Chat, 
    "joined":true, 
    "messages": [...Message],
    "users": [...User]
}
```

## API

#### /api/register
- Server will store user in db and send a email with validation link to email.
- Method: POST
- Request:
```json
{User, "email": "user@example.com", "password": "password"}
```

#### /api/verifyEmail/{username}/{activationToken}
- Server will set user activated status to true in database
- Method: GET

#### /api/authorize
- Server will check if user is in database, if user had verified email and if password is correct, then force logins user.
- Method: POST
- Request:
```json
{User, "password": "password"}
```
- Response: valid cookie with session token

#### /api/login
- Login for unregistered users, only if user with this name isn't already logged in.
- Method: POST
- Response: valid cookie with session token and
```json
User
```

### /api/loggedonly/*
- this routes are only for users with valid cookie, that can be acquired via upper routes

#### /api/loggedonly/login
- Force logins user.
- Method: POST
- Response: valid cookie with session token and
```json
User
```

#### /api/loggedonly/logout
- Logouts user.
- Method POST

#### /api/loggedonly/keepalive
- Updates user's last activity.
- Method: PUT

#### /api/loggedonly/poll
- Long poll for actions.
- Method: GET
- Response: Event
```json
{"type": "EventType", "event": Event}
```

#### /api/loggedonly/joinChat
- Joins chat
- Method: POST
- Request:
```json
Chat
```
- Response:
```json
ChatReport
```

#### /api/loggedonly/leaveChat
- Leaves chat
- Method: POST
- Request:
```json
Chat
```

#### /api/loggedonly/getChats
- Get all chats, joined status, messages for chats and users, if user permitted.
- Method: GET
- Response:
```json
[...ChatReport]
```

#### /api/loggedonly/getLastMessages
- Gets last messages for a chat if is permitted
- Method: POST
- Request:
```json
Chat
```
- Response:
```json
[...Message]
```

#### /api/loggedonly/postMessage
- Posts message in chat if permitted
- Method: POST
- Request:
```json
{User, "text": "text", Chat}
```

#### /api/loggedonly/getUsers
- Gets users in chat if permitted
- Method: POST
- Request:
```json
Chat
```
- Response:
```json
[...User]
```

### /api/adminonly/*
- this routes are only for admins

#### /api/adminonly/createChat
- Creates new chat
- Method: POST
- Request:
```json
Chat
```

#### /api/adminonly/removeChat
- Removes a chat
- Method: POST
- Request:
```json
Chat
```