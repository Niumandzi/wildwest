.PHONY: run

run:
	@echo "Loading environment variables..."; \
	. ./scripts/load_env.sh; \
	echo "Starting services..."; \
	docker-compose up -d

stop:
	@echo "Stopping services..."
	@docker-compose down

build:
	@echo "Building services..."
	@docker-compose build
