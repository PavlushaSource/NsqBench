.DEFAULT_GOAL=run

.PHONY: run
run:
	@go run .

init:
	sudo docker compose up -d