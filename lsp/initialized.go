package lsp

import "encoding/json"

import "taskfile-language-server/jsonrpc"

func (s *LSPServer) InitializedHandler(params json.RawMessage) (interface{}, *jsonrpc.ResponseError) {
	i, ok := s.impl.(ServerImplementation)
	if !ok {
		return nil, MethodNotFoundError("Initialized")
	}
	return nil, i.Initialized()
}
