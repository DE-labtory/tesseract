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
 */

package tesseract

import (
	"github.com/it-chain/tesseract/pb"
)

type ContainerID = string
type CallBack func(response *pb.Response, err error)

type ContainerFactory interface {
	Create(ContainerConfig) Container
}

type ContainerConfig struct {
	language string // language that icode use
	// todo ENUM 으로 대체하면 좋음 ( expect to change enum )

	Name string // container name

	ContainerImage ContainerImage // container docker image name to pull. example : 'golang:1.9'

	IP string // IP address that manage icode.

	Port string // port num to communicate container and icode

	StartCmd []string // start command to initiate icode. it must contain icodename. ex {"go","run.sh","testicode/icode.go","-p","4401}

	Network *Network // docker network option. if you don't use, nil

	Mount string // [Mount-src-path or volume name]:[Mount-dist-path] ex) /go/src/github.com/it-chain/learn-icode:/icode

	//Volume *Volume // docker volume option.
	//
	//ImgSrcRootPath string // input icode base root path.
	//// if language need specify path like go ( golang source file have to be
	//// in %GOPATH% to build ), please input GOPATH like "/go/src".
	//// if language do not use specify path to build, empty or "/"
	//HostICodeRoot string // if not use volume option, 호스트 아이코드 루트 디렉토리를 써주세용
}

type ContainerImage struct {
	Name string // language docker image name ( ex : go )
	Tag  string // language docker image Tag ( ex : 1.9 )
}

func (dc ContainerImage) GetFullName() string {
	return dc.Name + ":" + dc.Tag
}

type Container interface {
	// send request to container
	Request(req Request, callback CallBack) error

	// close container
	Close() error

	// get container ID
	GetID() string
}

type Request struct {
	Uuid     string
	TypeName string
	FuncName string
	Args     []string
}
