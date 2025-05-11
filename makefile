.PHONY: dev db-up db-down

# roda o servidor Go em modo dev
dev:
	go run cmd/server/main.go

# sobe apenas o serviço de banco de dados
db-up:
	docker-compose up -d db

# desliga todos os serviços
db-down:
	docker-compose down
