.DEFAULT_GOAL=run

.PHONY: run
run:
	@go run .

init:
	sudo docker compose up -d

run_requester:
	go run ./cmd/serviceRequest/main.go

run_responser:
	go run ./cmd/serviceResponse/main.go