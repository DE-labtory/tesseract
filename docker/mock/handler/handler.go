/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package handler

import (
	"strconv"

	"errors"

	"fmt"

	"github.com/it-chain/sdk"
	"github.com/it-chain/sdk/logger"
	"github.com/it-chain/sdk/pb"
)

type HandlerExample struct {
}

func (*HandlerExample) Name() string {
	return "sample"
}

func (*HandlerExample) Versions() []string {
	vers := make([]string, 0)
	vers = append(vers, "1.0")
	vers = append(vers, "1.2")
	return vers
}

func (*HandlerExample) Handle(request *pb.Request, cell *sdk.Cell) (*pb.Response, error) {
	switch request.Type {
	case "invoke":
		return handleInvoke(request, cell)
	case "query":
		return handleQuery(request, cell)
	case "test":
		fmt.Println("req : " + request.Uuid)
		if request.Uuid == "0" {
			cell.PutData("test", []byte("0"))
			return responseSuccess(request, []byte(string(0))), nil
		}
		data, err := cell.GetData("test")
		if err != nil {
			return responseError(request, err), err
		}
		if len(data) == 0 {
			err := errors.New("no data err")
			return responseError(request, err), err
		}
		strData := string(data)
		intData, err := strconv.Atoi(strData)
		if err != nil {
			return responseError(request, err), err
		}
		intData = intData + 1
		changeData := strconv.Itoa(intData)
		err = cell.PutData("test", []byte(changeData))
		if err != nil {
			return responseError(request, err), err
		}
		return responseSuccess(request, []byte(changeData)), nil
	default:
		logger.Fatal(nil, "unknown request type")
		err := errors.New("unknown request type")
		return responseError(request, err), err
	}
}
func handleQuery(request *pb.Request, cell *sdk.Cell) (*pb.Response, error) {
	switch request.FunctionName {
	case "getA":
		b, err := cell.GetData("A")
		if err != nil {
			return responseError(request, err), err
		}
		return responseSuccess(request, b), nil

	default:
		err := errors.New("unknown query method")
		return responseError(request, err), err
	}
}
func handleInvoke(request *pb.Request, cell *sdk.Cell) (*pb.Response, error) {
	switch request.FunctionName {
	case "initA":
		err := cell.PutData("A", []byte("0"))
		if err != nil {
			return responseError(request, err), err
		}
		return responseSuccess(request, nil), nil
	case "incA":
		data, err := cell.GetData("A")
		if err != nil {
			return responseError(request, err), err
		}
		if len(data) == 0 {
			err := errors.New("no data err")
			return responseError(request, err), err
		}
		strData := string(data)
		intData, err := strconv.Atoi(strData)
		if err != nil {
			return responseError(request, err), err
		}
		intData++
		changeData := strconv.Itoa(intData)
		err = cell.PutData("A", []byte(changeData))
		if err != nil {
			return responseError(request, err), err
		}
		return responseSuccess(request, nil), nil
	default:
		err := errors.New("unknown invoke method")
		return responseError(request, err), err
	}
}

func responseError(request *pb.Request, err error) *pb.Response {
	return &pb.Response{
		Uuid:   request.Uuid,
		Type:   request.Type,
		Result: false,
		Data:   nil,
		Error:  err.Error(),
	}
}

func responseSuccess(request *pb.Request, data []byte) *pb.Response {
	return &pb.Response{
		Uuid:   request.Uuid,
		Type:   request.Type,
		Result: true,
		Data:   data,
		Error:  "",
	}
}
