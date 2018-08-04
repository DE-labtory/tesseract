package tesseract

import (
	"context"
	"errors"
	"net"
	"strconv"
	"time"

	"log"

	"fmt"

	"github.com/it-chain/tesseract/docker"
	"github.com/it-chain/tesseract/pb"
	"github.com/it-chain/tesseract/rpc"
)

type ContainerID = string
type CallBack func(response *pb.Response, err error)

type Tesseract struct {
	Clients map[ContainerID]*rpc.ClientStream
}

type Request struct {
	Uuid     string
	TypeName string
	FuncName string
	Args     []string
}

type ICodeInfo struct {
	Name        string
	Directory   string
	DockerImage docker.Image
	language    string // ENUM 으로 대체하면 좋음
}

func New() *Tesseract {
	return &Tesseract{
		Clients: make(map[string]*rpc.ClientStream),
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
	res, err := docker.CreateContainer(
		docker.Image{Name: docker.DefaultImageName, Tag: docker.DefaultImageTag},
		iCodeInfo.Directory,
		port,
	)

	if err != nil {
		return "", err
	}
	fmt.Println("Container create finish")
	// StartContainer
	err = docker.StartContainer(res)

	if err != nil {
		return "", err
	}

	client, err := createClient()
	if err != nil {
		docker.CloseContainer(res.ID)
		return "", err
	}

	t.Clients[res.ID] = client

	return res.ID, nil
}

//1씩 증가 시키며 port를 확인한다
func getAvailablePort() (string, error) {
	portList, err := docker.GetUsingPorts()
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
			if portNumber == portInfo.PublicPort || portNumber == portInfo.PrivatePort {
				fmt.Println("already using port! : ", string(portNumber))
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

	//todo need to remove
	//todo client connect retry or ip check
	//todo maybe need ping operation
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
				log.Println("successfully connected")
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

func (t *Tesseract) Request(containerID string, req Request, callback CallBack) error {
	// Args : Transaction
	// Get Container handler using SmartContract ID
	// Send Query or Invoke massage
	// Receive result
	// Return result

	client := t.Clients[containerID]
	err := client.RunICode(&pb.Request{
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

func (t *Tesseract) StopContainers() {

	for name, client := range t.Clients {
		log.Println("rpc client [%s] is closing", name)
		client.Close()
		delete(t.Clients, name)
		docker.CloseContainer(name)
	}
}

func (t *Tesseract) StopContainerById(id ContainerID) error {
	client := t.Clients[id]
	if client == nil {
		return errors.New(fmt.Sprintf("no container with id : %s", id))
	}
	client.Close()
	delete(t.Clients, id)
	return docker.CloseContainer(id)
}
