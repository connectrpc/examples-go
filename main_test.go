package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bufbuild/connect-demo/internal/gen/connect-go/buf/connect/demo/eliza/v1/elizav1connect"
	elizav1 "github.com/bufbuild/connect-demo/internal/gen/go/buf/connect/demo/eliza/v1"
	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"
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
	t.Run("connect_say", func(t *testing.T) { // nolint: paralleltest
		result, err := connectClient.Say(context.Background(), connect.NewRequest(&elizav1.SayRequest{
			Sentence: "Hello",
		}))
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.NotNil(t, result)
		assert.NotNil(t, result.Msg)
		assert.True(t, len(result.Msg.Sentence) > 0)
	})
	t.Run("connect_converse", func(t *testing.T) { // nolint: paralleltest
		stream := connectClient.Converse(context.Background())

		for i := 0; i < 10; i++ {
			sentence := "hello"
			if i == 9 {
				sentence = "bye"
			}
			err := stream.Send(&elizav1.ConverseRequest{Sentence: sentence})
			if err != nil && errors.Is(err, io.EOF) {
				return
			} else if err != nil {
				assert.FailNow(t, err.Error())
			}
		}

		// Listen and count responses.
		var responseCount int
		for {
			result, err := stream.Receive()
			if err != nil && errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				assert.FailNow(t, err.Error())
			}
			if len(result.Sentence) > 0 {
				responseCount++
			}
		}
		err := stream.CloseReceive()
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		err = stream.CloseSend()
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, responseCount, 10)
	})
}
