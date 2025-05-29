// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package main implements a simple HTTP server.
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	gateway "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/unix4ever/fe-test-service/api/logs"
	"github.com/unix4ever/fe-test-service/internal/pkg/logging"
	"github.com/unix4ever/fe-test-service/internal/pkg/services"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("server run failed %s", err)
	}

	log.Printf("the server was stopped gracefully")
}

func run() error {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return err
	}

	defer lis.Close() //nolint:errcheck

	grpcServer := grpc.NewServer()

	logsService := services.NewLogsService()

	logs.RegisterLogsServiceServer(grpcServer, logsService)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	marshaller := &gateway.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:  true,
			UseEnumNumbers: true,
		},
	}

	runtimeMux := gateway.NewServeMux(
		gateway.WithMarshalerOption(gateway.MIMEWildcard, marshaller),
	)

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// we are proxying requests to ourselves, so we don't need to impose a limit
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)),
		grpc.WithSharedWriteBuffer(true),
	}

	if err = logs.RegisterLogsServiceHandlerFromEndpoint(ctx, runtimeMux, lis.Addr().String(), dialOpts); err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return grpcServer.Serve(lis)
	})

	gatewayServer := &http.Server{
		Addr:    "localhost:12000",
		Handler: logging.NewHandler(runtimeMux),
	}

	eg.Go(func() error {
		return gatewayServer.ListenAndServe()
	})

	log.Printf("the API server is running on the address 0.0.0.0:12000")

	<-ctx.Done()

	grpcServer.Stop()

	if err := gatewayServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown failed %w", err)
	}

	return nil
}
