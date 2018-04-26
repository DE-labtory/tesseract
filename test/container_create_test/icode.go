package main

import (
	"os/exec"

	"github.com/it-chain/tesseract/cellcode/cell"
)

type ICode struct {
	test string
}

func (ic *ICode) Query(cell cell.Cell) {
	cmd := exec.Command("touch", "/cellcode/query")
	cmd.Run()
}

func (ic *ICode) Invoke(cell cell.Cell) {
	cmd := exec.Command("touch", "/cellcode/invoke")
	cmd.Run()
}

var ICodeInstance ICode
