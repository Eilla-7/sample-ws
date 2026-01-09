# Sample Go JWT Web Service

Simple Go server with JWT auth.

## Server (`sample-ws`)
- `/auth` → login, returns JWT
- `/query` → protected, requires header: `Authorization: Bearer <token>`
- Uses MySQL to store users

Use `go run .` to start the server.