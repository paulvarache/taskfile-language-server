package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"yaml/extension"
	"yaml/jsonrpc"
	"yaml/lsp"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "try" {
		cmd := exec.Command(os.Args[0])
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		out, err := cmd.StdinPipe()
		if err != nil {
			panic(err)
		}
		err = cmd.Start()
		if err != nil {
			panic(err)
		}
		ctnt := `{"id": 12, "method": "try"}`
		out.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(ctnt), ctnt)))
		return
	}

	logfile := flag.String("logfile", "", "log to this file")
	traceEnabled := flag.Bool("trace", false, "print all requests and responses")

	flag.Parse()

	var output io.Writer = os.Stderr

	if *logfile != "" {
		f, err := os.OpenFile(*logfile, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		output = f
	}

	logger := log.New(output, "", 0)

	reader := os.Stdin
	writer := os.Stdout

	impl := &extension.TaskfileExtension{}

	s := jsonrpc.NewServer()

	if *traceEnabled {
		// Override the Discard output and provide the same output as the logger
		s.Logger.SetOutput(output)
	}

	_ = lsp.NewServer(s, impl, logger)

	s.AddHandler("extension/getTasks", impl.GetTasks)

	logger.Println("Listening to Stdin")
	s.Listen(reader, writer)
}
