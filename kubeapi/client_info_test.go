package kubeapi_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var once sync.Once

type ClientInfo struct {
	Host                string
	ResourceProcessFlag bool
	Enabled             bool
}

func (c *ClientInfo) SetEnabled(enabled bool) {
	c.Enabled = enabled
}

func GetClientInfoArray() []*ClientInfo {
	size := 2
	host_arr := []string{"192.168.0.138", "192.168.0.140"}

	var clientInfoArray []*ClientInfo
	once.Do(func() {
		clientInfoArray = make([]*ClientInfo, size)
		for i := 0; i < size; i++ {
			clientInfoArray[i] = &ClientInfo{
				Host:                host_arr[i],
				ResourceProcessFlag: false,
				Enabled:             true,
			}
		}
	})

	return clientInfoArray
}

func TestClientInfoSingleton(t *testing.T) {
	clientInfoArray := GetClientInfoArray()
	clientInfoArray[1].SetEnabled(false)

	assert.Equal(t, clientInfoArray[0].Enabled, true)
	assert.Equal(t, clientInfoArray[1].Enabled, false)
}
