package main

import (
	"os"
	"os/exec"
	"plugin"

	"github.com/it-chain/tesseract/cellcode/cell"
	"github.com/it-chain/tesseract/stream"
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

	cmd := exec.Command("touch", "/icode/1")
	cmd.Run()
	// Socket Connection
	serverStream := stream.NewDefaultServerStream(":50001")
	serverStream.Listen(func() {
		cmd := exec.Command("touch", "/icode/2")
		cmd.Run()

		// Setting Cell

		cell := cell.NewCell()

		if true {
			iCode.Query(*cell)
		}
		cmd = exec.Command("touch", "/icode/inHandler")
		cmd.Run()
	})

	cmd = exec.Command("touch", "/icode/end")
	cmd.Run()

}
