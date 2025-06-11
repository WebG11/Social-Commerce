.PHONY: gen tidy test lint all

all: gen tidy test lint

gen:
	sh ./scripts/gen.sh

tidy:
	sh ./scripts/tidy.sh

test:
	sh ./scripts/test.sh

lint:
	sh ./scripts/lint.sh

# 测试命令
test-unit:
	go test -v ./tests/unit/...

test-functional:
	go test -v ./tests/functional/...

test-all: test-unit test-functional

test-coverage:
	go test -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out

.PHONY: run
run: ## run {svc} server. example: make run svc=product
	@scripts/run.sh ${svc}

.PHONY: gensvc
gensvc: ## generate thrift code
	@scripts/gensvc.sh ${svc}