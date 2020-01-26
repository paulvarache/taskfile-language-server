package extension

import (
	"taskfile-language-server/jsonrpc"

	"github.com/sourcegraph/go-lsp"
)

func (t *TaskfileExtension) Initialize(params *lsp.InitializeParams) (*lsp.InitializeResult, *jsonrpc.ResponseError) {
	caps := lsp.ServerCapabilities{
		CompletionProvider: &lsp.CompletionOptions{ResolveProvider: true},
		TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
			Options: &lsp.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    lsp.TDSKFull,
			},
		},
	}
	return &lsp.InitializeResult{Capabilities: caps}, nil
}

func (t *TaskfileExtension) Initialized() *jsonrpc.ResponseError {
	return nil
}
