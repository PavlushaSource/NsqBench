/*
 * Copyright 2023 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package shmipc

import (
	"context"
	"fmt"
	"github.com/PavlushaSource/NsqBench/internal/core/port"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cloudwego/shmipc-go"
)

type ServiceRequest struct {
	manager *shmipc.SessionManager
}

func NewServiceRequest() (port.RequestService, error) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	udsPath := filepath.Join(dir, "../ipc_test.sock")

	// 1.create client session manager
	conf := shmipc.DefaultSessionManagerConfig()
	conf.ShareMemoryPathPrefix = "/dev/shm/client.ipc.shm"
	conf.Network = "unix"
	conf.Address = udsPath
	if runtime.GOOS == "darwin" {
		conf.ShareMemoryPathPrefix = "/tmp/client.ipc.shm"
		conf.QueuePath = "/tmp/client.ipc.shm_queue"
	}

	s, err := shmipc.NewSessionManager(conf)
	if err != nil {
		return nil, fmt.Errorf("new session manager failed %w", err)
	}

	return &ServiceRequest{s}, nil
}

func (s *ServiceRequest) Run(ctx context.Context, iterations int) error {
	stream, err := s.manager.GetStream()
	if err != nil {
		return fmt.Errorf("client get stream failed, %w", err)
	}
	for i := 0; i < iterations; i++ {
		err = s.Send(ctx, "hello", stream)
		if err != nil {
			return fmt.Errorf("client send failed, %w", err)
		}
	}

	defer s.manager.PutBack(stream)
	return nil
}

func (s *ServiceRequest) Send(ctx context.Context, msg string, stream *shmipc.Stream) error {
	writer := stream.BufferWriter()
	err := writer.WriteString(msg)
	if err != nil {
		return fmt.Errorf("buffer writeString failed %w", err)
	}

	fmt.Println("client stream send request:" + msg)
	err = stream.Flush(true)
	if err != nil {
		return fmt.Errorf("stream Flush failed %w", err)
	}

	reader := stream.BufferReader()
	respData, err := reader.ReadBytes(len("server hello world!!!"))
	if err != nil {
		return fmt.Errorf("respBuf ReadBytes failed %w", err)
	}

	fmt.Println("client stream receive response " + string(respData))
	return nil
}

func (s *ServiceRequest) Close() error {
	return s.manager.Close()
}
