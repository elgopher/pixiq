test:
	go test -race -v ./...

lint:
	golint -set_exit_status ./...

test-ci:
	./ci/run-tests.sh
