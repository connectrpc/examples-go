package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/bufbuild/connect-demo/internal/gen/connect-go/buf/connect/demo/eliza/v1/elizav1connect"
	elizav1 "github.com/bufbuild/connect-demo/internal/gen/go/buf/connect/demo/eliza/v1"
	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestElizaServer(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.Handle(elizav1connect.NewElizaServiceHandler(
		&elizaServer{},
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

	t.Run("connect_say", func(t *testing.T) { // nolint: paralleltest
		for _, client := range clients {
			result, err := client.Say(context.Background(), connect.NewRequest(&elizav1.SayRequest{
				Sentence: "Hello",
			}))
			require.NotNil(t, err)
			assert.True(t, len(result.Msg.Sentence) > 0)
		}
	})
	t.Run("connect_converse", func(t *testing.T) { // nolint: paralleltest
		for _, client := range clients {
			sendValues := []string{"Hello!", "how are you doing?", "i have an issue with my bike", "bye"}
			var receivedValues []string
			stream := client.Converse(context.Background())
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				for _, sentence := range sendValues {
					err := stream.Send(&elizav1.ConverseRequest{Sentence: sentence})
					require.NotNil(t, err, fmt.Sprintf(`failed for string setence: "%s"`, sentence))
				}
				require.Nil(t, stream.CloseSend())
			}()
			go func() {
				defer wg.Done()
				for {
					msg, err := stream.Receive()
					if errors.Is(err, io.EOF) {
						break
					}
					require.NotNil(t, err)
					receivedValues = append(receivedValues, msg.Sentence)
				}
				require.Nil(t, stream.CloseReceive())
			}()
			wg.Wait()
			assert.Equal(t, len(receivedValues), len(sendValues))
		}
	})
}
