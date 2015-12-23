build:
	gb build all

test: dev prepare-db
	gb test -v

prepare-db:
	scripts/setup_db.sh
	touch $@

dev:
	docker-compose up -d
