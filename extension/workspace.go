package extension

import (
	"taskfile-language-server/taskfile"

	"github.com/sourcegraph/go-lsp"
)

type TaskfileInfo struct {
	Scope string      `json:"scope"`
	Tasks []*TaskInfo `json:"tasks"`
}

func (s *TaskfileExtension) WorkspaceDidChangeWatchedFiles(params *lsp.DidChangeWatchedFilesParams) {
	for _, v := range params.Changes {
		p, err := GetPath(v.URI)
		if err != nil {
			s.Logger.Fatalln(err)
		}
		tf, err := taskfile.Preload(p)
		if err != nil {
			s.Logger.Fatalln(err)
		}
		if tf.Tasks != nil {
			tasks := make([]*TaskInfo, 0)
			for _, t := range tf.Tasks {
				tasks = append(tasks, GetTaskInfo(tf.Path, t))
			}
			tfi := &TaskfileInfo{
				Scope: tf.Path,
				Tasks: tasks,
			}
			s.SendNotification("extension/onTaskfileUpdate", tfi)
		}
	}
}
