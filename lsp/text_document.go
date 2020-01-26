package lsp

import (
	"encoding/json"
	"fmt"
	"os"
	"taskfile-language-server/jsonrpc"

	"github.com/sourcegraph/go-lsp"
)

func (s *LSPServer) TextDocumentOpen(params json.RawMessage) {
	parsed := &lsp.DidOpenTextDocumentParams{}
	err := json.Unmarshal(params, parsed)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}
	i, ok := s.impl.(TextDocumentSync)
	if !ok {
		return
	}
	i.TextDocumentDidOpen(parsed)
}

func (s *LSPServer) TextDocumentChange(params json.RawMessage) {
	parsed := &lsp.DidChangeTextDocumentParams{}
	err := json.Unmarshal(params, parsed)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}
	i, ok := s.impl.(TextDocumentSync)
	if !ok {
		return
	}
	i.TextDocumentDidChange(parsed)
}

func (s *LSPServer) TextDocumentClose(params json.RawMessage) {
	parsed := &lsp.DidCloseTextDocumentParams{}
	err := json.Unmarshal(params, parsed)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}
	i, ok := s.impl.(TextDocumentSync)
	if !ok {
		return
	}
	i.TextDocumentDidClose(parsed)
}

func (s *LSPServer) TextDocumentCompletion(params json.RawMessage) (interface{}, *jsonrpc.ResponseError) {
	parsed := &lsp.CompletionParams{}
	err := json.Unmarshal(params, parsed)
	if err != nil {
		return nil, jsonrpc.NewError(jsonrpc.ParseError, err.Error(), nil)
	}
	i, ok := s.impl.(TextDocumentCompletion)
	if !ok {
		return nil, MethodNotFoundError("TextDocumentCompletion")
	}
	return i.TextDocumentCompletion(parsed)
}

func (s *LSPServer) TextDocumentHover(params json.RawMessage) (interface{}, *jsonrpc.ResponseError) {
	parsed := &lsp.TextDocumentPositionParams{}
	err := json.Unmarshal(params, parsed)
	if err != nil {
		return nil, jsonrpc.NewError(jsonrpc.ParseError, err.Error(), nil)
	}
	i, ok := s.impl.(TextDocumentHover)
	if !ok {
		return nil, MethodNotFoundError("TextDocumentHover")
	}
	return i.TextDocumentHover(parsed)
}
