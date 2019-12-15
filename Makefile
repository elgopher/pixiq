test:
	go test -race -v ./...

lint:
	golint

test-ci:
	./ci/run-tests.sh
