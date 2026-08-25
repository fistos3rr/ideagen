# ================================================================== #
# HELPERS
# ================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

# ================================================================== #
# DEVELOPMENT
# ================================================================== #

## git/merge: merging dev in main
.PHONY: git/merge
git/merge: confirm
	git checkout dev
	git rebase -i main
	git push origin dev --force-with-lease
	git checkout main
	git merge dev
	git push origin main

## git/pull: pull from github
.PHONY: git/pull
git/pull:
	git checkout main
	git pull
	git checkout dev
	git fetch origin dev && git reset --hard origin/dev
	
## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}
	
## db/psql: connect to db in docker
.PHONY: db/psql
db/psql:
	psql -h localhost -p 5432 -U ideagenuser -d ideagendb

# ================================================================== #
# DOCKER
# ================================================================== #

## docker/up: up service
.PHONY: docker/up
docker/up:
	docker compose up -d

## docker/down: down service
.PHONY: docker/down
docker/down:
	docker compose down

## docker/clean: clean all data with db data
.PHONY: docker/clean
docker/clean:
	docker compose down -v