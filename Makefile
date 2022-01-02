COVERAGEDIR = .coverage

test:
	go test -cover ./... 
	golangci-lint run ./...

test-coverage:
	if [ ! -d $(COVERAGEDIR) ]; then mkdir $(COVERAGEDIR); fi
	go test -coverpkg ./... -coverprofile $(COVERAGEDIR)/request.coverprofile ./...
	go tool cover -html $(COVERAGEDIR)/request.coverprofile

test-bench:
	go test -bench=. -benchtime=1000000x -benchmem | tee benchmark_results.txt
