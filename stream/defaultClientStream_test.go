package stream

import (
	"fmt"
	"testing"

	"github.com/it-chain/tesseract/pb"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultClientStream(t *testing.T) {
	cs := NewDefaultClientStream("127.0.0.1", "50001")
	fmt.Println(cs)
}

func TestConnect(t *testing.T) {
	cs := NewDefaultClientStream("127.0.0.1", "50001")
	err := cs.Connect()

	assert.NoError(t, err)
}

func TestSendRequest(t *testing.T) {
	cs := NewDefaultClientStream("127.0.0.1", "50001")
	err := cs.Connect()
	assert.NoError(t, err)

	cs.SendRequest(&pb.Request{"testetset"})
}
