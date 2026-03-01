ifneq (,$(wildcard .env))
include .env
export
endif

MIGRATE_IMAGE := migrate/migrate:v4.18.3
MIGRATE := docker run --rm -v $(PWD)/migrations:/migrations $(MIGRATE_IMAGE)

.PHONY: migrate-up migrate-down migrate-version

migrate-up:
	@test -n "$(DATABASE_URL)" || (echo "DATABASE_URL is required (set in .env)"; exit 1)
	$(MIGRATE) -path=/migrations -database "$(DATABASE_URL)" up

migrate-down:
	@test -n "$(DATABASE_URL)" || (echo "DATABASE_URL is required (set in .env)"; exit 1)
	$(MIGRATE) -path=/migrations -database "$(DATABASE_URL)" down 1

migrate-version:
	@test -n "$(DATABASE_URL)" || (echo "DATABASE_URL is required (set in .env)"; exit 1)
	$(MIGRATE) -path=/migrations -database "$(DATABASE_URL)" version
