package extension

type TaskInfo struct {
	Task  TaskRef `json:"task"`
	Scope string  `json:"scope"`
}

type TaskRef struct {
	Value     string `json:"value"`
	StartLine int    `json:"startLine"`
	StartCol  int    `json:"startCol"`
	EndLine   int    `json:"endLine"`
	EndCol    int    `json:"endCol"`
}
