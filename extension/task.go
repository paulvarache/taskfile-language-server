package extension

import (
	"encoding/json"
	"path/filepath"
	"taskfile-language-server/jsonrpc"
	"taskfile-language-server/taskfile"
)

type GetTasksParams struct {
	FsPath string `json:"fsPath"`
}

func (t *TaskfileExtension) GetTasks(params json.RawMessage) (interface{}, *jsonrpc.ResponseError) {
	parsed := &GetTasksParams{}

	err := json.Unmarshal(params, parsed)
	if err != nil {
		return nil, jsonrpc.NewError(jsonrpc.ParseError, err.Error(), nil)
	}

	tf, ok := taskfile.Taskfiles[filepath.ToSlash(parsed.FsPath)]
	if !ok {
		return nil, jsonrpc.NewError(jsonrpc.ParseError, "Could not find taskfile", nil)
	}

	tasks := make([]*TaskInfo, 0)

	for _, t := range tf.Tasks {
		ti := &TaskInfo{
			Scope: parsed.FsPath,
			Task: TaskRef{
				Value:     t.Name,
				StartLine: t.Range[0],
				StartCol:  t.Range[1],
				EndLine:   t.Range[2],
				EndCol:    t.Range[3],
			},
		}
		tasks = append(tasks, ti)
	}

	return tasks, nil
}
