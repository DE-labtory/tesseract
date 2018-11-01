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
	"errors"
	"net"
	"strconv"
	"time"

	"docker.io/go-docker/api/types"
	"github.com/it-chain/iLogger"
	"github.com/it-chain/tesseract"
	"github.com/it-chain/tesseract/docker"
	"github.com/it-chain/tesseract/rpc"
)

var ErrFailedPullImage = errors.New("failed to pull image")
var defaultPort = "50001"

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

	if config.Network == nil {
		var port string
		if port, err = getAvailablePort(); err != nil {
			return DockerContainer{}, err
		}
		config.Port = port
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

func getAvailablePort() (string, error) {
	portList, err := docker.GetPorts()
	if err != nil {
		return "", err
	}

findLoop:
	for {
		portNumber, err := strconv.Atoi(defaultPort)
		if err != nil {
			return "", err
		}
		for _, portInfo := range portList {
			if portNumber == int(portInfo.PublicPort) || portNumber == int(portInfo.PrivatePort) {
				portNumber++
				defaultPort = strconv.Itoa(portNumber)
				continue findLoop
			}
		}

		lis, err := net.Listen("tcp", "127.0.0.1:"+defaultPort)

		if err == nil {
			lis.Close()
			return defaultPort, nil
		}

		portNumber++
		defaultPort = strconv.Itoa(portNumber)
	}
}
