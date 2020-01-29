package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"taskfile-language-server/extension"
	"taskfile-language-server/jsonrpc"
	"taskfile-language-server/lsp"
)

// These variables are provided at build time using ldflags
var (
	BuildVersion string = "dev"
	BuildHash    string = "dev"
)

func main() {
	version := flag.Bool("version", false, "display the Language Server version")
	logfile := flag.String("logfile", "", "log to this file")
	traceEnabled := flag.Bool("trace", false, "print all requests and responses")

	flag.Parse()

	// Display version if asked for it
	if *version {
		fmt.Printf("Taskfile Language Server version %s-%s", BuildVersion, BuildHash)
		return
	}

	var output io.Writer = os.Stderr
	// Output to a file if a path is provided
	if *logfile != "" {
		f, err := os.OpenFile(*logfile, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		output = f
	}

	logger := log.New(output, "", log.Ldate|log.Ltime)

	reader := os.Stdin
	writer := os.Stdout

	// Create the taskfile implementation of the LSP
	impl := extension.New()
	// Create the jsonrpc server
	s := jsonrpc.NewServer()

	// Override the Discard output and provide the same output as the logger
	if *traceEnabled {
		s.Logger.SetOutput(output)
		impl.Logger.SetOutput(output)
	}

	// Create the LSP Server
	_ = lsp.NewServer(s, impl, logger)
	s.Listen(reader, writer)
}
