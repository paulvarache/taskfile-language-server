package lsp

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sourcegraph/go-lsp"
)

func (s *LSPServer) DidChangeWatchedFiles(params json.RawMessage) {
	parsed := &lsp.DidChangeWatchedFilesParams{}
	err := json.Unmarshal(params, parsed)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}
	i, ok := s.impl.(WorkspaceSync)
	if !ok {
		return
	}
	i.WorkspaceDidChangeWatchedFiles(parsed)
}
