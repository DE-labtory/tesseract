package main

import (
	"os/exec"

	"strconv"

	"github.com/it-chain/tesseract/cellcode/cell"
	"github.com/it-chain/tesseract/pb"
)

type ICode struct {
	test string
}

func (ic *ICode) Query(cell cell.Cell) pb.Response {
	var r pb.Response

	cmd := exec.Command("touch", "/icode/query")
	cmd.Run()

	function, args := cell.GetFunctionAndParameters()
	if function == "getA" {
		r = ic.getA(cell, args)
	}

	return r
}

func (ic *ICode) Invoke(cell cell.Cell) pb.Response {
	var r pb.Response

	cmd := exec.Command("touch", "/icode/invoke")
	cmd.Run()

	function, args := cell.GetFunctionAndParameters()
	if function == "initA" {
		r = ic.initA(cell, args)
	} else if function == "incA" {
		r = ic.incA(cell, args)
	} else if function == "decA" {
		r = ic.decA(cell, args)
	}

	return r
}

func (ic *ICode) getA(cell cell.Cell, args []string) pb.Response {
	A, err := cell.GetData("A")
	if err != nil {
		return cell.Error(err.Error())
	}
	return cell.Success(map[string]string{"A": string(A[:])})
}

func (ic *ICode) initA(cell cell.Cell, args []string) pb.Response {
	err := cell.PutData("A", []byte("0"))
	if err != nil {
		return cell.Error(err.Error())
	}
	return cell.Success(nil)
}

func (ic *ICode) incA(cell cell.Cell, args []string) pb.Response {
	A, err := cell.GetData("A")
	if err != nil {
		return cell.Error(err.Error())
	}

	A_int, err := strconv.Atoi(string(A[:]))
	if err != nil {
		return cell.Error(err.Error())
	}

	err = cell.PutData("A", []byte(string(A_int+1)))
	if err != nil {
		return cell.Error(err.Error())
	}
	return cell.Success(nil)
}

func (ic *ICode) decA(cell cell.Cell, args []string) pb.Response {
	A, err := cell.GetData("A")
	if err != nil {
		return cell.Error(err.Error())
	}

	A_int, err := strconv.Atoi(string(A[:]))
	if err != nil {
		return cell.Error(err.Error())
	}

	err = cell.PutData("A", []byte(string(A_int-1)))
	if err != nil {
		return cell.Error(err.Error())
	}
	return cell.Success(nil)
}

var ICodeInstance ICode
