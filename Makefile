.SILENT: test1

test-update:
	@git fetch template
	@git checkout template/main .github

test1: 
	rm -f ./.runtime/test-incr1.log
	rm -f ./.runtime/server
	go build -o ./.runtime/server ./cmd/server/*.go
	metricstest -test.v -test.run=^TestIteration1$ -binary-path=./.runtime/server > ./.runtime/test-incr1.log 