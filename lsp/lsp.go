package lsp

import (
	"fmt"
	"log"
	"taskfile-language-server/jsonrpc"

	"github.com/sourcegraph/go-lsp"
)

type ServerImplementation interface {
	Initialize(*lsp.InitializeParams) (*lsp.InitializeResult, *jsonrpc.ResponseError)
	Initialized() *jsonrpc.ResponseError
}

type TextDocumentSync interface {
	TextDocumentDidOpen(*lsp.DidOpenTextDocumentParams)
	TextDocumentDidChange(*lsp.DidChangeTextDocumentParams)
	TextDocumentDidClose(*lsp.DidCloseTextDocumentParams)
}

type WorkspaceSync interface {
	WorkspaceDidChangeWatchedFiles(*lsp.DidChangeWatchedFilesParams)
}

type TextDocumentCompletion interface {
	TextDocumentCompletion(*lsp.CompletionParams) (*lsp.CompletionList, *jsonrpc.ResponseError)
}

type CompletionItemResolve interface {
	CompletionItemResolve(*lsp.CompletionItem) (*lsp.CompletionItem, *jsonrpc.ResponseError)
}

type TextDocumentHover interface {
	TextDocumentHover(*lsp.TextDocumentPositionParams) (*lsp.Hover, *jsonrpc.ResponseError)
}

type LSPServer struct {
	server *jsonrpc.Server
	impl   interface{}
	logger *log.Logger
}

func NewServer(s *jsonrpc.Server, impl interface{}, logger *log.Logger) *LSPServer {
	server := &LSPServer{server: s, impl: impl}
	server.logger = logger

	s.AddHandler("initialize", server.InitializeHandler)
	s.AddHandler("initialized", server.InitializedHandler)

	s.AddNotificationHandler("textDocument/didOpen", server.TextDocumentOpen)
	s.AddNotificationHandler("textDocument/didChange", server.TextDocumentChange)
	s.AddNotificationHandler("textDocument/didClose", server.TextDocumentClose)
	s.AddHandler("textDocument/completion", server.TextDocumentCompletion)
	s.AddHandler("completionItem/resolve", server.CompletionItemResolve)
	s.AddNotificationHandler("workspace/didChangeWatchedFiles", server.DidChangeWatchedFiles)
	s.AddHandler("textDocument/hover", server.TextDocumentHover)

	return server
}

func MethodNotFoundError(name string) *jsonrpc.ResponseError {
	return jsonrpc.NewError(jsonrpc.MethodNotFound, fmt.Sprintf("Server is not implementing the %s method", name), nil)
}
