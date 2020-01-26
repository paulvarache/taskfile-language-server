package lsp

import (
	"encoding/json"
	"taskfile-language-server/jsonrpc"

	"github.com/sourcegraph/go-lsp"
)

func (s *LSPServer) InitializeHandler(params json.RawMessage) (interface{}, *jsonrpc.ResponseError) {
	parsed := &lsp.InitializeParams{}
	err := json.Unmarshal(params, parsed)
	if err != nil {
		return nil, jsonrpc.NewError(jsonrpc.ParseError, err.Error(), nil)
	}
	i, ok := s.impl.(ServerImplementation)
	if !ok {
		return nil, MethodNotFoundError("Initialize")
	}
	return i.Initialize(parsed)
}
