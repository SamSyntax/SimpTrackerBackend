gooseup:
	goose -dir internal/db/migrations postgres "postgres://postgres:postgres@localhost:5432/stuleja?sslmode=disable" up 

goosedown:
	goose -dir internal/db/migrations postgres "postgres://postgres:postgres@localhost:5432/stuleja?sslmode=disable" down 

build:
	go build -o ./backend cmd/main.go

run:
	go run cmd/main.go
