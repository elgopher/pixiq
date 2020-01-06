test:
	go test -race -v ./...

lint:
	./scripts/gofmt-check.sh
	golint -set_exit_status ./...

test-ci:
	./ci/run-tests.sh
