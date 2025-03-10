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
	"fmt"
	"github.com/PavlushaSource/NsqBench/internal/core/port"
	"github.com/cloudwego/shmipc-go"
	"net"
	"os"
	"path/filepath"
	"syscall"
)

type ServiceResponse struct {
	server     *shmipc.Session
	con        net.Conn
	iterations int
}

func (s *ServiceResponse) Run() error {
	for i := 0; i < s.iterations; i++ {
		stream, err := s.server.AcceptStream()
		if err != nil {
			return fmt.Errorf("accept stream failed %w", err)
		}

		reader := stream.BufferReader()
		reqData, err := reader.ReadBytes(len("client say hello world!!!"))
		if err != nil {
			return fmt.Errorf("reqBuf readData failed %w", err)
		}
		fmt.Println("server receive request message:" + string(reqData))

		respMsg := "server hello world!!!"
		writer := stream.BufferWriter()
		err = writer.WriteString(respMsg)
		if err != nil {
			return fmt.Errorf("respBuf WriteString failed %w", err)
		}

		err = stream.Flush(true)
		if err != nil {
			return fmt.Errorf("stream flush failed %w", err)
		}
		stream.Close()
	}

	return nil
}

func (s *ServiceResponse) Close() error {
	if err := s.server.Close(); err != nil {
		return fmt.Errorf("close server failed %w", err)
	}
	return s.con.Close()
}

func NewServiceResponse(iter int) (port.ResponseService, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getwd failed %w", err)
	}

	udsPath := filepath.Join(dir, "../ipc_test.sock")
	_ = syscall.Unlink(udsPath)
	ln, err := net.ListenUnix("unix", &net.UnixAddr{Name: udsPath, Net: "unix"})
	if err != nil {
		return nil, fmt.Errorf("listen unix failed %w", err)
	}

	conn, err := ln.Accept()
	if err != nil {
		return nil, fmt.Errorf("accept unix failed %w", err)
	}

	conf := shmipc.DefaultConfig()
	s, err := shmipc.Server(conn, conf)
	if err != nil {
		return nil, fmt.Errorf("new server failed %w", err)
	}

	return &ServiceResponse{s, conn, iter}, nil
}
