package extension

import "github.com/sourcegraph/go-lsp"

import "yaml/jsonrpc"

func (t *TaskfileExtension) CompletionItemResolve(item *lsp.CompletionItem) (*lsp.CompletionItem, *jsonrpc.ResponseError) {
	return item, nil
}
