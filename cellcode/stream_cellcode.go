package main

/*

import (
	"os"
	"os/exec"
	"plugin"

	"fmt"
	"log"
	"net"

	"github.com/it-chain/tesseract/cellcode/cell"
	"github.com/it-chain/tesseract/pb"
	"github.com/it-chain/tesseract/stream"
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

	cmd := exec.Command("touch", "/icode/1")
	cmd.Run()
	// Socket Connection
	serverStream := stream.NewDefaultServerStream(port, func() {
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
	Listen(serverStream)

	cmd = exec.Command("touch", "/icode/end")
	cmd.Run()

}

func Listen(s *stream.DefaultServerStream) {
	lis, err := net.Listen("tcp", ":50003")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, s)
	reflection.Register(server)
	fmt.Println(s.Port + "in Listen")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
*/
