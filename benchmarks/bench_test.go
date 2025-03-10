package benchmarks

import (
	"context"
	"fmt"
	"github.com/PavlushaSource/NsqBench/internal/adapters/monolith/channels"
	"github.com/PavlushaSource/NsqBench/internal/adapters/monolith/shmipc"
	"github.com/PavlushaSource/NsqBench/internal/adapters/nsq/optimizationNsq"
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
	//{
	//	name: "Sending 10 message",
	//	iter: 10,
	//},
	{
		name: "Sending 10^3 message",
		iter: 100,
	},
	////{
	//	name: "Sending 10^4 message",
	//	iter: 10e4,
	//},
}

func RunNsq(requestService port.RequestService, responseService port.ResponseService, iter int) func(b *testing.B) {

	return func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		b.ResetTimer()

		defer responseService.Close()

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			if err := responseService.Run(); err != nil {
				b.Error(err)
			}
			defer wg.Done()
		}()

		if err := requestService.Run(ctx, iter); err != nil {
			b.Error(err)
		}
		_ = cancel
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
		requester, err := optimizationNsq.NewServiceRequest("127.0.0.1:4161", "127.0.0.1:4150")
		if err != nil {
			b.Error(err)
		}
		responser, err := rawNsq.NewServiceResponse("127.0.0.1:4161", "127.0.0.1:4150", tc.iter)
		if err != nil {
			b.Error(err)
		}

		f := RunNsq(requester, responser, tc.iter)
		b.Run(tc.name, f)

		requester.Close()
	}
}

func RunMonolith(requestService port.RequestService, responseService port.ResponseService, iter int) func(b *testing.B) {
	return func(b *testing.B) {
		b.ResetTimer()

		go func() {
			if err := responseService.Run(); err != nil {
				b.Error(err)
			}
		}()

		if err := requestService.Run(context.Background(), iter); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkMonolithChannels(b *testing.B) {
	for _, tc := range TestTable {
		responseCh := make(chan domain.Message)
		requestCh := make(chan domain.Message)

		requester := channels.NewServiceRequest(responseCh, requestCh)
		responser := channels.NewServiceResponse(requestCh, responseCh)

		f := RunMonolith(requester, responser, tc.iter)
		b.Run(tc.name, f)

		// not need close channels
		//responser.Close()
		//requester.Close()
	}
}

func BenchmarkShmipc(b *testing.B) {
	for _, tc := range TestTable {
		responser, err := shmipc.NewServiceResponse(tc.iter)
		if err != nil {
			b.Error(err)
		}
		fmt.Println("HEY")

		requester, err := shmipc.NewServiceRequest()
		if err != nil {
			b.Error(err)
		}

		f := RunMonolith(requester, responser, tc.iter)
		b.Run(tc.name, f)
	}
}
