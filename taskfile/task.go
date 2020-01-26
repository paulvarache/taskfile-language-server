package taskfile

import (
	"fmt"

	"github.com/goccy/go-yaml/ast"
)

type Task struct {
	Name        string          `json:"name"`
	Range       Range           `json:"range"`
	Vars        map[string]*Var `json:"vars"`
	Expressions []Expr          `json:"expressions"`
}

func (t *Task) ExpressionAtPosition(line int, col int) *Expr {
	for _, e := range t.Expressions {
		if IsInRange(line, col, e.Range) {
			return &e
		}
	}
	return nil
}

func GetTasks(node *ast.MappingValueNode) (map[string]*Task, error) {
	sn, ok := node.Key.(*ast.StringNode)
	if !ok {
		return nil, fmt.Errorf("OOPS")
	}
	if sn.Value == "tasks" {
		switch tasksNode := node.Value.(type) {
		case *ast.MappingValueNode:
			key, val := ExtractTaskFromMappingValueNode(tasksNode)
			tasks := make(map[string]*Task)
			tasks[key] = val
			return tasks, nil
		case *ast.MappingNode:
			return ExtractTasksFromMappingNode(tasksNode)
		}
	}
	return nil, nil
}

func ExtractTasksFromMappingNode(node *ast.MappingNode) (map[string]*Task, error) {
	tasks := make(map[string]*Task)
	for _, v := range node.Values {
		key, val := ExtractTaskFromMappingValueNode(v)
		tasks[key] = val
	}
	return tasks, nil
}

func ExtractTaskFromMappingValueNode(node *ast.MappingValueNode) (string, *Task) {
	name := node.Key.(*ast.StringNode).Value
	res := Analyze(node)
	var r Range = []int{
		node.Key.GetToken().Position.Line - 1,
		node.Key.GetToken().Position.Column - 1,
		res.LastToken.Position.Line - 1,
		res.LastToken.Position.Column + len(res.LastToken.Value) - 1,
	}
	task := &Task{Name: name, Range: r, Expressions: res.Expressions}
	varsNode, ok := node.Value.(*ast.MappingNode)
	if ok {
		vars, _ := ExtractTaskVarsFromMappingNode(varsNode)
		task.Vars = vars
	}
	return name, task
}

func ExtractTaskVarsFromMappingNode(node *ast.MappingNode) (map[string]*Var, error) {
	for _, v := range node.Values {
		vars, err := GetVars(v)
		if err != nil {
			return nil, err
		}
		if vars != nil {
			return vars, nil
		}
	}
	return nil, nil
}
