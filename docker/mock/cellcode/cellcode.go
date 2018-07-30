package main

import (
	"os"
	"os/exec"
	"plugin"

	"github.com/it-chain/tesseract/cellcode/cell"
)

type ICode interface {
	Query(cell.Cell)
	Invoke(cell.Cell)
}

func main() {

	if len(os.Args) != 2 {
		os.Exit(1)
	}

	iCodePath := os.Args[1]

	plug, err := plugin.Open(iCodePath)
	if err != nil {
		os.Exit(1)
	}

	iCodePlugin, err := plug.Lookup("ICodeInstance")
	if err != nil {
		os.Exit(1)
	}

	iCode := iCodePlugin.(ICode)

	print(iCode)

	cmd := exec.Command("touch", "/sh/main")
	cmd.Run()
}
