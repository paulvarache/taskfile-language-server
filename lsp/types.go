package lsp

import "github.com/sourcegraph/go-lsp"

/*
Sourcegraph's protocol is missing a few types from the latest LSP spec
This was taken from go's internals and adapted to co-exist with the sourcegraph implementation
*/

type MarkupKind string

const (
	PlainText MarkupKind = "plaintext"
	Markdown  MarkupKind = "markdown"
)

type MarkupContent struct {
	/**
	 * The type of the Markup
	 */
	Kind MarkupKind `json:"kind"`
	/**
	 * The content itself
	 */
	Value string `json:"value"`
}

type Hover struct {
	/**
	 * The hover's content
	 */
	Contents MarkupContent `json:"contents"`
	/**
	 * An optional range
	 */
	Range lsp.Range `json:"range,omitempty"`
}
