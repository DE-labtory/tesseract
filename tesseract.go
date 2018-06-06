package tesseract

import (
	"net"
	"time"

	"strconv"

	"log"

	"github.com/it-chain/tesseract/docker"
	"github.com/it-chain/tesseract/rpc"
	"github.com/pkg/errors"
)

type ContainerID = string

type Tesseract struct {
	Config  Config
	Clients map[ContainerID]rpc.Client
}

type Config struct {
	ShPath string
}

type ICodeInfo struct {
	Name        string
	Directory   string
	DockerImage docker.Image
	language    string // ENUM 으로 대체하면 좋음
}

func New(c Config) *Tesseract {
	return &Tesseract{
		Config:  c,
		Clients: make(map[string]rpc.Client),
	}
}

var ErrFailedPullImage = errors.New("failed to pull image")
var defaultPort = "50001"
var ipAddress = "127.0.0.1"

// Deploy create Docker Container with running ShimCode and copying SmartContract.
// todo sh를 실행시키는데 시간이 많이 걸려서 client connect전에 time을 걸어놓았음 다르게 처리할 방법이 필요함
func (t *Tesseract) SetupContainer(iCodeInfo ICodeInfo) (string, error) {

	// Todo : port 선정 기준은? (포트 번호 생성 함수 필요?)
	if iCodeInfo.DockerImage.Name == "" {
		iCodeInfo.DockerImage.Name = docker.DefaultImageName
		iCodeInfo.DockerImage.Tag = docker.DefaultImageTag
	}

	var port string
	var err error

	if port, err = getAvailablePort(); err != nil {
		return "", err
	}

	if err = pullImage(iCodeInfo.DockerImage.GetFullName()); err != nil {
		return "", ErrFailedPullImage
	}

	// Create Docker
	res, err := docker.CreateContainerWithCellCode(
		docker.Image{Name: docker.DefaultImageName, Tag: docker.DefaultImageTag},
		iCodeInfo.Directory,
		t.Config.ShPath,
		port,
	)

	if err != nil {
		return "", err
	}

	// StartContainer
	err = docker.StartContainer(res)

	if err != nil {
		return "", err
	}

	client, err := createClient(res.ID)

	if err != nil {
		return "", err
	}

	t.Clients[res.ID] = client

	return res.ID, nil
}

func createClient(containerID ContainerID) (rpc.Client, error) {

	//todo need to remove
	//todo client connect retry or ip check
	//todo maybe need ping operation
	time.Sleep(60 * time.Second)

	client, err := rpc.Connect(ipAddress + ":" + defaultPort)

	if err != nil {
		return nil, err
	}

	return client, nil
}

//1씩 증가 시키며 port를 확인한다
func getAvailablePort() (string, error) {

	for {
		lis, err := net.Listen("tcp", "127.0.0.1:"+defaultPort)

		if err == nil {
			lis.Close()
			return defaultPort, nil
		}

		portNumber, err := strconv.Atoi(defaultPort)

		if err != nil {
			return "", err
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

func (t *Tesseract) QueryOrInvoke() {
	// args : Transaction
	// Get Container handler using SmartContract ID
	// Send Query or Invoke massage
	// Receive result
	// Return result
}

func (t *Tesseract) StopContainer() {

	for name, client := range t.Clients {
		log.Println("rpc client [%s] is closing", name)
		client.Close()
	}
}
