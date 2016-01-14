package WhatIsMyIp_test

import (
	"github.com/philwinder/WhatIsMyIp"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestValidIp(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "127.0.0.1:8080",
	}
	res, _ := WhatIsMyIp.GetIp(req)
	assert.True(t, res.Ip == "127.0.0.1", "Ip address not parsed")
	assert.True(t, res.Port == "8080", "Port not parsed")
}

func TestInvalidIp(t *testing.T) {
	req := &http.Request{
		RemoteAddr: "1aa.43.fss:8080",
	}
	res, err := WhatIsMyIp.GetIp(req)
	assert.True(t, res == nil, "Ip address not parsed")
	assert.False(t, err == nil, "There should have been an error")
}
