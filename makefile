.PHONY: setup
setup: datafiles
	docker compose up -d

datafiles:
	docker run --security-opt seccomp=unconfined -v $(pwd)/data:/data ghcr.io/tigerbeetle/tigerbeetle format --cluster=0 --replica=0 --replica-count=3 /data/0_0.tigerbeetle
	docker run --security-opt seccomp=unconfined -v $(pwd)/data:/data ghcr.io/tigerbeetle/tigerbeetle format --cluster=0 --replica=1 --replica-count=3 /data/0_1.tigerbeetle
	docker run --security-opt seccomp=unconfined -v $(pwd)/data:/data ghcr.io/tigerbeetle/tigerbeetle format --cluster=0 --replica=2 --replica-count=3 /data/0_2.tigerbeetle

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
