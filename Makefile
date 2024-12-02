.SILENT: test1, test2

test-update:
	@git fetch template
	@git checkout template/main .github

test1: 
	rm -f ./.runtime/test-incr1.log
	rm -f ./.runtime/server
	go build -o ./.runtime/server ./cmd/server/*.go
	metricstest -test.v -test.run=^TestIteration1$ -binary-path=./.runtime/server > ./.runtime/test-incr1.log 

test2: 
	rm -f ./.runtime/test-incr2.log
	rm -f ./.runtime/agent
	go build -o ./.runtime/agent ./cmd/agent/*.go
	metricstest -test.v -test.run=^TestIteration2$ -agent-binary-path=./.runtime/agent > ./.runtime/test-incr2.log 

test3: 
	rm -f ./.runtime/test-incr3.log
	rm -f ./.runtime/agent
	rm -f ./.runtime/server
	go build -o ./.runtime/agent ./cmd/agent/*.go
	go build -o ./.runtime/server ./cmd/server/*.go
	metricstest -test.v -test.run=^TestIteration3$ -agent-binary-path=./.runtime/agent -binary-path=./.runtime/server > ./.runtime/test-incr3.log 

test4: 
	rm -f ./.runtime/test-incr4.log
	rm -f ./.runtime/agent
	rm -f ./.runtime/server
	go build -o ./.runtime/agent ./cmd/agent/*.go
	go build -o ./.runtime/server ./cmd/server/*.go
	metricstest -test.v -test.run=^TestIteration4$ -agent-binary-path=./.runtime/agent -binary-path=./.runtime/server > ./.runtime/test-incr4.log 