
build: test

fmt:
	go vet ./...

cover: test
	go tool cover -html=./cover.out

test: fmt
	rm -rf ./cover.out
	go test -race -coverpkg=. -coverprofile=./cover.out ./...

fast:
	go test -failfast ./...

bench:
	go test -bench=.

