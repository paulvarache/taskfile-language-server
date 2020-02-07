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

func GetTaskInfo(scope string, task *taskfile.Task) *TaskInfo {
	return &TaskInfo{
		Scope: scope,
		Task: TaskRef{
			Value:     task.Name,
			StartLine: task.Range[0],
			StartCol:  task.Range[1],
			EndLine:   task.Range[2],
			EndCol:    task.Range[3],
		},
	}
}

func (t *TaskfileExtension) GetTasks(params json.RawMessage) (interface{}, *jsonrpc.ResponseError) {
	parsed := &GetTasksParams{}

	err := json.Unmarshal(params, parsed)
	if err != nil {
		return nil, jsonrpc.NewError(jsonrpc.ParseError, err.Error(), nil)
	}

	path := filepath.ToSlash(parsed.FsPath)

	tf := taskfile.GetParsedTaskfile(path)
	if tf == nil {
		return nil, jsonrpc.NewError(jsonrpc.ParseError, "Could not find taskfile", nil)
	}

	tasks := make([]*TaskInfo, 0)

	if tf.Tasks == nil {
		return tasks, nil
	}

	for _, t := range tf.Tasks {
		ti := GetTaskInfo(path, t)
		tasks = append(tasks, ti)
	}

	return tasks, nil
}
