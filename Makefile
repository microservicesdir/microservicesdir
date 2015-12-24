build: clean
	gb build all

clean:
	rm -rf bin

test: dev
	gb test -v

prepare-db: dev
	scripts/setup_db.sh
	touch $@

dev:
	docker-compose up -d
