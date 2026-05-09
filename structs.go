package main

type ToolCall struct {
	Tool      string `json:"tool"`
	Arguments struct {
		Filename string `json:"filename,omitempty"`
		Content  string `json:"content,omitempty"`
		Command  string `json:"command,omitempty"`
	} `json:"arguments"`
}
