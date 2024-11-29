gooseup:
	goose -dir internal/db/migrations postgres "postgres://postgres:postgres@localhost:5432/dev_stuleja?sslmode=disable" up 

goosedown:
	goose -dir internal/db/migrations postgres "postgres://postgres:postgres@localhost:5432/dev_stuleja?sslmode=disable" down 

build:
	go build -o ./backend cmd/main.go

run:
	go run cmd/main.go

gooseup-extern:
	goose -dir internal/db/migrations postgres "postgres://postgres:g5HxxLc9n7ZXPdm1AB2FqyRPipsMTdr8pe6LjH6SOlPuStK1MoNmzezaViuJDOrP@88.198.203.75:5432/simptracker?sslmode=disable" up 

goosedown-extern:
	goose -dir internal/db/migrations postgres "postgres://postgres:g5HxxLc9n7ZXPdm1AB2FqyRPipsMTdr8pe6LjH6SOlPuStK1MoNmzezaViuJDOrP@88.198.203.75:5432/simptracker?sslmode=disable" down 
