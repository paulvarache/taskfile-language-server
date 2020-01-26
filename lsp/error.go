package lsp

import "yaml/jsonrpc"

const (
	RequestCancelled jsonrpc.ErrorCode = -32800
	ContentModified  jsonrpc.ErrorCode = -32801
)
