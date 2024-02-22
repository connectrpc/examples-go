// Copyright 2022-2024 The Connect Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	elizav1 "connect-examples-go/internal/gen/connectrpc/eliza/v1"
	"connect-examples-go/internal/gen/connectrpc/eliza/v1/elizav1connect"
)

func TestElizaServer(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.Handle(elizav1connect.NewElizaServiceHandler(
		NewElizaServer(0),
	))
	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	connectClient := elizav1connect.NewElizaServiceClient(
		server.Client(),
		server.URL,
	)
	grpcClient := elizav1connect.NewElizaServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPC(),
	)
	clients := []elizav1connect.ElizaServiceClient{connectClient, grpcClient}

	t.Run("say", func(t *testing.T) {
		for _, client := range clients {
			result, err := client.Say(context.Background(), connect.NewRequest(&elizav1.SayRequest{
				Sentence: "Hello",
			}))
			require.NoError(t, err)
			assert.NotEmpty(t, result.Msg.GetSentence())
		}
	})
	t.Run("converse", func(t *testing.T) {
		for _, client := range clients {
			sendValues := []string{"Hello!", "How are you doing?", "I have an issue with my bike", "bye"}
			var receivedValues []string
			stream := client.Converse(context.Background())
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				for _, sentence := range sendValues {
					err := stream.Send(&elizav1.ConverseRequest{Sentence: sentence})
					if assert.ErrorContainsf(t, err, `failed for string sentence: "%s"`, sentence) { 
						break
					}
				}
				assert.NoError(t, stream.CloseRequest()) 
			}()
			go func() {
				defer wg.Done()
				for {
					msg, err := stream.Receive()
					if errors.Is(err, io.EOF) {
						break
					}
					if !assert.NoError(t, err) {
						break
					}
					assert.NotEmpty(t, msg.GetSentence())
					receivedValues = append(receivedValues, msg.GetSentence())
				}
				assert.NoError(t, stream.CloseResponse()) 
			}()
			wg.Wait()
			assert.Equal(t, len(receivedValues), len(sendValues))
		}
	})
	t.Run("introduce", func(t *testing.T) {
		total := 0
		for _, client := range clients {
			request := connect.NewRequest(&elizav1.IntroduceRequest{
				Name: "Ringo",
			})
			stream, err := client.Introduce(context.Background(), request)
			require.NoError(t, err)
			for stream.Receive() {
				total++
			}
			require.NoError(t, stream.Err())
			require.NoError(t, stream.Close())
			assert.Greater(t, total, 0)
		}
	})
}
