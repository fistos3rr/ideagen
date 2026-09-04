include .env

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
	#psql -h localhost -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME)
	psql "postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)"

## redis/cli: connect to redis cli
.PHONY: redis/cli
redis/cli:
	docker compose exec -it ideagen-redis redis-cli -a $(REDIS_PASSWORD)

# ================================================================== #
# DOCKER
# ================================================================== #

## docker/rebuild: rebuild service
.PHONY: docker/rebuild
docker/rebuild:
	docker compose up -d --build

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
