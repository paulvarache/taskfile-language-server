package taskfile

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
)

type Range []int

var Vars map[string]*Var

func init() {
	Vars = make(map[string]*Var)
	Vars["OS"] = &Var{Name: "OS"}
	Vars["ARCH"] = &Var{Name: "ARCH"}
	Vars["exeExt"] = &Var{Name: "exeExt"}
}

type Taskfile struct {
	Path  string           `json:"path"`
	Tasks map[string]*Task `json:"tasks"`
	Vars  map[string]*Var  `json:"vars"`
}

func IsInRange(line int, col int, r Range) bool {
	return (r[0] < line && r[2] > line) || ((r[0] == line || r[2] == line) && (r[1] <= col && r[3] >= col))
}

func (t *Taskfile) TaskAtPosition(line int, col int) *Task {
	for _, t := range t.Tasks {
		if IsInRange(line, col, t.Range) {
			return t
		}
	}
	return nil
}

func Parse(doc *ast.Document) (*Taskfile, error) {
	m, ok := doc.Body.(*ast.MappingNode)
	if !ok {
		return nil, fmt.Errorf("OOPS")
	}
	taskfile := &Taskfile{}
	for _, v := range m.Values {
		tasks, err := GetTasks(v)
		if err != nil {
			return nil, err
		}
		vars, err := GetVars(v)
		if err != nil {
			return nil, err
		}
		if tasks != nil {
			taskfile.Tasks = tasks
		}
		if vars != nil {
			taskfile.Vars = vars
		}
	}
	return taskfile, nil
}

// PreloadWithBytes will parse a yaml file and extract
// the Taskfile specific information like tasks, variables and expressions
// This will ignore parsing errors
func PreloadWithBytes(path string, contents []byte) {
	f, err := parser.ParseBytes(contents, parser.ParseComments)
	if err != nil {
		Taskfiles[path] = nil
		return
	}
	tf, err := Parse(f.Docs[0])
	tf.Path = path
	if err != nil {
		Taskfiles[path] = nil
		return
	}
	Taskfiles[path] = tf
	return
}

func Preload(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return PreloadWithBytes(path, bytes)
}

func IsTemplateKey(key string) bool {
	switch key {
	case "desc":
		return true
	default:
		return false
	}
}

type Expr struct {
	Range Range
	Value string
}

type ExprInString struct {
	Indices [2]int
	Value   string
}

type Result struct {
	LastToken   *token.Token
	Expressions []Expr
}

func GetAllExpr(src string) []ExprInString {
	r := regexp.MustCompile(`{{(.*?)}}`)
	f := r.FindAllSubmatchIndex([]byte(src), 12)
	items := make([]ExprInString, 0)
	ru := []rune(src)
	for _, i := range f {
		expr := ExprInString{Value: string(ru[i[2]:i[3]]), Indices: [2]int{i[2], i[3]}}
		items = append(items, expr)
	}
	return items
}

func Analyze(node ast.Node) *Result {
	switch n := node.(type) {
	case *ast.MappingValueNode:
		return Analyze(n.Value)
	case ast.ScalarNode:
		expressions := make([]Expr, 0)
		if sn, ok := n.(*ast.StringNode); ok {
			exps := GetAllExpr(sn.Value)
			for _, exp := range exps {
				rang := []int{
					sn.Token.Position.Line - 1,
					sn.Token.Position.Column + exp.Indices[0] - 1,
					sn.Token.Position.Line - 1,
					sn.Token.Position.Column + exp.Indices[1] - 1,
				}
				expressions = append(expressions, Expr{Value: exp.Value, Range: rang})
			}
		}
		return &Result{
			LastToken:   n.GetToken(),
			Expressions: expressions,
		}
	case *ast.SequenceNode:
		var t *token.Token
		expressions := make([]Expr, 0)
		for i, v := range n.Values {
			a := Analyze(v)
			expressions = append(expressions, a.Expressions...)
			if i == (len(n.Values) - 1) {
				t = a.LastToken
			}
		}
		return &Result{
			LastToken:   t,
			Expressions: expressions,
		}
	case *ast.MappingNode:
		var t *token.Token
		expressions := make([]Expr, 0)
		for i, v := range n.Values {
			a := Analyze(v)
			expressions = append(expressions, a.Expressions...)
			if i == (len(n.Values) - 1) {
				t = a.LastToken
			}
		}
		return &Result{
			LastToken:   t,
			Expressions: expressions,
		}
	default:
		return nil
	}
}
