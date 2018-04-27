package main

import (
	"log"
	"net"
	"os"
	"os/exec"
	"plugin"

	"github.com/it-chain/tesseract/cellcode/cell"
	"github.com/it-chain/tesseract/pb"
	"github.com/it-chain/tesseract/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ICode interface {
	Query(cell.Cell)
	Invoke(cell.Cell)
}

func main() {
	if len(os.Args) != 3 {
		os.Exit(1)
	}

	iCodePath := os.Args[1]
	port := os.Args[2]

	plug, err := plugin.Open(iCodePath)
	if err != nil {
		os.Exit(1)
	}

	iCodePlugin, err := plug.Lookup("ICodeInstance")
	if err != nil {
		os.Exit(1)
	}

	iCode := iCodePlugin.(ICode)

	// Socket Connection
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := rpc.NewDefaultRpcServer(port, func() {
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
	server := grpc.NewServer()
	pb.RegisterDefaultServiceServer(server, s)
	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	cmd := exec.Command("touch", "/icode/end")
	cmd.Run()

}
