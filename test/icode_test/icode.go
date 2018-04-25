package main

import (
	"os/exec"

	"github.com/it-chain/tesseract/cellcode/cell"
)

type ICode struct {
	test string
}

func (ic *ICode) Query(cell cell.Cell) {
	cmd := exec.Command("touch", "/icode/query")
	cmd.Run()
}

func (ic *ICode) Invoke(cell cell.Cell) {
	cmd := exec.Command("touch", "/icode/invoke")
	cmd.Run()
}

var ICodeInstance ICode
