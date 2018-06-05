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
	DockerImage docker.Image
	language    string // ENUM 으로 대체하면 좋음
}

func NewTesseract(c Config) *Tesseract {
	return &Tesseract{Config: c}
}

// Deploy create Docker Container with running ShimCode and copying SmartContract.
func (t *Tesseract) SetupContainer(iCodeInfo ICodeInfo) error {
	// Todo : port 선정 기준은? (포트 번호 생성 함수 필요?)

	if iCodeInfo.DockerImage.Name == "" {
		iCodeInfo.DockerImage.Name = docker.DefaultImageName
		iCodeInfo.DockerImage.Tag = docker.DefaultImageTag
	}

	port := "50001"

	// Docker IMAGE pull
	r, err := docker.HasImage(iCodeInfo.DockerImage.GetFullName())
	if err != nil {
		return err
	}
	if !r {
		docker.PullImage(iCodeInfo.DockerImage.GetFullName())
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

func (t *Tesseract) QueryOrInvoke() {
	// args : Transaction
	// Get Container handler using SmartContract ID
	// Send Query or Invoke massage
	// Receive result
	// Return result
}
