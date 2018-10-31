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
	"context"

	"docker.io/go-docker/api/types"

	"errors"

	"net"
	"time"

	"github.com/it-chain/iLogger"
	"github.com/it-chain/tesseract"
	"github.com/it-chain/tesseract/docker"
	"github.com/it-chain/tesseract/rpc"
)

var ErrFailedPullImage = errors.New("failed to pull image")

func Create(config tesseract.ContainerConfig) (DockerContainer, error) {

	iLogger.Info(nil, "[Tesseract] creating container")
	containerImage := config.ContainerImage

	// checking port
	lis, err := net.Listen("tcp", config.HostIp+":"+config.Port)
	lis.Close()
	if err != nil {
		return DockerContainer{}, err
	}

	if err = pullImage(containerImage.GetFullName()); err != nil {
		return DockerContainer{}, ErrFailedPullImage
	}

	// Create Docker
	res, err := docker.CreateContainer(
		config,
	)

	if err != nil {
		return DockerContainer{}, err
	}

	containerInfo, err := docker.StartContainer(res)
	if err != nil {
		iLogger.Errorf(nil, "[Tesseract] fail to create container: %s", err.Error())
		docker.RemoveContainer(res.ID)
		return DockerContainer{}, err
	}

	var client *rpc.ClientStream
	if config.Network == nil {
		client, err = createClient(config.ContainerIp, config.Port)
	} else {
		ipAddress := retrieveNetworkIpAddress(config.Network.Name, containerInfo)
		client, err = createClient(ipAddress, config.Port)
	}

	if err != nil {
		iLogger.Errorf(nil, "[Tesseract] closing container %s", res.ID)
		docker.KillContainer(res.ID)
		docker.RemoveContainer(res.ID)
		return DockerContainer{}, err
	}

	return NewDockerContainer(res.ID, client, config), nil
}

func retrieveNetworkIpAddress(network string, containerInfo types.ContainerJSON) string {
	return containerInfo.NetworkSettings.Networks[network].IPAddress
}

func pullImage(ImageFullName string) error {

	// Docker IMAGE pull
	r, err := docker.HasImage(ImageFullName)
	if err != nil {
		return err
	}

	if !r {
		docker.PullImage(ImageFullName)
	}

	return nil
}

func createClient(ipAddress string, port string) (*rpc.ClientStream, error) {
	return retryConnectWithTimeOut(ipAddress, port, 120*time.Second)
}

func retryConnectWithTimeOut(ipAddress string, port string, timeout time.Duration) (*rpc.ClientStream, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c := make(chan *rpc.ClientStream, 1)

	go func() {

		ticker := time.NewTicker(2 * time.Second)
		for _ = range ticker.C {
			client, err := rpc.NewClientStream(ipAddress + ":" + port)
			if err != nil {
				continue
			}

			_, err = client.Ping()

			if err == nil {
				client.SetHandler(rpc.NewDefaultHandler())
				client.StartHandle()
				c <- client
				return
			}
		}
	}()

	select {

	case <-ctx.Done():
		//timeoutted body
		return nil, ctx.Err()

	case client := <-c:
		//okay body
		return client, nil
	}
}
