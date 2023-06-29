package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IpCheck(t *testing.T) {
	var CampusIPs = []string{
		"10.10.10.0/8",
		"192.168.0.0/16",
		"172.16.0.0/12",
		"202.120.92.60",
	}
	assert.True(t, IPCheck("10.100.100.1", CampusIPs))
	assert.True(t, IPCheck("192.168.100.1", CampusIPs))
	assert.True(t, IPCheck("172.24.100.1", CampusIPs))
	assert.True(t, IPCheck("202.120.92.60", CampusIPs))
	assert.False(t, IPCheck("114.114.114.114", CampusIPs))
	assert.False(t, IPCheck("172.12.114.114", CampusIPs))
}
