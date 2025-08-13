.PHONY: setup
setup: datafiles
	docker compose up -d

datafiles:
	./scripts/setup.sh

.PHONY: teardown
teardown:
	-docker compose down --remove-orphans
	-docker compose rm -f -v

.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	go test ./...
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	COVER_THRESHOLD=0.0 ./lint-project.sh
endif
