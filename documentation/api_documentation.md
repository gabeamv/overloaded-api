# Overloaded API Documentation

Base URL: `http://localhost:8080`

This document covers every active route registered in [main.go](../main.go). Commented-out routes are intentionally excluded.

## Authentication

Most protected routes require this header (exercises, progression rules, tokens, sets, exercise stats, users, workouts):

```http
Authorization: Bearer <token>
```

There are two token types in this API:

- Access token: returned from `POST /api/users/login`, used for most authenticated endpoints.
- Refresh token: returned from `POST /api/users/login`, used only with `POST /api/refresh` and `POST /api/revoke`.

## Response Conventions

- Successful responses are JSON.
- Error responses use this shape:

```json
{
  "error": "error message"
}
```

## Admin

### DELETE `/admin/reset`

Deletes all users. This route is intended for development use only and requires `PLATFORM=dev`.

Auth: none

Success response:

```json
"All users deleted."
```

Common status codes:

- `200 OK`
- `500 Internal Server Error`

## Users

### POST `/api/users/register`

Creates a new user account.

Auth: none

Request body:

```json
{
  "email": "user@example.com",
  "username": "gabeamv",
  "password": "password123"
}
```

Success response:

```json
{
  "id": "05cf60c3-b730-4e40-a495-232d3b587637",
  "username": "gabeamv",
  "email": "user@example.com",
  "created_at": "2026-04-25T08:00:00Z",
  "updated_at": "2026-04-25T08:00:00Z"
}
```

Common status codes:

- `201 Created`
- `500 Internal Server Error`

### POST `/api/users/login`

Authenticates a user and returns an access token plus a refresh token.

Auth: none

Request body:

```json
{
  "username": "gabeamv",
  "password": "password123"
}
```

Success response:

```json
{
  "id": "05cf60c3-b730-4e40-a495-232d3b587637",
  "username": "gabeamv",
  "email": "user@example.com",
  "created_at": "2026-04-25T08:00:00Z",
  "updated_at": "2026-04-25T08:00:00Z",
  "token": "<access-token>",
  "refresh_token": "<refresh-token>"
}
```

Common status codes:

- `202 Accepted`
- `401 Unauthorized`
- `500 Internal Server Error`

### PUT `/api/users/username`

Updates the authenticated user's username.

Auth: access token required

Request body:

```json
{
  "username": "new_username"
}
```

Success response:

```json
{
  "id": "05cf60c3-b730-4e40-a495-232d3b587637",
  "username": "new_username",
  "email": "user@example.com",
  "created_at": "2026-04-25T08:00:00Z",
  "updated_at": "2026-04-25T08:05:00Z"
}
```

Common status codes:

- `202 Accepted`
- `401 Unauthorized`
- `409 Conflict`
- `500 Internal Server Error`

### PUT `/api/users/email`

Updates the authenticated user's email.

Auth: access token required

Request body:

```json
{
  "email": "new_email@example.com"
}
```

Success response:

```json
{
  "id": "05cf60c3-b730-4e40-a495-232d3b587637",
  "username": "gabeamv",
  "email": "new_email@example.com",
  "created_at": "2026-04-25T08:00:00Z",
  "updated_at": "2026-04-25T08:05:00Z"
}
```

Common status codes:

- `202 Accepted`
- `401 Unauthorized`
- `409 Conflict`
- `500 Internal Server Error`

### PUT `/api/users/password`

Updates the authenticated user's password.

Auth: access token required

Request body:

```json
{
  "password": "newPassword123"
}
```

Success response:

```json
{
  "id": "05cf60c3-b730-4e40-a495-232d3b587637",
  "username": "gabeamv",
  "email": "user@example.com",
  "created_at": "2026-04-25T08:00:00Z",
  "updated_at": "2026-04-25T08:05:00Z"
}
```

Common status codes:

- `202 Accepted`
- `401 Unauthorized`
- `500 Internal Server Error`

## Token Management

### POST `/api/refresh`

Generates a new access token using a refresh token.

Auth: refresh token required

Header:

```http
Authorization: Bearer <refresh-token>
```

Success response:

```json
{
  "token": "<new-access-token>"
}
```

Common status codes:

- `200 OK`
- `401 Unauthorized`
- `500 Internal Server Error`

### POST `/api/revoke`

Revokes a refresh token.

Auth: refresh token required

Header:

```http
Authorization: Bearer <refresh-token>
```

Success response:

```json
{}
```

Common status codes:

- `204 No Content`
- `500 Internal Server Error`

## Exercises

### POST `/api/exercises`

Creates a new exercise for the authenticated user.

Auth: access token required

Request body:

```json
{
  "name": "bench press"
}
```

Success response:

```json
{
  "id": "aaa9d580-9cb9-4b82-a43e-999da8f63ada",
  "name": "bench press",
  "is_custom": true,
  "user_id": "05cf60c3-b730-4e40-a495-232d3b587637"
}
```

