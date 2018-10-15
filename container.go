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

import "github.com/it-chain/tesseract/pb"

type ContainerID = string
type CallBack func(response *pb.Response, err error)

type ContainerFactory interface {
	Create(ContainerConfig) Container
}

type ContainerConfig struct {
	Name           string
	Url            string
	Directory      string
	ContainerImage ContainerImage
	UseDocker      bool
	language       string // ENUM 으로 대체하면 좋음
}

const (
	DefaultImageName = "golang"
	DefaultImageTag  = "1.9"
)

type ContainerImage struct {
	Name string
	Tag  string
}

func GetDefaultImage() ContainerImage {
	return ContainerImage{
		Name: DefaultImageName,
		Tag:  DefaultImageTag,
	}
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

type Port struct {
	// Host IP address that the container's port is mapped to
	IP string

	// Port on the container
	PrivatePort int

	// Port exposed on the host
	PublicPort int
}
