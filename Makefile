.PHONY: start build clean frontend-dev frontend-build

binary_name = expense-tracker
binary_path = ./bin/$(binary_name)
main_file = ./cmd/app/main.go

start: $(binary_path)
	@echo "Start the apllication..."
	@$(binary_path)

build:
	@mkdir -p ./bin
	@go build -o $(binary_path) $(main_file)

clean:
	@rm -rf ./bin

frontend-dev:
	@cd web && npm run dev

frontend-build:
	@cd web && npm run build

frontend-install:
	@cd web && npm install

dev: build start

full-dev:
	@make dev && make frontend-dev

help:
	@echo "Available commands:"
	@echo ""
	@echo "Backend:"
	@echo "  make build    - Build the application"
	@echo "  make start    - Start the application" 
	@echo "  != make run      - Run without building"
	@echo "  make dev      - Build and start"
	@echo "  make clean    - Clean build files"
	@echo ""
	@echo "Frontend:"
	@echo "  make frontend-dev     - Start frontend dev server"
	@echo "  make frontend-build   - Build frontend for production"
	@echo "  make frontend-install - Install frontend dependencies"
	@echo ""
	@echo "Full-stack:"
	@echo "  make full-dev - Start both backend and frontend"
	@echo ""
	@echo "  make help     - Show this help"