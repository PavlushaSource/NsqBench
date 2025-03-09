# NsqBench

go clean -testcache && go test -bench=BenchmarkNSQMessageNow -benchmem -count=1 -benchtime=1x benchmarks/bench_test.go

go clean -testcache && go test -bench=BenchmarkNSQMessageOpt -benchmem -count=1 -benchtime=1x benchmarks/bench_test.go


<!-- Старт NSQ -->
sudo docker compose up -d

<!-- Остановка NSQ -->
sudo docker compose down


<!-- ping nsqlookupd -->
curl http://0.0.0.0:4151/ping

<!-- Запусн бенчмарков -->
go clean -testcache && go test -bench=BenchmarkNSQProducerConsumerBig -benchmem -count=1 -benchtime=1x tests/nsq_benchmarks_test.go
go clean -testcache && go test -bench=BenchmarkNSQProducerConsumerMedium -benchmem -count=1 -benchtime=1x tests/nsq_benchmarks_test.go
go clean -testcache && go test -bench=BenchmarkNSQProducerConsumerSmall -benchmem -count=1 -benchtime=1x tests/nsq_benchmarks_test.go

<!-- админка -->
http://localhost:4171/