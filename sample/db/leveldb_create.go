package main

import (
	"fmt"

	"os"

	"github.com/it-chain/leveldb-wrapper"
)

func main() {
	GOPATH := os.Getenv("GOPATH")
	path := GOPATH + "/src/github.com/it-chain/tesseract/sample/db/leveldb"
	dbProvider := leveldbwrapper.CreateNewDBProvider(path)
	//defer os.RemoveAll(path)

	studentDB := dbProvider.GetDBHandle("Student")
	studentDB.Put([]byte("20164403"), []byte("JUN"), true)

	name, _ := studentDB.Get([]byte("20164403"))

	fmt.Printf("%s", name)
}
