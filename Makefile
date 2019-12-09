test:
	go test -race -v ./...

test-ci:
	./ci/run-tests.sh
