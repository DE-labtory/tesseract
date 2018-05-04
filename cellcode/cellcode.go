package main

import (
	"log"
	"net"
	"os"
	"os/exec"
	"plugin"

	"fmt"

	"github.com/it-chain/leveldb-wrapper"
	"github.com/it-chain/tesseract/cellcode/cell"
	"github.com/it-chain/tesseract/pb"
	"github.com/it-chain/tesseract/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ICode interface {
	Query(cell.Cell) pb.Response
	Invoke(cell.Cell) pb.Response
}

func main() {
	fmt.Println(os.Args)
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

	//init DB
	dbHandler := InitDB("wsdb")

	// Socket Connection
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := rpc.NewDefaultRpcServer(port, func(tx *cell.TxInfo) pb.Response {
		var res pb.Response

		// Setting Cell
		cell := cell.NewCell(tx, dbHandler)

		if cell.Tx.Method == "query" {
			fmt.Println("before query")
			res = iCode.Query(*cell)
		} else if cell.Tx.Method == "invoke" {
			fmt.Println("before invoke")
			res = iCode.Invoke(*cell)
		}

		return res
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

func InitDB(dbName string) *leveldbwrapper.DBHandle {
	path := "./wsdb"
	dbProvider := leveldbwrapper.CreateNewDBProvider(path)
	return dbProvider.GetDBHandle(dbName)
}
