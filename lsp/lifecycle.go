package lsp

import (
	"encoding/json"
	"os"
	"taskfile-language-server/jsonrpc"
)

type LifecycleShutdown interface {
	OnShutdown() *jsonrpc.ResponseError
}

type LifecycleExit interface {
	OnExit()
}

// ShutdownHandler notifies the implementation that the server was requested to shutdown
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#shutdown
func (s *LSPServer) ShutdownHandler(params json.RawMessage) (interface{}, *jsonrpc.ResponseError) {
	// A server can't be shutdpwn twice
	if s.wasShutdown {
		return nil, jsonrpc.NewError(jsonrpc.InvalidRequest, "Shutdown was already sent", nil)
	}
	// Marked as shutdown
	s.wasShutdown = true
	// Notifies the implementation if it supports it
	i, ok := s.impl.(LifecycleShutdown)
	if ok {
		err := i.OnShutdown()
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// ExitHandler notifies the implementation that the server will exit
// This is the last thing that will happen as it calls os.Exit
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#exit
func (s *LSPServer) ExitHandler(params json.RawMessage) {
	i, ok := s.impl.(LifecycleExit)
	if ok {
		i.OnExit()
	}
	// Exit with error if shutdown wasn't received first
	if s.wasShutdown {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
