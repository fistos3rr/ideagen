include .envrc

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
