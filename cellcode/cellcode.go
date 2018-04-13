package main

import (
	"os"
	"plugin"
)

type ICode interface {
	Query()
	Invoke()
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

	// ToDo : Socket Connection
	// ToDo : Setting Cell

	if(true) {
		iCode.Query()
	} else {
		iCode.Invoke()
	}
}