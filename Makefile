.PHONY: run run-inventory run-order run-payment build tidy

run: ## Run all three services in parallel (Ctrl+C stops all)
	$(MAKE) -j3 run-inventory run-order run-payment

run-inventory:
	go run ./inventory/cmd

run-order:
	go run ./order/cmd

run-payment:
	go run ./payment/cmd

build:
	go build -o bin/inventory ./inventory/cmd
	go build -o bin/order ./order/cmd
	go build -o bin/payment ./payment/cmd

tidy:
	cd inventory && go mod tidy
	cd order && go mod tidy
	cd payment && go mod tidy
	cd shared && go mod tidy
