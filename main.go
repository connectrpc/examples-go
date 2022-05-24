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
	"net/http"

	"github.com/bufbuild/connect-demo/internal/eliza"
	"github.com/bufbuild/connect-demo/internal/gen/connect-go/buf/connect/demo/eliza/v1/elizav1connect"
	elizav1 "github.com/bufbuild/connect-demo/internal/gen/go/buf/connect/demo/eliza/v1"
	connect "github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// ElizaServer implements some trivial business logic. The Protobuf
// definition for this API is in proto/buf/connect/demo/eliza/v1/eliza.proto.
type ElizaServer struct {
	elizav1connect.ElizaServiceHandler
}

// Say is a unary request demo. This method should allow for a one sentence
// response given a one sentence request.
func (e *ElizaServer) Say(ctx context.Context, req *connect.Request[elizav1.SayRequest]) (*connect.Response[elizav1.SayResponse], error) {
	sentenceInput := req.Msg.Sentence
	reply, _ := eliza.Reply(sentenceInput)
	res := connect.NewResponse(&elizav1.SayResponse{
		Sentence: reply,
	})
	return res, nil
}

// Converse is a bi-directional request demo. This method should allow for
// many requests and many responses.
func (e *ElizaServer) Converse(ctx context.Context, stream *connect.BidiStream[elizav1.ConverseRequest, elizav1.ConverseResponse]) error {
	for {
		receive, err := stream.Receive()
		if err != nil {
			return err
		}
		sentenceInput := receive.Sentence
		reply, endSession := eliza.Reply(sentenceInput)
		err = stream.Send(&elizav1.ConverseResponse{
			Sentence: reply,
		})
		if err != nil {
			return err
		}
		if endSession {
			return nil
		}
	}
}

func main() {
	// The business logic here is trivial, but the rest of the example is meant
	// to be somewhat realistic. This server has basic timeouts configured, and
	// it also exposes gRPC's server reflection and health check APIs.

	// protoc-gen-connect-go generates constructors that return plain net/http
	// Handlers, so they're compatible with most Go HTTP routers and middleware
	// (for example, net/http's StripPrefix). Each handler automatically supports
	// the Connect, gRPC, and gRPC-Web protocols.
	mux := http.NewServeMux()
	elizaServiceHandler := &ElizaServer{} // our business logic
	path, handler := elizav1connect.NewElizaServiceHandler(elizaServiceHandler)
	mux.Handle(path, handler)
	_ = http.ListenAndServe(
		"localhost:9000", // TODO: use PORT from terraform or default to 8080
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
