package main

type ToolCall struct {
	Tool      string `json:"tool"`
	Arguments struct {
		Filename string `json:"filename"`
		Content  string `json:"content,omitempty"`
	} `json:"arguments"`
}
