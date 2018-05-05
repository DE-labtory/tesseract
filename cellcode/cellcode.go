package main

import (
	"log"
	"net"
	"os"
	"plugin"

	"fmt"

	"github.com/it-chain/leveldb-wrapper"
	"github.com/it-chain/tesseract/cellcode/cell"
	"github.com/it-chain/tesseract/pb"
	"github.com/it-chain/tesseract/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var iCode ICode
var dbHandler *leveldbwrapper.DBHandle

type ICode interface {
	Query(cell.Cell) pb.Response
	Invoke(cell.Cell) pb.Response
}

func main() {

	log.Println("Cellcode is Starting")

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

	iCode = iCodePlugin.(ICode)

	log.Println("Icode is initiated")

	//init DB
	//todo 외부로 부터 wsdb이름 받아오기
	dbHandler = InitDB("wsdb")

	// Socket Connection
	lis, err := net.Listen("tcp", ":"+port)
	defer lis.Close()

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := rpc.NewDefaultRpcServer(port, Handle)

	server := grpc.NewServer()
	defer server.Stop()

	pb.RegisterDefaultServiceServer(server, s)
	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		server.Stop()
		lis.Close()
		log.Fatalf("failed to serve: %v", err)
	}
}

func InitDB(dbName string) *leveldbwrapper.DBHandle {

	path := "./wsdb"
	dbProvider := leveldbwrapper.CreateNewDBProvider(path)
	return dbProvider.GetDBHandle(dbName)
}

func Handle(tx *cell.TxInfo) pb.Response {

	log.Printf("Run Icode [%s]", tx)

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

	log.Printf("End Icode [%s]", res)
	return res
}
