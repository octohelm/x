fmt:
	goimports -w -l .

test: tidy
	CGO_ENABLED=0 go test -v -failfast ./...

test.bench:
	CGO_ENABLED=0 go test -test.bench=. -test.benchmem ./...

cover:
	CGO_ENABLED=1 go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

tidy:
	go mod tidy

dep:
	go get -u -t ./...