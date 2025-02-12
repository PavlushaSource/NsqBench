# NsqBench

go clean -testcache && go test -bench=BenchmarkNSQMessageNow -benchmem -count=1 -benchtime=1x benchmarks/bench_test.go