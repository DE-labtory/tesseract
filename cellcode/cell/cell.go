package cell

import (
	"encoding/json"

	"github.com/it-chain/leveldb-wrapper"
	"github.com/it-chain/tesseract/pb"
)

type Params struct {
	Type     int
	Function string
	Args     []string
}

type TxInfo struct {
	Method string
	ID     string
	Params Params
}

type Cell struct {
	DBHandler *leveldbwrapper.DBHandle
	Tx        *TxInfo
}

func NewCell(tx *TxInfo, dbHandler *leveldbwrapper.DBHandle) *Cell {
	return &Cell{Tx: tx, DBHandler: dbHandler}
}

func (c Cell) GetFunctionAndParameters() (string, []string) {
	return c.Tx.Params.Function, c.Tx.Params.Args
}

func (c Cell) PutData(key string, value []byte) error {
	return c.DBHandler.Put([]byte(key), value, true)
}

func (c Cell) GetData(key string) ([]byte, error) {
	value, err := c.DBHandler.Get([]byte(key))
	return value, err
}

func (c Cell) Error(err string) pb.Response {
	return pb.Response{Result: "Error", Method: "", Data: nil, Error: err}
}

func (c Cell) Success(data map[string]string) pb.Response {
	dataBytes, _ := json.Marshal(data)
	return pb.Response{Result: "Success", Method: "", Data: dataBytes, Error: ""}
}
