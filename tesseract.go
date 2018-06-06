package tesseract

import (
	"net"

	"strconv"

	"github.com/it-chain/tesseract/docker"
	"github.com/it-chain/tesseract/rpc"
	"github.com/pkg/errors"
)

type Tesseract struct {
	Config  Config
	clients map[string]rpc.Client
}

type Config struct {
	shPath string
}

type ICodeInfo struct {
	Name        string
	Directory   string
	DockerImage docker.Image
	language    string // ENUM 으로 대체하면 좋음
}

func NewTesseract(c Config) *Tesseract {
	return &Tesseract{
		Config:  c,
		clients: make(map[string]rpc.Client),
	}
}

var ErrFailedPullImage = errors.New("failed to pull image")
var defaultPort = "50001"

// Deploy create Docker Container with running ShimCode and copying SmartContract.
func (t *Tesseract) SetupContainer(iCodeInfo ICodeInfo) error {

	// Todo : port 선정 기준은? (포트 번호 생성 함수 필요?)
	if iCodeInfo.DockerImage.Name == "" {
		iCodeInfo.DockerImage.Name = docker.DefaultImageName
		iCodeInfo.DockerImage.Tag = docker.DefaultImageTag
	}

	var port string
	var err error

	if port, err = getAvailablePort(); err != nil {
		return err
	}

	if err := pullImage(iCodeInfo.DockerImage.GetFullName()); err != nil {
		return ErrFailedPullImage
	}

	// Create Docker
	res, err := docker.CreateContainerWithCellCode(
		docker.Image{Name: docker.DefaultImageName, Tag: docker.DefaultImageTag},
		iCodeInfo.Directory,
		t.Config.shPath,
		port,
	)

	if err != nil {
		return err
	}

	// StartContainer
	err = docker.StartContainer(res)
	if err != nil {
		return err
	}

	// Get Container handler

	return nil
}

//1씩 증가 시키며 port를 확인한다
func getAvailablePort() (string, error) {

	for {
		lis, err := net.Listen("tcp", defaultPort)

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
