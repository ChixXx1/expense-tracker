.PHONY: start build

binary_name = expense-tracker
binary_path = ./bin/$(binary_name)
main_file = ./cmd/app/main.go

start: $(binary_path)
	@echo "Start the apllication..."
	@$(binary_path)

build:
	@mkdir -p ./bin
	@go build -o $(binary_path) $(main_file)