Common status codes:

- `201 Created`
- `401 Unauthorized`
- `500 Internal Server Error`

### GET `/api/exercises`

Returns the authenticated user's exercises.

Auth: access token required

Success response:

```json
[
  {
    "id": "aaa9d580-9cb9-4b82-a43e-999da8f63ada",
    "name": "bench press",
    "is_custom": true,
    "user_id": "05cf60c3-b730-4e40-a495-232d3b587637"
  }
]
```

Common status codes:

- `202 Accepted`
- `401 Unauthorized`
- `500 Internal Server Error`

### GET `/api/exercises/{exercise_id}`

Returns a single exercise by ID.

Auth: access token required

Path params:

- `exercise_id`: UUID of the exercise

Behavior:

- Custom exercises can only be accessed by their owner.
- Non-custom exercises can be read by authenticated users.

Success response:

```json
{
  "id": "aaa9d580-9cb9-4b82-a43e-999da8f63ada",
  "name": "bench press",
  "is_custom": true,
  "user_id": "05cf60c3-b730-4e40-a495-232d3b587637"
}
```

Common status codes:

- `202 Accepted`
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

### PUT `/api/exercises/{exercise_id}`

Updates the name of a custom exercise.

Auth: access token required

Path params:

- `exercise_id`: UUID of the exercise

Request body:

```json
{
  "name": "incline bench press"
}
```

Behavior:

- Only the owner of a custom exercise can update it.
- Non-custom exercises cannot be updated through this route.

Success response:

```json
{
  "id": "aaa9d580-9cb9-4b82-a43e-999da8f63ada",
  "name": "incline bench press",
  "is_custom": true,
  "user_id": "05cf60c3-b730-4e40-a495-232d3b587637"
}
```

Common status codes:

- `202 Accepted`
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

### DELETE `/api/exercises/{exercise_id}`

Deletes a custom exercise.

Auth: access token required

Path params:

- `exercise_id`: UUID of the exercise

Behavior:

- Only the owner of a custom exercise can delete it.
- Non-custom exercises cannot be deleted through this route.

Success response:

```json
{}
```

Common status codes:

- `204 No Content`
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

## Workouts

### POST `/api/workouts`

Creates a workout for the authenticated user.

Auth: access token required

Request body:

```json
{
  "started_at": "2026-04-08T16:30:00Z",
  "ended_at": "2026-04-08T16:40:00Z",
  "volume": 0,
  "prs": 0
}
```

Success response:

```json
{
  "id": "9d667d95-9c83-4c51-add4-83d0166914ef",
  "user_id": "05cf60c3-b730-4e40-a495-232d3b587637",
  "started_at": "2026-04-08T16:30:00Z",
  "ended_at": "2026-04-08T16:40:00Z",
  "volume": 0,
  "prs": 0,
  "is_completed": false
}
```

Common status codes:

- `201 Created`
- `401 Unauthorized`
- `500 Internal Server Error`

Note: In order to add a batch of sets, a workout must be created first.

### GET `/api/workouts`

Returns all workouts belonging to the authenticated user.

Auth: access token required

Success response:

```json
[
  {
    "id": "9d667d95-9c83-4c51-add4-83d0166914ef",
    "user_id": "05cf60c3-b730-4e40-a495-232d3b587637",
    "started_at": "2026-04-08T16:30:00Z",
    "ended_at": "2026-04-08T16:40:00Z",
    "volume": 1100,
    "prs": 2,
    "is_completed": true
  }
]
```

Common status codes:

- `201 Created`
- `401 Unauthorized`
- `500 Internal Server Error`

Note:

- The current implementation returns `201 Created` instead of the more typical `200 OK`.

### GET `/api/workouts/{workout_id}`

Returns a single workout by ID.

Auth: access token required

Path params:

- `workout_id`: UUID of the workout

Success response:

```json
{
  "id": "9d667d95-9c83-4c51-add4-83d0166914ef",
  "user_id": "05cf60c3-b730-4e40-a495-232d3b587637",
  "started_at": "2026-04-08T16:30:00Z",
  "ended_at": "2026-04-08T16:40:00Z",
  "volume": 1100,
  "prs": 2,
  "is_completed": true
}
```

Common status codes:

- `202 Accepted`
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

### DELETE `/api/workouts/{workout_id}`

Deletes a workout by ID.

Auth: access token required

Path params:

- `workout_id`: UUID of the workout

Success response:

```json
{}
```

Common status codes:

- `204 No Content`
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

## Workout Sets

### POST `/api/workouts/{workout_id}/sets/batch`

Creates a batch of workout sets for a workout, updates workout PR count and volume, and marks the workout as completed.

Auth: access token required

Path params:

- `workout_id`: UUID of the workout

Request body:

