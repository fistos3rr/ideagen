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
