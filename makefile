include ./.env

DBURL=postgres://$(DBUSER):$(DBPASS)@$(DBHOST):$(DBPORT)/$(DBNAME)?sslmode=disable
MIGRATIONPATH=db/migrations
SEEDPATH=db/seeds

# bikin migration baru
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONPATH) -seq create_$(NAME)_table

# jalankan semua migration
migrate-up:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) up

# rollback 1 step
migrate-down:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) down 1

# isi data dummy dari seed.sql
insert-seed:
	psql $(DBURL) -f $(SEEDPATH)/seed.sql

# full migration (up + seed)
migrate-full:
	make migrate-up insert-seed


