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

	"github.com/it-chain/tesseract"
	"github.com/it-chain/tesseract/docker"
	"github.com/it-chain/tesseract/logger"
	"github.com/it-chain/tesseract/rpc"
)

var ErrFailedPullImage = errors.New("failed to pull image")
var defaultPort = "50001"
var ipAddress = "127.0.0.1"

func Create(config tesseract.ContainerConfig) (DockerContainer, error) {

	logger.Info(nil, "[Tesseract] creating container")
	containerImage := tesseract.GetDefaultImage()

	var port string
	var err error

	if port, err = getAvailablePort(); err != nil {
		return DockerContainer{}, err
	}

	if err = pullImage(containerImage.GetFullName()); err != nil {
		return DockerContainer{}, ErrFailedPullImage
	}

	if config.LogFileName == "" {
		config.LogFileName = tesseract.DefaultLogFileName
	}

	// Create Docker
	res, err := docker.CreateContainer(
		containerImage,
		config.Directory,
		port,
		normalizeLogFileName(config.LogFileName),
	)

	if err != nil {
		return DockerContainer{}, err
	}

	err = docker.StartContainer(res)

	if err != nil {
		logger.Errorf(nil, "[Tesseract] fail to create container: %s", err.Error())
		docker.RemoveContainer(res.ID)
		return DockerContainer{}, err
	}

	client, err := createClient()

	if err != nil {
		logger.Errorf(nil, "[Tesseract] closing container %d", res.ID)
		docker.KillContainer(res.ID)
		docker.RemoveContainer(res.ID)
		return DockerContainer{}, err
	}

	return NewDockerContainer(res.ID, client, port), nil
}

//1씩 증가 시키며 port를 확인한다
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

func createClient() (*rpc.ClientStream, error) {
	return retryConnectWithTimeOut(120 * time.Second)
}

func retryConnectWithTimeOut(timeout time.Duration) (*rpc.ClientStream, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c := make(chan *rpc.ClientStream, 1)

	go func() {

		ticker := time.NewTicker(2 * time.Second)

		for _ = range ticker.C {
			client, err := rpc.NewClientStream(ipAddress + ":" + defaultPort)

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

func normalizeLogFileName(name string) string {
	if name == "" {
		return tesseract.DefaultLogFileName
	}
	return name
}