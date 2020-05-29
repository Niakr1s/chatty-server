## API

#### /api/register
- Server will store user in db and send a email with validation link to email.
- Method: POST
- Request:
```json
{"user": "user", "email": "user@example.com", "password": "password"}
```

#### /api/verifyEmail/{username}/{activationToken}
- Server will set user activated status to true in database
- Method: GET

#### /api/authorize
- Server will check if user is in database, if user had verified email and if password is correct, then force logins user.
- Method: POST
- Request:
```json
{"user": "user", "password": "password"}
```
- Response: valid cookie with session token

#### /api/login
- Login for unregistered users, only if user with this name isn't already logged in.
- Method: POST
- Response: valid cookie with session token and
```json
{"user": "user"}
```

### /api/loggedonly/*
- this routes are only for users with valid cookie, that can be acquired via upper routes

#### /api/loggedonly/login
- Force logins user.
- Method: POST
- Response: valid cookie with session token and
```json
{"user": "user"}
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

#### /api/loggedonly/joinChat
- Joins chat
- Method: POST
 - Request:
```json
{"chat": "chat"}
```

#### /api/loggedonly/leaveChat
- Leaves chat
- Method: POST
- Request:
```json
{"chat": "chat"}
```

#### /api/loggedonly/getChats
- Get all chats, joined status, messages for chats and users, if user permitted.
- Method: GET
- Response:
```json
[{
    "chat":"chat", 
    "joined":true, 
    "messages": [
        {"user": "user", "chat": "chat", "id": "id", "text": "text", "time": "time"}
    ],
    "users": [
        {"user": "user"}
    ]}]
```

#### /api/loggedonly/getLastMessages
- Gets last messages for a chat if is permitted
- Method: POST
- Request:
```json
{"chat": "chat"}
```
- Response:
```json
[{"user": "user", "chat": "chat", "id": "id", "text": "text", "time": "time"}]
```

#### /api/loggedonly/postMessage
- Posts message in chat if permitted
- Method: POST
- Request:
```json
{"user": "user", "text": "text", "chat": "chat"}
```

#### /api/loggedonly/getUsers
- Gets users in chat if permitted
- Method: POST
- Request:
```json
{"chat": "chat"}
```
- Response:
```json
[{"user": "user"}]
```