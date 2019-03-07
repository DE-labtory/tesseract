/*
 * Copyright 2018 DE-labtory
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

package main

import (
	"os"
	"time"

	"github.com/DE-labtory/tesseract"
	"github.com/DE-labtory/tesseract/docker"
)

func main() {

	volume := os.Args[1]
	if volume == "" {
		panic("Volume name is missing")
	}

	testGolangImg := tesseract.ContainerImage{
		Name: "golang",
		Tag:  "1.9",
	}

	// when
	res, err := docker.CreateContainer(
		tesseract.ContainerConfig{
			Name:           "container_mock",
			ContainerImage: testGolangImg,
			HostIp:         "127.0.0.1",
			Port:           "50002",
			//StartCmd:       []string{"go", "run", "/go/src/github.com/DE-labtory/tesseract/mock/test-volume/main.go"},
			StartCmd: []string{"sleep", "1000"},
			Network:  nil,
			Mount: []string{
				volume + ":" + "/go/src",
			},
		},
	)

	if err != nil {
		panic(err)
	}

	err = docker.StartContainer(res)
	if err != nil {
		panic(err)
	}

	time.Sleep(60 * time.Second)

	if _, err := os.Stat("/go/src/test"); os.IsNotExist(err) {
		panic("fail to bind volume")
	}
}
