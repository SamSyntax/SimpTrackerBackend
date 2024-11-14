gooseup:
	goose -dir internal/db/migrations postgres "postgres://postgres:postgres@localhost:5432/stuleja?sslmode=disable" up 

goosedown:
	goose -dir internal/db/migrations postgres "postgres://postgres:postgres@localhost:5432/stuleja?sslmode=disable" down 

build:
	go build -o ./backend cmd/main.go

run:
	go run cmd/main.go

gooseup-extern:
	goose -dir internal/db/migrations postgres "postgres://postgres:NddjW8Aj2fXdamACj4VrokZDGmUT0WhqM3cvR2qnLxf4IDLNhkr30YuFYqMhatsC@88.198.203.75:5432/postgres?sslmode=disable" up 

goosedown-extern:
	goose -dir internal/db/migrations postgres "postgres://postgres:NddjW8Aj2fXdamACj4VrokZDGmUT0WhqM3cvR2qnLxf4IDLNhkr30YuFYqMhatsC@88.198.203.75:5432/postgres?sslmode=disable" down 
