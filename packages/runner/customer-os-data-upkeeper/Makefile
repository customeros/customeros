.PHONY: build migrate run stop tidy

build:
	go build -o app .

migrate:
	go run main.go migrate

run:
	go run main.go

stop:
	@-pkill -f "customer-os-data-upkeeper"

tidy:
	go mod tidy

# Helper to show all available commands
help:
	@echo "Available commands:"
	@echo "  build    - Build the application"
	@echo "  migrate  - Run database migrations"
	@echo "  run      - Run the application"
	@echo "  stop     - Stop the running application"
	@echo "  tidy     - Run go mod tidy"
