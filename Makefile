fmt:
	goimports -w -l .

test: tidy
	go test -race -failfast ./...

test.bench:
	go test -test.bench=. -test.benchmem ./...

cover:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

tidy:
	go mod tidy

dep:
	go get -u -t ./...