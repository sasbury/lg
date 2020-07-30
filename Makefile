
build: test

cover: test
	go tool cover -html=./cover.out

test:
	go vet ./...
	rm -rf ./cover.out
	go test -race -coverpkg=. -coverprofile=./cover.out ./...

fast:
	go test -failfast ./...

bench:
	go test -bench=.

