package extension

import (
	"fmt"
	"taskfile-language-server/jsonrpc"
	"taskfile-language-server/taskfile"

	"github.com/sourcegraph/go-lsp"
)

func (t *TaskfileExtension) TextDocumentDidOpen(params *lsp.DidOpenTextDocumentParams) {
	err := reloadTaskfile(params.TextDocument.URI, params.TextDocument.Text)
	if err != nil {
		t.Logger.Panicf(err.Error())
	}
}

func (t *TaskfileExtension) TextDocumentDidChange(params *lsp.DidChangeTextDocumentParams) {
	err := reloadTaskfile(params.TextDocument.URI, params.ContentChanges[0].Text)
	if err != nil {
		t.Logger.Panicf(err.Error())
	}
}

func (t *TaskfileExtension) TextDocumentDidClose(params *lsp.DidCloseTextDocumentParams) {}

func CompletionItemFromVar(v *taskfile.Var, scoped bool) lsp.CompletionItem {
	data := v.Name
	if scoped {
		data = fmt.Sprintf(".%s", v.Name)
	}
	return lsp.CompletionItem{Label: v.Name, Kind: lsp.CIKVariable, InsertText: data}
}

func CompletionItemsFromVars(vars map[string]*taskfile.Var, scoped bool) []lsp.CompletionItem {
	items := make([]lsp.CompletionItem, 0)
	for _, v := range vars {
		items = append(items, CompletionItemFromVar(v, scoped))
	}
	return items
}

func (t *TaskfileExtension) TextDocumentCompletion(params *lsp.CompletionParams) (*lsp.CompletionList, *jsonrpc.ResponseError) {
	p, err := GetPath(params.TextDocument.URI)
	if err != nil {
		return nil, jsonrpc.NewError(jsonrpc.InternalError, err.Error(), nil)
	}
	empty := &lsp.CompletionList{Items: []lsp.CompletionItem{}, IsIncomplete: false}
	tf := taskfile.GetParsedTaskfile(p)
	if tf == nil {
		t.Logger.Printf("Could not find node tree for %s", p)
		// No taskfile means the parsing went wrong. Maybe the user is stil typing
		return empty, nil
	}
	task := tf.TaskAtPosition(params.Position.Line, params.Position.Character)
	if task == nil {
		t.Logger.Println("Cursor is not in task")
		return empty, nil
	}
	exp := task.ExpressionAtPosition(params.Position.Line, params.Position.Character)
	if exp == nil {
		t.Logger.Println("Cursor is not in expression")
		return empty, nil
	}
	items := make([]lsp.CompletionItem, 0)
	// Add local variables
	items = append(items, CompletionItemsFromVars(task.Vars, true)...)
	// Add taskfile variables
	items = append(items, CompletionItemsFromVars(tf.Vars, true)...)
	// Add global variables
	items = append(items, CompletionItemsFromVars(taskfile.Vars, false)...)

	return &lsp.CompletionList{Items: items, IsIncomplete: false}, nil
}
