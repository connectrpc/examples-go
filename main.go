// Copyright 2022 Buf Technologies, Inc.
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
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bufbuild/connect-demo/internal/eliza"
	elizav1 "github.com/bufbuild/connect-demo/internal/gen/buf/connect/demo/eliza/v1"
	"github.com/bufbuild/connect-demo/internal/gen/buf/connect/demo/eliza/v1/elizav1connect"
	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/rs/cors"
	"github.com/spf13/pflag"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type elizaServer struct {
	// The time to sleep between sending responses on a stream
	streamDelay time.Duration
}

// NewElizaServer returns a new elizaServer.  streamDelay applies to server-streaming and will delay the responses
// sent on a stream by the given duration.
func NewElizaServer(streamDelay time.Duration) elizav1connect.ElizaServiceHandler {
	return &elizaServer{
		streamDelay: streamDelay,
	}
}

func (e *elizaServer) Say(
	ctx context.Context,
	req *connect.Request[elizav1.SayRequest],
) (*connect.Response[elizav1.SayResponse], error) {
	reply, _ := eliza.Reply(req.Msg.Sentence) // ignore end-of-conversation detection
	return connect.NewResponse(&elizav1.SayResponse{
		Sentence: reply,
	}), nil
}

func (e *elizaServer) Converse(
	ctx context.Context,
	stream *connect.BidiStream[elizav1.ConverseRequest, elizav1.ConverseResponse],
) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		request, err := stream.Receive()
		if err != nil && errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return fmt.Errorf("receive request: %w", err)
		}
		reply, endSession := eliza.Reply(request.Sentence)
		if err := stream.Send(&elizav1.ConverseResponse{Sentence: reply}); err != nil {
			return fmt.Errorf("send response: %w", err)
		}
		if endSession {
			return nil
		}
	}
}

func newCORS() *cors.Cors {
	// To let web developers play with the demo service from browsers, we need a
	// very permissive CORS setup.
	return cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowOriginFunc: func(origin string) bool {
			// Allow all origins, which effectively disables CORS.
			return true
		},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			// Content-Type is in the default safelist.
			"Accept",
			"Accept-Encoding",
			"Accept-Post",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Content-Encoding",
			"Grpc-Accept-Encoding",
			"Grpc-Encoding",
			"Grpc-Message",
			"Grpc-Status",
			"Grpc-Status-Details-Bin",
		},
		// Let browsers cache CORS information for longer, which reduces the number
		// of preflight requests. Any changes to ExposedHeaders won't take effect
		// until the cached data expires. FF caps this value at 24h, and modern
		// Chrome caps it at 2h.
		MaxAge: int(2 * time.Hour / time.Second),
	})
}

func main() {
	helpArg := pflag.BoolP("help", "h", false, "")
	streamDelayArg := pflag.DurationP(
		"server-stream-delay",
		"d",
		0,
		"The duration to delay sending responses on the server stream.",
	)
	pflag.Parse()

	if *helpArg {
		pflag.PrintDefaults()
		return
	}

	if *streamDelayArg < 0 {
		log.Printf("Server stream delay cannot be negative.")
		return
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/",
		http.RedirectHandler("https://connect.build/demo", http.StatusFound),
	)
	compress1KB := connect.WithCompressMinBytes(1024)
	mux.Handle(elizav1connect.NewElizaServiceHandler(
		NewElizaServer(*streamDelayArg),
		compress1KB,
	))
	mux.Handle(grpchealth.NewHandler(
		grpchealth.NewStaticChecker(elizav1connect.ElizaServiceName),
		compress1KB,
	))
	mux.Handle(grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(elizav1connect.ElizaServiceName),
		compress1KB,
	))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(elizav1connect.ElizaServiceName),
		compress1KB,
	))
	addr := "localhost:8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	srv := &http.Server{
		Addr: addr,
		Handler: h2c.NewHandler(
			newCORS().Handler(mux),
			&http2.Server{},
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP listen and serve: %v", err)
		}
	}()

	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown: %v", err) //nolint:gocritic
	}
}

func (e *elizaServer) Introduce(
	ctx context.Context,
	req *connect.Request[elizav1.IntroduceRequest],
	stream *connect.ServerStream[elizav1.IntroduceResponse],
) error {
	name := req.Msg.Name
	if name == "" {
		name = "Anonymous User"
	}
	intros := eliza.GetIntroResponses(name)
	var ticker *time.Ticker
	if e.streamDelay > 0 {
		ticker = time.NewTicker(e.streamDelay)
		defer ticker.Stop()
	}
	for _, resp := range intros {
		if ticker != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
			}
		}
		if err := stream.Send(&elizav1.IntroduceResponse{Sentence: resp}); err != nil {
			return err
		}
	}
	return nil
}
