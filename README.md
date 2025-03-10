## Repository for benchmarks of different mechanisms for synchronous message sending 

Для запуска бенчмарков, использующих NSQ
```bash
make init
```

## NSQ панель

#### ping nsqlookupd
```bash
curl http://0.0.0.0:4151/ping
```

### admin panel
```bash
http://localhost:4171/
```
---


## Запуск бенчмарков

#### Raw NSQ 
```bash
go clean -testcache && go test ./benchmarks/. -bench=Raw -benchtime=1x -count=1
```
#### Optimization NSQ
```bash
go clean -testcache && go test ./benchmarks/. -bench=Opt -benchtime=1x -count=1
```

#### Monolith go channels
```bash
go clean -testcache && go test ./benchmarks/. -bench=MonolithChannels -benchtime=1x -count=1
```

#### Monolith shmipc-go
```bash
go clean -testcache && go test ./benchmarks/. -bench=Shmipc -benchtime=1x -count=1
```

---

