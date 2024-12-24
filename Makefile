.SILENT: test1, test2

test-update:
	@git fetch template
	@git checkout template/main .github

test-cover:
	@go test -cover ./...

test7: 
	rm -f ./.runtime/agent
	rm -f ./.runtime/server
	go build -o ./.runtime/agent ./cmd/agent/*.go
	go build -o ./.runtime/server ./cmd/server/*.go
	@metricstest -test.v -test.run=^TestIteration7$\ -agent-binary-path=.runtime/agent -binary-path=.runtime/server -source-path=./ > ./.runtime/test.log

test8: 
	rm -f ./.runtime/agent
	rm -f ./.runtime/server
	go build -o ./.runtime/agent ./cmd/agent/*.go
	go build -o ./.runtime/server ./cmd/server/*.go
	@metricstest -test.v -test.run=^TestIteration8$\ -server-port=8080 -agent-binary-path=.runtime/agent -binary-path=.runtime/server -source-path=./ > ./.runtime/test.log

test9: 
	rm -f ./.runtime/agent
	rm -f ./.runtime/server
	go build -o ./.runtime/agent ./cmd/agent/*.go
	go build -o ./.runtime/server ./cmd/server/*.go
	@metricstest -test.v -test.run=^TestIteration9$\ -file-storage-path=./BACKUP -server-port=8080 -agent-binary-path=.runtime/agent -binary-path=.runtime/server -source-path=./ > ./.runtime/test.log

lint:
	@golangci-lint run