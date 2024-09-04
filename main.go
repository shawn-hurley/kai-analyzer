package main

import (
	"log/slog"
	"net/rpc"
	"os"

	"github.com/go-logr/logr"
	"github.com/konveyor/kai-analyzer/pkg/codec"
	"github.com/konveyor/kai-analyzer/pkg/service"
)

func main() {

	// In the future add cobra for flags maybe
	// create log file in working directory for now.

	file, err := os.Create("kai-analyzer.log")
	if err != nil {
		panic(err)
	}
	logger := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.Level(-100),
	})

	l := logr.FromSlogHandler(logger)
	// We need to start up the JSON RPC server and start listening for messages
	analyzerService, err := service.NewAnalyzer(10000, 10, 10, "/Users/shurley/repos/kai/demo-apps/coolstore", "", []string{"/Users/shurley/repos/MTA/rulesets/default/generated"}, l)
	if err != nil {
		panic(err)
	}
	server := rpc.NewServer()
	err = server.RegisterName("analysis_engine", analyzerService)
	if err != nil {
		panic(err)
	}

	codec := codec.NewCodec(codec.Connection{Input: os.Stdin, Output: os.Stdout}, l)
	l.Info("Starting Server")
	server.ServeRequest(codec)
}
