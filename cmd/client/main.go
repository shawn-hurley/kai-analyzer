package main

import (
	"fmt"
	"log/slog"
	"net/rpc"
	"os"
	"os/exec"

	"github.com/go-logr/logr"
	"github.com/konveyor/kai-analyzer/pkg/codec"
	"github.com/konveyor/kai-analyzer/pkg/service"
)

func main() {
	// Create a connnection to the server

	file, err := os.Create("testing.log")
	if err != nil {
		panic(err)
	}
	logger := slog.NewJSONHandler(file, nil)
	l := logr.FromSlogHandler(logger)

	cmd := exec.Command("./analyzer-lsp-jsonrpc-server", "-source-directory", "/Users/shurley/repos/kai/demo-apps/coolstore", "-rules-directory", "/Users/shurley/repos/MTA/rulesets/default/generated")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", cmd.Args)

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	c := rpc.NewClientWithCodec(codec.NewCodec(codec.Connection{Input: stdout, Output: stdin}, l))

	response := service.Response{}
	err = c.Call("analysis_engine.Analyze", service.Args{
		LabelSelector: "konveyor.io/target=cloud-readiness",
	}, &response)
	if err != nil {
		panic(err)
	}

	l.Info("got result", "r", response)

}
