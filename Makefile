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

## run/idea: run the idea generator service
.PHONY: run/idea
run/idea:
	go run ./cmd/idea

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

## db/psl: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${POSTGRES_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${POSTGRES_DB_DSN} up

# ================================================================== #
# BUILD
# ================================================================== #

## build/idea: build the cmd/idea application
.PHONY: build/idea
build/idea:
	@echo 'Building cmd/idea...'
	go build -ldflags='-s' -o=./bin/idea ./cmd/idea
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/idea ./cmd/idea
