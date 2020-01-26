package lsp

import (
	"encoding/json"
	"yaml/jsonrpc"

	"github.com/sourcegraph/go-lsp"
)

func (s *LSPServer) CompletionItemResolve(params json.RawMessage) (interface{}, *jsonrpc.ResponseError) {
	item := &lsp.CompletionItem{}
	err := json.Unmarshal(params, item)
	if err != nil {
		return nil, jsonrpc.NewError(jsonrpc.ParseError, err.Error(), nil)
	}
	i, ok := s.impl.(CompletionItemResolve)
	if !ok {
		return nil, MethodNotFoundError("CompletionItemResolve")
	}
	return i.CompletionItemResolve(item)
}
