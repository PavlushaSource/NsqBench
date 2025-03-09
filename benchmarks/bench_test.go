package benchmarks

import (
	"github.com/PavlushaSource/NsqBench/internal/adapters/nsq/optNsq"
	"github.com/PavlushaSource/NsqBench/internal/adapters/nsq/rawNsq"
	"github.com/PavlushaSource/NsqBench/internal/core/domain"
	"github.com/PavlushaSource/NsqBench/internal/core/port"
	"sync"
	"testing"
)

type testCase struct {
	name string
	iter int
}

var TestTable = []testCase{
	{
		name: "Sending 10 message",
		iter: 10,
	},
	//{
	//	name: "Sending 10^2 message",
	//	iter: 100,
	//},
	//{
	//	name: "Sending 10^7 message",
	//	iter: 10e7,
	//},
}

func RunNsq(requestService port.RequestService, responseService port.ResponseService, iter int) func(b *testing.B) {
	return func(b *testing.B) {
		b.ResetTimer()

		//defer requestService.Close()
		defer responseService.Close()

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			if err := responseService.Run(); err != nil {
				b.Error(err)
			}
			defer wg.Done()
		}()

		if err := requestService.Run(iter); err != nil {
			b.Error(err)
		}
		wg.Wait()
	}
}

func BenchmarkRawNsq(b *testing.B) {
	for _, tc := range TestTable {
		requester, err := rawNsq.NewServiceRequest("127.0.0.1:4161", "127.0.0.1:4150")
		if err != nil {
			b.Error(err)
		}
		responser, err := rawNsq.NewServiceResponse("127.0.0.1:4161", "127.0.0.1:4150", tc.iter)
		if err != nil {
			b.Error(err)
		}

		defer requester.Close()
		f := RunNsq(requester, responser, tc.iter)
		b.Run(tc.name, f)
	}
}

func BenchmarkOptNsq(b *testing.B) {
	for _, tc := range TestTable {
		requester, err := optNsq.NewServiceRequestOpt("127.0.0.1:4161", "127.0.0.1:4150", domain.ResponseTopic, domain.ResponseChannel)
		if err != nil {
			b.Error(err)
		}
		responser, err := rawNsq.NewServiceResponse("127.0.0.1:4161", "127.0.0.1:4150", tc.iter)
		if err != nil {
			b.Error(err)
		}

		f := RunNsq(requester, responser, tc.iter)
		b.Run(tc.name, f)
	}
}
