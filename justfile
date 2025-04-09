fmt:
    go tool gofumpt -w -l .

test:
    CGO_ENABLED=0 go test -v -failfast ./...

test-bench:
    CGO_ENABLED=0 go test -test.bench=. -test.benchmem ./...

test-cover:
    CGO_ENABLED=1 go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

dep:
    go mod tidy

update:
    go get -u -t ./...
