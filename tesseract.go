package tesseract

import (
	"github.com/it-chain/tesseract/docker"
)

type Tesseract struct {
	Config Config
}

type Config struct {
	shPath string
}

type ICodeInfo struct {
	Name        string
	Directory   string
	DockerImage docker.DockerImage
	language    string // ENUM 으로 대체하면 좋음
}

func NewTesseract(c Config) *Tesseract {
	return &Tesseract{Config: c}
}

// Deploy create Docker Container with running ShimCode and copying SmartContract.
func (t *Tesseract) SetupContainer(iCodeInfo ICodeInfo) error {
	// Todo : ImageName과 ImageTag를 parameter로 받아와야 함 (ICodeInfo에 포함될 수도 있음)
	// Todo : shPath는 어디서 받아올 것인가? (config 쪽으로, default는 있고)
	// Todo : port 선정 기준은? (포트 번호 생성 함수 필요?)

	port := "50001"

	// Docker IMAGE pull
	r, err := docker.HasImage(docker.DefaultImageName + ":" + docker.DefaultImageTag)
	if err != nil {
		return err
	}
	if r {
		docker.PullImage(docker.DefaultImageName + ":" + docker.DefaultImageTag)
	}

	// Create Docker
	res, err := docker.CreateContainerWithCellCode(
		docker.DockerImage{docker.DefaultImageName, docker.DefaultImageTag},
		iCodeInfo,
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

	// (Connect socket)
	// Get Container handler

	return nil
}

func (t *Tesseract) QueryOrInvoke() {
	// args : Transaction
	// Get Container handler using SmartContract ID
	// Send Query or Invoke massage
	// Receive result
	// Return result
}