```json
{
  "sets": [
    {
      "exercise_id": "aaa9d580-9cb9-4b82-a43e-999da8f63ada",
      "progress_track": "d5ed01a7-a710-44cc-ad2f-7850d1916fe7",
      "weight_in_lbs": 250,
      "reps": 5,
      "time_in_seconds": 0
    },
    {
      "exercise_id": "8fb51970-46a9-4fb9-bb95-95571b728a53",
      "progress_track": "d5ed01a7-a710-44cc-ad2f-7850d1916fe7",
      "weight_in_lbs": 250,
      "reps": 8,
      "time_in_seconds": 0
    }
  ]
}
```

Behavior:

- All sets are created within a transaction.
- The route calculates estimated one-rep max and per-set volume to update workout PR count and total volume.
- Once the batch succeeds, the workout is marked completed.
- The current implementation rejects adding more sets to an already completed workout.

Success response:

```json
[
  {
    "id": "5e764b8b-f371-4d6b-806f-325de6392372",
    "workout_id": "9d667d95-9c83-4c51-add4-83d0166914ef",
    "exercise_id": "aaa9d580-9cb9-4b82-a43e-999da8f63ada",
    "progress_track": "d5ed01a7-a710-44cc-ad2f-7850d1916fe7",
    "weight_in_lbs": 250,
    "reps": 5,
    "time_in_seconds": 0
  }
]
```

Common status codes:

- `201 Created`
- `401 Unauthorized`
- `404 Not Found`
- `409 Conflict`
- `500 Internal Server Error`

Note: In order to add a batch of sets, a workout must be created first.

### GET `/api/workouts/{workout_id}/sets`

Returns all sets belonging to a workout.

Auth: access token required

Path params:

- `workout_id`: UUID of the workout

Success response:

```json
[
  {
    "id": "5e764b8b-f371-4d6b-806f-325de6392372",
    "workout_id": "9d667d95-9c83-4c51-add4-83d0166914ef",
    "exercise_id": "aaa9d580-9cb9-4b82-a43e-999da8f63ada",
    "progress_track": "d5ed01a7-a710-44cc-ad2f-7850d1916fe7",
    "weight_in_lbs": 250,
    "reps": 5,
    "time_in_seconds": 0
  }
]
```

Common status codes:

- `202 Accepted`
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

### GET `/api/sets/{set_id}`

Returns a single set by ID.

Auth: access token required

Path params:

- `set_id`: UUID of the workout set

Success response:

```json
{
  "id": "5e764b8b-f371-4d6b-806f-325de6392372",
  "workout_id": "9d667d95-9c83-4c51-add4-83d0166914ef",
  "exercise_id": "aaa9d580-9cb9-4b82-a43e-999da8f63ada",
  "progress_track": "d5ed01a7-a710-44cc-ad2f-7850d1916fe7",
  "weight_in_lbs": 250,
  "reps": 5,
  "time_in_seconds": 0
}
```

Common status codes:

- `202 Accepted`
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

Note: 'progress_track' is the id of the progression rule used in the set.

## Progression Rules

### GET `/api/progression_rules`

Returns all progression rules.

Auth: access token required

Success response:

```json
[
  {
    "id": "d5ed01a7-a710-44cc-ad2f-7850d1916fe7",
    "label": "p",
    "rule": "progress",
    "description": "Progress and go up in weight next session."
  }
]
```

Common status codes:

- `202 Accepted`
- `401 Unauthorized`
- `500 Internal Server Error`

### POST `/api/progression_rules`

Seeds the built-in progression rules. This route is intended for development use only and requires `PLATFORM=dev`.

Auth: none

Rules created:

- `p` -> `progress`
- `s` -> `stay`
- `t` -> `tag on form`

Success response:

```json
{}
```

Common status codes:

- `204 No Content`
- `500 Internal Server Error`

### DELETE `/api/progression_rules`

Deletes all progression rules. This route is intended for development use only and requires `PLATFORM=dev`.

Auth: none

Success response:

- No response body is written by the current handler.

Common status codes:

- Typically succeeds with an empty response
- `500 Internal Server Error`

## Stats

### GET `/api/exercises/{exercise_id}/1rm`

Returns the authenticated user's estimated one-rep max for an exercise using the Epley formula.

Auth: access token required

Path params:

- `exercise_id`: UUID of the exercise

Success response:

```json
{
  "one_rep_max": 291.67
}
```

Common status codes:

- `202 Accepted`
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`

### GET `/api/exercises/{exercise_id}/volume`

Returns the authenticated user's highest recorded set volume for an exercise.

Auth: access token required

Path params:

- `exercise_id`: UUID of the exercise

Success response:

```json
{
  "volume_max": 2000
}
```

Common status codes:

- `202 Accepted`
- `401 Unauthorized`
- `404 Not Found`
- `500 Internal Server Error`
