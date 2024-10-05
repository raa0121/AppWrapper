package main

type Setting struct {
	Cmd map[string]Cmd
}

type Cmd struct {
	Command string
	Env map[string]string
}
