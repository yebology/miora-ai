SHELL := /bin/bash

# --- Development (local) ---

run-fe:
	cd frontend && npm run dev

run-be:
	cd backend && go run main.go

run-agent:
	cd agent && python main.py

run-all:
	@make run-be & make run-fe & make run-agent

# --- Database ---

db-reset:
	cd backend && go run cmd/reset/main.go

db-seed:
	cd backend && go run cmd/seed/main.go

# --- EAS ---

register-schema:
	cd backend && go run cmd/register-schema/main.go

# --- MockUSDT ---

deploy-musdt:
	cd contracts && source .env && forge script script/DeployMockUSDT.s.sol --rpc-url $$BASE_SEPOLIA_RPC_URL --broadcast

# --- Agent setup ---

setup-agent:
	cd agent && pip install -r requirements.txt

# --- Docker (all services) ---

docker-up:
	docker compose build --no-cache
	docker compose up -d --build

docker-down:
	docker compose down

docker-reset:
	docker compose down -v && docker compose up -d

docker-be:
	docker compose up -d --build backend

docker-fe:
	docker compose up -d --build frontend

docker-agent:
	docker compose up -d --build agent

docker-db:
	docker compose up -d db

start:
	docker compose up -d --build

# --- Git ---

commit:
	@git add .
	@git status
	@read -p "Commit message: " msg; \
	git commit -m "$msg"
