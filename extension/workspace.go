package extension

import (
	"taskfile-language-server/taskfile"

	"github.com/sourcegraph/go-lsp"
)

func (s *TaskfileExtension) WorkspaceDidChangeWatchedFiles(params *lsp.DidChangeWatchedFilesParams) {
	for _, v := range params.Changes {
		p, err := GetPath(v.URI)
		if err != nil {
			s.Logger.Fatalln(err)
			continue
		}
		err = taskfile.Preload(p)
		if err != nil {
			s.Logger.Fatalln(err)
		}
	}
}
