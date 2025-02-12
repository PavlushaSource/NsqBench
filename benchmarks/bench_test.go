package benchmarks

import (
	"fmt"
	"github.com/PavlushaSource/NsqBench/src/services/requestAdapters/serviceRequestNow"
	"github.com/PavlushaSource/NsqBench/src/services/responseAdapters/serviceResponseNow"
	"sync"
	"testing"
)

type testCase struct {
	name string
	iter int
}

var TestTable = []testCase{
	{
		name: "Sending 10^2 message",
		iter: 3,
	},
	//{
	//	name: "Sending 10^5 message",
	//	iter: 10e5,
	//},
	//{
	//	name: "Sending 10^7 message",
	//	iter: 10e7,
	//},
}

func BenchmarkNSQMessageNow(b *testing.B) {
	for _, tc := range TestTable {
		fmt.Println("TEST " + tc.name)
		b.Run(tc.name, func(b *testing.B) {
			requester, err := serviceRequestNow.NewServiceRequest("127.0.0.1:4161", "127.0.0.1:4150")
			if err != nil {
				b.Error(err)
			}
			responser, err := serviceResponseNow.NewServiceResponse("127.0.0.1:4161", "127.0.0.1:4150", tc.iter)
			if err != nil {
				b.Error(err)
			}

			b.ResetTimer()

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				if err := responser.Run(); err != nil {
					b.Error(err)
				}
				defer responser.Close()
				defer wg.Done()
			}()

			if err = requester.Run(tc.iter); err != nil {
				b.Error(err)
			}
			wg.Wait()
		})
	}
}
