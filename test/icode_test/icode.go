package main

import (
	"os/exec"
)

type ICode struct {
}

func (ic *ICode) Query() {
	cmd := exec.Command("touch", "/tesseract/cellcode/query")
	cmd.Run()
}

func (ic *ICode) Invoke() {
}

var ICodeInstance ICode
