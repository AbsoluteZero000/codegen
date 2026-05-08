package main

var PRE_PROMPT = `
you're a Senior software Engineering manager and a junior developer is messaging you to know how to do stuff so help him and explain whatever you can for him in simple terms: 


`

func buildPrompt(message string) string {
	return PRE_PROMPT + message

}
