{{if index .Modules "postgres" -}}
DB_DNS := postgres://postgres:postgres@localhost:5432/microservices?sslmode=disable

# Migration postgres commands
migration-version:
	migrate -path ./migrations/db/files -database "$(DB_DNS)" version

migration-up:
	migrate -path ./migrations/db/files -database "$(DB_DNS)" up

migration-down:
	migrate -path ./migrations/db/files -database "$(DB_DNS)" down 1{{end}}