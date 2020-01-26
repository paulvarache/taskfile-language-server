package taskfile

import (
	"fmt"

	"github.com/goccy/go-yaml/ast"
)

type Var struct {
	Name  string `json:"name"`
	Range Range  `json:"range"`
}

func GetVars(node *ast.MappingValueNode) (map[string]*Var, error) {
	sn, ok := node.Key.(*ast.StringNode)
	if !ok {
		return nil, fmt.Errorf("OOPS")
	}
	if sn.Value == "vars" {
		switch varsNode := node.Value.(type) {
		case *ast.MappingValueNode:
			key, val := ExtractVarFromMappingValueNode(varsNode)
			vars := make(map[string]*Var)
			vars[key] = val
			return vars, nil
		case *ast.MappingNode:
			return ExtractVarsFromMappingNode(varsNode)
		}
	}
	return nil, nil
}

func ExtractVarsFromMappingNode(node *ast.MappingNode) (map[string]*Var, error) {
	vars := make(map[string]*Var)
	for _, v := range node.Values {
		key, val := ExtractVarFromMappingValueNode(v)
		vars[key] = val
	}
	return vars, nil
}

func ExtractVarFromMappingValueNode(node *ast.MappingValueNode) (string, *Var) {
	name := node.Key.(*ast.StringNode).Value
	res := Analyze(node)
	var r Range = []int{
		node.Key.(*ast.StringNode).Token.Position.Line - 1,
		node.Key.(*ast.StringNode).Token.Position.Column - 1,
		res.LastToken.Position.Line - 1,
		res.LastToken.Position.Column + len(name) - 1,
	}
	return node.Key.(*ast.StringNode).Value, &Var{Name: name, Range: r}
}
