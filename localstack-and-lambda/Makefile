init: stop
	docker-compose -f docker-compose.yml up -d
	sleep 1 # grace time
	make migrate

start:
	docker-compose up

stop:
	docker-compose -f docker-compose.yml down --remove-orphans

migrate:
	PGPASSWORD="Don't look!" psql -h localhost -p 5432 -U postgres -f migration.sql