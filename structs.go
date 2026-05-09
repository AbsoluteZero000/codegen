package main

type ToolCall struct {
	Tool      string `json:"tool"`
	Arguments struct {
		Filename string `json:"filename"`
	} `json:"arguments"`
}
