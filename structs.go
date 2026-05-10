package main

import (
	"fmt"
	"time"
)

type ToolCall struct {
	CallID    string `json:"call_id"`
	Tool      string `json:"tool"`
	Arguments struct {
		Filename string `json:"filename,omitempty"`
		Content  string `json:"content,omitempty"`
		Command  string `json:"command,omitempty"`
		Pattern  string `json:"pattern,omitempty"`
		Path     string `json:"path,omitempty"`
	} `json:"arguments"`
}

func newToolCallID() string {
	return fmt.Sprintf("call_%d", time.Now().UnixNano())
}
