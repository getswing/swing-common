# Go Echo Microservice + RabbitMQ + Postgres

Komponen:
- API Service (Echo) dengan GORM Postgres, publish event ke RabbitMQ saat membuat Todo.
- Worker Service (consumer) untuk memproses event `todos.created`.

## RabbitMQ
- AMQP: amqp://localhost:5672
- Management UI: http://localhost:15672 (user/pass default: `guest`/`guest` atau sesuai .env)
- definitions.json otomatis meload:
  - Exchange: `todos` (topic)
  - Queue: `todo_created`
  - Binding: routing key `todos.created`

## Auth (JWT)
- Register: POST http://localhost:8080/auth/register
  - Body: {"email":"user@example.com","password":"secret123"}
- Login: POST http://localhost:8080/auth/login
  - Body: {"email":"user@example.com","password":"secret123"}
  - Response: {"access_token":"...","token_type":"Bearer","expires_in":900}

Gunakan header Authorization: Bearer <token> untuk mengakses endpoint di bawah `/api`.

Contoh:
- GET/POST http://localhost:8080/api/todos (protected)
- GET http://localhost:8080/health (public)

## DB Migrations
Migrations are applied automatically by the `migrate` one-shot service when you run `docker-compose up`.

- Directory: `db/migrations`
- Naming: `NNNN_description.up.sql` and `NNNN_description.down.sql` (e.g., `0002_add_index.up.sql`)
- Image: `migrate/migrate:4`

Manual usage examples:
- Run all up migrations: `docker compose run --rm migrate up`
- Roll back one step: `docker compose run --rm migrate down 1`
- Force version (use with caution): `docker compose run --rm migrate force <version>`

The migrate service uses the in-network DB URL:
`postgres://postgres:postgres@postgres:5432/microservice?sslmode=disable`

Jalankan dengan Docker:
1) Salin `.env.example` menjadi `.env` (opsional untuk lokal).
2) `docker-compose up --build -d`
3) API: http://localhost:8080/health, POST `http://localhost:8080/api/todos` dengan JSON `{"title":"example"}`
4) RabbitMQ Management: http://localhost:15672 (guest/guest)

Struktur kode memisahkan config, db, mq, models, repository, handlers untuk maintainability.



Grpc Setup:

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

brew install protobuf

export PATH="$PATH:$(go env GOPATH)/bin"

Try: 
protoc --go_out=. --go-grpc_out=. hello.proto
protoc --go_out=. --go-grpc_out=. services/grpc/proto/hello.proto


protoc \
  --go_out=paths=source_relative:. \
  --go-grpc_out=paths=source_relative:. \
  proto/player.proto
