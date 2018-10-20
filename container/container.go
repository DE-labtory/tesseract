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

package container

import (
	"github.com/it-chain/tesseract"
	"github.com/it-chain/tesseract/docker"
	"github.com/it-chain/tesseract/pb"
	"github.com/it-chain/tesseract/rpc"
)

type DockerContainer struct {
	ID     tesseract.ContainerID
	Client *rpc.ClientStream
	config tesseract.ContainerConfig
}

func NewDockerContainer(id tesseract.ContainerID, client *rpc.ClientStream, config tesseract.ContainerConfig) DockerContainer {
	return DockerContainer{
		ID:     id,
		Client: client,
		config: config,
	}
}

// send request
func (d DockerContainer) Request(req tesseract.Request, callback tesseract.CallBack) error {

	// Args : Transaction
	// Get Container handler using SmartContract ID
	// Send Query or Invoke massage
	// Receive result
	// Return result

	err := d.Client.RunICode(&pb.Request{
		Uuid:         req.Uuid,
		Type:         req.TypeName,
		FunctionName: req.FuncName,
		Args:         req.Args,
	}, callback)

	if err != nil {
		return err
	}

	return nil
}

// close docker container
func (d DockerContainer) Close() error {

	d.Client.Close()

	if err := docker.KillContainer(d.ID); err != nil {
		return err
	}
	if err := docker.RemoveContainer(d.ID); err != nil {
		return err
	}

	return nil
}

func (d DockerContainer) GetID() string {
	return d.ID
}
