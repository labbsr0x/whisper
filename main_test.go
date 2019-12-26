package main

import (
	"testing"
)

var executed = false

func MockExecute() {
	executed = true
}

func TestMain(t *testing.T) {
	_exec = MockExecute
	main()

	if !executed {
		t.Error("cmd.Exec did not run")
	}
}
