package taskfile

type Memory map[string]*Taskfile

var Taskfiles Memory

func init() {
	Taskfiles = make(Memory)
}
